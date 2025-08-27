package idempotence

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/hredis"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

type IdempotencyTool interface {
	// Check 检查消息是否已经处理过
	Check(ctx context.Context, topic, channel, messageID string) (bool, error)

	// MarkProcessed 标记消息已处理
	MarkProcessed(ctx context.Context, channel, group, messageID string) error
}

type MessageIdempotence struct {
	data database.IDataBase
	rdb  *redis.Client
}

func NewIdempotencyTool(data database.IDataBase, rdb *hredis.RedisClient) IdempotencyTool {
	// 同步表
	tables := []interface{}{
		&IdempotencyRecord{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables to mysql error: %v", err)
	}
	return &MessageIdempotence{
		data: data,
		rdb:  rdb.GetClient(),
	}
}

// Check 幂等性检查
func (r *MessageIdempotence) Check(ctx context.Context, topic, channel, messageID string) (bool, error) {
	redisKey := r.generateRedisKey(topic, channel, messageID)

	// 1. 先检查 Redis
	exists, err := r.rdb.Exists(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return true, nil // 消息已处理
	}
	// 2. 检查 PostgreSQL
	var record IdempotencyRecord
	result := r.data.DB(ctx).Where("topic = ? AND channel = ? AND message_id = ?", topic, channel, messageID).First(&record)
	if result.RowsAffected > 0 {
		// 将记录缓存到 Redis
		r.rdb.Set(ctx, redisKey, 1, 24*time.Hour)
		return true, nil
	}

	return false, nil // 消息未处理
}

// MarkProcessed 消息处理函数
func (r *MessageIdempotence) MarkProcessed(ctx context.Context, topic, channel, messageID string) error {
	redisKey := r.generateRedisKey(topic, channel, messageID)

	// 1. 将记录存入 PostgreSQL
	record := IdempotencyRecord{
		Topic:     topic,
		Channel:   channel,
		MessageID: messageID,
		Status:    "success",
		CreatedAt: utils.GetTimeNow(),
	}
	if err := r.data.DB(ctx).Create(&record).Error; err != nil {
		return err
	}

	// 2. 缓存到 Redis
	r.rdb.Set(ctx, redisKey, 1, 24*time.Hour)

	return nil
}

// 生成 Redis key
func (r *MessageIdempotence) generateRedisKey(topic, group, messageID string) string {
	return fmt.Sprintf("message:topic:%s:group:%s:msgid:%s", topic, group, messageID)
}
