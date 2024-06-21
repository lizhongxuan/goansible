package ansible

import "fmt"

func ExampleWhen() {
	// 假设的变量映射，用于模拟Ansible变量
	var vars = map[string]interface{}{
		"A": 1,
		"B": "redhat",
		"C": "kylin",
		"D": 2,
	}
	fmt.Printf("vars: %+v \n", vars)

	// 定义条件表达式
	conditions := []string{
		"A == 1 or (B == 'redhat' and C == 'kylin')",
		"A == 1 or (B == 'ubuntu2' and C == 'kylin')",
		"A == 1 and (B == 'ubuntu2' and C == 'kylin')",
		"A == 2 or (B == 'ubuntu' and C == 'kylin')",
		"D == 2 and B == 'redhat'",
		"D == 3 and B == 'redhat'",
	}

	// 测试用例
	for _, condition := range conditions {
		fmt.Printf("Testing condition: %s\n", condition)
		t := &Task{
			When: condition,
		}
		result := t.WhenFunc(vars)
		fmt.Println("Condition result:", result)
	}
	/*
		Testing condition: A == 1 or (B == 'redhat' and C == 'kylin')
		Condition result: true
		Testing condition: A == 1 or (B == 'ubuntu2' and C == 'kylin')
		Condition result: true
		Testing condition: A == 1 and (B == 'ubuntu2' and C == 'kylin')
		Condition result: false
		Testing condition: A == 2 or (B == 'ubuntu' and C == 'kylin')
		Condition result: false
		Testing condition: D == 2 and B == 'redhat'
		Condition result: true
		Testing condition: D == 3 and B == 'redhat'
		Condition result: false
	*/

}
