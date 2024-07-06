package ansible

import (
	"context"
	"errors"
	"fmt"
	"goansible/module"
	"goansible/work"
	"strconv"
	"strings"
	"time"
)

// Task represents a single task in an Ansible playbook.
type Task struct {
	Name         string                 `yaml:"name"` // 任务名称
	ShellBody    string                 `yaml:"shell_body,omitempty"`
	Module       module.Module          `yaml:",inline"`
	Args         map[string]interface{} `yaml:"args,omitempty"`        // 模块参数，内联以适应各种模块参数
	Register     string                 `yaml:"register,omitempty"`    // 任务输出结果的变量名
	DelegateTo   string                 `yaml:"delegate_to,omitempty"` // 将任务委派给其他主机执行
	IgnoreErrors bool                   `yaml:"ignore_errors"`         // 是否忽略任务错误
	Retry        int64                  `yaml:"retry,omitempty"`       // 重试次数, 默认为0
	Delay        int64                  `yaml:"delay,omitempty"`       // 重试之间的延迟时间
	Until        string                 `yaml:"until,omitempty"`       // 重试的条件
	WithItems    []interface{}          `yaml:"with_items,omitempty"`  // 循环关键字，用于循环执行任务,生成参数{{ item }}
	Loop         []interface{}          `yaml:"loop,omitempty"`        // 循环关键字,通WithItems
	Notify       []string               `yaml:"notify,omitempty"`      // 通知处理程序
	When         string                 `yaml:"when,omitempty"`        // 条件语句，用于有条件地执行任务
	ShowShell    bool                   `yaml:"show_shell,omitempty"`  // 是否打印shell
	Become       bool                   `yaml:"become"`                // 是否提升权限（类似于 sudo）
	Timeout      int64                  `yaml:"timeout,omitempty"`
	PreProcess   *Process               `yaml:"out_put,omitempty"`
	NotWait      bool                   `yaml:"not_wait,omitempty"`
	ModuleObject module.ModuleInterface
	Worker       work.Worker
}

type Process struct {
	PID       int64 // 最后一个shell的执行进程ID
	StateCode int
	Stdout    string
	Register  map[string]string
}

func (t *Task) run(ctx context.Context, vars map[string]interface{}) error {
	if t == nil {
		return errors.New("task is nil")
	}

	register := make(map[string]string)
	args := make(map[string]interface{})
	if t.PreProcess != nil {
		args["pre_pid"] = t.PreProcess.PID
		args["pre_state_code"] = t.PreProcess.StateCode
		args["pre_stdout"] = t.PreProcess.Stdout

		if t.PreProcess.Register != nil {
			register = t.PreProcess.Register
			for k, v := range t.PreProcess.Register {
				args[k] = v
			}
		}
	}
	if vars != nil {
		for k, v := range vars {
			args[k] = v
		}
	}
	if t.Args != nil {
		for k, v := range t.Args {
			args[k] = v
		}
	}

	if !t.WhenFunc(args) {
		PrintfMsg(ctx, "ignore task, when:%s \n", t.When)
		return nil
	}

	// 根据不同的moduleName构建不同的shell命令
	sh, err := t.ModuleObject.StringShell(args)
	if err != nil {
		PrintfMsg(ctx, "error:%s ,args:%+v \n", err.Error(), args)
		return err
	}
	PrintfMsg(ctx, "module:%s sh:%s \n", t.ModuleObject.Show(), sh)

	sudoPassword := ""

	if t.NotWait {
		// 不等待命令结果,返回:pid,错误
		pid, err := t.Worker.Start(sh,
			work.WithTimeOut(time.Duration(t.Timeout)*time.Second),
			work.WithSudoPassword(sudoPassword),
		)
		if err != nil {
			PrintError(ctx, err)
			return err
		}
		PrintfMsg(ctx, "pid:%d \n", pid)
		register[fmt.Sprintf("%s.pid", t.Name)] = strconv.Itoa(pid)
		t.PreProcess = &Process{
			Register: register,
		}
		return nil
	}

	// 等待命令结果,返回:状态码,输出,错误
	stateCode, stdout, err := t.Worker.RunOutput(sh,
		work.WithTimeOut(time.Duration(t.Timeout)*time.Second),
		work.WithSudoPassword(sudoPassword),
	)
	if err != nil {
		PrintfMsg(ctx, "stateCode:%d, stdout:%s", stateCode, stdout)
		PrintError(ctx, err)
		return err
	} else {
		PrintfMsg(ctx, "stateCode:%d, stdout:%s \n", stateCode, stdout)
	}

	if t.Register != "" {
		registerName := fmt.Sprintf("%s.stdout", t.Register)
		register[registerName] = stdout
	}
	t.PreProcess = &Process{
		StateCode: stateCode,
		Stdout:    stdout,
		Register:  register,
	}
	return nil
}

func (t Task) trimSpace() Task {
	t.Name = strings.TrimSpace(t.Name)
	if t.Args != nil {
		args := make(map[string]interface{})
		for key, value := range t.Args {
			switch v := value.(type) {
			case string:
				args[strings.TrimSpace(key)] = strings.TrimSpace(v)
			default:
				args[strings.TrimSpace(key)] = value
			}
		}
		t.Args = args
	}
	t.Register = strings.TrimSpace(t.Register)
	t.DelegateTo = strings.TrimSpace(t.DelegateTo)
	if len(t.Notify) != 0 {
		for i, _ := range t.Notify {
			t.Notify[i] = strings.TrimSpace(t.Notify[i])
		}
	}
	t.When = strings.TrimSpace(t.When)
	t.ShellBody = strings.TrimSpace(t.ShellBody)
	return t
}

// When 模块，接受一个条件函数并执行相应的操作
func (t *Task) WhenFunc(vars map[string]interface{}) bool {
	if strings.TrimSpace(t.When) == "" {
		return true
	}
	return evaluateCondition(vars, t.When)
}

// 递归解析和判断条件表达式
func evaluateCondition(vars map[string]interface{}, condition string) bool {
	// 去除两边的空格
	condition = strings.TrimSpace(condition)

	// 检查是否有括号
	if strings.HasPrefix(condition, "(") && strings.HasSuffix(condition, ")") {
		return evaluateCondition(vars, condition[1:len(condition)-1])
	}

	// 检查是否有 "or" 或 "and"
	if strings.Contains(condition, " or ") {
		parts := strings.SplitN(condition, " or ", 2)
		return evaluateCondition(vars, parts[0]) || evaluateCondition(vars, parts[1])
	}

	if strings.Contains(condition, " and ") {
		parts := strings.SplitN(condition, " and ", 2)
		return evaluateCondition(vars, parts[0]) && evaluateCondition(vars, parts[1])
	}

	// 检查是否有 "=="
	if strings.Contains(condition, "==") {
		eqParts := strings.Split(condition, "==")
		if len(eqParts) == 2 {
			left := strings.TrimSpace(eqParts[0])
			right := strings.TrimSpace(eqParts[1])

			v, ok := vars[left]
			if !ok {
				fmt.Printf("when var:%s not found. \n", left)
				return false
			}

			switch v.(type) {
			case int, int64, int32, uint, uint64, uint32:
				fmt.Printf("key:%s int left:%d right:%d, when:%v\n", left, v, parseInt(right), v == parseInt(right))
				return v == parseInt(right)
			default:
				fmt.Printf("key:%s string left:%s right:%s, when:%v\n", left, v, parseString(right), v == parseString(right))
				return v == parseString(right)
			}
		}
	}

	return false
}

// 辅助函数，用于解析整数
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// 辅助函数，用于解析字符串
func parseString(s string) string {
	return strings.Trim(s, ` "'`)
}
