package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const cgroupRoot = "/sys/fs/cgroup/cage"

type Config struct {
	MemoryLimit int64
	CPUQuota    int64
	CPUPeriod   int64
	MaxPIDs     int64
}

func DefaultConfig() *Config {
	return &Config{
		MemoryLimit: 512 * 1024 * 1024,
		CPUQuota:    50000,
		CPUPeriod:   100000,
		MaxPIDs:     100,
	}
}

type Cgroup struct {
	id   string
	path string
}

func New(containerID string) *Cgroup {
	return &Cgroup{
		id:   containerID,
		path: filepath.Join(cgroupRoot, containerID),
	}
}

func EnableControllers() error {
	controllers := []string{
		"/sys/fs/cgroup/cgroup.subtree_control",
		"/sys/fs/cgroup/cage/cgroup.subtree_control",
	}
	for _, path := range controllers {
		err := os.WriteFile(path, []byte("+memory +cpu +pids"), 0700)
		if err != nil {
			return fmt.Errorf("enable controllers at %s: %w", path, err)
		}
	}
	return nil
}

func (c *Cgroup) Create(cfg *Config) error {
	if err := os.MkdirAll(cgroupRoot, 0755); err != nil {
		return fmt.Errorf("create cgroup root: %w", err)
	}

	if err := EnableControllers(); err != nil {
		return fmt.Errorf("enable controllers: %w", err)
	}

	if err := os.MkdirAll(c.path, 0755); err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	if err := c.setMemory(cfg.MemoryLimit); err != nil {
		return err
	}
	if err := c.setCPU(cfg.CPUQuota, cfg.CPUPeriod); err != nil {
		return err
	}
	if err := c.setPIDs(cfg.MaxPIDs); err != nil {
		return err
	}

	return nil
}

func (c *Cgroup) AddProcess(pid int) error {
	return c.write("cgroup.procs", strconv.Itoa(pid))
}

func (c *Cgroup) Destroy() error {
	if err := os.RemoveAll(c.path); err != nil {
		return fmt.Errorf("destroy cgroup: %w", err)
	}
	return nil
}

func (c *Cgroup) setMemory(limitBytes int64) error {
	if limitBytes == 0 {
		return nil
	}
	return c.write("memory.max", strconv.FormatInt(limitBytes, 10))
}

func (c *Cgroup) setCPU(quota, period int64) error {
	if quota == 0 {
		return nil
	}
	val := fmt.Sprintf("%d %d", quota, period)
	return c.write("cpu.max", val)
}

func (c *Cgroup) setPIDs(max int64) error {
	if max == 0 {
		return nil
	}
	return c.write("pids.max", strconv.FormatInt(max, 10))
}

func (c *Cgroup) write(file, value string) error {
	path := filepath.Join(c.path, file)
	if err := os.WriteFile(path, []byte(value), 0700); err != nil {
		return fmt.Errorf("cgroup write %s: %w", file, err)
	}
	return nil
}
