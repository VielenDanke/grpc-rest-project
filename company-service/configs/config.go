package configs

type Config struct {
	HTTP *ServerHTTP     `json:"http"`
	GRPC *ServerGRPC     `json:"grpc"`
	DB   *DatabaseConfig `json:"database"`
}

func NewConfig() *Config {
	return &Config{
		HTTP: &ServerHTTP{},
		GRPC: &ServerGRPC{},
		DB:   &DatabaseConfig{},
	}
}

type DatabaseConfig struct {
	URL string `json:"url"`
}

type ServerHTTP struct {
	Addr string `json:"addr"`
}

type ServerGRPC struct {
	Addr string `json:"addr"`
}
