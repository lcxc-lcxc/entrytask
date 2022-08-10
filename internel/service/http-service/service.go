package http_service

import (
	"entrytask/global"
	"entrytask/internel/dao"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

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
