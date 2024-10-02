package domain

type (
	Config struct {
		Env   string `yaml:"env" env:"ENV" default:"prod"`
		Port  int    `yaml:"port" env:"PORT" default:"8080"`
		Redis struct {
			Addr      string `yaml:"addr" env:"REDIS_ADDR" default:"localhost:6379"`
			NameSpace string `yaml:"name_space" env:"REDIS_NAMESPACE" default:"onetimepw"`
		}
	}
)
