package controllers

import (
	"database/sql"
	"net/http"
	"ppob/app/configs"
	"ppob/app/domain"
	"ppob/app/services"

	"github.com/labstack/echo/v4"
)

type onboardingController struct {
	db *sql.DB
	onboardingService services.OnboardingService
}

func NewOnboardingController(db *sql.DB, onboardingService services.OnboardingService) *onboardingController {
	return &onboardingController{
		db: db,
		onboardingService: onboardingService,
	}
}


func (oc *onboardingController) Check(ctx echo.Context) error {

	payload := new(domain.CheckPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 7581,
			Message: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			Message: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	status, code, data, err := oc.onboardingService.CheckAccount(ctx.Request().Context(), payload);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	if code != 0 {
		response := &configs.Response{
			Status: status,
			Code: code,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		Data: data,
	}
	return response.ResponseMiddleware(ctx);
}