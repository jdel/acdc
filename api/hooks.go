package api // import jdel.org/acdc/api

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"jdel.org/acdc/cfg"
	"jdel.org/acdc/util"
)

// Hook represents a hook
type Hook struct {
	Name     string
	APIKey   string
	FileName string
}

// Hooks is a slice of hook
type Hooks []*Hook

// GetHook returns a hook based on the imput name
func (key Key) GetHook(name string) *Hook {
	hook := Hook{
		Name:     name,
		APIKey:   key.Unique,
		FileName: func() string { return filepath.Join(cfg.GetComposeDir(), key.Unique, name) }(),
	}

	if !util.FileExists(fmt.Sprintf("%s.yml", hook.FileName)) {
		return nil
	}
	return &hook
}

func (hook *Hook) callMethod(m string) (string, error) {
	allowedActions := []string{
		"pull",
		"up",
		"down",
		"ps",
		"logs",
		"start",
		"stop",
		"restart",
		"build",
		"config",
		"kill",
	}

	if !util.IsStringInSlice(m, allowedActions) {
		return "", fmt.Errorf("Invalid action %s", m)
	}

	method := reflect.ValueOf(hook).MethodByName(strings.Title(m))
	if method.Kind() == reflect.Invalid {
		return "", fmt.Errorf("Unknown method for action %s", m)
	}
	out := method.Call([]reflect.Value{})
	cmd := out[0].Interface().(*exec.Cmd)
	ticket := util.NextTicket()
	util.InputQueue <- util.QueueItem{Key: hook.APIKey, Ticket: ticket, Cmd: cmd}
	// o, err := cmd.CombinedOutput()
	return strconv.Itoa(ticket), nil
}

func (hook *Hook) callMethodNow(m string) (string, error) {
	allowedActions := []string{
		"pull",
		"up",
		"down",
		"ps",
		"logs",
		"start",
		"stop",
		"restart",
		"build",
		"config",
		"kill",
	}

	if !util.IsStringInSlice(m, allowedActions) {
		return "", fmt.Errorf("Invalid action %s", m)
	}

	method := reflect.ValueOf(hook).MethodByName(strings.Title(m))
	if method.Kind() == reflect.Invalid {
		return "", fmt.Errorf("Unknown method for action %s", m)
	}
	out := method.Call([]reflect.Value{})
	cmd := out[0].Interface().(*exec.Cmd)
	o, err := cmd.CombinedOutput()
	return string(o), err
}

// ExecuteSequentially will queue the commands sequentually and return ticket numbers
func (hook *Hook) ExecuteSequentially(actions ...string) (string, error) {
	var err error
	var o string
	var output bytes.Buffer
	for _, a := range actions {
		logAPI.Debugf("Queing %s on %s for key %s", a, hook.Name, hook.APIKey)
		o, err = hook.callMethod(a)
		output.WriteString(fmt.Sprintf("%s\n", o))
	}

	return output.String(), err
}

// ExecuteSequentiallyNow will execute the commands immediately and return output
func (hook *Hook) ExecuteSequentiallyNow(actions ...string) (string, error) {
	var err error
	var o string
	var output bytes.Buffer
	for _, a := range actions {
		logAPI.Debugf("Queing %s on %s for key %s", a, hook.Name, hook.APIKey)
		o, err = hook.callMethodNow(a)
		output.WriteString(o)
	}

	return output.String(), err
}

// Pull pulls images for the hook
func (hook *Hook) Pull() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"pull",
		)...,
	)
}

// Up brings hook up
func (hook *Hook) Up() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"up",
			"-d",
			"--no-build",
		)...,
	)
}

// Down brings hook down
func (hook *Hook) Down() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"down",
		)...,
	)
}

// Ps executes docker-compose ps
func (hook *Hook) Ps() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"ps",
		)...,
	)
}

// Logs return hook logs
func (hook *Hook) Logs() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"logs",
		)...,
	)
}

// Start starts hook
func (hook *Hook) Start() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"start",
		)...,
	)
}

// Stop stops hook
func (hook *Hook) Stop() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"stop",
		)...,
	)
}

// Restart restarts hook
func (hook *Hook) Restart() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"restart",
		)...,
	)
}

// Build builds hook
func (hook *Hook) Build() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"build",
			"--no-cache",
			"--memory",
			cfg.GetBuildMemLimit(),
		)...,
	)
}

// Config shows hook compose file
func (hook *Hook) Config() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"config",
		)...,
	)
}

// Kill kills hook's containers
func (hook *Hook) Kill() *exec.Cmd {
	return exec.Command(cfg.GetDockerComposeLocation(),
		append(hook.composeCommonArgs(),
			"kill",
		)...,
	)
}

// Delete deletes the hook
func (hook *Hook) Delete() error {
	filePath := fmt.Sprintf("%s.yml", hook.FileName)
	return os.Remove(filePath)
}

func (hook *Hook) composeCommonArgs() []string {
	return []string{
		"-f",
		fmt.Sprintf("%s.yml", hook.FileName),
		"-p",
		hook.APIKey,
	}
}
