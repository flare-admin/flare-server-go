package mq

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"golang.org/x/exp/rand"
)

// GetDeadLetterQueueName 获取死信队列ID
func GetDeadLetterQueueName(topic, group string) string {
	return fmt.Sprintf("%s_%s_deadletter_queue", topic, group)
}

func GetSubscribeId(topic, group string) string {
	// 创建一个新的随机数生成器，使用当前时间作为种子
	r := rand.New(rand.NewSource(uint64(utils.GetTimeNow().UnixNano())))
	// 生成一个随机数，范围可以根据需要调整
	randomNum := r.Intn(1000000)
	// 返回带有随机数的字符串
	return fmt.Sprintf("%s_%s_%d", topic, group, randomNum)
}

func GetSubscribeBaseId(topic, group string) string {
	return fmt.Sprintf("%s_%s", topic, group)
}
