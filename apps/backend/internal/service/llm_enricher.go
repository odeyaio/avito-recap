package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"avito-recap/internal/engine"
)

type LLMConfig struct {
	Enabled bool          `env:"LLM_NEXT_ACTION_ENABLED" env-default:"false"`
	APIKey  string        `env:"LLM_API_KEY"`
	BaseURL string        `env:"LLM_BASE_URL" env-default:"https://api.groq.com/openai/v1"`
	Model   string        `env:"LLM_MODEL" env-default:"llama-3.3-70b-versatile"`
	Timeout time.Duration `env:"LLM_TIMEOUT" env-default:"8s"`
}

type LLMTextEnricher struct {
	config LLMConfig
	client *http.Client
}

func NewLLMTextEnricher(config LLMConfig) *LLMTextEnricher {
	return &LLMTextEnricher{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

const systemPrompt = `Ты формируешь персонализированный итог года для пользователя маркетплейса Avito.
Тебе дают JSON со сработавшими поведенческими сценариями, достижениями, топ-категориями и ключевыми метриками активности за год — это уже проверенные факты, не выдумывай ничего сверх них.
Ответь СТРОГО в формате JSON без markdown и пояснений: {"description": string, "actionText": string, "recommendedActionCode": string}.
"description" — 1-2 коротких дружелюбных предложения на русском, характеризующих пользователя на основе фактов.
"actionText" — 1 короткое предложение, персонализированная формулировка предложенного действия.
"recommendedActionCode" — один из кодов в поле allowedActionCodes, тот что лучше всего подходит; если не уверен, верни primaryBehavior.code.
Не придумывай цифры и категории, которых нет в предоставленных данных.`

type promptContext struct {
	PrimaryBehavior    promptBehavior      `json:"primaryBehavior"`
	OtherBehaviors     []promptBehavior    `json:"otherBehaviors,omitempty"`
	Achievements       []promptAchievement `json:"achievements,omitempty"`
	TopCategories      []promptCategory    `json:"topCategories,omitempty"`
	NewCategories      []promptCategory    `json:"newCategories,omitempty"`
	MonthlyActivity    [12]int64           `json:"monthlyActivity"`
	KeyMetrics         map[string]any      `json:"keyMetrics"`
	AllowedActionCodes []string            `json:"allowedActionCodes"`
}

type promptBehavior struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type promptAchievement struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type promptCategory struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Share float64 `json:"share"`
}

var promptMetricCodes = []string{
	"activity.total_actions",
	"activity.active_months",
	"activity.returned_after_long_gap",
	"interests.distinct_categories",
	"interests.top_category_share",
	"interests.preferred_price_band",
	"intent.contact_to_deal_conversion",
	"intent.completed_deals",
	"intent.cancelled_after_contact",
	"marketplace.sales",
	"marketplace.purchases",
	"marketplace.listing_edits",
	"community.average_rating",
	"community.reviews_left",
	"features.listings_shared",
}

func buildPromptContext(result engine.Result) (promptContext, bool) {
	primary, ok := primaryBehavior(result.Behaviors)
	if !ok {
		return promptContext{}, false
	}

	ctx := promptContext{
		PrimaryBehavior: promptBehavior{
			Code: primary.Definition.Code, Name: primary.Definition.Name, Description: primary.Definition.Description,
		},
		MonthlyActivity: result.MonthlyActivity,
		KeyMetrics:      make(map[string]any, len(promptMetricCodes)),
	}

	allowed := []string{primary.Definition.DefaultAction.Code}
	for _, match := range result.Behaviors {
		if match.IsPrimary {
			continue
		}
		ctx.OtherBehaviors = append(ctx.OtherBehaviors, promptBehavior{
			Code: match.Definition.Code, Name: match.Definition.Name, Description: match.Definition.Description,
		})
		allowed = append(allowed, match.Definition.DefaultAction.Code)
	}
	ctx.AllowedActionCodes = allowed

	for _, match := range result.Achievements {
		ctx.Achievements = append(ctx.Achievements, promptAchievement{
			Code: match.Definition.Code, Name: match.Definition.Name, Description: match.Definition.Description,
		})
	}

	for _, stat := range result.TopCategories {
		ctx.TopCategories = append(ctx.TopCategories, promptCategory{Code: stat.Code, Name: stat.Name, Share: stat.Share})
	}
	for _, stat := range result.NewCategories {
		ctx.NewCategories = append(ctx.NewCategories, promptCategory{Code: stat.Code, Name: stat.Name, Share: stat.Share})
	}

	for _, code := range promptMetricCodes {
		if value, exists := result.Metrics.Get(code); exists {
			ctx.KeyMetrics[code] = value
		}
	}

	return ctx, true
}

type chatCompletionRequest struct {
	Model          string            `json:"model"`
	Messages       []chatMessage     `json:"messages"`
	Temperature    float64           `json:"temperature"`
	MaxTokens      int               `json:"max_tokens"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

type modelOutput struct {
	Description           string `json:"description"`
	ActionText            string `json:"actionText"`
	RecommendedActionCode string `json:"recommendedActionCode"`
}

func (e *LLMTextEnricher) Enrich(ctx context.Context, result engine.Result) (string, string, error) {
	if !e.config.Enabled {
		return "", "", nil
	}
	if e.config.APIKey == "" {
		return "", "", fmt.Errorf("llm enricher: LLM_API_KEY is not set")
	}

	promptCtx, ok := buildPromptContext(result)
	if !ok {
		return "", "", fmt.Errorf("llm enricher: no primary behavior in result")
	}
	payload, err := json.Marshal(promptCtx)
	if err != nil {
		return "", "", fmt.Errorf("llm enricher: marshal prompt context: %w", err)
	}

	requestBody, err := json.Marshal(chatCompletionRequest{
		Model: e.config.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: string(payload)},
		},
		Temperature:    0.4,
		MaxTokens:      300,
		ResponseFormat: map[string]string{"type": "json_object"},
	})
	if err != nil {
		return "", "", fmt.Errorf("llm enricher: marshal request: %w", err)
	}

	requestCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(
		requestCtx, http.MethodPost, strings.TrimRight(e.config.BaseURL, "/")+"/chat/completions", bytes.NewReader(requestBody),
	)
	if err != nil {
		return "", "", fmt.Errorf("llm enricher: build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+e.config.APIKey)

	httpResponse, err := e.client.Do(httpRequest)
	if err != nil {
		return "", "", fmt.Errorf("llm enricher: call model: %w", err)
	}
	defer func() { _ = httpResponse.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(httpResponse.Body, 1<<20))
	if err != nil {
		return "", "", fmt.Errorf("llm enricher: read response: %w", err)
	}
	if httpResponse.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("llm enricher: unexpected status %d: %s", httpResponse.StatusCode, string(body))
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(body, &completion); err != nil {
		return "", "", fmt.Errorf("llm enricher: decode completion: %w", err)
	}
	if len(completion.Choices) == 0 {
		return "", "", fmt.Errorf("llm enricher: completion has no choices")
	}

	var output modelOutput
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &output); err != nil {
		return "", "", fmt.Errorf("llm enricher: decode model output: %w", err)
	}

	if !slices.Contains(promptCtx.AllowedActionCodes, output.RecommendedActionCode) {
		return "", "", fmt.Errorf("llm enricher: model returned action code %q not in allowlist", output.RecommendedActionCode)
	}

	return truncate(output.Description, 240), truncate(output.ActionText, 160), nil
}

func truncate(value string, limit int) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return strings.TrimSpace(string(runes[:limit])) + "…"
}
