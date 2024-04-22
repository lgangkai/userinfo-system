package dao

import (
	"context"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestProfileDao_GetProfileById(t *testing.T) {
	dbMaster, err := sql.Open("mysql", "root:qwer1234@tcp(127.0.0.1:13306)/userinfo")
	dbSlave, err := sql.Open("mysql", "root:qwer1234@tcp(127.0.0.1:23306)/userinfo")
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"127.0.0.1:6371", "127.0.0.1:6372", "127.0.0.1:6373", "127.0.0.1:6374", "127.0.0.1:6375", "127.0.0.1:6376"},
	})
	err = rdb.Set(context.Background(), "q", "1", 0).Err()
	if err != nil {
		t.Errorf(err.Error())
		//return
	}
	profileDao := NewProfileDao(&DBMaster{dbMaster}, &DBSlave{dbSlave}, rdb)
	profile, err := profileDao.GetProfileById(context.Background(), 1)
	if err != nil {
		t.Errorf("call profileDao.GetProfile failed!")
		return
	}
	if profile.Id != 1 || profile.UserId != 1 {
		t.Errorf("got unexpected result!")
		return
	}
}
