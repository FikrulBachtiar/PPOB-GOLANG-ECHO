package repository

import (
	"database/sql"
	"fmt"
	"os"
	"ppob/app/domain"
	"strconv"
)

type OtpRepository interface {
	GetUserAttempt(Msisdn string, Type string, DeviceID string) (*domain.GetUserAttempt, error)
	GetDurationExpiredOTP() (*domain.ExpiredDurationOTP, error)
	GetParamValue(code string, currentDate string) (*domain.ParamValue, error)
	InsertRequestOTP(data *domain.InsertRequestOTP) (string, error)
	UpdateOTPActive(idUser int, Type string, DeviceID string, updatedOn string) error
	UpdateOTPAttempt(IdUser int, Otp string, UpdatedOn string) error
	ResetOTP(IdUser int, GroupId string) error
	BlockUser(IdUser int, currentDate string) error
}

type otpRepository struct {
	db *sql.DB
}

func NewOtpRepository(db *sql.DB) OtpRepository {
	return &otpRepository{
		db: db,
	}
}

func (otpRepo *otpRepository) GetUserAttempt(Msisdn string, Type string, DeviceID string) (*domain.GetUserAttempt, error) {
	var result domain.GetUserAttempt
	sqlQuery := fmt.Sprintf(`
	SELECT
		tu.id_user,
		tu.user_status_code,
		COALESCE(to1.attempt_request_otp, 0) AS attempt_request_otp,
		CASE
			WHEN to1.attempt_request_otp > 0 THEN to1.group_id
			ELSE NULL
		END AS group_id,
		to1.attempt,
		to1.otp,
		to1.expired_on
	FROM
		user_management.t_users tu
	LEFT JOIN (
		SELECT
			id_user,
			COUNT(*) AS attempt_request_otp,
			MAX(CASE WHEN currently_active = TRUE THEN group_id END) AS group_id,
			MAX(CASE WHEN currently_active = TRUE THEN attempt END) AS attempt,
			MAX(CASE WHEN currently_active = TRUE THEN otp END) AS otp,
			MAX(CASE WHEN currently_active = TRUE THEN expired_on END) AS expired_on
		FROM
			security.t_otp
		WHERE
			TRIM("type") = TRIM('%s')
			AND status = 1
			AND device_id = '%s'
		GROUP BY
			id_user
	) to1 ON tu.id_user = to1.id_user
	WHERE
		tu.msisdn = '%s'
	`, DeviceID, Type, Msisdn);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result.IdUser, &result.UserStatusCode, &result.AttemptRequestOtp, &result.GroupId, &result.AttemptVerification, &result.Otp, &result.ExpiredOn);
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

func (otpRepo *otpRepository) GetParamValue(code string, currentDate string) (*domain.ParamValue, error) {
	var result domain.ParamValue;
	sqlQuery := fmt.Sprintf("SELECT * FROM master.sp_get_param_value('%s', '%s')", currentDate, code);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result.Value, &result.Measure);
	if err != nil {
		return nil, err;
	}


	return &result, nil;
}

func (otpRepo *otpRepository) InsertRequestOTP(data *domain.InsertRequestOTP) (string, error) {
	var result string
	sqlQuery := fmt.Sprintf("INSERT INTO security.t_otp (id_user, type, otp, created_on, expired_on, group_id, device_id) VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s') RETURNING otp", data.IdUser, data.Type, data.Otp, data.CratedOn, data.ExpiredOn, data.GroupId, data.DeviceID);

	err := otpRepo.db.QueryRow(sqlQuery).Scan(&result);
	if err != nil {
		return "", err;
	}
	
	return result, nil;
}

func (otpRepo *otpRepository) UpdateOTPActive(idUser int, Type string, DeviceID string, updatedOn string) error {
	sqlQuery := fmt.Sprintf(`UPDATE security.t_otp SET currently_active = false, updated_on = '%s' WHERE id_user = %d AND type = '%s' AND device_id = '%s' AND currently_active = true`, updatedOn, idUser, Type, DeviceID);

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
	sqlQuery := fmt.Sprintf("UPDATE security.t_otp SET attempt = attempt + 1, updated_on = '%s' WHERE id_user = %d AND otp = '%s' AND currently_active = true", UpdatedOn, IdUser, Otp);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		return err;
	}

	return nil;
}

func (otpRepo *otpRepository) ResetOTP(IdUser int, GroupId string) error {
	sqlQuery := fmt.Sprintf("UPDATE security.t_otp SET status = 0, currently_active WHERE id_user = %d AND group_id = '%s'", IdUser, GroupId);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		return err;
	}

	return nil;
}

func (otpRepo *otpRepository) BlockUser(IdUser int, currentDate string) error {
	sqlQuery := fmt.Sprintf("UPDATE user_management.t_users SET user_status_code = 'T', blocked_on = '%s' WHERE id_user = %d", currentDate, IdUser);

	_, err := otpRepo.db.Exec(sqlQuery);
	if err != nil {
		return err;
	}

	return nil;
}