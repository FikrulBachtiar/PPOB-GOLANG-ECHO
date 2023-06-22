package repository

import (
	"database/sql"
	"fmt"
	"ppob/app/domain"
)

type UserRepo interface {
	GetUserDetail(IdUser int) (*domain.UserDetailResponse, error)
	GetBalanceUser(IdUser int) (*domain.SpGetBalance, error)
	GetPointUser(IdUser int) (*domain.SpGetPoint, error)
	IsKtpVerify(IdUser int) (bool, error)
	GetPointType() ([]domain.PointTypeResponse, error)
	GetListPoint(IdUser int, IdPointType int) ([]domain.ListPointResponse, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	};
}

func (repo *userRepo) GetUserDetail(IdUser int) (*domain.UserDetailResponse, error) {
	var result domain.UserDetailResponse;

	sqlQuery := fmt.Sprintf(`
	SELECT
		tu.msisdn,
		tu.fullname,
		COALESCE(tue.email, '') AS email,
		CASE WHEN tu.is_email_verified THEN 'Terverifikasi' ELSE 'Belum Terverifikasi' END AS email_status,
		CASE WHEN tu.is_email_verified THEN 'VERIFIED' ELSE 'UNVERIFIED' END AS email_status_en,
		CASE WHEN tu.is_ktp_verified THEN 'Terverifikasi' ELSE 'Belum Terverifikasi' END AS ktp_status,
		CASE WHEN tu.is_ktp_verified THEN 'VERIFIED' ELSE 'UNVERIFIED' END AS ktp_status_en,
		tu.uuid,
		tu.is_question,
		COALESCE(tup.url_pro_pic, '') AS url_pro_pic
	FROM
		user_management.t_users tu
	LEFT JOIN
		"security".t_user_email tue ON tue.id_user = tu.id_user
	LEFT JOIN
		user_management.t_user_picture tup ON tup.id_user = tu.id_user
	WHERE
		tu.id_user = %d`,
	IdUser);

	err := repo.db.QueryRow(sqlQuery).Scan(&result.Msisdn, &result.FullName, &result.Email, &result.EmailStatus, &result.EmailStatusEn, &result.KtpStatus, &result.KtpStatusEn, &result.Uuid, &result.Question, &result.UrlProPic);
	if err != nil {
		return nil, err;
	}

	return &result, nil;

}

func (repo *userRepo) GetBalanceUser(IdUser int) (*domain.SpGetBalance, error) {
	var result domain.SpGetBalance;
	sqlQuery := fmt.Sprintf("SELECT * FROM user_management.sp_get_balance_user(%d)", IdUser);

	err := repo.db.QueryRow(sqlQuery).Scan(&result.Currency, &result.Value);
	if err != nil {
		return nil, err;
	}
	
	return &result, nil;
}

func (repo *userRepo) GetPointUser(IdUser int) (*domain.SpGetPoint, error) {
	var result domain.SpGetPoint;
	sqlQuery := fmt.Sprintf("SELECT * FROM user_management.sp_get_point_user(%d)", IdUser);

	err := repo.db.QueryRow(sqlQuery).Scan(&result.Currency, &result.Value, &result.IsTap);
	if err != nil {
		return nil, err;
	}
	
	return &result, nil;
}

func (repo *userRepo) IsKtpVerify(IdUser int) (bool, error) {
	var isVerify bool
	sqlQuery := fmt.Sprintf("SELECT is_ktp_verified FROM user_management.t_users WHERE id_user = %d", IdUser);

	err := repo.db.QueryRow(sqlQuery).Scan(&isVerify);
	if err != nil {
		return false, err;
	}

	return isVerify, nil;
}

func (repo *userRepo) GetPointType() ([]domain.PointTypeResponse, error) {
	sqlQuery := "SELECT id_point_type, point_type_name, point_type_name_en FROM master.t_point_type WHERE status = 1";

	rows, err := repo.db.Query(sqlQuery);
	if err != nil {
		return nil, err;
	}

	defer rows.Close();
	var results []domain.PointTypeResponse;

	for rows.Next() {
		var result domain.PointTypeResponse;
		err := rows.Scan(&result.IdPointType, &result.PointTypeName, &result.PointTypeNameEn);
		if err != nil {
			return nil, err;
		}

		results = append(results, result);
	}

	return results, nil;
}

func (repo *userRepo) GetListPoint(IdUser int, IdPointType int) ([]domain.ListPointResponse, error) {
	sqlQuery := fmt.Sprintf(`
	select
		ttp.product_category,
		tpt.direction,
		tc.currency,
		tup.point_total,
		ttp.created_on
	from
		user_management.t_user_point tup
	left join master.t_currency tc on
		tc.id_currency = tup.id_currency
	left join master.t_point_type tpt on
		tpt.id_point_type = tup.id_point_type
	left join "transaction".t_trx_payment ttp on
		ttp.invoice = tup.transaction_invoice
	where
		tup.id_user = %d
		and tup.status = 1
		and tup.id_point_type = %d
	order by
		ttp.created_on desc`,
	IdUser, IdPointType);

	rows, err := repo.db.Query(sqlQuery);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();

	var results []domain.ListPointResponse;

	for rows.Next() {
		var result domain.ListPointResponse;
		err := rows.Scan(&result.ProductCategory, &result.Direction, &result.Currency, &result.PointTotal, &result.TransactionDate);
		if err != nil {
			return nil, err;
		}

		results = append(results, result);
	}
	
	return results, nil;

}