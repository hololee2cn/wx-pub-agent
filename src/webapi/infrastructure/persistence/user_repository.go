package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/hololee2cn/pkg/ginx"
	"github.com/hololee2cn/wxpub/v1/src/webapi/domain/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB    *gorm.DB
	Redis *redis.UniversalClient
}

var defaultUserRepo *UserRepo

func NewUserRepo() {
	if defaultUserRepo == nil {
		defaultUserRepo = &UserRepo{
			DB:    CommonRepositories.DB,
			Redis: CommonRepositories.Redis,
		}
	}
}

func DefaultUserRepo() *UserRepo {
	return defaultUserRepo
}

func (a *UserRepo) IsExistUserMsgFromDB(ctx context.Context, fromUserName string, createTime int64) (bool, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("IsExistUserMsgFromDB traceID:%s", traceID)
	var user entity.User
	err := a.DB.Where("open_id = ? AND create_time = ?", fromUserName, createTime).First(&user).Error
	if err != nil {
		// 不存在记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("IsExistUserMsgFromDB record is not found,traceID:%s,err:%+v", traceID, err)
			return false, nil
		}
		log.Errorf("IsExistUserMsgFromDB failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	return true, nil
}

func (a *UserRepo) IsExistUserFromDB(ctx context.Context, fromUserName string) (bool, error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("IsExistUserFromDB traceID:%s", traceID)
	var user entity.User
	err := a.DB.Where("open_id = ?", fromUserName).First(&user).Error
	if err != nil {
		// 不存在记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("IsExistUserFromDB record is not found,traceID:%s,err:%+v", traceID, err)
			return false, nil
		}
		log.Errorf("IsExistUserFromDB failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	return true, nil
}

func (a *UserRepo) SaveUser(ctx context.Context, user entity.User, isUpdateAll bool) error {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("SaveUser traceID:%s", traceID)
	// 先查看是否有这用户，如果没有则创建，否则将创建时间和删除时间更新
	exist, err := a.IsExistUserFromDB(ctx, user.OpenID)
	if err != nil {
		log.Errorf("SaveUser IsExistUserFromDB failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	if !exist {
		if err = a.DB.Create(&user).Error; err != nil {
			log.Errorf("SaveUser create user failed,traceID:%s,err:%+v", traceID, err)
			return err
		}
	}

	if !isUpdateAll {
		err = a.UpdateUserTime(ctx, user)
	} else {
		err = a.DB.Model(&entity.User{}).Where("open_id=?", user.OpenID).Updates(&user).Error
	}
	if err != nil {
		log.Errorf("SaveUser UpdateUser failed, traceID: %s,err: %+v", traceID, err)
		return err
	}

	return nil
}

func (a *UserRepo) DelUser(ctx context.Context, user entity.User) error {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("DelUser traceID:%s", traceID)
	user.DeleteTime = time.Now().Unix()
	if err := a.DB.Model(&entity.User{}).Where("open_id = ?", user.OpenID).Updates(&user).Error; err != nil {
		log.Errorf("DelUser delete user failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *UserRepo) UpdateUserTime(ctx context.Context, user entity.User) error {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("UpdateUserTime traceID:%s", traceID)
	if err := a.DB.Model(&entity.User{}).Where("open_id = ?", user.OpenID).Updates(map[string]interface{}{
		"create_time": user.CreateTime,
		"delete_time": user.DeleteTime,
	}).Error; err != nil {
		log.Errorf("UpdateUserTime update user failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *UserRepo) GetUserByOpenID(ctx context.Context, openID string) (user entity.User, err error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("GetUserByID traceID:%s", traceID)
	if err = a.DB.Where("open_id = ?", openID).First(&user).Error; err != nil {
		log.Errorf("GetUserByOpenID get user by open_id failed,traceID: %s, err: %+v", traceID, err)
		return
	}
	return
}

func (a *UserRepo) ListUserByPhones(ctx context.Context, phones []string) (users []entity.User, err error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("ListUserByPhones traceID:%s", traceID)
	if err = a.DB.Where("phone IN (?)", phones).Find(&users).Error; err != nil {
		log.Errorf("ListUserByPhones get list user by phones failed,traceID:%s,err:%+v", traceID, err)
		return
	}
	return
}

func (a *UserRepo) IsExistPhone(ctx context.Context, phone string) (exist bool, err error) {
	traceID := ginx.ShouldGetTraceID(ctx)
	log.Debugf("IsExistPhone traceID:%s", traceID)
	var user entity.User
	if err = a.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		// 不存在记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("IsExistPhone find user by phone is not found,traceID:%s,err:%+v", traceID, err)
			return false, nil
		}
		log.Errorf("IsExistPhone find user by phone failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	return true, nil
}
