package grpc_service

import (
	"context"
	"encoding/json"
	"entrytask/global"
	"entrytask/internel/constant"
	"entrytask/internel/dao"
	pb "entrytask/internel/proto"
	"entrytask/pkg/utils"
	"errors"
	"fmt"
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

type UserService struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.RedisClient
	pb.UnimplementedUserServiceServer
}

func NewUserService(ctx context.Context) UserService {
	return UserService{
		ctx:   ctx,
		dao:   dao.NewDBClient(global.DBEngine),
		cache: dao.NewCache(global.RedisClient),
	}
}

func (svc UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	_, err := svc.dao.GetUserByName(req.GetUsername())

	if err == nil || !errors.Is(gorm.ErrRecordNotFound, err) { //代表找到数据

		return nil, fmt.Errorf("Register : user already exists")
	}

	hash := utils.Hash(req.Password)
	_, err = svc.dao.CreateUser(req.Username, hash)
	if err != nil {
		log.Printf("Register : create user failed: %v", err)
		return nil, err
	}
	return &pb.RegisterReply{}, nil

}

func (svc UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	//1 从数据库获取user
	dbUser, err := svc.dao.GetUserByName(req.Username)
	if err != nil { //用户不存在或其他错误
		return nil, err
	}
	//2 校验数据库里面的密码和输入的密码
	verify := utils.HashVerify(dbUser.Password, req.Password)
	if !verify { //密码错误
		return nil, err
	}
	// 3 生成缓存的结构
	cacheUser := pb.AuthReply{
		ID:       uint64(dbUser.ID),
		Username: dbUser.Username,
	}
	cacheUserJson, err := json.Marshal(cacheUser)
	if err != nil {
		log.Printf("cache session_id failed : %v", err)
		return nil, err
	}
	// 4 生成session_id并存进redis
	sessionId := constant.SESSION_ID + "_" + uuid.NewString()
	err = cache.New[string](svc.cache.RedisStore).Set(svc.ctx, sessionId, string(cacheUserJson), store.WithExpiration(time.Hour))
	if err != nil {
		log.Printf("cache session_id failed : %v", err)
		return nil, err
	}

	//5 返回
	return &pb.LoginReply{
		Username:  dbUser.Username,
		SessionId: sessionId,
	}, nil
}

func (svc UserService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthReply, error) {
	cacheUserJson, err := cache.New[string](svc.cache.RedisStore).Get(svc.ctx, req.SessionId)
	if err != nil {
		log.Println("auth failed : get redis session message failed")
		return nil, err
	}
	var authReply pb.AuthReply
	err = json.Unmarshal([]byte(cacheUserJson), &authReply)
	if err != nil {
		log.Println("auth failed : unmarshal session message failed")
		return nil, err
	}
	return &authReply, nil

}
