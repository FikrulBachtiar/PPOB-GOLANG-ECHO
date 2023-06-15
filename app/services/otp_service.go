package services

import (
	"net/http"
	"os"
	"ppob/app/domain"
	"ppob/app/repository"
	"ppob/app/utils"
	"strconv"
	"time"
)

type OtpService interface {
	CreateOTP(payload *domain.RequestOtpPayload, header *domain.RequestOtpHeader) (int, int, *domain.RequestOtpResponse, error)
	VerificationOtp(payload *domain.VerificationOtpPayload, header *domain.RequestOtpHeader) (int, int, error)
}

type otpService struct {
	otpRepo repository.OtpRepository
}

var codeMaxRequestOTP, codeMaxValidateOTP, codeDurationExpiredOTP, codeAccountBlockTemporary, codeAccountBlock string

func NewOtpService(otpRepo repository.OtpRepository) OtpService {

	codeMaxRequestOTP = os.Getenv("CODE_MAX_REQUEST_OTP");
	codeMaxValidateOTP = os.Getenv("CODE_MAX_VALIDATION_OTP");
	codeDurationExpiredOTP = os.Getenv("CODE_DURATION_EXPIRED_OTP");
	codeAccountBlockTemporary = os.Getenv("STATUS_ACCOUNT_TEMPORARY_BLOCK");
	codeAccountBlock = os.Getenv("STATUS_ACCOUNT_BLOCKED");

	return &otpService{
		otpRepo: otpRepo,
	}
}

func (otpService *otpService) CreateOTP(payload *domain.RequestOtpPayload, header *domain.RequestOtpHeader) (int, int, *domain.RequestOtpResponse, error) {
	var response domain.RequestOtpResponse;
	payload.Msisdn = utils.FormatPhoneNumber(payload.Msisdn);

	now := time.Now();
	currentDate := now.Local().Format("2006-01-02 15:04:05.999");

	// get max attempt request
	maxRequest, err := otpService.otpRepo.GetParamValue(codeMaxRequestOTP, currentDate);
	if err != nil {
		return http.StatusInternalServerError, 1013, nil, err;
	}

	if maxRequest.Value == nil {
		return http.StatusInternalServerError, 3607, nil, err;
	}

	maxRequestInt, err := strconv.Atoi(*maxRequest.Value);
	if err != nil {
		return http.StatusInternalServerError, 2089, nil, err;
	}

	// get user attempt
	users, err := otpService.otpRepo.GetUserAttempt(payload.Msisdn, header.DeviceID, payload.Type);
	if err != nil {
		return http.StatusInternalServerError, 5707, nil, err;
	}

	// check user
	if users == nil {
		return http.StatusBadRequest, 2203, nil, nil;
	}

	if users.UserStatusCode == codeAccountBlockTemporary {
		return http.StatusForbidden, 8567, nil, nil;
	}

	if users.UserStatusCode == codeAccountBlock {
		return http.StatusForbidden, 6492, nil, nil;
	}

	// check can request again?
	if users.AttemptRequestOtp >= maxRequestInt {
		return http.StatusTooManyRequests, 4736, nil, nil;
	}
	
	// create groupID otp
	var GroupId string
	if users.GroupId == nil {
		GroupId, err = utils.RandomString(30);
		if err != nil {
			return http.StatusInternalServerError, 5525, nil, err;
		}
	} else {
		GroupId = *users.GroupId;
	}

	otp := strconv.Itoa(utils.RandomNumber(6));

	duration, err := otpService.otpRepo.GetParamValue(codeDurationExpiredOTP, currentDate);
	if err != nil {
		return http.StatusInternalServerError, 1013, nil, err;
	}

	if duration.Value == nil {
		return http.StatusInternalServerError, 1093, nil, err;
	}

	var durationType time.Duration
	if duration.Measure != nil {
		durationType, err = utils.ConvertStringToDurationType(*duration.Measure);
		if err != nil {
			return http.StatusInternalServerError, 1871, nil, err;
		}
	}

	durationInt, err := strconv.Atoi(*duration.Value);
	if err != nil {
		return http.StatusInternalServerError, 2089, nil, err;
	}

	expiredOtp := now.Add(time.Duration(durationInt) * durationType).Local().Format("2006-01-02 15:04:05.999");

	if err := otpService.otpRepo.UpdateOTPActive(users.IdUser, payload.Type, header.DeviceID, currentDate); err != nil {
		return http.StatusInternalServerError, 5029, nil, err;
	}

	var data domain.InsertRequestOTP;
	data.IdUser = users.IdUser;
	data.Type = payload.Type;
	data.Otp = otp;
	data.CratedOn = currentDate;
	data.ExpiredOn = expiredOtp;
	data.GroupId = GroupId;
	data.DeviceID = header.DeviceID;
	dataOtp, err := otpService.otpRepo.InsertRequestOTP(&data);
	if err != nil {
		return http.StatusInternalServerError, 1656, nil, err;
	}

	response.Otp = dataOtp;
	response.Attempt = users.AttemptRequestOtp + 1;
	response.MaxAttempt = maxRequestInt;
	return http.StatusOK, 0, &response, nil;
}

func (otpService *otpService) VerificationOtp(payload *domain.VerificationOtpPayload, header *domain.RequestOtpHeader) (int, int, error) {
	payload.Msisdn = utils.FormatPhoneNumber(payload.Msisdn);

	now := time.Now();
	currentDate := now.Local().Format("2006-01-02 15:04:05.999");

	maxValidate, err := otpService.otpRepo.GetParamValue(codeMaxValidateOTP, currentDate);
	if err != nil {
		return http.StatusInternalServerError, 1013, err;
	}

	if maxValidate.Value == nil {
		return http.StatusInternalServerError, 3607, err;
	}

	maxValidateInt, err := strconv.Atoi(*maxValidate.Value);
	if err != nil {
		return http.StatusInternalServerError, 2089, err;
	}

	users, err := otpService.otpRepo.GetUserAttempt(payload.Msisdn, header.DeviceID, payload.Type);
	if err != nil {
		return http.StatusInternalServerError, 5707, err;
	}

	if users == nil {
		return http.StatusBadRequest, 2203, nil;
	}

	if users.UserStatusCode == codeAccountBlockTemporary {
		return http.StatusForbidden, 8567, nil;
	}

	if users.UserStatusCode == codeAccountBlock {
		return http.StatusForbidden, 6492, nil;
	}

	if users.Otp == nil && users.ExpiredOn == nil && users.GroupId == nil && users.AttemptVerification == nil {
		return http.StatusNotFound, 1748, nil;
	}

	if users.AttemptVerification != nil && *users.AttemptVerification >= maxValidateInt {
		err := otpService.otpRepo.BlockUser(users.IdUser, currentDate);
		if err != nil {
			return http.StatusInternalServerError, 4421, err;
		}

		return http.StatusTooManyRequests, 8761, nil; 
	}
	
	if users.ExpiredOn != nil {
		ExpiredOn, err := time.Parse("2006-01-02T15:04:05.999Z", *users.ExpiredOn);
		if err != nil {
			return http.StatusInternalServerError, 2175, err;
		}
		
		if now.After(ExpiredOn) {
			return http.StatusBadRequest, 2476, nil;
		}
	}

	if users.Otp != nil && *users.Otp != payload.OtpToken {
		updatedOn := now.Local().Format("2006-01-02 15:04:05.999");
		err := otpService.otpRepo.UpdateOTPAttempt(users.IdUser, *users.Otp, updatedOn);
		if err != nil {
			return http.StatusInternalServerError, 7365, err;
		}
		return http.StatusBadRequest, 1799, nil;
	}

	return http.StatusOK, 0, nil;
}