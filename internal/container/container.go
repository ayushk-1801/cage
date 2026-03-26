package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/ayushk-1801/cage/internal/cgroup"
	"github.com/ayushk-1801/cage/internal/namespace"
	"github.com/ayushk-1801/cage/internal/seccomp"
)

type Container struct {
	ID          string
	Name        string
	Hostname    string
	Status      string
	PID         int
	ExitCode    int
	IP          string
	ImageID     string
	Rootfs      string
	Cmd         []string
	Env         []string
	Mounts      []string
	Ports       []PortMap
	HealthCheck HealthCheck
	CreatedAt   time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
	Labels      map[string]string
}

type HealthCheck struct {
	Test     []string
	Interval time.Duration
	Timeout  time.Duration
	Retries  int
}

type PortMap struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

func New(id, rootfs string) *Container {
	return &Container{
		ID:       id,
		Rootfs:   rootfs,
		Hostname: id,
	}
}

func (c *Container) Run(args []string) error {
	c.Status = "created"
	c.CreatedAt = time.Now()
	c.Cmd = args

	cg := cgroup.New(c.ID)
	if err := cg.Create(cgroup.DefaultConfig()); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	cmd := namespace.NewParentProcess(args)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start process: %w", err)
	}

	c.PID = cmd.Process.Pid
	c.Status = "running"
	c.StartedAt = time.Now()

	if err := cg.AddProcess(c.PID); err != nil {
		return fmt.Errorf("add to cgroup: %w", err)
	}

	err := cmd.Wait()

	c.FinishedAt = time.Now()
	c.Status = "stopped"

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				c.ExitCode = ws.ExitStatus()
			}
		} else {
			c.ExitCode = 1
		}
	} else {
		c.ExitCode = 0
	}

	if err := cg.Destroy(); err != nil {
		return fmt.Errorf("destroy cgroup: %w", err)
	}

	return nil
}

func (c *Container) Child(args []string) error {
	if err := namespace.SetupNamespace(c.Hostname, c.Rootfs); err != nil {
		return fmt.Errorf("setup namespace: %w", err)
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		path = args[0]
	}

	if err := seccomp.Apply(); err != nil {
		return fmt.Errorf("apply seccomp: %w", err)
	}

	env := []string{
		"PATH=/bin:/usr/bin",
		"TERM=xterm",
	}

	if err := syscall.Exec(path, args, env); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
