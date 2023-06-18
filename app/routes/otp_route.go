package routes

import (
	"database/sql"
	"ppob/app/controllers"
	"ppob/app/repository"
	"ppob/app/services"

	"github.com/labstack/echo/v4"
)

func OtpRoute(db *sql.DB, app *echo.Group) *echo.Group {
	otp := app.Group("/otp");
	otpRepo := repository.NewOtpRepository(db);
	otpService := services.NewOtpService(otpRepo);
	otpController := controllers.NewOtpController(db, otpService);

	otp.POST("", otpController.RequestOtp);
	otp.POST("/verification", otpController.VerificationOtp);

	return otp;
}