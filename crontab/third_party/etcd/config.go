package _etcd

// EtcdConfig etcd 配置信息
type EtcdConfig struct {
	Server struct {
		Endpoints   []string `json:"endpoints"`
		DialTimeout int64    `json:"dialTimeout"`
	} `json:"server"`
}
