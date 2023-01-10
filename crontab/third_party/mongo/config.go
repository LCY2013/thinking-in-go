package _mongo

// MongoDBConfig mongodb 配置
type MongoDBConfig struct {
	Uri               string `json:"uri"`
	ConnectionTimeout int    `json:"connectionTimeout"`
}
