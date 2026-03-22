package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/ayushk-1801/cage/internal/cgroup"
	"github.com/ayushk-1801/cage/internal/namespace"
)

type Container struct {
	ID       string
	Rootfs   string
	Hostname string
}

func New(id, rootfs string) *Container {
	return &Container{
		ID:       id,
		Rootfs:   rootfs,
		Hostname: id,
	}
}

func (c *Container) Run(args []string) error {
	cg := cgroup.New(c.ID)
	if err := cg.Create(cgroup.DefaultConfig()); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}
	cmd := namespace.NewParentProcess(args)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start process: %w", err)
	}
	if err := cg.AddProcess(cmd.Process.Pid); err != nil {
		return fmt.Errorf("add to cgroup: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait: %w", err)
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

	path, err := exec.LookPath(args[0])
	if err != nil {
		path = args[0]
	}

	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
