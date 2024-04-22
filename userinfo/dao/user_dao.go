package dao

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"user-server/model"
)

const TAB_NAME_USER = "user_tab"

type UserDao struct {
	db *DBMaster
}

func NewUserDao(db *DBMaster) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) GetUserByEmail(email string) (*model.User, error) {
	logger.Info("Call UserDao.GetUserByEmail, email: ", email)
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
		logger.Error("Fail to scan data, err: ", err.Error())
		return nil, err
	}
	logger.Info("Get user done, user: ", user)
	return user, nil
}

func (d *UserDao) Insert(user *model.User) error {
	logger.Info("Call UserDao.Insert, user: ", user)
	updateFields, args := user.UpdateFields()
	sqlString := user.InsertSql(updateFields, TAB_NAME_USER)
	_, err := d.db.Query(sqlString, args...)
	logger.Debug("sql: ", sqlString)
	if err != nil {
		logger.Error("Fail to insert into sql DB, err: ", err.Error())
		return err
	}
	logger.Info("Insert user into sql DB succeed.")
	return nil
}
