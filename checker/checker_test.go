package checker

import (
	"goansible/work"
	"testing"
)

// 模拟 Worker 接口
type mockWorker struct{}

func (m *mockWorker) RunOutput(cmd string) (int, []byte, error) {
	// 根据需要模拟不同的输出和错误
	return 0, []byte("mock output"), nil
}

func TestCheck(t *testing.T) {
	// 初始化测试数据
	CheckerPoor = map[string]*Checker{
		"test1": {
			Name:  "Test1",
			Shell: "echo hello",
			Regex: "hello",
		},
		"test2": {
			Name:  "Test2",
			Shell: "echo world",
			Regex: "world",
		},
	}

	mockW := &work.LocalCmd{}

	tests := []struct {
		name    string
		list    []string
		wantErr bool
	}{
		{
			name:    "Valid checks",
			list:    []string{"test1", "test2"},
			wantErr: false,
		},
		{
			name:    "Non-existent check",
			list:    []string{"test1", "non-existent"},
			wantErr: false,
		},
		// 可以添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Check(mockW, tt.list)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChecker_check(t *testing.T) {
	mockW := &work.LocalCmd{}
	tests := []struct {
		name    string
		checker Checker
		wantErr bool
	}{
		{
			name: "Successful check",
			checker: Checker{
				Shell: "echo hello",
				Regex: "hello",
				Work:  mockW,
			},
			wantErr: false,
		},
		{
			name: "Failed regex match",
			checker: Checker{
				Shell: "echo hello",
				Regex: "world",
				Work:  mockW,
			},
			wantErr: true,
		},
		// 可以添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.checker.check()
			if (err != nil) != tt.wantErr {
				t.Errorf("Checker.check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
