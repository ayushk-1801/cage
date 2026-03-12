package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/ayushk-1801/cage/internal/namespace"
)

type Container struct {
	ID string
	Rootfs string
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
	cmd := namespace.NewParentProcess(args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}
	return nil
}

func (c *Container) Child(args []string) error {
	if err := namespace.SetupNamespace(c.Hostname, c.Rootfs); err != nil {
		return fmt.Errorf("failed to setup namespace: %w", err)
	}
	defer syscall.Unmount("/proc", 0)
	
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command in container: %w", err)
	}
	
	return nil
}