package repository

import (
	"database/sql"
	"fmt"
	"os"
	"ppob/app/domain"
	"strconv"
)

type OtpRepository interface {
	GetUserAttempt(Msisdn string, Type string) (*domain.GetUserAttempt, error)
	GetDurationExpiredOTP() (*domain.ExpiredDurationOTP, error)
	InsertRequestOTP(data *domain.InsertRequestOTP) (string, error)
	UpdateOTPActive(idUser int, Type string, updatedOn string) error
	UpdateOTPAttempt(IdUser int, Otp string, UpdatedOn string) error
}

type otpRepository struct {
	db *sql.DB
}

func NewOtpRepository(db *sql.DB) OtpRepository {
	return &otpRepository{
		db: db,
	}
}

func (otpRepo *otpRepository) GetUserAttempt(Msisdn string, Type string) (*domain.GetUserAttempt, error) {
	var result domain.GetUserAttempt
	sqlQuery := fmt.Sprintf(`
	SELECT
		tu.id_user,
		(SELECT COUNT(*) FROM security.t_otp to2 WHERE to2.id_user = tu.id_user AND TRIM("type") = TRIM('%s') AND status = 1) AS attempt_request_otp,
		case WHEN
		(SELECT COUNT(*) FROM security.t_otp to2 WHERE to2.id_user = tu.id_user AND TRIM("type") = TRIM('%s') AND status = 1) > 0
		THEN (SELECT group_id FROM security.t_otp to2 WHERE to2.id_user = tu.id_user AND TRIM("type") = TRIM('%s') AND status = 1 ORDER BY created_on DESC LIMIT 1)
		ELSE
			NULL
		END AS group_id,
		(SELECT to2.attempt FROM "security".t_otp to2 WHERE currently_active = TRUE AND TRIM("type") = TRIM('%s') AND status = 1 LIMIT 1),
		(SELECT to2.otp FROM "security".t_otp to2 WHERE currently_active = TRUE AND TRIM("type") = TRIM('REGISTER') AND status = 1 LIMIT 1)
	FROM
		user_management.t_users tu
	WHERE
		tu.msisdn = '%s'
	`, Type, Type, Type, Type, Msisdn);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result.IdUser, &result.AttemptRequestOtp, &result.GroupId, &result.AttemptVerification, &result.Otp);
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
	sqlQuery := fmt.Sprintf("INSERT INTO security.t_otp (id_user, type, otp, created_on, expired_on, group_id) VALUES (%d, '%s', '%s', '%s', '%s', '%s') RETURNING otp", data.IdUser, data.Type, data.Otp, data.CratedOn, data.ExpiredOn, data.GroupId);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result);
	if err != nil {
		return "", err;
	}
	
	return result, nil;
}

func (otpRepo *otpRepository) UpdateOTPActive(idUser int, Type string, updatedOn string) error {
	sqlQuery := fmt.Sprintf(`UPDATE security.t_otp SET currently_active = false, updated_on = '%s' WHERE id_user = %d AND type = '%s' AND currently_active = true`, updatedOn, idUser, Type);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		if err == sql.ErrNoRows {
			return nil;
		}
		return err;
	}

	return nil;
}

func (otpRepo *otpRepository) UpdateOTPAttempt(IdUser int, Otp string, UpdatedOn string) error {
	sqlQuery := fmt.Sprintf("UPDATE security.t_otp SET attempt = CASE WHEN attempt + 1 > 3 THEN attempt ELSE attempt + 1 END, updated_on = '%s' WHERE id_user = %d AND otp = '%s'", UpdatedOn, IdUser, Otp);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		return err;
	}

	return nil;
}