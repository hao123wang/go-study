package discover

// 服务实例的元信息
type InstanceInfo struct {
	ID      string         `json:"id"`                // 服务实例id
	Name    string         `json:"name"`              // 服务名称
	Service string         `json:"service,omitempty"` // 服务发现时返回的服务名
	Address string         `json:"address"`           // ip地址
	Port    int            `json:"port"`              // 端口
	Check   `json:"check"` // 健康检查相关配置
}

type Check struct {
	HTTP     string `json:"http"`               // 健康检查接口
	Interval string `json:"interval,omitempty"` // consul主动调用健康检查接口的时间间隔，例如"10s"
	TTL      string `json:"ttl,omitempty"`      // 服务实例方主动上报心跳的时间间隔
}
