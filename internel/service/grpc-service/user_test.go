package grpc_service

import (
	"context"
	"entrytask/internel/dao"
	"entrytask/internel/model"
	pb "entrytask/internel/proto"
	"github.com/agiledragon/gomonkey/v2"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	svc := NewUserService(context.Background())

	username := "test12345678"
	password := "test12345678"

	//Input
	request := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}

	t.Run("normal register", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "CreateUser", func(_ *dao.Dao, _, _ string) (*model.User, error) {
			return nil, nil
		})

		_, err := svc.Register(svc.ctx, request)
		if err != nil {
			t.Errorf("TestUserService_Register got error %v", err)
		}
	})

	t.Run("invalid register", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ string) (*model.User, error) {
			return nil, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "CreateUser", func(_, _ string) (*model.User, error) {
			return nil, nil
		})

		_, err := svc.Register(svc.ctx, request)
		if err == nil {
			t.Errorf("TestUserService_Register should return error but didn't")
		}

	})

}
