package model

type Host struct {
	IP           string `yaml:"ip"`
	Port         string `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	SudoPassword string `yaml:"sudo_password"`
}
