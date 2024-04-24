package dao

import (
	"context"
	"fmt"
	"loggers"
	"user-server/model"
)

const TAB_NAME_USER = "user_tab"

type UserDao struct {
	db     *DBMaster
	logger *logger.Logger
}

func NewUserDao(db *DBMaster, logger *logger.Logger) *UserDao {
	return &UserDao{
		db:     db,
		logger: logger,
	}
}

func (d *UserDao) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	d.logger.Info(ctx, "Call UserDao.GetUserByEmail, email: ", email)
	user := &model.User{}
	sqlString := fmt.Sprintf("SELECT id, name, password, email, status"+
		" FROM %v WHERE email = ?", TAB_NAME_USER)
	row := d.db.QueryRow(sqlString, email)

	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Status,
	)
	if err != nil {
		d.logger.Error(ctx, "Fail to scan data, err: ", err.Error())
		return nil, err
	}
	d.logger.Info(ctx, "Get user done, user: ", *user)
	return user, nil
}

func (d *UserDao) Insert(ctx context.Context, user *model.User) error {
	d.logger.Info(ctx, "Call UserDao.Insert, user: ", user)
	updateFields, args := user.UpdateFields()
	sqlString := user.InsertSql(updateFields, TAB_NAME_USER)
	rows, err := d.db.Query(sqlString, args...)
	defer rows.Close()
	d.logger.Debug(ctx, "sql: ", sqlString)
	if err != nil {
		d.logger.Error(ctx, "Fail to insert into sql DB, err: ", err.Error())
		return err
	}
	d.logger.Info(ctx, "Insert user into sql DB succeed.")
	return nil
}
