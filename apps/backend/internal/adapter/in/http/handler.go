package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"avito-recap/internal/adapter/in/http/generated"
	"avito-recap/internal/engine"
	"avito-recap/internal/repository"
	"avito-recap/internal/service"
)

const healthcheckTimeout = time.Second

type databasePinger interface {
	Ping(context.Context) error
}

type Handler struct {
	database databasePinger
	profiles profileService
	recaps   recapService
}

var _ generated.StrictServerInterface = (*Handler)(nil)

func NewHandler(database databasePinger, profiles profileService, recaps recapService) *Handler {
	return &Handler{database: database, profiles: profiles, recaps: recaps}
}

func RegisterHandlers(
	router generated.EchoRouter,
	database databasePinger,
	profiles profileService,
	recaps recapService,
) {
	handler := NewHandler(database, profiles, recaps)
	strictHandler := generated.NewStrictHandler(handler, nil)
	generated.RegisterHandlers(router, strictHandler)
}

func (h *Handler) ListProfiles(
	ctx context.Context,
	_ generated.ListProfilesRequestObject,
) (generated.ListProfilesResponseObject, error) {
	profiles, err := h.profiles.ListProfiles(ctx)
	if err != nil {
		problem := internalProblem()
		//nolint:nilerr // The service error is deliberately translated into a typed HTTP response.
		return generated.ListProfiles500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: generated.InternalErrorApplicationProblemPlusJSONResponse(problem),
		}, nil
	}

	items := make([]generated.ProfileSummary, 0, len(profiles))
	for _, profile := range profiles {
		items = append(items, profileResponse(profile))
	}
	return generated.ListProfiles200JSONResponse{Items: items}, nil
}

func (h *Handler) GenerateRecap(
	ctx context.Context,
	request generated.GenerateRecapRequestObject,
) (generated.GenerateRecapResponseObject, error) {
	if request.Body == nil {
		problem := newProblem(http.StatusBadRequest, "bad_request", "Некорректный запрос", "Тело запроса обязательно.")
		return generated.GenerateRecap400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: generated.BadRequestApplicationProblemPlusJSONResponse(problem),
		}, nil
	}

	recap, err := h.recaps.GenerateRecap(ctx, request.ProfileID, request.Body.Year)
	if err != nil {
		return generateRecapError(err), nil
	}
	response, err := recapResponse(recap)
	if err != nil {
		problem := internalProblem()
		//nolint:nilerr // The mapping error is deliberately translated into a typed HTTP response.
		return generated.GenerateRecap500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: generated.InternalErrorApplicationProblemPlusJSONResponse(problem),
		}, nil
	}
	return generated.GenerateRecap200JSONResponse(response), nil
}

func (h *Handler) GetRecap(
	ctx context.Context,
	request generated.GetRecapRequestObject,
) (generated.GetRecapResponseObject, error) {
	recap, err := h.recaps.GetRecap(ctx, request.RecapID)
	if errors.Is(err, repository.ErrRecapNotFound) {
		problem := newProblem(http.StatusNotFound, "recap_not_found", "Recap не найден", "Запрошенный recap не существует.")
		return generated.GetRecap404ApplicationProblemPlusJSONResponse{
			RecapNotFoundApplicationProblemPlusJSONResponse: generated.RecapNotFoundApplicationProblemPlusJSONResponse(problem),
		}, nil
	}
	if err != nil {
		problem := internalProblem()
		//nolint:nilerr // The service error is deliberately translated into a typed HTTP response.
		return generated.GetRecap500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: generated.InternalErrorApplicationProblemPlusJSONResponse(problem),
		}, nil
	}
	response, err := recapResponse(recap)
	if err != nil {
		problem := internalProblem()
		//nolint:nilerr // The mapping error is deliberately translated into a typed HTTP response.
		return generated.GetRecap500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: generated.InternalErrorApplicationProblemPlusJSONResponse(problem),
		}, nil
	}
	return generated.GetRecap200JSONResponse(response), nil
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

func internalProblem() generated.Problem {
	return newProblem(
		http.StatusInternalServerError,
		"internal_error",
		"Внутренняя ошибка сервиса",
		"Не удалось обработать запрос.",
	)
}

func generateRecapError(err error) generated.GenerateRecapResponseObject {
	switch {
	case errors.Is(err, repository.ErrProfileNotFound):
		problem := newProblem(http.StatusNotFound, "profile_not_found", "Профиль не найден", "Тестовый профиль не существует.")
		return generated.GenerateRecap404ApplicationProblemPlusJSONResponse{
			ProfileNotFoundApplicationProblemPlusJSONResponse: generated.ProfileNotFoundApplicationProblemPlusJSONResponse(problem),
		}
	case errors.Is(err, engine.ErrNoActivity), errors.Is(err, service.ErrBehaviorNotMatched):
		problem := newProblem(
			http.StatusUnprocessableEntity,
			"insufficient_activity",
			"Недостаточно активности",
			"Для выбранного года недостаточно данных для содержательного recap.",
		)
		return generated.GenerateRecap422ApplicationProblemPlusJSONResponse{
			InsufficientActivityApplicationProblemPlusJSONResponse: generated.InsufficientActivityApplicationProblemPlusJSONResponse(problem),
		}
	case errors.Is(err, repository.ErrCatalogUnavailable):
		problem := newProblem(
			http.StatusServiceUnavailable,
			"generation_unavailable",
			"Генерация временно недоступна",
			"Каталог правил ещё не загружен.",
		)
		return generated.GenerateRecap503ApplicationProblemPlusJSONResponse{
			GenerationUnavailableApplicationProblemPlusJSONResponse: generated.GenerationUnavailableApplicationProblemPlusJSONResponse(problem),
		}
	default:
		problem := internalProblem()
		return generated.GenerateRecap500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: generated.InternalErrorApplicationProblemPlusJSONResponse(problem),
		}
	}
}
