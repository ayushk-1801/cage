package seccomp

import (
	"fmt"
	"syscall"

	libseccomp "github.com/seccomp/libseccomp-golang"
	"golang.org/x/sys/unix"
)

var allowedSyscalls = []string{
	// Process lifecycle
	"execve", "execveat", "exit", "exit_group",
	"fork", "vfork", "clone", "clone3", "wait4", "waitid",

	// File I/O
	"read", "write", "readv", "writev", "pread64", "pwrite64",
	"open", "openat", "openat2", "close", "close_range", "creat",
	"stat", "fstat", "lstat", "statx", "newfstatat",
	"lseek", "dup", "dup2", "dup3", "pipe", "pipe2",
	"access", "faccessat", "faccessat2",
	"readlink", "readlinkat",
	"getcwd", "chdir", "fchdir",
	"mkdir", "mkdirat", "rmdir",
	"unlink", "unlinkat",
	"rename", "renameat", "renameat2",
	"chmod", "fchmod", "fchmodat",
	"chown", "fchown", "lchown", "fchownat",
	"truncate", "ftruncate",
	"getdents", "getdents64",
	"fsync", "fdatasync",
	"sendfile",

	// Memory
	"mmap", "mprotect", "munmap", "madvise", "brk",
	"mremap", "msync", "mincore", "mlock", "munlock",

	// Signals
	"rt_sigaction", "rt_sigprocmask", "rt_sigreturn",
	"rt_sigpending", "rt_sigsuspend", "rt_sigqueueinfo", "rt_tgsigqueueinfo",
	"kill", "tkill", "tgkill", "sigaltstack",
	"pause",

	// Identity / scheduling
	"getpid", "getppid", "gettid",
	"getuid", "getgid", "geteuid", "getegid",
	"getgroups", "setgroups",
	"setuid", "setgid", "setreuid", "setregid",
	"setresuid", "setresgid", "getresuid", "getresgid",
	"capget", "capset",
	"sched_getaffinity", "sched_setaffinity",
	"sched_yield", "sched_getparam", "sched_setparam",
	"sched_getscheduler", "sched_setscheduler",
	"nanosleep", "clock_nanosleep", "clock_gettime", "clock_getres",
	"gettimeofday", "time",
	"getitimer", "setitimer",
	"alarm",

	// Networking
	"socket", "socketpair", "bind", "listen", "accept", "accept4",
	"connect", "getsockname", "getpeername",
	"sendto", "sendmsg", "sendmmsg",
	"recvfrom", "recvmsg", "recvmmsg",
	"setsockopt", "getsockopt",
	"shutdown",
	"poll", "ppoll", "select", "pselect6",
	"epoll_create", "epoll_create1", "epoll_ctl", "epoll_wait", "epoll_pwait",
	"eventfd", "eventfd2",
	"timerfd_create", "timerfd_settime", "timerfd_gettime",
	"inotify_init", "inotify_init1", "inotify_add_watch", "inotify_rm_watch",

	// Go runtime + glibc internals
	"arch_prctl", "prctl", "seccomp",
	"futex", "futex_waitv",
	"rseq",
	"set_tid_address", "set_robust_list", "get_robust_list",
	"uname", "sysinfo",
	"getrandom",
	"memfd_create",
	"copy_file_range",
	"ioctl",
	"fcntl",
	"umask",
	"getrlimit", "setrlimit", "prlimit64",
	"getrusage",
}

func Apply() error {
	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		return fmt.Errorf("seccomp: set no_new_privs: %w", err)
	}

	filter, err := libseccomp.NewFilter(libseccomp.ActErrno.SetReturnCode(int16(syscall.ENOSYS)))
	if err != nil {
		return fmt.Errorf("seccomp: create filter: %w", err)
	}

	for _, name := range allowedSyscalls {
		sc, err := libseccomp.GetSyscallFromName(name)
		if err != nil {
			continue 
		}
		if err := filter.AddRule(sc, libseccomp.ActAllow); err != nil {
			filter.Release()
			return fmt.Errorf("seccomp: add rule for %s: %w", name, err)
		}
	}

	if err := filter.Load(); err != nil {
		filter.Release()
		return fmt.Errorf("seccomp: load filter: %w", err)
	}

	filter.Release()
	return nil
}
