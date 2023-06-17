package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"ppob/app/domain"
	"ppob/app/repository"
	"ppob/app/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type OnboardingService interface {
	CheckAccount(c context.Context, payload *domain.CheckPayload) (int, int, *domain.CheckResponse, error)
	LoginAccount(c context.Context, payload *domain.LoginPayload, header *domain.LoginHeader) (int, int, *domain.LoginResponse, error)
}

type onboardingService struct {
	onboardingRepo repository.OnboardingRepo
}

var user_status_code, statusAccountBlocked, passphrase, vector, appName, expiredTokenDuration, expiredTokenType string
var tokenKey string

func NewOnboardingService(onboardingRepo repository.OnboardingRepo) OnboardingService {

	user_status_code = os.Getenv("STATUS_ACCOUNT_NORMAL");
	statusAccountBlocked = os.Getenv("STATUS_ACCOUNT_BLOCKED");
	passphrase = os.Getenv("KEY_AES");
	vector = os.Getenv("VECTOR_AES");
	appName = os.Getenv("JWT_ISSUER");
	expiredTokenDuration = os.Getenv("JWT_EXPIRED_DURATION");
	expiredTokenType = os.Getenv("JWT_EXPIRED_DURATION_TYPE");
	tokenKey = os.Getenv("JWT_KEY");

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

func (services *onboardingService) LoginAccount(c context.Context, payload *domain.LoginPayload, header *domain.LoginHeader) (int, int, *domain.LoginResponse, error) {
	var response domain.LoginResponse
	payload.Msisdn = utils.FormatPhoneNumber(payload.Msisdn);

	users, err := services.onboardingRepo.GetUserByNumber(payload.Msisdn);
	if err != nil {
		return http.StatusInternalServerError, 9438, nil, err;
	}

	if users == nil {
		return http.StatusForbidden, 2203, nil, nil;
	}

	if users.UserStatusCode == statusAccountBlocked {
		return http.StatusForbidden, 6492, nil, nil;
	}

	if users.Pin == nil {
		return http.StatusForbidden, 2203, nil, nil;
	}

	pinPayloadDecrypt, err := utils.DecryptAES(payload.Pin, passphrase, vector);
	if err != nil {
		return http.StatusInternalServerError, 8753, nil, err;
	}
	
	pinUserDecrypt, err := utils.DecryptAES(*users.Pin, passphrase, vector);
	if err != nil {
		return http.StatusInternalServerError, 5706, nil, err;
	}

	if string(pinUserDecrypt) != string(pinPayloadDecrypt) {
		return http.StatusBadRequest, 6229, nil, nil;
	}

	expiredTokenInt, err := strconv.Atoi(expiredTokenDuration);
	if err != nil {
		return http.StatusInternalServerError, 7629, nil, err;
	}

	expiredTokenType, err := utils.ConvertStringToDurationType(expiredTokenType);
	if err != nil {
		return http.StatusInternalServerError, 7362, nil, err;
	}

	now := time.Now().Local();
	expiredDate := now.Add(time.Duration(expiredTokenInt) * expiredTokenType);
	MyClaims := &domain.ClaimsJWT{
		StandardClaims: jwt.StandardClaims{
			Issuer: appName,
			ExpiresAt: expiredDate.Unix(),
		},
		Msisdn: payload.Msisdn,
		DeviceID: header.DeviceID,
		NotificationID: payload.NotificationID,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		MyClaims,
	);

	signedToken, err := token.SignedString([]byte(tokenKey));
	if err != nil {
		return http.StatusInternalServerError, 4852, nil, err;
	}

	uniqueCode, err := utils.RandomString(30);
	if err != nil {
		return http.StatusInternalServerError, 3460, nil, err;
	}

	refreshToken := utils.Base64EncodeToString([]byte(uniqueCode));

	dataInsertToken := &domain.InsertToken{
		IdUser: users.IdUser,
		Msisdn: payload.Msisdn,
		Token: signedToken,
		RefreshToken: refreshToken,
		NotificationID: payload.NotificationID,
		ExpiredDate: expiredDate.Format("2006-01-02 15:04:05.999"),
		CreatedOn: now.Format("2006-01-02 15:04:05.999"),
		DeviceID: header.DeviceID,
	};

	dataInsertLogin := &domain.InsertLogin{
		IdUser: users.IdUser,
		OsName: header.OsName,
		OsVersion: header.OsVersion,
		DeviceModel: header.DeviceModel,
		DeviceId: header.DeviceID,
		AppVersion: header.AppVersion,
		Longitude: header.Longitude,
		Latitude: header.Latitude,
		NotificationID: payload.NotificationID,
		Pin: payload.Pin,
		LoginOn: now.Format("2006-01-02 15:04:05.999"),
		StatusCode: 0,
		Fullname: users.Fullname,
		Msisdn: payload.Msisdn,
		Uuid: users.Uuid,
	};

	if err := services.onboardingRepo.InsertLogin(dataInsertToken, dataInsertLogin); err != nil {
		return http.StatusInternalServerError, 8031, nil, err;
	}

	response.FullName = users.Fullname;
	response.Token = signedToken;
	response.RefreshToken = refreshToken;

	return http.StatusOK, 0, &response, nil;
}