package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"loggers"
	"math/rand"
	"time"
	"user-server/model"
)

const TAB_NAME_PROFILE = "profile_tab"
const (
	REDIS_KEY_GET_PROFILE_PREFIX           = "userinfo:get_profile:"
	REDIS_KEY_GET_PROFILE_EXPIRE_BASE      = time.Second * 60
	REDIS_KEY_GET_PROFILE_EXPIRE_MAX_SHIFT = 30
)

type ProfileDao struct {
	dbMaster *DBMaster
	dbSlave  *DBSlave
	dbRedis  *redis.ClusterClient
	logger   *logger.Logger
}

func NewProfileDao(dbMaster *DBMaster, dbSlave *DBSlave, dbRedis *redis.ClusterClient, logger *logger.Logger) *ProfileDao {
	return &ProfileDao{
		dbMaster: dbMaster,
		dbSlave:  dbSlave,
		dbRedis:  dbRedis,
		logger:   logger,
	}
}

func (d *ProfileDao) GetProfileById(ctx context.Context, userId uint64) (*model.Profile, error) {
	d.logger.Info(ctx, "Call ProfileDao.GetProfile.")
	profile := &model.Profile{}

	// 1. try to get value from redis first.
	rKey := fmt.Sprintf("%v%d", REDIS_KEY_GET_PROFILE_PREFIX, userId)
	profileStr, err := d.dbRedis.Get(ctx, rKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			d.logger.Info(ctx, "Can not find in cache, go to sql DB.")
		} else {
			d.logger.Error(ctx, "Can not get from cache, err: ", err.Error(), ". Go to sql DB")
		}
	} else {
		d.logger.Info(ctx, "Get profile json from cache, profile: ", profileStr)
		err = json.Unmarshal([]byte(profileStr), profile)
		if err != nil {
			d.logger.Error(ctx, "json.Unmarshal failed, err: ", err.Error(), ". Go to sql DB")
		} else {
			d.logger.Info(ctx, "Get profile from cache succeeded.")
			return profile, nil
		}
	}

	// 2. get value from mysql-slave if not found in redis.
	sqlString := fmt.Sprintf("SELECT id, user_id, username, birthday, email, avatar_url"+
		" FROM %v WHERE user_id = ?", TAB_NAME_PROFILE)
	row := d.dbSlave.QueryRow(sqlString, userId)

	err = row.Scan(
		&profile.Id,
		&profile.UserId,
		&profile.Username,
		&profile.Birthday,
		&profile.Email,
		&profile.AvatarUrl,
	)
	if err != nil {
		d.logger.Error(ctx, "Fail to scan data, err: ", err.Error())
		return nil, err
	}
	d.logger.Info(ctx, "Get profile done, profile: ", profile)

	// 3. write profile as json string back to cache.
	pBytes, err := json.Marshal(profile)
	if err != nil {
		d.logger.Error(ctx, "json.Marshal failed, err: ", err.Error(), ". It will not be saved to cache.")
		return profile, nil
	}
	// set key expiration time as base time plus random time to avoid cache avalanche.
	randExp := time.Duration(rand.Intn(REDIS_KEY_GET_PROFILE_EXPIRE_MAX_SHIFT)) * time.Second
	err = d.dbRedis.Set(ctx, rKey, string(pBytes), REDIS_KEY_GET_PROFILE_EXPIRE_BASE+randExp).Err()
	if err != nil {
		d.logger.Error(ctx, "redis set failed, err: ", err.Error(), ". It will not be saved to cache.")
		return profile, nil
	}
	d.logger.Info(ctx, "Write back to redis done. key: ", rKey, ", value: ", string(pBytes),
		", expiration: ", REDIS_KEY_GET_PROFILE_EXPIRE_BASE+randExp)

	return profile, nil
}

func (d *ProfileDao) Update(ctx context.Context, userId uint64, profile *model.Profile) error {
	// use cache aside pattern to update DB and then delete from cache.
	// 1. update data to mysql-master.
	d.logger.Info(ctx, "Call ProfileDao.Update.")
	updateFields, args := profile.UpdateFields()
	sqlString := profile.UpdateSql(updateFields, TAB_NAME_PROFILE)
	d.logger.Debug(ctx, "sql: ", sqlString)
	_, err := d.dbMaster.Exec(sqlString, append(args, userId)...)
	if err != nil {
		d.logger.Error(ctx, "Fail to update to sql DB, err: ", err.Error())
		return err
	}
	d.logger.Info(ctx, "Update profile to sql DB succeed.")

	// 2. delete data from redis.
	d.deleteFromCache(ctx, userId)

	return nil
}

func (d *ProfileDao) Delete(ctx context.Context, userId uint64) error {
	// use cache aside pattern to delete from DB and then delete from cache.
	// 1. delete data from mysql-master.
	d.logger.Info(ctx, "Call ProfileDao.Delete.")
	sqlString := fmt.Sprintf("DELETE FROM %v WHERE user_id = ?", TAB_NAME_PROFILE)
	_, err := d.dbMaster.Exec(sqlString, userId)
	if err != nil {
		d.logger.Error(ctx, "Fail to delete from sql DB, err: ", err.Error())
		return err
	}
	d.logger.Info(ctx, "Delete profile from sql DB succeed.")

	// 2. delete data from redis.
	d.deleteFromCache(ctx, userId)

	return nil
}

func (d *ProfileDao) Insert(ctx context.Context, profile *model.Profile) error {
	// we don't operate cache in insert. Cache data will be load when read.
	d.logger.Info(ctx, "Call ProfileDao.Insert, profile: ", profile)
	updateFields, args := profile.UpdateFields()
	sqlString := profile.InsertSql(updateFields, TAB_NAME_PROFILE)
	_, err := d.dbMaster.Exec(sqlString, args...)
	d.logger.Debug(ctx, "sql: ", sqlString)
	if err != nil {
		d.logger.Error(ctx, "Fail to insert into sql DB, err: ", err.Error(), " sql: ", sqlString, " args: ", args)
		return err
	}
	d.logger.Info(ctx, "Insert profile into sql DB succeed.")

	return nil
}

func (d *ProfileDao) deleteFromCache(ctx context.Context, userId uint64) {
	rKey := fmt.Sprintf("%v%d", REDIS_KEY_GET_PROFILE_PREFIX, userId)
	err := d.dbRedis.Del(ctx, rKey).Err()
	// if delete failed, other process might read the dirty data.
	// since we have set an expiration time on it, it'll be eventually consist.
	if err != nil {
		d.logger.Error(ctx, "Fail to delete from cache, err: ", err.Error())
	}
	d.logger.Info(ctx, "Delete profile from cache succeed.")
}
