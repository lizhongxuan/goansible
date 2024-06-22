package ansible

import (
	"context"
	"errors"
	"fmt"
	"go-ansible/module"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

// Playbook represents an Ansible playbook structure.
type Playbook struct {
	Name         string                 `yaml:"name"`          // Playbook 名称
	Hosts        string                 `yaml:"hosts"`         // 任务运行的主机或主机组
	Middlewares  string                 `yaml:"middlewares"`   // 中间件组
	Become       bool                   `yaml:"become"`        // 是否提升权限（类似于 sudo）
	Vars         map[string]interface{} `yaml:"vars"`          // 变量定义
	Tasks        []Task                 `yaml:"tasks"`         // 任务列表
	Handlers     []Task                 `yaml:"handlers"`      // 处理程序列表
	Roles        []string               `yaml:"roles"`         // 角色列表
	Strategy     string                 `yaml:"strategy"`      // 执行策略：linear 是默认策略，按顺序执行每个任务；free 策略允许主机并行执行任务，不必等待其他主机完成当前任务。
	IgnoreErrors bool                   `yaml:"ignore_errors"` // 是否忽略任务错误
}

type ListPlaybook struct {
	List           []*Playbook
	HostMaps       map[string][]*Host
	MiddlewareMaps map[string][]*Middleware
}

func GeneratePlaybook(options ...AnsiblePlaybookOptionsFunc) (*ListPlaybook, error) {
	listPlaybook := &ListPlaybook{
		List:           make([]*Playbook, 0),
		HostMaps:       make(map[string][]*Host),
		MiddlewareMaps: make(map[string][]*Middleware),
	}
	for _, option := range options {
		if err := option(listPlaybook); err != nil {
			return nil, err
		}
	}
	if listPlaybook.List == nil || len(listPlaybook.List) == 0 {
		return nil, errors.New("no playbook")
	}
	return listPlaybook, nil
}

type AnsiblePlaybookOptionsFunc func(*ListPlaybook) error

func WithPlaybooks(data []byte) AnsiblePlaybookOptionsFunc {
	return func(lpb *ListPlaybook) error {
		var pbList []*Playbook
		if err := yaml.Unmarshal(data, &pbList); err != nil {
			return err
		}
		lpb.List = append(lpb.List, pbList...)
		return nil
	}
}
func WithHosts(hosts map[string][]*Host) AnsiblePlaybookOptionsFunc {
	return func(lpb *ListPlaybook) error {
		lpb.HostMaps = hosts
		return nil
	}
}
func WithMiddlewares(middlewares map[string][]*Middleware) AnsiblePlaybookOptionsFunc {
	return func(lpb *ListPlaybook) error {
		lpb.MiddlewareMaps = middlewares
		return nil
	}
}

func (pbList *ListPlaybook) Run() error {
	if pbList == nil {
		return errors.New("ListPlaybook is nil error")
	}
	if pbList.HostMaps == nil {
		pbList.HostMaps = make(map[string][]*Host)
		pbList.HostMaps[""] = make([]*Host, 0)
	}
	if pbList.MiddlewareMaps == nil {
		pbList.MiddlewareMaps = make(map[string][]*Middleware)
		pbList.MiddlewareMaps[""] = make([]*Middleware, 0)
	}

	for i, _ := range pbList.List {
		pb := pbList.List[i]
		//pb.trimSpace()
		if err := pb.verify(); err != nil {
			log.Printf("end playbook, because verify playbook-%d:%s error. \n", i, pb.Name)
			return err
		}

		if pb.Vars == nil {
			pb.Vars = make(map[string]interface{})
		}
		if hosts, ok := pbList.HostMaps[pb.Hosts]; ok {
			pb.Vars["hosts"] = hosts
		}
		if middlewares, ok := pbList.MiddlewareMaps[pb.Middlewares]; ok {
			pb.Vars["middlewares"] = middlewares
		}

		printPlaybook(pb)

		if err := pb.run(); err != nil {
			if !pb.IgnoreErrors {
				log.Printf("end playbook, because run playbook-%d:%s error. \n", i, pb.Name)
				return err
			}
		}
	}
	return nil
}

// TODO trimSpace 去除所有字符串的前后空格,好像有点多余,先注释掉
func (pb *Playbook) trimSpace() {
	if pb == nil {
		return
	}
	pb.Name = strings.TrimSpace(pb.Name)
	pb.Hosts = strings.TrimSpace(pb.Hosts)
	pb.Middlewares = strings.TrimSpace(pb.Middlewares)
	pb.Strategy = strings.TrimSpace(pb.Strategy)
	if len(pb.Tasks) != 0 {
		for i, _ := range pb.Tasks {
			pb.Tasks[i] = pb.Tasks[i].trimSpace()
		}
	}
	if len(pb.Handlers) != 0 {
		for i, _ := range pb.Handlers {
			pb.Handlers[i] = pb.Handlers[i].trimSpace()
		}
	}
	if len(pb.Roles) != 0 {
		for i, _ := range pb.Roles {
			pb.Roles[i] = strings.TrimSpace(pb.Roles[i])
		}
	}

	if pb.Vars != nil {
		vars := make(map[string]interface{})
		for key, value := range pb.Vars {
			switch v := value.(type) {
			case string:
				vars[strings.TrimSpace(key)] = strings.TrimSpace(v)
			default:
				vars[strings.TrimSpace(key)] = value
			}
		}
		pb.Vars = vars
	}
}

func (pb *Playbook) verify() error {
	if pb.Strategy == "" {
		pb.Strategy = "linear"
	}
	if pb.Name == "" {
		return errors.New("There's a playbook that doesn't have a name.")
	}
	if len(pb.Tasks) == 0 {
		return errors.New("playbook tasks is nil")
	}
	for i, t := range pb.Tasks {
		if t.Name == "" {
			return errors.New("There's a task that doesn't have a name.")
		}
		mn, err := module.ModuleVerify(t.Module)
		if err != nil {
			return err
		}
		fmt.Println("mn:", mn)
		pb.Tasks[i].ModuleName = mn
	}
	return nil
}

func (pb *Playbook) run() error {
	ctx := context.Background()
	var preTask *Task
	var errGroup errgroup.Group
	for j, _ := range pb.Tasks {
		task := pb.Tasks[j]
		if task.ShowShell {
			log.Printf("run playbook:%s task-%d:%s %s shell:%s \n", pb.Name, j, task.Name, pb.Strategy, task.Module.Shell)
		} else {
			log.Printf("run playbook:%s task-%d:%s %s \n", pb.Name, j, task.Name, pb.Strategy)
		}

		// 执行任务
		if pb.Strategy == "free" {
			errGroup.Go(func() error {
				return pb.runTask(ctx, &task, &Task{})
			})
		} else {
			if err := pb.runTask(ctx, &task, preTask); err != nil {
				return err
			}
			preTask = &task
		}

	}
	return errGroup.Wait()
}

func (pb *Playbook) runTask(ctx context.Context, task *Task, preTask *Task) error {
	items := []interface{}{
		"flag",
	}
	items = append(items, task.WithItems...)
	items = append(items, task.Loop...)

	for i, t := range items {
		item := t
		if len(items) != 1 {
			if i == 0 {
				continue
			}
			pb.Vars["item"] = item
		}

		if err := task.run(ctx, pb.Vars); err != nil {
			if !task.IgnoreErrors {
				log.Printf("playbook:%s task:%s error:%+v \n", pb.Name, task.Name, err)
				return err
			}
		}
	}

	// Notify
	if err := pb.runNotify(ctx, task.Notify, task.PreProcess); err != nil {
		if !task.IgnoreErrors {
			return err
		}
	}
	return nil
}

func (pb *Playbook) runNotify(ctx context.Context, notifys []string, preProcess *Process) error {
	if len(notifys) == 0 || len(pb.Handlers) == 0 {
		return nil
	}
	for _, noti := range notifys {
		for _, hand := range pb.Handlers {
			if hand.Name != noti {
				continue
			}
			hand.PreProcess = preProcess
			if err := hand.run(ctx, pb.Vars); err != nil {
				if !hand.IgnoreErrors {
					return err
				}
			}
			break
		}
	}
	return nil
}

func printPlaybook(pb *Playbook) {
	if pb == nil {
		return
	}
	if pb.Tasks != nil {
		for _, t := range pb.Tasks {
			fmt.Printf("Playbook:%s Task: %+v\n", pb.Name, t)
		}
	}
	if pb.Handlers != nil {
		for _, t := range pb.Handlers {
			fmt.Printf("Playbook:%s Handler: %+v\n", pb.Name, t)
		}
	}
}
