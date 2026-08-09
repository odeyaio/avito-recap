package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"avito-recap/internal/adapter/in/http/generated"
	"avito-recap/internal/catalog"
	"avito-recap/internal/engine"
	"avito-recap/internal/model"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

var (
	errPrimaryBehaviorMissing = errors.New("primary behavior is missing")
	errNextActionMissing      = errors.New("next action is missing")
)

// categoryMetric is the decode-side counterpart of engine's categoryStatJSON
// (see internal/engine/serialize.go) — connected only via matching JSON
// tags, not a shared Go type, to keep the http adapter decoupled from the
// engine's internal representation.
type categoryMetric struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Actions      int64   `json:"actions"`
	ActiveMonths int64   `json:"active_months"`
	Share        float64 `json:"share"`
}

type snapshotMetrics struct {
	Activity struct {
		Views                   int64   `json:"views"`
		UniqueListingsViewed    int64   `json:"unique_listings_viewed"`
		Searches                int64   `json:"searches"`
		ActiveDays              int32   `json:"active_days"`
		ActiveMonths            int32   `json:"active_months"`
		LongestActiveStreakDays int32   `json:"longest_active_streak_days"`
		MostActiveMonth         string  `json:"most_active_month"`
		FavoriteHour            int32   `json:"favorite_hour"`
		MonthlyActions          [12]int64 `json:"monthly_actions"`
	} `json:"activity"`
	Interests struct {
		TopCategories          []categoryMetric `json:"top_categories"`
		NewCategories          []categoryMetric `json:"new_categories"`
		MostConsistentCategory *categoryMetric  `json:"most_consistent_category"`
		PreferredPriceBand     string           `json:"preferred_price_band"`
	} `json:"interests"`
	Intent struct {
		RepeatViews             int64   `json:"repeat_viewed_listings"`
		FavoritesAdded          int64   `json:"favorites_added"`
		ActiveFavorites         int64   `json:"active_favorites"`
		Contacts                int64   `json:"contacts"`
		CompletedDeals          int64   `json:"completed_deals"`
		ContactToDealConversion float64 `json:"contact_to_deal_conversion"`
		CancelledAfterContact   int64   `json:"cancelled_after_contact"`
	} `json:"intent"`
	Marketplace struct {
		Purchases         int64 `json:"purchases"`
		Sales             int64 `json:"sales"`
		DeliveryDeals     int64 `json:"delivery_deals"`
		PublishedListings int64 `json:"published_listings"`
		ClosedListings    int64 `json:"closed_listings"`
		ListingViews      int64 `json:"listing_views"`
		ListingContacts   int64 `json:"listing_contacts"`
		ListingEdits      int64 `json:"listing_edits"`
	} `json:"marketplace"`
	Community struct {
		ReviewsLeft     int64   `json:"reviews_left"`
		ReviewsReceived int64   `json:"reviews_received"`
		FiveStarRatings int64   `json:"five_star_ratings"`
		AverageRating   float64 `json:"average_rating"`
	} `json:"community"`
	Features struct {
		NotificationOpens int64 `json:"notification_opens"`
		PromotionUses     int64 `json:"promotion_uses"`
		ListingsShared    int64 `json:"listings_shared"`
	} `json:"features"`
}

func profileResponse(profile model.ProfileSummary) generated.ProfileSummary {
	return generated.ProfileSummary{
		ID:             profile.User.ID,
		DisplayName:    profile.User.DisplayName,
		Region:         profile.User.Region,
		RegisteredAt:   profile.User.RegisteredAt,
		AvailableYears: profile.AvailableYears,
		LatestRecapID:  profile.LatestRecapID,
	}
}

func recapResponse(recap model.Recap) (generated.Recap, error) {
	metrics, err := metricsResponse(recap.Snapshot.Metrics)
	if err != nil {
		return generated.Recap{}, err
	}

	primary, traits, err := behaviorResponse(recap.Behaviors)
	if err != nil {
		return generated.Recap{}, err
	}
	achievements, cards, err := achievementResponse(recap.Achievements)
	if err != nil {
		return generated.Recap{}, err
	}
	nextAction, err := nextActionResponse(recap.NextAction, recap.Behaviors)
	if err != nil {
		return generated.Recap{}, err
	}

	shareCard := shareCardResponse(primary, achievements)

	year := recap.Snapshot.PeriodStart.Year()
	return generated.Recap{
		ID:      recap.Snapshot.ID,
		Profile: profileResponse(recap.Profile),
		Period: generated.RecapPeriod{
			Year:     year,
			StartsAt: openapitypes.Date{Time: recap.Snapshot.PeriodStart},
			EndsAt:   openapitypes.Date{Time: recap.Snapshot.PeriodEnd},
		},
		Behavior:     generated.BehaviorProfile{Primary: primary, Traits: traits},
		Metrics:      metrics,
		Achievements: achievements,
		Story: generated.RecapStory{
			Headline: "Ваш тип года — «" + primary.Name + "»",
			Summary:  primary.Description,
			Cards:    cards,
		},
		NextAction:  nextAction,
		ShareCard:   &shareCard,
		GeneratedAt: recap.Snapshot.GeneratedAt,
	}, nil
}

// shareCardResponse builds the compact, safe-to-share summary of a recap.
// Only already-public/shareable-flagged data goes in here: the primary
// behavior's own name/description (always shown to the user anyway) and, if
// any achievement was marked shareable by the catalog, its icon. Previously
// nothing populated generated.Recap.ShareCard at all — the field existed in
// the contract but recapResponse never set it, so API responses always had
// shareCard: null.
func shareCardResponse(primary generated.BehaviorMatch, achievements []generated.Achievement) generated.ShareCard {
	card := generated.ShareCard{
		Title:    fmt.Sprintf("Мой тип года — «%s»", primary.Name),
		Subtitle: primary.Description,
	}
	for _, achievement := range achievements {
		if !achievement.Shareable {
			continue
		}
		imageURL := achievement.Image.URL
		card.ImageURL = &imageURL
		break
	}
	return card
}

func metricsResponse(raw json.RawMessage) (generated.RecapMetrics, error) {
	var value snapshotMetrics
	if err := json.Unmarshal(raw, &value); err != nil {
		return generated.RecapMetrics{}, err
	}

	activity := generated.ActivityMetrics{
		Views:                   value.Activity.Views,
		UniqueListingsViewed:    value.Activity.UniqueListingsViewed,
		Searches:                value.Activity.Searches,
		ActiveDays:              value.Activity.ActiveDays,
		ActiveMonths:            value.Activity.ActiveMonths,
		LongestActiveStreakDays: value.Activity.LongestActiveStreakDays,
	}
	if value.Activity.MostActiveMonth != "" {
		activity.MostActiveMonth = &value.Activity.MostActiveMonth
	}
	activity.FavoriteHour = &value.Activity.FavoriteHour

	intent := generated.IntentMetrics{
		RepeatViews:     value.Intent.RepeatViews,
		FavoritesAdded:  value.Intent.FavoritesAdded,
		ActiveFavorites: value.Intent.ActiveFavorites,
		Contacts:        value.Intent.Contacts,
		CompletedDeals:  value.Intent.CompletedDeals,
	}
	intent.ContactToDealConversion = &value.Intent.ContactToDealConversion

	community := generated.CommunityMetrics{
		ReviewsLeft:     value.Community.ReviewsLeft,
		ReviewsReceived: value.Community.ReviewsReceived,
		FiveStarRatings: value.Community.FiveStarRatings,
	}
	if value.Community.ReviewsReceived > 0 {
		community.AverageRating = &value.Community.AverageRating
	}

	interests := generated.InterestMetrics{
		TopCategories: categoryMetricsResponse(value.Interests.TopCategories),
		NewCategories: categoryMetricsResponse(value.Interests.NewCategories),
	}
	if value.Interests.MostConsistentCategory != nil {
		consistent := categoryMetricResponse(*value.Interests.MostConsistentCategory)
		interests.MostConsistentCategory = &consistent
	}

	return generated.RecapMetrics{
		Activity:  activity,
		Interests: interests,
		Intent:    intent,
		Marketplace: generated.MarketplaceMetrics{
			Purchases: value.Marketplace.Purchases, Sales: value.Marketplace.Sales,
			DeliveryDeals: value.Marketplace.DeliveryDeals, PublishedListings: value.Marketplace.PublishedListings,
			ClosedListings: value.Marketplace.ClosedListings, ListingViews: value.Marketplace.ListingViews,
			ListingContacts: value.Marketplace.ListingContacts,
		},
		Community: community,
		Features: generated.FeatureMetrics{
			NotificationOpens: value.Features.NotificationOpens,
			PromotionUses:     value.Features.PromotionUses,
		},
	}, nil
}

func categoryMetricsResponse(items []categoryMetric) []generated.CategoryMetric {
	result := make([]generated.CategoryMetric, 0, len(items))
	for _, item := range items {
		result = append(result, categoryMetricResponse(item))
	}
	return result
}

func categoryMetricResponse(item categoryMetric) generated.CategoryMetric {
	share := item.Share
	return generated.CategoryMetric{
		Code:         item.Code,
		Name:         item.Name,
		Actions:      item.Actions,
		ActiveMonths: int32(item.ActiveMonths),
		Share:        &share,
	}
}

func behaviorResponse(
	behaviors []model.StoredBehavior,
) (generated.BehaviorMatch, []generated.BehaviorMatch, error) {
	traits := make([]generated.BehaviorMatch, 0, len(behaviors))
	var primary generated.BehaviorMatch
	primaryFound := false
	for _, behavior := range behaviors {
		evidence, err := evidenceResponse(behavior.Match.Evidence)
		if err != nil {
			return generated.BehaviorMatch{}, nil, err
		}
		value := generated.BehaviorMatch{
			Code: behavior.Definition.Code, Name: behavior.Definition.Name,
			Description: behavior.Definition.Description, Explanation: explanation(evidence),
			Score: behavior.Match.Score, Evidence: evidence,
		}
		if behavior.Match.IsPrimary {
			primary = value
			primaryFound = true
			continue
		}
		traits = append(traits, value)
	}
	if !primaryFound {
		return generated.BehaviorMatch{}, nil, errPrimaryBehaviorMissing
	}
	return primary, traits, nil
}

func achievementResponse(
	achievements []model.StoredAchievement,
) ([]generated.Achievement, []generated.StoryCard, error) {
	values := make([]generated.Achievement, 0, len(achievements))
	cards := make([]generated.StoryCard, 0, len(achievements))
	for _, achievement := range achievements {
		evidence, err := evidenceResponse(achievement.Match.Evidence)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, generated.Achievement{
			Code: achievement.Definition.Code, Name: achievement.Definition.Name,
			Description: achievement.Definition.Description, Explanation: explanation(evidence),
			Image: generated.AchievementImage{
				URL: achievement.Definition.IconKey,
				Alt: achievement.Definition.Name,
			},
			AchievedAt: achievement.Match.AchievedAt,
			Shareable:  achievement.Match.IsShareable,
		})
		metricCodes := make([]string, 0, len(evidence))
		for _, item := range evidence {
			metricCodes = append(metricCodes, item.MetricCode)
		}
		cards = append(cards, generated.StoryCard{
			ID: achievement.Definition.Code, Kind: "achievement",
			Title: achievement.Definition.Name, Text: achievement.Definition.Description,
			MetricCodes: metricCodes, Shareable: achievement.Match.IsShareable,
		})
	}
	return values, cards, nil
}

func nextActionResponse(
	action *model.RecapNextAction,
	behaviors []model.StoredBehavior,
) (generated.NextAction, error) {
	if action == nil {
		return generated.NextAction{}, errNextActionMissing
	}
	var target generated.ActionTarget
	if err := json.Unmarshal(action.Target, &target); err != nil {
		return generated.NextAction{}, err
	}

	var definition catalog.DefaultAction
	for _, behavior := range behaviors {
		if !behavior.Match.IsPrimary {
			continue
		}
		if err := json.Unmarshal(behavior.Definition.DefaultAction, &definition); err != nil {
			return generated.NextAction{}, err
		}
		break
	}
	return generated.NextAction{
		Code: action.Code, Title: definition.Title, Text: definition.Title,
		Href: action.Href, Target: target,
	}, nil
}

func evidenceResponse(raw json.RawMessage) ([]generated.Evidence, error) {
	var stored []engine.Evidence
	if err := json.Unmarshal(raw, &stored); err != nil {
		return nil, err
	}
	result := make([]generated.Evidence, 0, len(stored))
	for _, item := range stored {
		description := item.MetricCode
		if item.Actual != nil {
			description = fmt.Sprintf("%s: %v", item.MetricCode, item.Actual)
		}
		result = append(result, generated.Evidence{MetricCode: item.MetricCode, Description: description})
	}
	return result, nil
}

func explanation(evidence []generated.Evidence) string {
	parts := make([]string, 0, len(evidence))
	for _, item := range evidence {
		parts = append(parts, item.Description)
	}
	return strings.Join(parts, "; ")
}
