package repository

import (
	"database/sql"
	"fmt"
	"os"
	"ppob/app/domain"
	"strconv"
)

type OtpRepository interface {
	GetUserAttempt(msisdn string) (*domain.GetUserAttempt, error)
	GetDurationExpiredOTP() (*domain.ExpiredDurationOTP, error)
	InsertRequestOTP(data *domain.InsertRequestOTP) (string, error)
	UpdateOTPActive(idUser int, Type string, updatedOn string) error
	UpdateAttemptUser(IdUser int) error
}

type otpRepository struct {
	db *sql.DB
}

func NewOtpRepository(db *sql.DB) OtpRepository {
	return &otpRepository{
		db: db,
	}
}

func (otpRepo *otpRepository) GetUserAttempt(msisdn string) (*domain.GetUserAttempt, error) {
	var result domain.GetUserAttempt
	sqlQuery := fmt.Sprintf(`SELECT id_user, attempt_request_otp FROM user_management.t_users WHERE msisdn = '%s'`, msisdn);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result.IdUser, &result.AttemptRequestOtp);
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil;
		}

		return nil, err;
	}

	return &result, nil;
}

func (otpRepo *otpRepository) GetDurationExpiredOTP() (*domain.ExpiredDurationOTP, error) {
	var result domain.ExpiredDurationOTP;
	sqlQuery := "SELECT duration, duration_type FROM master.t_expired_duration_otp WHERE status = 1";

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result.Duration, &result.DurationType);
	if err != nil {
		if err == sql.ErrNoRows {
			result.Duration, err = strconv.Atoi(os.Getenv("DEFAULT_EXPIRED_DURATION_OTP"));
			if err != nil {
				return nil, err;
			}
			result.DurationType = os.Getenv("DEFAULT_EXPIRED_DURATION_TYPE_OTP");
			return &result, nil;
		}

		return nil, err;
	}


	return &result, nil;
}

func (otpRepo *otpRepository) InsertRequestOTP(data *domain.InsertRequestOTP) (string, error) {
	var result string
	sqlQuery := fmt.Sprintf("INSERT INTO security.t_otp (id_user, type, otp, created_on, expired_on) VALUES (%d, '%s', '%s', '%s', '%s') RETURNING otp", data.IdUser, data.Type, data.Otp, data.CratedOn, data.ExpiredOn);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result);
	if err != nil {
		return "", err;
	}
	
	return result, nil;
}

func (otpRepo *otpRepository) UpdateOTPActive(idUser int, Type string, updatedOn string) error {
	sqlQuery := fmt.Sprintf(`UPDATE security.t_otp SET status = 0, updated_on = '%s' WHERE id_user = %d AND type = '%s'`, updatedOn, idUser, Type);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		if err == sql.ErrNoRows {
			return nil;
		}
		return err;
	}

	return nil;
}

func (otpRepo *otpRepository) UpdateAttemptUser(IdUser int) error {
	sqlQuery := fmt.Sprintf("UPDATE user_management.t_users SET attempt_request_otp = CASE WHEN attempt_request_otp + 1 > 3 THEN attempt_request_otp ELSE attempt_request_otp + 1 END WHERE id_user = %d", IdUser);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		return err;
	}
	
	return nil;
}