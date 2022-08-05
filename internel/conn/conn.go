package conn

import (
	"entrytask/pkg/setting"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func NewRPCClient(setting *setting.RpcClientSetting) (*grpc.ClientConn, error) {
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx, setting.RPCHost, opts...)
	return conn, err

}

func NewCacheClient(cacheSetting *setting.CacheSetting) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cacheSetting.Host,
		DB:   cacheSetting.DBIndex,
	})
	return client, nil
}

func NewDBEngine(dbSetting *setting.DBSetting) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		dbSetting.Username,
		dbSetting.Password,
		dbSetting.Host,
		dbSetting.DBName,
		dbSetting.Charset,
		dbSetting.ParseTime,
	)))
	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB()
	//连接可复用的最大时间
	sqlDb.SetConnMaxIdleTime(time.Minute * 3)
	sqlDb.SetMaxIdleConns(dbSetting.MaxIdleConns)
	sqlDb.SetMaxOpenConns(dbSetting.MaxOpenConns)
	return db, nil
}
