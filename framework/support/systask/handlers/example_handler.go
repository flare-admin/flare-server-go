package handlers

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

type ExampleHandler struct{}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) Execute(args map[string]string) error {
	message := args["message"]
	if message == "" {
		message = "默认消息"
	}

	fmt.Printf("[%s] 执行示例任务: %s\n", utils.GetTimeNow().Format("2006-01-02 15:04:05"), message)
	return nil
}
