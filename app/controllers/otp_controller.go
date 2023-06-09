package controllers

import (
	"database/sql"
	"net/http"
	"ppob/app/configs"
	"ppob/app/domain"
	"ppob/app/services"

	"github.com/labstack/echo/v4"
)

type otpController struct {
	otpService services.OtpService
	db *sql.DB
}

func NewOtpController(db *sql.DB, otpService services.OtpService) *otpController {
	return &otpController{
		db: db,
		otpService: otpService,
	}
}

func (controller *otpController) RequestOtp(ctx echo.Context) error {

	payload := new(domain.RequestOtpPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 7581,
			DB: controller.db,
			Type: 2,
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			DB: controller.db,
			Type: 2,
		}
		return response.ResponseMiddleware(ctx);
	}

	headerData := &domain.RequestOtpHeader{
		DeviceID: ctx.Request().Header.Get("device_id"),
	}

	status, code, data, err := controller.otpService.CreateOTP(payload, headerData);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			Error: err.Error(),
			DB: controller.db,
			Type: 2,
		}
		return response.ResponseMiddleware(ctx);
	}

	if code != 0 {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Type: 2,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: http.StatusOK,
		Code: 0,
		Data: data,
		DB: controller.db,
		Type: 2,
	}
	return response.ResponseMiddleware(ctx);
}

func (controller *otpController) VerificationOtp(ctx echo.Context) error {

	payload := new(domain.VerificationOtpPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 7581,
			DB: controller.db,
			Type: 3,
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 6708,
			DB: controller.db,
			Type: 3,
		}
		return response.ResponseMiddleware(ctx);
	}

	headerData := &domain.RequestOtpHeader{
		DeviceID: ctx.Request().Header.Get("device_id"),
	}

	status, code, err := controller.otpService.VerificationOtp(payload, headerData);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			Error: err.Error(),
			DB: controller.db,
			Type: 3,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Type: 3,
	}
	return response.ResponseMiddleware(ctx);
}