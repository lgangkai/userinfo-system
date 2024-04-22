package model

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"time"
)

type Profile struct {
	Id        uint64 `json:"id"`
	UserId    uint64 `json:"user_id"`
	Username  string `json:"username"`
	Birthday  string `json:"birthday"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

func (p *Profile) UpdateFields() ([]string, []any) {
	fields := make([]string, 0)
	args := make([]any, 0)
	if p.Id != 0 {
		fields = append(fields, "id")
		args = append(args, p.Id)
	}
	if p.UserId != 0 {
		fields = append(fields, "user_id")
		args = append(args, p.UserId)
	}
	if p.Username != "" {
		fields = append(fields, "username")
		args = append(args, p.Username)
	}
	if p.Birthday != "" {
		fields = append(fields, "birthday")
		date, err := time.Parse("2006-01-02", p.Birthday)
		if err != nil {
			logger.Warn("Parse birthday failed, using default date. Err: ", err.Error())
			date = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
		}
		args = append(args, date)
	}
	if p.Email != "" {
		fields = append(fields, "email")
		args = append(args, p.Email)
	}
	if p.AvatarUrl != "" {
		fields = append(fields, "avatar_url")
		args = append(args, p.AvatarUrl)
	}
	return fields, args
}

func (p *Profile) UpdateSql(fields []string, tabName string) string {
	setSql := "SET "
	for i, field := range fields {
		setSql = setSql + field + "=?"
		if i < len(fields)-1 {
			setSql += ","
		}
	}
	sqlString := fmt.Sprintf("UPDATE %v %v WHERE user_id=?", tabName, setSql)
	return sqlString
}

func (p *Profile) InsertSql(fields []string, tabName string) string {
	fieldsSql := ""
	for i, field := range fields {
		fieldsSql += field
		if i < len(fields)-1 {
			fieldsSql += ","
		}
	}
	valuesSql := ""
	for i := range fields {
		valuesSql += "?"
		if i < len(fields)-1 {
			valuesSql += ","
		}
	}
	sqlString := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", tabName, fieldsSql, valuesSql)
	return sqlString
}
