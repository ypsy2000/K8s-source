# Lab 1

## 1.1 Container runtime

The goal of this first lab is to install the container runtime on your server
and check that it's working fine.

## 1.1.1 Docker

Docker is pre-installed on all the node

- Check that Docker is installed correctly and running:

```shell
docker run helloworld
```

- You should see this:

```shell
Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

## 1.1.2 Containerd & co (Optional)

As an alternative to Docker, you can use Kubernetes with containerd.
Containerd is a minimal container runtime, as such it relies on other tools
for many functionalities:

- runc: CLI tools for spawning and running containers
- CNI: Container Network Interface, to configure containers network interfaces
- crictl: CLI tool to interact with CRI (Container Runtime Interface)
  compatible runtimes

To use containerd, just see the following steps:

- Download containerd and dependencies:

```shell
curl -LO https://github.com/containerd/containerd/releases/download/v1.2.2/containerd-1.2.2.linux-amd64.tar.gz
curl -LO https://github.com/opencontainers/runc/releases/download/v1.0.0-rc6/runc.amd64
curl -LO https://github.com/containernetworking/plugins/releases/download/v0.7.4/cni-plugins-amd64-v0.7.4.tgz
curl -LO https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.13.0/crictl-v1.13.0-linux-amd64.tar.gz
```

- Install containerd binaries to path:

```shell
tar xf containerd-1.2.2.linux-amd64.tar.gz -C /usr/local/bin
```

- Install runc binary to path:

```shell
cp runc.amd64 /usr/local/bin/runc
chmod +x /usr/local/bin/runc
```

- Install cni-plugins to path:

```shell
mkdir -p /opt/cni/bin
tar xf ~/cni-plugins-amd64-v0.7.4.tgz -C /opt/cni/bin
```

- Generate containerd configuration:

```shell

```

- Configure CNI plugin:

```shell
mkdir -p /etc/cni/net.d/
cat <<EOF > /etc/cni/net.d/99-loopback.conf
{
    "cniVersion": "0.3.1",
    "type": "loopback"
}
EOF
```

- Declare containerd systemd service by creating the following file:
  - The file is available in your workspace as `./Lab1/containerd.service`

```service
/etc/systemd/system/containerd.service

[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target

[Service]
ExecStartPre=/sbin/modprobe overlay
ExecStart=/usr/local/bin/containerd
Restart=always
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity

[Install]
WantedBy=multi-user.target
```

- Reload `systemctl` daemons and start containerd.service:

```shell
systemctl daemon-reload
systemctl start containerd
```

- Test containerd:

```shell
# Download an image
crictl pull busybox

# Create a pod (busybox-pod.yaml is in ./Lab1/busybox-pod.yaml)
crictl runp busybox-pod.yaml

# Create a container in the pod (busybox-container.yaml is in ./Lab1/busybox-container.yaml)
crictl create $(crictl pods -q) busybox-container.yaml busybox-pod.yaml

# Start the container
crictl start $(crictl ps -aq)

# Execute a command in the container
crictl exec $(crictl ps -q) ps aux

# Stop and remove the pod
crictl stopp $(crictl pods -q)
crictl rmp $(crictl pods -q)
```

## 1.2 Kubelet

### Install kubelet

Kubernetes binaries are distributed:

- As distribution packages available on specific repositories
- As binaries joined to Release Notes, and grouped by scope:
  - Client binaries
  - Server binaries
  - Node binaries
- Directly, as single binaries
  - [https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/${ARCH}/kubelet](https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/${ARCH}/kubelet)
  - [https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/${ARCH}/kubectl](https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/${ARCH}/kubectl)

- Download kubelet and kubectl binaries:

```shell
export RELEASE=v1.11.2
export ARCH=amd64
export DOWNLOAD_SITE=https://storage.googleapis.com/kubernetes-release
curl -LO ${DOWNLOAD_SITE}/release/${RELEASE}/bin/linux/${ARCH}/kubectl
curl -LO ${DOWNLOAD_SITE}/release/${RELEASE}/bin/linux/${ARCH}/kubelet
```

- Move kubelet and kubectl binaries to the path:

```shell
chmod 755 {kubelet,kubectl}
mv kubectl kubelet /usr/local/bin/
```

- Check the kubelet command parameters: `kubelet --help`

Many configuration options of the kubelet can be set by command line
parameters but this is not the recommended way anymore. You should see
a `--config` parameter.

Example `kubelet-config.yaml`:

```yaml
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
evictionHard:
    memory.available:  "200Mi"
```

Documentation on these parameters is available [here](https://godoc.org/k8s.io/kubelet/config/v1beta1#KubeletConfiguration) and [here](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/apis/config/types.go)

- Start the kubelet with the command `kubelet`
- Does it work?
  >! - Yes, the kubelet can run as standalone
- Check kubelet status with:

```shell
curl -k https://localhost:10250/healthz
```

- Check the directory `/var/lib/kubelet`
  - `pki` stores the self-signed certificate and private key of the kubelet
  - `pods` is the default location for the node static pods
  - `plugins`, `device-plugins`, `container-plugins` are for plugins (CSI,
    advertising node custom resources, ...)
- Stop the kubelet with `killall kubelet`
- Define a static pod with the following descriptor, by placing it in `/var/lib/kubelet/pods`:
  - The file is available in your workspace as `./Lab1/busybox-pod.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: busybox
  namespace: default
spec:
  containers:
  - image: busybox
    command:
      - sleep
      - "3600"
    imagePullPolicy: IfNotPresent
    name: busybox
  restartPolicy: Always
```

- Start kubelet with the command:

```shell
kubelet --pod-manifest-path=/var/lib/kubelet/pods &
```

- Check created containers with: `docker container ls`
- What do you see?
  >! - Containers get created even without an api-server
- Remove the busybox container with: `docker container rm -f <CONTAINER_ID>`
- What happens?
  >! - The container is recreated, according to the pod restartPolicy
  >! - We can see that the kubelet is responsible of containers running in its
  >!   pods

- Stop the kubelet: `killall kubelet`
- List running containers: `docker container ls`
- What do you see?
  >! - Containers are still up
  >! - Your workloads stays online even without the kubelet
- Remove the containers: `docker rm -f <CONTAINERS_IDS>`
- List running containers: `docker container ls`
- What do you see?
  >! - However this time they are not recreated
- Remove the created directories: `rm -Rf /var/lib/kubelet`

<div class="pb"></div>
