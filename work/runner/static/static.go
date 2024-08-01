package static

type SandboxGlobalConfigurations struct {
	App struct {
		Port  int    `yaml:"port"`
		Debug bool   `yaml:"debug"`
		Key   string `yaml:"key"`
	} `yaml:"app"`
	MaxWorkers    int  `yaml:"max_workers"`
	MaxRequests   int  `yaml:"max_requests"`
	WorkerTimeout int  `yaml:"worker_timeout"`
	EnableNetwork bool `yaml:"enable_network"`
	Proxy         struct {
		Socks5 string `yaml:"socks5"`
		Https  string `yaml:"https"`
		Http   string `yaml:"http"`
	} `yaml:"proxy"`
}

var difySandboxGlobalConfigurations SandboxGlobalConfigurations

// avoid global modification, use value copy instead
func GetDifySandboxGlobalConfigurations() SandboxGlobalConfigurations {
	return difySandboxGlobalConfigurations
}
