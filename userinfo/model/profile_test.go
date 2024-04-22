package model

import "testing"

func TestProfile_UpdateFields(t *testing.T) {
	p := &Profile{
		Id:        1,
		UserId:    2,
		Username:  "",
		Birthday:  "123",
		Email:     "",
		AvatarUrl: "",
	}
	fields, args := p.UpdateFields()
	t.Logf("fields: %v", fields)
	t.Log(args)
}

func TestProfile_UpdateSql(t *testing.T) {
	p := &Profile{
		Id:        1,
		UserId:    2,
		Username:  "",
		Birthday:  "123",
		Email:     "",
		AvatarUrl: "",
	}
	fields, _ := p.UpdateFields()
	t.Log(p.UpdateSql(fields, "profile_tab"))
}

func TestProfile_InsertSql(t *testing.T) {
	p := &Profile{
		Id:        1,
		UserId:    2,
		Username:  "",
		Birthday:  "123",
		Email:     "",
		AvatarUrl: "",
	}
	fields, _ := p.UpdateFields()
	t.Log(p.InsertSql(fields, "profile_tab"))
}
