package domain

type (
	Config struct {
		Env            string `yaml:"env" env:"ENV" default:"prod"`
		Port           int    `yaml:"port" env:"PORT" default:"8080"`
		RedisAddr      string `yaml:"redis_addr" env:"REDIS_ADDR" default:"redis://:redispwd@localhost:6379"`
		NameSpace      string `yaml:"name_space" env:"NAMESPACE" default:"onetimepw"`
		MemoryCapacity int    `yaml:"memory_capacity" env:"MEMORY_CAPACITY" default:"10000"`
	}
)
