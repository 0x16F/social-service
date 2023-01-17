package config

type Config struct {
	HTTP     HTTP     `yaml:"http"`
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`
	CF       CF       `yaml:"cf"`
}

type HTTP struct {
	Port int `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
	SSL      string `yaml:"ssl"`
}

type JWT struct {
	AccessSecret  string `yaml:"access"`
	RefreshSecret string `yaml:"refresh"`
}

type CF struct {
	Secret string `yaml:"secret"`
}
