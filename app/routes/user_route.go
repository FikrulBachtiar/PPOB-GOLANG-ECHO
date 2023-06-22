package routes

import (
	"database/sql"
	"ppob/app/controllers"
	"ppob/app/middleware"
	"ppob/app/repository"
	"ppob/app/services"

	"github.com/labstack/echo/v4"
)

func UserRoute(db *sql.DB, app *echo.Group) *echo.Group {
	user := app.Group("/user");
	headers := &middleware.HeaderValidator{DB: db}
	user.Use(headers.JWTValidation);
	userRepo := repository.NewUserRepo(db);
	userService := services.NewUserService(userRepo);
	userController := controllers.NewUserController(db, userService);

	user.GET("/detail", userController.Detail);
	user.GET("/balance", userController.Balance);
	user.GET("/point", userController.Point);
	user.GET("/point/type", userController.PointType);
	user.GET("/point/type/:IdPointType", userController.ListPoint);

	return user;
}