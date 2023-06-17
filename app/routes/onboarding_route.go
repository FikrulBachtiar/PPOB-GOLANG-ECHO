package routes

import (
	"database/sql"
	"ppob/app/controllers"
	"ppob/app/repository"
	"ppob/app/services"

	"github.com/labstack/echo/v4"
)

func OnboardingRoute(db *sql.DB, app *echo.Group) *echo.Group {
	onboard := app.Group("/onboard");
	onboardRepo := repository.NewOnboardingRepo(db);
	onboardService := services.NewOnboardingService(onboardRepo);
	onboardController := controllers.NewOnboardingController(db, onboardService);

	onboard.POST("/check", onboardController.Check);
	onboard.POST("/login", onboardController.Login);

	return onboard;
}