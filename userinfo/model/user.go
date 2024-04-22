package model

import "fmt"

type User struct {
	Id       uint64
	Name     string
	Password string
	Email    string
	Status   uint8
}

func (u *User) UpdateFields() ([]string, []any) {
	fields := make([]string, 0)
	args := make([]any, 0)
	if u.Id != 0 {
		fields = append(fields, "id")
		args = append(args, u.Id)
	}
	if u.Name != "" {
		fields = append(fields, "name")
		args = append(args, u.Name)
	}
	if u.Password != "" {
		fields = append(fields, "password")
		args = append(args, u.Password)
	}
	if u.Email != "" {
		fields = append(fields, "email")
		args = append(args, u.Email)
	}
	if u.Status != 0 {
		fields = append(fields, "status")
		args = append(args, u.Status)
	}
	return fields, args
}

func (p *User) InsertSql(fields []string, tabName string) string {
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
