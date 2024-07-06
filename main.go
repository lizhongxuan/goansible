package main

import (
	"fmt"
	"goansible/ansible"
	"os"
)

func main() {
	data, err := os.ReadFile("./example/playbook.yaml")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	fmt.Println("yaml: \n", string(data))
	p, err := ansible.GeneratePlaybook(
		ansible.WithPlaybooks(data),
	)
	if err != nil {
		fmt.Println("GeneratePlaybook err:", err)
		return
	}

	if err = p.Run(); err != nil {
		fmt.Println("Run err:", err)
		return
	}
	fmt.Println("Success.")
}
