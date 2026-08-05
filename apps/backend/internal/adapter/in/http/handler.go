package http

import (
	"context"
	"net/http"
	"time"

	"avito-recap/internal/adapter/in/http/generated"
)

const healthcheckTimeout = time.Second

type databasePinger interface {
	Ping(context.Context) error
}

type Handler struct {
	database databasePinger
}

var _ generated.StrictServerInterface = (*Handler)(nil)

func NewHandler(database databasePinger) *Handler {
	return &Handler{database: database}
}

func RegisterHandlers(router generated.EchoRouter, database databasePinger) {
	handler := NewHandler(database)
	strictHandler := generated.NewStrictHandler(handler, nil)
	generated.RegisterHandlers(router, strictHandler)
}

func (h *Handler) ListProfiles(
	_ context.Context,
	_ generated.ListProfilesRequestObject,
) (generated.ListProfilesResponseObject, error) {
	return generated.ListProfiles200JSONResponse{
		Items: []generated.ProfileSummary{},
	}, nil
}

func (h *Handler) GenerateRecap(
	_ context.Context,
	_ generated.GenerateRecapRequestObject,
) (generated.GenerateRecapResponseObject, error) {
	problem := newProblem(
		http.StatusServiceUnavailable,
		"generation_unavailable",
		"Генерация recap пока недоступна",
		"Обработчик генерации ещё не подключён к use case.",
	)

	return generated.GenerateRecap503ApplicationProblemPlusJSONResponse{
		GenerationUnavailableApplicationProblemPlusJSONResponse: generated.GenerationUnavailableApplicationProblemPlusJSONResponse(problem),
	}, nil
}

func (h *Handler) GetRecap(
	_ context.Context,
	_ generated.GetRecapRequestObject,
) (generated.GetRecapResponseObject, error) {
	problem := newProblem(
		http.StatusNotFound,
		"recap_not_found",
		"Recap не найден",
		"Хранилище recap ещё не подключено.",
	)

	return generated.GetRecap404ApplicationProblemPlusJSONResponse{
		RecapNotFoundApplicationProblemPlusJSONResponse: generated.RecapNotFoundApplicationProblemPlusJSONResponse(problem),
	}, nil
}

func (h *Handler) CheckHealth(
	ctx context.Context,
	_ generated.CheckHealthRequestObject,
) (generated.CheckHealthResponseObject, error) {
	pingCtx, cancel := context.WithTimeout(ctx, healthcheckTimeout)
	defer cancel()

	if err := h.database.Ping(pingCtx); err != nil {
		problem := newProblem(
			http.StatusServiceUnavailable,
			"database_unavailable",
			"Сервис временно недоступен",
			"Не удалось подключиться к базе данных.",
		)

		return generated.CheckHealth503ApplicationProblemPlusJSONResponse{
			ServiceUnavailableApplicationProblemPlusJSONResponse: generated.ServiceUnavailableApplicationProblemPlusJSONResponse(problem),
		}, nil
	}

	return generated.CheckHealth200JSONResponse{Status: "ok"}, nil
}

func newProblem(status int32, code, title, detail string) generated.Problem {
	return generated.Problem{
		Type:   "https://avito-recap.example/problems/" + code,
		Title:  title,
		Status: status,
		Code:   code,
		Detail: &detail,
	}
}
