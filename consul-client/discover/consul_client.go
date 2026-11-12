package discover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ConsulClient struct {
	ServiceHost string `json:"host"` // consul服务端的ip
	ServicePort int    `json:"port"` // consul服务端的端口
}

func NewConsulClient(host string, port int) *ConsulClient {
	return &ConsulClient{
		ServiceHost: host,
		ServicePort: port,
	}
}

func (c *ConsulClient) RegisterWithHealthCheck(serviceName, instanceID, healthCheckURL string, instanceHost string, instancePort int, meta map[string]any) bool {
	// 封装服务实例元数据
	instanceInfo := &InstanceInfo{
		ID:      instanceID,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Check: Check{
			HTTP:     fmt.Sprintf("http://%s:%d%s", instanceHost, instancePort, healthCheckURL),
			Interval: "10s",
		},
	}

	bytesData, err := json.Marshal(instanceInfo)
	if err != nil {
		log.Printf("serialization instance info err: %v", err)
		return false
	}

	// 向Consul服务端发送注册请求
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/v1/agent/service/register", c.ServiceHost, c.ServicePort), bytes.NewReader(bytesData))
	if err != nil {
		log.Printf("http new request err: %v", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("register service %s err: %v", instanceID, err)
		return false
	}
	if resp.StatusCode == 200 {
		log.Printf("register service %s success", instanceID)
		return true
	} else {
		log.Printf("register service %s failed, code: %d", instanceID, resp.StatusCode)
		return false
	}
}

func (c *ConsulClient) DeRegister(instanceID string) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/v1/agent/service/deregister/%s", c.ServiceHost, c.ServicePort, instanceID), nil)
	if err != nil {
		log.Printf("http new request err: %v", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client do req err: %v", err)
		return
	}
	if resp.StatusCode == 200 {
		log.Printf("deregister service %s success", instanceID)
		return
	} else {
		log.Printf("deregister service %s failed", instanceID)
		return
	}
}

func (c *ConsulClient) Discover(serviceName string) []any {
	// 同上，根据服务名，发送对应的http请求获取服务列表
	// url: http://ip:port/v1/health/service/serviceName
	return nil
}
