package routes

import "github.com/labstack/echo/v4"

func Version1Route(app *echo.Echo) *echo.Group {
	version := app.Group("/api/v1");

	return version;
}