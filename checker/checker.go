package checker

type Checker struct {
	Name   string
	CMD    string
	Output string
	Regex  string
	Hooks  map[string]string
}
