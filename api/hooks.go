package api

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jdel/acdc/cfg"
	"github.com/jdel/acdc/util"
)

//Hook represents a hook
type Hook struct {
	Name     string
	APIKey   string
	FileName string
}

// GetHook returns a hook based on the imput name
func (key Key) GetHook(name string) Hook {
	hook := Hook{
		Name:     name,
		APIKey:   key.Unique,
		FileName: func() string { return filepath.Join(cfg.GetComposeDir(), key.Unique, name) }(),
	}

	if !util.FileExists(fmt.Sprintf("%s.yml", hook.FileName)) {
		return Hook{}
	}
	return hook
}

func (hook Hook) callMethod(m string) (string, error) {
	logAPI.Debugf("method: %+v", strings.Title(m))
	method := reflect.ValueOf(&hook).MethodByName(strings.Title(m))
	logAPI.Debugf("reflected method: %+v", method)
	logAPI.Debugf("reflected method type %s", method.Kind())
	if method.Kind() == reflect.Invalid {
		return "", fmt.Errorf("Ignored action %s", m)
	}
	out := method.Call([]reflect.Value{})
	cmd := out[0].Interface().(*exec.Cmd)
	o, err := cmd.CombinedOutput()
	logAPI.Debugf("output: %v", string(o))
	return string(o), err
}

// ExecuteSequentially will execute the commands sequentually and return output
func (hook Hook) ExecuteSequentially(actions ...string) string {
	logAPI.Debugf("Sequential Actions: %+v", actions)
	var output bytes.Buffer
	for _, a := range actions {
		logAPI.Debugf("Processing: %+v", a)
		hook.callMethod(a)
	}
	return output.String()
}

// Pull pulls images for the hook
func (hook Hook) Pull() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"pull",
		)...,
	)
}

// Up brings hook up
func (hook Hook) Up() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"up",
			"-d",
		)...,
	)
}

// Down brings hook down
func (hook Hook) Down() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"down",
		)...,
	)
}

// Ps executes docker-compose ps
func (hook Hook) Ps() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"ps",
		)...,
	)
}

// Logs return hook logs
func (hook Hook) Logs() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"logs",
		)...,
	)
}

// Restart restarts hook
func (hook Hook) Restart() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"restart",
		)...,
	)
}

// Start starts hook
func (hook Hook) Start() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"start",
		)...,
	)
}

// Stop stops hook
func (hook Hook) Stop() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"stop",
		)...,
	)
}

// Build builds hook
func (hook Hook) Build() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"build",
			"--no-cache",
		)...,
	)
}

// Delete deletes the hook
func (hook Hook) Delete() error {
	filePath := fmt.Sprintf("%s.yml", hook.FileName)
	return os.Remove(filePath)
}

func (hook Hook) composeCommonArgs() []string {
	return []string{
		"-f",
		fmt.Sprintf("%s.yml", hook.FileName),
		"-p",
		hook.APIKey,
	}
}
