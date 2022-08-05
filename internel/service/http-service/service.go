package http_service

import (
	"entrytask/global"
	"entrytask/internel/dao"
	"entrytask/internel/redisCache"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Service struct {
	ctx       context.Context
	dao       *dao.Dao
	cache     *redisCache.RedisClient
	rpcClient *grpc.ClientConn
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx:       ctx,
		dao:       dao.NewDBClient(global.DBEngine),
		cache:     redisCache.NewCache(global.RedisClient),
		rpcClient: global.GRPCClient,
	}
}
