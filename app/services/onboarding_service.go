package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"ppob/app/domain"
	"ppob/app/repository"
	"ppob/app/utils"
	"time"
)

type OnboardingService interface {
	CheckAccount(c context.Context, payload *domain.CheckPayload) (int, int, *domain.CheckResponse, error)
}

type onboardingService struct {
	onboardingRepo repository.OnboardingRepo
}

var user_status_code, statusAccountBlocked string

func NewOnboardingService(onboardingRepo repository.OnboardingRepo) OnboardingService {

	user_status_code = os.Getenv("STATUS_ACCOUNT_NORMAL");
	statusAccountBlocked = os.Getenv("STATUS_ACCOUNT_BLOCKED");

	return &onboardingService{
		onboardingRepo: onboardingRepo,
	}
}

func (onboardService *onboardingService) CheckAccount(c context.Context, payload *domain.CheckPayload) (int, int, *domain.CheckResponse, error) {
	var response domain.CheckResponse
	
	payload.Msisdn = utils.FormatPhoneNumber(payload.Msisdn);
	
	user, err := onboardService.onboardingRepo.GetUserByNumber(payload.Msisdn);
	if err != nil {
		return http.StatusInternalServerError, 4007, nil, err;
	}

	fmt.Println(user)
	
	if user == nil {

		uuid := utils.GenerateMD5(payload.Msisdn);
		created_on := time.Now().Local().Format("2006-01-02 15:04:05.999");
		user_code, err := utils.RandomString(30);
		if err != nil {
			return http.StatusInternalServerError, 3925, nil, err;
		}
		
		serial, err := utils.SerialNumberString(utils.RandomNumber(10));
		if err != nil {
			return http.StatusInternalServerError, 7146, nil, err;
		}
		fullName := fmt.Sprintf("User%s", serial);
		
		result, err := onboardService.onboardingRepo.CreateUser(payload.Msisdn, uuid, fullName, user_code, created_on, user_status_code);
		if err != nil {
			return http.StatusInternalServerError, 1738, nil, err;
		}
		
		response.EmailStatus = utils.EmailParsingString(result.IsEmailVerified);
		response.IsQuestion = result.IsQuestion;
		response.IsBlocked = false;
		response.State = "REGISTER";

		return http.StatusOK, 0, &response, nil;
	}

	if user.UserStatusCode == statusAccountBlocked {
		return http.StatusForbidden, 7803, nil, nil;
	}

	var isLoginCode int = 1;
	if user.LoginStatus != nil && *user.LoginStatus == isLoginCode {
		response.EmailStatus = utils.EmailParsingString(user.EmailStatus);
		response.IsQuestion = user.IsQuestion;
		response.IsBlocked = false;
		response.State = "OTP";
		return http.StatusOK, 0, &response, nil;
	}

	response.EmailStatus = utils.EmailParsingString(user.EmailStatus);
	response.IsQuestion = user.IsQuestion;
	response.IsBlocked = false;
	response.State = "LOGIN";

	return http.StatusOK, 0, &response, nil;
}