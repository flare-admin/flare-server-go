package net

import (
	"fmt"
	"net"
)

// CheckPortAvailability 检查端口是否被占用
func CheckPortAvailability(port int) bool {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return false // 端口被占用
	}
	ln.Close()  // 关闭监听器
	return true // 端口可用
}

// GetAnAvailablePort 获取一个可用的端口
func GetAnAvailablePort(port int) int {
	i := port
	for {
		if CheckPortAvailability(i) {
			return i
		}
		i++
	}
}

// GetWorkingDirectory 获取下一个接口
func GetFreePort(instanceCounter int) int {
	instanceCounter++
	prot := 1012 + instanceCounter
	return GetAnAvailablePort(prot)
}
