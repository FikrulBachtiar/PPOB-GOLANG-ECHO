package routes

import (
	"database/sql"
	"ppob/app/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitRoutes(db *sql.DB) *echo.Echo {
	app := echo.New();
	app.Use(middleware.HeaderContentType);

	v := validator.New()
	app.Validator = &middleware.PayloadValidator{Validator: v};

	version_route := Version1Route(app);
	validHeader := &middleware.HeaderValidator{DB: db};
	version_route.Use(validHeader.HeaderValidator);
	OnboardingRoute(db, version_route);
	OtpRoute(db, version_route);
	UserRoute(db, version_route);

	return app;
}