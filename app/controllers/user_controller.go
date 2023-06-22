package controllers

import (
	"database/sql"
	"net/http"
	"ppob/app/configs"
	"ppob/app/domain"
	"ppob/app/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userController struct {
	db *sql.DB
	userService services.UserService
}

func NewUserController(db *sql.DB, userService services.UserService) *userController {
	return &userController{
		db: db,
		userService: userService,
	};
}

func (controller *userController) Detail(ctx echo.Context) error {

	IdUser := ctx.Get("id_user").(int);

	status, code, data, err := controller.userService.UserDetail(ctx.Request().Context(), IdUser);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Error: err.Error(),
			Type: 6,
		};
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Data: data,
		Type: 6,
	};
	return response.ResponseMiddleware(ctx);
}

func (controller *userController) Balance(ctx echo.Context) error {
	
	IdUser := ctx.Get("id_user").(int);

	status, code, data, err := controller.userService.UserBalance(ctx.Request().Context(), IdUser);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Error: err.Error(),
			Type: 7,
		};
		return response.ResponseMiddleware(ctx);
	}
	
	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Data: data,
		Type: 7,
	};
	return response.ResponseMiddleware(ctx);
}

func (controller *userController) Point(ctx echo.Context) error {
	
	IdUser := ctx.Get("id_user").(int);

	status, code, data, err := controller.userService.UserPoint(ctx.Request().Context(), IdUser);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Error: err.Error(),
			Type: 8,
		};
		return response.ResponseMiddleware(ctx);
	}
	
	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Data: data,
		Type: 8,
	};
	return response.ResponseMiddleware(ctx);
}

func (controller *userController) PointType(ctx echo.Context) error {

	status, code, data, err := controller.userService.PointType(ctx.Request().Context());
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Error: err.Error(),
			Type: 9,
		};
		return response.ResponseMiddleware(ctx);
	}
	
	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Data: data,
		Type: 9,
	};
	return response.ResponseMiddleware(ctx);
}

func (controller *userController) ListPoint(ctx echo.Context) error {

	var pageNumber int
	var limitNumber int
	var IdPointType int

	if ctx.QueryParam("page") == "" {
		pageNumber = 1;
	} else {
		pageParse, err := strconv.Atoi(ctx.QueryParam("page"));
		if err != nil {
			response := &configs.Response{
				Status: http.StatusBadRequest,
				Code: 3896,
				DB: controller.db,
				Type: 10,
			};
			return response.ResponseMiddleware(ctx);
		}
		pageNumber = pageParse;
	}

	if ctx.QueryParam("limit") == "" {
		limitNumber = 10;
	} else {
		limitParse, err := strconv.Atoi(ctx.QueryParam("limit"));
		if err != nil {
			response := &configs.Response{
				Status: http.StatusBadRequest,
				Code: 3896,
				DB: controller.db,
				Type: 10,
			};
			return response.ResponseMiddleware(ctx);
		}
		limitNumber = limitParse;
	}

	IdPointTypeParam := ctx.Param("IdPointType");
	if IdPointTypeParam == ":IdPointType" {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 3896,
			DB: controller.db,
			Type: 10,
		};
		return response.ResponseMiddleware(ctx);
	} else {
		IdPointTypeParse, err := strconv.Atoi(IdPointTypeParam);
		if err != nil {
			response := &configs.Response{
				Status: http.StatusBadRequest,
				Code: 3896,
				DB: controller.db,
				Type: 10,
			};
			return response.ResponseMiddleware(ctx);
		}
		IdPointType = IdPointTypeParse;
	}

	IdUser := ctx.Get("id_user").(int);
	
	payload := &domain.ListPointRequest{
		IdUser: IdUser,
		Page: pageNumber,
		Limit: limitNumber,
		IdPointType: IdPointType,
	};
	
	status, code, data, err := controller.userService.ListPointByType(ctx.Request().Context(), payload);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			DB: controller.db,
			Error: err.Error(),
			Type: 10,
		};
		return response.ResponseMiddleware(ctx);
	}
	
	response := &configs.Response{
		Status: status,
		Code: code,
		DB: controller.db,
		Data: data,
		Type: 10,
	};
	return response.ResponseMiddleware(ctx);
}