package ansible

type Middleware struct {
	Kind     string                 `yaml:"kind"`
	Hosts    []*Host                `yaml:"hosts"`
	Vip      string                 `yaml:"vip"`
	Port     string                 `yaml:"port"`
	Username string                 `yaml:"username"`
	Password string                 `yaml:"password"`
	Args     map[string]interface{} `yaml:"args"`
}
