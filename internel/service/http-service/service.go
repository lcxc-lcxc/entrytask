package http_service

import (
	"entrytask/global"
	"entrytask/internel/dao"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Service 对context和dao以及redis-client的封装，方便上层调用
type Service struct {
	ctx context.Context
	dao *dao.Dao

	rpcClient *grpc.ClientConn
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx:       ctx,
		dao:       dao.NewDAO(global.DBEngine, global.RedisClient),
		rpcClient: global.GRPCClient,
	}
}
