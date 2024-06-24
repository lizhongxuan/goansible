package ansible

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func setCtxPlaybook(ctx context.Context, pb string, i ...int) context.Context {
	if len(i) != 0 {
		ctx = context.WithValue(ctx, "playbookIndex", strconv.Itoa(i[0]+1))
	}
	return context.WithValue(ctx, "playbookName", pb)
}

func setCtxTask(ctx context.Context, task string, i ...int) context.Context {
	if len(i) != 0 {
		ctx = context.WithValue(ctx, "taskIndex", strconv.Itoa(i[0]+1))
	}
	return context.WithValue(ctx, "taskName", task)
}

func setCtxTaskShell(ctx context.Context, shell string) context.Context {
	return context.WithValue(ctx, "taskShell", shell)
}

func contextKey(ctx context.Context) string {
	arr := make([]string, 0)
	playbookName, ok := ctx.Value("playbookName").(string)
	if ok {
		playbookIndex := ctx.Value("playbookIndex").(string)
		if playbookIndex != "" {
			arr = append(arr, fmt.Sprintf("%s-playbook:%s", playbookIndex, playbookName))
		} else {
			arr = append(arr, fmt.Sprintf("playbook:%s", playbookName))
		}
	}
	taskName, ok := ctx.Value("taskName").(string)
	if ok {
		taskIndex := ctx.Value("taskIndex").(string)
		if taskIndex != "" {
			arr = append(arr, fmt.Sprintf("%s-task:%s", taskIndex, taskName))
		} else {
			arr = append(arr, fmt.Sprintf("task:%s", taskName))
		}
	}
	taskShell, ok := ctx.Value("taskShell").(string)
	if ok {
		arr = append(arr, fmt.Sprintf("shell:%s", taskShell))
	}
	return strings.Join(arr, " ")
}

func PrintError(ctx context.Context, err error) {
	log.Printf("[%s] error:%+v \n", contextKey(ctx), err)
}

func PrintfMsg(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.Printf("[%s] msg:%s \n", contextKey(ctx), msg)
}
