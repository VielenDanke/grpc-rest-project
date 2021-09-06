package configs

type Config struct {
	HTTP *ServerHTTP `json:"http"`
	GRPC *ServerGRPC `json:"grpc"`
}

func NewConfig() *Config {
	return &Config{
		HTTP: &ServerHTTP{},
		GRPC: &ServerGRPC{},
	}
}

type ServerHTTP struct {
	Addr string `json:"addr"`
}

type ServerGRPC struct {
	Addr string `json:"addr"`
}
