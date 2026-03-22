# Cage — TODO 

## Phase 1 — Core Runtime

### Namespace Isolation
- [x] Create UTS, PID, Mount, Network namespaces
- [x] Re-exec trick using `/proc/self/exe`
- [x] Set custom hostname inside container
- [x] Chroot into busybox rootfs
- [x] Mount `/proc` inside container
- [x] Extract namespace logic into its own package
- [x] Extract container logic into its own package

### Cgroups
- [x] Create cgroup directory per container
- [x] Set memory limit
- [x] Set CPU limit
- [x] Set max process limit
- [x] Add container PID to cgroup
- [x] Clean up cgroup on container exit

### Container ID + State
- [ ] Generate unique container ID
- [ ] Create Container struct with ID, status, PID, IP, rootfs, created time
- [ ] Track containers in an in-memory map
- [ ] Implement get, save, delete, list on the store
- [ ] Add status transitions: created → running → stopped → deleted

### OverlayFS
- [ ] Create lower, upper, work, merged dirs per container
- [ ] Mount overlayfs combining all layers
- [ ] Point container chroot at merged dir
- [ ] Unmount overlayfs on container stop
- [ ] Delete upper and work dirs on container delete

### Image Management (local only)
- [ ] Define image directory structure
- [ ] Write image metadata to metadata.json
- [ ] Import a rootfs tarball as a named image
- [ ] Implement image list
- [ ] Implement image delete

### Networking
- [ ] Create cage0 bridge interface on host
- [ ] Assign 172.20.0.1/16 to bridge
- [ ] Create veth pair per container
- [ ] Move one veth end into container network namespace
- [ ] Assign IP from pool to container eth0
- [ ] Set default route inside container
- [ ] Write /etc/resolv.conf inside container
- [ ] Add iptables NAT rule for outbound internet
- [ ] Track and release IPs from pool

---

## Phase 2 — CLI + Daemon

### Wire Everything Together
- [ ] Call image, overlay, cgroup, network, namespace in correct order on run
- [ ] Implement full teardown on stop (kill, network, cgroup, overlay)
- [ ] Implement delete (remove state, dirs, release IP)

### CLI Commands
- [ ] `cage run <image> <cmd>`
- [ ] `cage ps` — list running containers
- [ ] `cage ps -a` — list all containers
- [ ] `cage stop <id>`
- [ ] `cage rm <id>`
- [ ] `cage logs <id>`
- [ ] `cage exec <id> <cmd>`
- [ ] `cage images`
- [ ] `cage import <tarball> <name>`
- [ ] `cage rmi <id>`
- [ ] `cage inspect <id>`

### Persist State to Disk
- [ ] Write container state to JSON file on create
- [ ] Load all state files on startup
- [ ] Update state file on status change
- [ ] Delete state file on container remove

### gRPC Daemon
- [ ] Define gRPC proto (create, start, stop, delete, list, logs, exec)
- [ ] Implement gRPC server as a daemon
- [ ] Create unix socket at `/var/run/cage/cage.sock`
- [ ] Update CLI to talk to daemon over socket instead of direct calls

---

## Phase 3 — Advanced Features

### Logs
- [ ] Redirect container stdout/stderr to a log file
- [ ] Implement `cage logs <id>` to print log file
- [ ] Implement `cage logs -f <id>` to stream/follow logs

### Exec into Running Container
- [ ] Look up container PID from state
- [ ] Use nsenter to join container namespaces
- [ ] Run command inside the container namespaces

### Image Pull from Registry
- [ ] Pull image manifest from OCI registry
- [ ] Download each layer
- [ ] Unpack and merge layers into image rootfs
- [ ] `cage pull ubuntu:22.04`

### Port Forwarding
- [ ] Accept `-p hostPort:containerPort` flag on run
- [ ] Add iptables DNAT rule on container start
- [ ] Remove iptables rule on container stop

---

## Phase 4 — Orchestrator (separate repo)

- [ ] Define node agent that calls cage gRPC API
- [ ] Build control plane API server
- [ ] Implement Pod spec (image, cmd, resources)
- [ ] Build scheduler — place pods on nodes
- [ ] Build deployment controller — replicas, rolling updates
- [ ] Add service discovery and internal DNS
- [ ] Add overlay network across nodes (VXLAN or WireGuard)
