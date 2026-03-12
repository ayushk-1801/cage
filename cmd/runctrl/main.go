package main

import (
	"fmt"
	"os"

	"github.com/ayushk-1801/cage/internal/container"
)

const defaultRootfs = "/tmp/busybox-rootfs"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cage run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cage run <command>")
			os.Exit(1)
		}
		run(os.Args[2:])

	case "child":
		// Called by re-exec, runs inside namespaces
		child(os.Args[2:])

	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func run(args []string) {
	c := container.New("cage-001", defaultRootfs)
	if err := c.Run(args); err != nil {
		fmt.Println("[run] error:", err)
		os.Exit(1)
	}
}

func child(args []string) {
	c := container.New("cage-001", defaultRootfs)
	if err := c.Child(args); err != nil {
		fmt.Println("[child] error:", err)
		os.Exit(1)
	}
}