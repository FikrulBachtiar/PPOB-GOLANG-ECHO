package repository

import (
	"database/sql"
	"fmt"
	"ppob/app/domain"
)

type OnboardingRepo interface {
	GetUserByNumber(msisdn string) (*domain.GetUserByNumber, error)
	CreateUser(msisdn string, uuid string, fullName string, user_code string, created_on string, user_status_code string) (*domain.InsertUser, error)
	CheckLoginDevice(IdUser int) (int, error)
	InsertLogin(tokens *domain.InsertToken, logins *domain.InsertLogin) error
	UpdateDataUserLogin(loginStatus int, lastLogin string, updatedOn string) error
	UpdateUserLogout(currentDate string, DeviceID string, Msisdn string) error
}

type onboardingRepo struct {
	db *sql.DB
}

func NewOnboardingRepo(db *sql.DB) OnboardingRepo {
	return &onboardingRepo{
		db: db,
	}
}
	
func (repo *onboardingRepo) GetUserByNumber(msisdn string) (*domain.GetUserByNumber, error) {
	var result domain.GetUserByNumber
	sqlQuery := fmt.Sprintf(`
	SELECT
		tu.id_user,
		tu.fullname,
		tu.uuid,
		tu.is_question,
		tu.login_status,
		tu.last_login,
		tu.user_status_code,
		coalesce(tue.is_verified, false) as email_status,
		tup.pin
	FROM
		user_management.t_users tu
	LEFT JOIN
		"security".t_user_email tue
	ON
		tue.id_user = tu.id_user
	LEFT JOIN
		"security".t_user_pin tup
	ON
		tup.id_user = tu.id_user
	WHERE
		tu.msisdn = '%s'
		AND tup.status = 1`,
	msisdn);

	err := repo.db.QueryRow(sqlQuery).Scan(&result.IdUser, &result.Fullname, &result.Uuid, &result.IsQuestion, &result.LoginStatus, &result.LastLogin, &result.UserStatusCode, &result.EmailStatus, &result.Pin);
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil;
		}
		return nil, err;
	}

	return &result, nil;
}
  
func (repo *onboardingRepo) CreateUser(msisdn string, uuid string, fullName string, user_code string, created_on string, user_status_code string) (*domain.InsertUser, error) {
	var result domain.InsertUser;
	sqlQuery := fmt.Sprintf(`INSERT INTO user_management.t_users (msisdn, uuid, fullname, user_code, created_on, user_status_code) VALUES ('%s', '%s', '%s', '%s', '%s', '%s') RETURNING is_question, is_email_verified`, msisdn, uuid, fullName, user_code, created_on, user_status_code);

	err := repo.db.QueryRow(sqlQuery).Scan(&result.IsQuestion, &result.IsEmailVerified);
	if err != nil {
		return nil, err;
	}

	return &result, nil;
}

func (repo *onboardingRepo) CheckLoginDevice(IdUser int) (int, error) {
	var resultNumber int
	sqlQuery := fmt.Sprintf("SELECT count(*) as resultNumber FROM transaction.t_login_app WHERE id_user = %d AND is_logout = TRUE", IdUser);

	err := repo.db.QueryRow(sqlQuery).Scan(&resultNumber);
	if err != nil {
		return 0, err;
	}

	return resultNumber, nil;
}

func (repo *onboardingRepo) InsertLogin(tokens *domain.InsertToken, logins *domain.InsertLogin) error {
	trx, err := repo.db.Begin();
	if err != nil {
		return err;
	}

	sqlQueryUpdateLogin := fmt.Sprintf("UPDATE transaction.t_login_app SET is_logout = true, logout_on = '%s'", logins.LoginOn);
	
	_, err = trx.Query(sqlQueryUpdateLogin);
	if err != nil {
		return err;
	}

	sqlQueryToken := fmt.Sprintf("INSERT INTO security.t_session_app (id_user, msisdn, token, refresh_token, notification_id, expired_date, created_on, device_id) VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s') ON CONFLICT ON CONSTRAINT t_session_app_un DO UPDATE SET token = '%s', refresh_token = '%s', msisdn = '%s', expired_date = '%s', updated_on = '%s'",
	tokens.IdUser, tokens.Msisdn, tokens.Token, tokens.RefreshToken, tokens.NotificationID, tokens.ExpiredDate, tokens.CreatedOn, tokens.DeviceID, tokens.Token, tokens.RefreshToken, tokens.Msisdn, tokens.ExpiredDate, tokens.CreatedOn);

	_, err = trx.Query(sqlQueryToken);
	if err != nil {
		trx.Rollback();
		return err;
	}

	sqlQueryInsert := fmt.Sprintf("INSERT INTO transaction.t_login_app (id_user, os_name, os_version, device_model, device_id, app_version, longitude, latitude, notification_id, pin, login_on, status_code, fullname, msisdn, uuid) VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, '%s', '%s', '%s')",
	logins.IdUser, logins.OsName, logins.OsVersion, logins.DeviceModel, logins.DeviceId, logins.AppVersion, logins.Longitude, logins.Latitude, logins.NotificationID, logins.Pin, logins.LoginOn, logins.StatusCode, logins.Fullname, logins.Msisdn, logins.Uuid);

	_, err = trx.Query(sqlQueryInsert);
	if err != nil {
		trx.Rollback();
		return err;
	}

	sqlQueryUpdateUser := fmt.Sprintf("UPDATE user_management.t_users SET login_status = %d, last_login = '%s', updated_on = '%s'", 1, logins.LoginOn, logins.LoginOn);
	
	_, err = trx.Query(sqlQueryUpdateUser);
	if err != nil {
		return err;
	}

	err = trx.Commit();
	if err != nil {
		return err
	}

	return nil;
}

func (repo *onboardingRepo) UpdateDataUserLogin(loginStatus int, lastLogin string, updatedOn string) error {
	sqlQuery := fmt.Sprintf("UPDATE user_management.t_users SET login_status = %d, last_login = '%s', updated_on = '%s'", loginStatus, lastLogin, updatedOn);
	
	_, err := repo.db.Query(sqlQuery);
	if err != nil {
		return err;
	}

	return nil;
}

func (repo *onboardingRepo) UpdateUserLogout(currentDate string, DeviceID string, Msisdn string) error {
	var IdUser int
	trx, err := repo.db.Begin();
	if err != nil {
		return err;
	}

	sqlQueryUser := fmt.Sprintf("UPDATE user_management.t_users SET login_status = 0, updated_on = '%s' WHERE msisdn = '%s' RETURNING id_user", currentDate, Msisdn);
	
	err = repo.db.QueryRow(sqlQueryUser).Scan(&IdUser);
	if err != nil {
		trx.Rollback();
		return err;
	}

	sqlQueryLogin := fmt.Sprintf("UPDATE transaction.t_login_app SET is_logout = true, logout_on = '%s' WHERE id_user = %d AND device_id = '%s' AND is_logout = false", currentDate, IdUser, DeviceID);
	
	_, err = repo.db.Query(sqlQueryLogin);
	if err != nil {
		trx.Rollback();
		return err;
	}

	sqlQuerySession := fmt.Sprintf("UPDATE security.t_session_app SET status = 0, updated_on = '%s' WHERE id_user = %d AND device_id = '%s'", currentDate, IdUser, DeviceID);
	
	_, err = repo.db.Query(sqlQuerySession);
	if err != nil {
		trx.Rollback();
		return err;
	}

	trx.Commit();
	return nil;
}