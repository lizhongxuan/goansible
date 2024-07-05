package ansible

import "go-ansible/model"

type Middleware struct {
	Kind     string                 `yaml:"kind"`
	Hosts    []*model.Host          `yaml:"hosts"`
	Vip      string                 `yaml:"vip"`
	Port     string                 `yaml:"port"`
	Username string                 `yaml:"username"`
	Password string                 `yaml:"password"`
	Args     map[string]interface{} `yaml:"args"`
}
