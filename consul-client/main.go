package main

import (
	"consul-client/discover"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	instanceID := "123456"
	// 创建Consul客户端
	consulClient := discover.NewConsulClient("127.0.0.1", 8500)
	ok := consulClient.RegisterWithHealthCheck("http服务", instanceID, "/health", "192.168.124.6", 8080, nil)
	if !ok {
		return
	}
	go http.ListenAndServe(":8080", nil)

	// 监听程序退出信号
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh
	consulClient.DeRegister(instanceID)
}
