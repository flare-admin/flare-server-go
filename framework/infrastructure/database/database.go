package database

import (
	"github.com/dtm-labs/rockscache"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/database/cache"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/snowflake_id"
	"github.com/flare-admin/flare-server-go/framework/pkg/hredis"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"time"
)

var ProviderSet = wire.NewSet(
	snowflake_id.NewSnowIdGen,
	NewHdbClient,
	NewDataBase,
	NewRc,
	database.NewDb,
	database.NewData,
	cache.NewCache,
	NewRedisClient,
	cache.NewCacheDecorator,
	NewTransactional,
)

func NewDataBase(data *database.Data) database.IDataBase {
	return data
}

func NewTransactional(data *database.Data) database.ITransactional {
	return data
}

func NewHdbClient(conf *configs.Data) (*hredis.RedisClient, func(), error) {
	return hredis.NewRedisClient(hredis.Option{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       int(conf.Redis.Db),
		Timeout:  conf.Redis.WriteTimeout,
	})
}

func NewRedisClient(hc *hredis.RedisClient) *redis.Client {
	return hc.GetClient()
}

func NewRc(rdb *hredis.RedisClient) *rockscache.Client {
	// 强一致性缓存，当一个key被标记删除，其他请求线程会被锁住轮询直到新的key生成，适合各种同步的拉取, 如果弱一致可能导致拉取还是老数据，毫无意义
	options := rockscache.NewDefaultOptions()
	options.StrongConsistency = true
	options.Delay = time.Millisecond * 1
	Rc := rockscache.NewClient(rdb.GetClient(), options)
	Rc.Options.StrongConsistency = true
	return Rc
}
