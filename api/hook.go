package api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/jdel/acdc/cfg"
	"github.com/jdel/acdc/lgr"
	"github.com/jdel/acdc/util"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

//Hook represents a hook
type Hook struct {
	Name     string
	APIKey   string
	FileName string
}

var c = make(chan []byte, 9999)

var myLogger = lgr.ACDCLogger{
	Output: c,
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

// Pull pulls images for the hook
func (hook Hook) Pull() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"pull",
		)...,
	)
}

// NewUp is the new PS using libcompose
func (hook Hook) NewUp() (string, error) {
	project, err := hook.composeProject()
	if err != nil {
		return "", err
	}
	// fmt.Printf("%+v", project)
	err = project.Up(context.Background(), options.Up{})
	if err != nil {
		return "ok", err
	}
	return "", err
}

// Up brings hook up
func (hook Hook) Up() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"up",
			"-d",
			"--remove-orphans",
		)...,
	)
}

// Down brings hook down
func (hook Hook) Down() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"down",
			"--remove-orphans",
		)...,
	)
}

// NewPs is the new PS using libcompose
func (hook Hook) NewPs() (string, error) {
	project, err := hook.composeProject()
	if err != nil {
		return "", err
	}

	info, err := project.Ps(context.Background())

	if err != nil {
		return "", err
	}
	return info.String([]string{"Name", "Command", "State", "Ports"}, true), err

	// project := project.NewProject(&project.Context{
	// 	ComposeFiles: []string{fmt.Sprintf("%s.yml", hook.FileName)},
	// 	ProjectName:  fmt.Sprintf("%s%s", hook.APIKey, hook.Name),
	// }, nil, nil)

	// err := project.Parse()
	// if err != nil {
	// 	return "", err
	// }

	// for _, key := range project.ServiceConfigs.Keys() {
	// 	if svc, ok := project.ServiceConfigs.Get(key); ok {
	// 		fmt.Println("=== " + key)

	// 		for _, env := range svc.Environment {
	// 			fmt.Println(env)
	// 		}
	// 	}
	// }

	// info, err := project.Ps(context.Background())
	// fmt.Printf("%+v", info)
	// return "", err
}

// Ps executes docker-compose ps
func (hook Hook) Ps() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"ps",
		)...,
	)
}

// NewLogs is the new PS using libcompose
func (hook Hook) NewLogs() (string, error) {
	project, err := hook.composeProject()
	if err != nil {
		return "", err
	}

	err = project.Log(context.Background(), false)
	if err != nil {
		return "", err
	}
	return string(<-c), err
}

// Logs return hook logs
func (hook Hook) Logs() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"logs",
		)...,
	)
}

// Restart restarts hook
func (hook Hook) Restart() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"restart",
		)...,
	)
}

// Start starts hook
func (hook Hook) Start() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"start",
		)...,
	)
}

// Stop stops hook
func (hook Hook) Stop() *exec.Cmd {
	return exec.Command(cfg.DockerComposePath,
		append(hook.composeCommonArgs(),
			"stop",
		)...,
	)
}

// Delete deletes the hook
func (hook Hook) Delete() error {
	filePath := fmt.Sprintf("%s.yml", hook.FileName)
	return os.Remove(filePath)
}

func (hook Hook) composeProject() (project.APIProject, error) {
	return docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles:  []string{fmt.Sprintf("%s.yml", hook.FileName)},
			ProjectName:   fmt.Sprintf("%s%s", hook.APIKey, hook.Name),
			LoggerFactory: &myLogger,
		},
	}, nil)
}

func (hook Hook) composeCommonArgs() []string {
	return []string{
		"-f",
		fmt.Sprintf("%s.yml", hook.FileName),
		"-p",
		fmt.Sprintf("%s%s", hook.APIKey, hook.Name),
	}
}
