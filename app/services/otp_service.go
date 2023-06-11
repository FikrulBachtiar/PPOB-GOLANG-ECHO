package services

import (
	"net/http"
	"ppob/app/domain"
	"ppob/app/repository"
	"ppob/app/utils"
	"strconv"
	"time"
)

type OtpService interface {
	CreateOTP(payload *domain.RequestOtpPayload) (int, int, *domain.RequestOtpResponse, error)
	VerificationOtp(payload *domain.VerificationOtpPayload) (int, int, error)
}

type otpService struct {
	otpRepo repository.OtpRepository
}

func NewOtpService(otpRepo repository.OtpRepository) OtpService {
	return &otpService{
		otpRepo: otpRepo,
	}
}

func (otpService *otpService) CreateOTP(payload *domain.RequestOtpPayload) (int, int, *domain.RequestOtpResponse, error) {
	var response domain.RequestOtpResponse;
	payload.Msisdn = utils.FormatPhoneNumber(payload.Msisdn);
	
	otp := strconv.Itoa(utils.RandomNumber(6));
	users, err := otpService.otpRepo.GetUserAttempt(payload.Msisdn, payload.Type);
	if err != nil {
		return http.StatusInternalServerError, 5707, nil, err;
	}

	if users == nil {
		return http.StatusBadRequest, 2203, nil, nil;
	}

	if users.IdUser == 0 {
		return http.StatusBadRequest, 2203, nil, nil;
	}

	if users.AttemptRequestOtp >= 3 {
		return http.StatusTooManyRequests, 4736, nil, nil;
	}

	var GroupId string
	if users.GroupId == nil {
		GroupId, err = utils.RandomString(15);
		if err != nil {
			return http.StatusInternalServerError, 5525, nil, err;
		}
	} else {
		GroupId = *users.GroupId;
	}

	now := time.Now();
	duration, err := otpService.otpRepo.GetDurationExpiredOTP();
	if err != nil {
		return http.StatusInternalServerError, 1013, nil, err;
	}

	durationType, err := utils.ConvertStringToDurationType(duration.DurationType);
	if err != nil {
		return http.StatusInternalServerError, 1871, nil, err;
	}

	expiredOtp := now.Add(time.Duration(duration.Duration) * durationType).Local().Format("2006-01-02 15:04:05.999");
	createdOn := now.Local().Format("2006-01-02 15:04:05.999");

	if err := otpService.otpRepo.UpdateOTPActive(users.IdUser, payload.Type, createdOn); err != nil {
		return http.StatusInternalServerError, 5029, nil, err;
	}

	var data domain.InsertRequestOTP;
	data.IdUser = users.IdUser;
	data.Type = payload.Type;
	data.Otp = otp;
	data.CratedOn = createdOn;
	data.ExpiredOn = expiredOtp;
	data.GroupId = GroupId;
	dataOtp, err := otpService.otpRepo.InsertRequestOTP(&data);
	if err != nil {
		return http.StatusInternalServerError, 1656, nil, err;
	}

	response.Otp = dataOtp;
	return http.StatusOK, 0, &response, nil;
}

func (otpService *otpService) VerificationOtp(payload *domain.VerificationOtpPayload) (int, int, error) {
	payload.Msisdn = utils.FormatPhoneNumber(payload.Msisdn);

	users, err := otpService.otpRepo.GetUserAttempt(payload.Msisdn, payload.Type);
	if err != nil {
		return http.StatusInternalServerError, 5707, err;
	}

	if users == nil {
		return http.StatusBadRequest, 2203, nil;
	}

	if users.IdUser == 0 {
		return http.StatusBadRequest, 2203, nil;
	}

	if users.AttemptVerification != nil && *users.AttemptVerification >= 3 {
		return http.StatusTooManyRequests, 8761, nil;
	}

	if users.Otp != nil && *users.Otp != payload.OtpToken {
		now := time.Now();
		updatedOn := now.Local().Format("2006-01-02 15:04:05.999");
		err := otpService.otpRepo.UpdateOTPAttempt(users.IdUser, *users.Otp, updatedOn);
		if err != nil {
			return http.StatusInternalServerError, 7365, err;
		}
		return http.StatusBadRequest, 1799, nil;
	}

	return http.StatusOK, 0, nil;
}