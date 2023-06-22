package services

import (
	"context"
	"math"
	"net/http"
	"ppob/app/domain"
	"ppob/app/repository"
)

type UserService interface{
	UserDetail(ctx context.Context, IdUser int) (int, int, *domain.UserDetailResponse, error)
	UserBalance(ctx context.Context, IdUser int) (int, int, *domain.UserBalanceResponse, error)
	UserPoint(ctx context.Context, IdUser int) (int, int, *domain.UserPointResponse, error)
	PointType(ctx context.Context) (int, int, []domain.PointTypeResponse, error)
	ListPointByType(ctx context.Context, payload *domain.ListPointRequest) (int, int, *domain.PaginationDefault, error)
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	};
}

func (service *userService) UserDetail(ctx context.Context, IdUser int) (int, int, *domain.UserDetailResponse, error) {
	
	result, err := service.userRepo.GetUserDetail(IdUser);
	if err != nil {
		return http.StatusInternalServerError, 2977, nil, err;
	}

	return http.StatusOK, 0, result, nil;
}

func (service *userService) UserBalance(ctx context.Context, IdUser int) (int, int, *domain.UserBalanceResponse, error) {
	UserBalance, err := service.userRepo.GetBalanceUser(IdUser);
	if err != nil {
		return http.StatusInternalServerError, 6911, nil, err;
	}

	isKtpVerify, err := service.userRepo.IsKtpVerify(IdUser);
	if err != nil {
		return http.StatusInternalServerError, 4460, nil, err;
	}

	response := &domain.UserBalanceResponse{
		Currency: UserBalance.Currency,
		Value: UserBalance.Value,
		KtpStatus: isKtpVerify,
	};

	return http.StatusOK, 0, response, nil;
}

func (service *userService) UserPoint(ctx context.Context, IdUser int) (int, int, *domain.UserPointResponse, error) {
	UserPoint, err := service.userRepo.GetPointUser(IdUser);
	if err != nil {
		return http.StatusInternalServerError, 9348, nil, err;
	}

	response := &domain.UserPointResponse{
		Currency: UserPoint.Currency,
		TotalPoint: UserPoint.Value,
		IsTap: UserPoint.IsTap,
	};

	return http.StatusOK, 0, response, nil;
}

func (service *userService) PointType(ctx context.Context) (int, int, []domain.PointTypeResponse, error) {
	results, err := service.userRepo.GetPointType();
	if err != nil {
		return http.StatusInternalServerError, 5361, nil, err;
	}

	return http.StatusOK, 0, results, nil;
}

func (service *userService) ListPointByType(ctx context.Context, payload *domain.ListPointRequest) (int, int, *domain.PaginationDefault, error) {
	
	var totalPage, totalData int;

	results, err := service.userRepo.GetListPoint(payload.IdUser, payload.IdPointType);
	if err != nil {
		return http.StatusInternalServerError, 2293, nil, err;
	}

	var list_data []domain.ListPointResponse

	if len(results) <= 0{
		list_data = make([]domain.ListPointResponse, 0);
	} else {
		totalData = len(results);
		totalPage = int(math.Ceil(float64(totalData) / float64(payload.Limit)));
		offset := (payload.Page - 1) * payload.Limit;

		if totalPage < payload.Page {
			return http.StatusNotFound, 2623, nil, nil;
		}

		for x := offset; x < totalData && x < offset+payload.Limit; x++ {
			list_data = append(list_data, results[x]);
		}
	}

	response := &domain.PaginationDefault{
		TotalPages: totalPage,
		CurrentPage: payload.Page,
		TotalData: totalData,
		ListData: list_data,
	};
	
	return http.StatusOK, 0, response, nil;
}