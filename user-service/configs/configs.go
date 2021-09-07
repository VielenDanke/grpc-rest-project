package configs

type DatabaseConfig struct {
	URL string `json:"url"`
}

type Service struct {
	ConnUrl string `json:"conn_url"`
	Name    string `json:"name"`
}

type Config struct {
	HTTP     *ServerHTTP     `json:"http"`
	GRPC     *ServerGRPC     `json:"grpc"`
	DB       *DatabaseConfig `json:"database"`
	Services []*Service      `json:"services"`
}

func NewConfig() *Config {
	return &Config{
		HTTP:     &ServerHTTP{},
		GRPC:     &ServerGRPC{},
		DB:       &DatabaseConfig{},
		Services: make([]*Service, 0),
	}
}

type ServerHTTP struct {
	Addr string `json:"addr"`
}

type ServerGRPC struct {
	Addr string `json:"addr"`
}
