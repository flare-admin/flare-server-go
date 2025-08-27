package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	ts "github.com/flare-admin/flare-server-go/framework/support/systask/service"
)

type SysCronService struct {
	tm ts.ITaskManager
}

func NewSysCronService(tm ts.ITaskManager) (*SysCronService, func(), error) {

	// 启动任务管理器
	err := tm.Initialize()
	if err != nil {
		hlog.Fatalf("task manager initialization failed: %v", err)
		return nil, nil, err
	}
	clumpfunc := func() {
		tm.Stop()
	}
	tm.Start()
	return &SysCronService{
		tm: tm,
	}, clumpfunc, nil
}
func (s *SysCronService) Start() {
	// 启动cron调度器
	// s.cron.Start()
	s.register()
}

func (s *SysCronService) Test(data map[string]string) error {
	hlog.Debugf("test task")
	return nil
}

func (s *SysCronService) register() {
	s.tm.RegisterHandler("test", s.Test)
}
