package config

type SystemEnv struct {
	Env               string `yaml:"env"`
	HealthCheckPeriod int    `yaml:"health_check_period"`
	ProxyHttp         struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Name string `yaml:"name"`
	} `yaml:"proxy_http"`
	Redis struct {
		Host     string  `yaml:"host"`
		Port     int     `yaml:"port"`
		Password *string `yaml:"password"`
		DB       int     `yaml:"db"`
	} `yaml:"redis"`
}

type ProxyMapping struct {
	Algorithm    string `yaml:"algorithm"`
	ProxyServers []Host `yaml:"proxy_servers"`
}

type Host struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Health   string `yaml:"health"`
	Protocol string `yaml:"protocol"`
}
