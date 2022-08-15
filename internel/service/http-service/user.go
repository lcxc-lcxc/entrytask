package http_service

import (
	pb "entrytask/internel/proto"
	"errors"
	"strings"
)

type UserRegisterRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=8,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=32"`
}

type UserRegisterResponse struct {
}

type UserLoginRequest struct {
	Username string `form:"username" json:"username" binding:"required,min=8,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=32"`
}

type UserLoginResponse struct {
	Username  string `json:"username"`
	SessionId string `json:"session_id"`
}

type UserAuthRequest struct {
	SessionId string
}

type AuthResponse struct {
	UserID   uint
	Username string
}

// GetCleanErr
// 由于 从grpc返回的错误都带有 rpc error = Unknown Desc = 前缀
// 该函数就是删除这个前缀
func GetCleanErr(err error) error {
	if err == nil {
		return err
	}
	errStr := err.Error()
	tmp1 := errStr[strings.Index(errStr, "=")+1 : len(errStr)]
	tmp2 := tmp1[strings.Index(tmp1, "=")+1 : len(tmp1)]
	return errors.New(tmp2)
}

// UserRegister 负责调用grpc 的service
func (svc *Service) UserRegister(request *UserRegisterRequest) (*UserRegisterResponse, error) {
	_, err := svc.GetUserRpcClient().Register(svc.ctx, &pb.RegisterRequest{
		Username: request.Username,
		Password: request.Password,
	})
	return &UserRegisterResponse{}, GetCleanErr(err)
}

// UserLogin 负责调用grpc 的service
func (svc *Service) UserLogin(request *UserLoginRequest) (*UserLoginResponse, error) {
	loginReply, err := svc.GetUserRpcClient().Login(svc.ctx, &pb.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, GetCleanErr(err)
	}
	return &UserLoginResponse{
		Username:  loginReply.Username,
		SessionId: loginReply.SessionId,
	}, nil
}

// AuthUser 负责调用grpc 的service
func (svc *Service) AuthUser(request *UserAuthRequest) (*AuthResponse, error) {
	authReply, err := svc.GetUserRpcClient().Auth(svc.ctx, &pb.AuthRequest{
		SessionId: request.SessionId,
	})
	if err != nil {
		return nil, GetCleanErr(err)
	}
	return &AuthResponse{
		UserID:   uint(authReply.ID),
		Username: authReply.Username,
	}, nil

}

var rpcUserServiceClient pb.UserServiceClient

// GetUserRpcClient 懒加载获取grpc的client
func (svc *Service) GetUserRpcClient() pb.UserServiceClient {
	if rpcUserServiceClient == nil {
		rpcUserServiceClient = pb.NewUserServiceClient(svc.rpcClient)
	}
	return rpcUserServiceClient

}
