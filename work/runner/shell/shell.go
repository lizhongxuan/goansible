package shell

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"goansible/work/runner"
	"goansible/work/runner/static"
	"goansible/work/runner/types"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	LIB_PATH = "/var/sandbox/sandbox-shell"
)

type ShellRunner struct {
	runner.TempDirRunner
}

var (
	REQUIRED_FS = []string{
		path.Join(LIB_PATH),
	}
)
var sandbox_fs []byte

func (p *ShellRunner) Run(
	code string,
	timeout time.Duration,
	stdin []byte,
	preload string,
	options *types.RunnerOptions,
) (chan []byte, chan []byte, chan bool, error) {
	configuration := static.GetDifySandboxGlobalConfigurations()

	// initialize the environment
	untrusted_code_path, key, err := p.InitializeEnvironment(code, preload, options)
	if err != nil {
		return nil, nil, nil, err
	}

	// capture the output
	output_handler := runner.NewOutputCaptureRunner()
	output_handler.SetTimeout(timeout)

	err = p.WithTempDir(LIB_PATH, REQUIRED_FS, func(root_path string) error {
		// cleanup
		output_handler.SetAfterExitHook(func() {
			os.RemoveAll(root_path)
			os.Remove(root_path)
		})

		// create a new process
		cmd := exec.Command(
			"bash",
			untrusted_code_path,
			LIB_PATH,
			key,
		)
		cmd.Env = []string{}

		if configuration.Proxy.Socks5 != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", configuration.Proxy.Socks5))
			cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PROXY=%s", configuration.Proxy.Socks5))
		} else if configuration.Proxy.Https != "" || configuration.Proxy.Http != "" {
			if configuration.Proxy.Https != "" {
				cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", configuration.Proxy.Https))
			}
			if configuration.Proxy.Http != "" {
				cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PROXY=%s", configuration.Proxy.Http))
			}
		}

		err = output_handler.CaptureOutput(cmd)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return output_handler.GetStdout(), output_handler.GetStderr(), output_handler.GetDone(), nil
}

func (p *ShellRunner) InitializeEnvironment(code string, preload string, options *types.RunnerOptions) (string, string, error) {
	// create a tmp dir and copy the python script
	temp_code_name := strings.ReplaceAll(uuid.New().String(), "-", "_")
	temp_code_name = strings.ReplaceAll(temp_code_name, "/", ".")

	script := string(sandbox_fs)
	if options.EnableNetwork {
		script = strings.Replace(
			script,
			"{{enable_network}}", "1", 1,
		)
	} else {
		script = strings.Replace(
			script,
			"{{enable_network}}", "0", 1,
		)
	}

	script = strings.Replace(
		script,
		"{{preload}}",
		fmt.Sprintf("%s\n", preload),
		1,
	)

	// generate a random 512 bit key
	key_len := 64
	key := make([]byte, key_len)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", err
	}

	// encrypt the code
	encrypted_code := make([]byte, len(code))
	for i := 0; i < len(code); i++ {
		encrypted_code[i] = code[i] ^ key[i%key_len]
	}

	// encode code using base64
	code = base64.StdEncoding.EncodeToString(encrypted_code)
	// encode key using base64
	encoded_key := base64.StdEncoding.EncodeToString(key)

	code = strings.Replace(
		script,
		"{{code}}",
		code,
		1,
	)

	untrusted_code_path := fmt.Sprintf("%s/tmp/%s.py", LIB_PATH, temp_code_name)
	err = os.MkdirAll(path.Dir(untrusted_code_path), 0755)
	if err != nil {
		return "", "", err
	}
	err = os.WriteFile(untrusted_code_path, []byte(code), 0755)
	if err != nil {
		return "", "", err
	}

	return untrusted_code_path, encoded_key, nil
}