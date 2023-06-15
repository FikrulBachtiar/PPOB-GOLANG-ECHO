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
}

type onboardingRepo struct {
	db *sql.DB
}

func NewOnboardingRepo(db *sql.DB) OnboardingRepo {
	return &onboardingRepo{
		db: db,
	}
}
	
func (onboardRepo *onboardingRepo) GetUserByNumber(msisdn string) (*domain.GetUserByNumber, error) {
	var result domain.GetUserByNumber
	sqlQuery := fmt.Sprintf(`
	SELECT
		tu.id_user,
		tu.fullname,
		tu.is_question,
		tu.login_status,
		tu.last_login,
		tu.user_status_code,
		coalesce(tue.is_verified, false) as email_status
	FROM
		user_management.t_users tu
	LEFT JOIN
		"security".t_user_email tue
	ON
		tue.id_user = tu.id_user
	WHERE
		tu.msisdn = '%s'`,
	msisdn);

	err := onboardRepo.db.QueryRow(sqlQuery).Scan(&result.IdUser, &result.Fullname, &result.IsQuestion, &result.LoginStatus, &result.LastLogin, &result.UserStatusCode, &result.EmailStatus);
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil;
		}
		return nil, err;
	}

	return &result, nil;
}
  
func (onboardRepo *onboardingRepo) CreateUser(msisdn string, uuid string, fullName string, user_code string, created_on string, user_status_code string) (*domain.InsertUser, error) {
	var result domain.InsertUser;
	sqlQuery := fmt.Sprintf(`INSERT INTO user_management.t_users (msisdn, uuid, fullname, user_code, created_on, user_status_code) VALUES ('%s', '%s', '%s', '%s', '%s', '%s') RETURNING is_question, is_email_verified`, msisdn, uuid, fullName, user_code, created_on, user_status_code);

	err := onboardRepo.db.QueryRow(sqlQuery).Scan(&result.IsQuestion, &result.IsEmailVerified);
	if err != nil {
		return nil, err;
	}

	return &result, nil;
}

func (onboardRepo *onboardingRepo) CheckLoginDevice(IdUser int) (int, error) {
	var resultNumber int
	sqlQuery := fmt.Sprintf("SELECT count(*) as resultNumber FROM transaction.t_login_app WHERE id_user = %d AND is_logout = TRUE", IdUser);

	err := onboardRepo.db.QueryRow(sqlQuery).Scan(&resultNumber);
	if err != nil {
		return 0, err;
	}

	return resultNumber, nil;
}