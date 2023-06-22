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


func (controller *onboardingController) Check(ctx echo.Context) error {

	payload := new(domain.CheckPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 7581,
			DB: controller.db,
			Type: 1,
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			DB: controller.db,
			Type: 1,
		}
		return response.ResponseMiddleware(ctx);
	}

	status, code, data, err := controller.onboardingService.CheckAccount(ctx.Request().Context(), payload);
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
			DB: controller.db,
			Type: 1,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		Data: data,
		DB: controller.db,
		Type: 1,
	}
	return response.ResponseMiddleware(ctx);
}

func (controller *onboardingController) Login(ctx echo.Context) error {

	payload := new(domain.LoginPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusInternalServerError,
			Code: 7581,
			Error: err.Error(),
			DB: controller.db,
			Type: 4,
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			DB: controller.db,
			Type: 4,
		}
		return response.ResponseMiddleware(ctx);
	}

	header := &domain.OnboardingHeader{
		DeviceID: ctx.Request().Header.Get("device_id"),
		OsName: ctx.Request().Header.Get("os_name"),
		OsVersion: ctx.Request().Header.Get("os_version"),
		DeviceModel: ctx.Request().Header.Get("device_model"),
		AppVersion: ctx.Request().Header.Get("app_version"),
		Longitude: ctx.Request().Header.Get("longitude"),
		Latitude: ctx.Request().Header.Get("latitude"),
		NotificationID: ctx.Request().Header.Get("notification_id"),
	}

	if err := ctx.Validate(header); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			DB: controller.db,
			Type: 4,
		}
		return response.ResponseMiddleware(ctx);
	}

	status, code, data, err := controller.onboardingService.LoginAccount(ctx.Request().Context(), payload, header);
	if err != nil {
		response := &configs.Response{
			Code: code,
			Status: status,
			Error: err.Error(),
			DB: controller.db,
			Type: 4,
		}
		return response.ResponseMiddleware(ctx);
	}

	if code != 0 {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Type: 4,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Data: data,
		Type: 4,
	}
	return response.ResponseMiddleware(ctx);
}

func (controller *onboardingController) Logout(ctx echo.Context) error {
	payload := new(domain.LogoutPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusInternalServerError,
			Code: 7581,
			Error: err.Error(),
			DB: controller.db,
			Type: 5,
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			DB: controller.db,
			Type: 5,
		}
		return response.ResponseMiddleware(ctx);
	}

	header := &domain.OnboardingHeader{
		DeviceID: ctx.Request().Header.Get("device_id"),
		OsName: ctx.Request().Header.Get("os_name"),
		OsVersion: ctx.Request().Header.Get("os_version"),
		DeviceModel: ctx.Request().Header.Get("device_model"),
		AppVersion: ctx.Request().Header.Get("app_version"),
		Longitude: ctx.Request().Header.Get("longitude"),
		Latitude: ctx.Request().Header.Get("latitude"),
		NotificationID: ctx.Request().Header.Get("notification_id"),
	}

	status, code, err := controller.onboardingService.LogoutAccount(ctx.Request().Context(), payload, header);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			Error: err.Error(),
			DB: controller.db,
			Type: 5,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Type: 5,
	}
	return response.ResponseMiddleware(ctx);
}