# Lab 2

## 2.1: Certificates

Communications in a Kubernetes cluster should be encrypted and peers generally
use TLS mutual authentication.
In order to set this up, there are many certificates to generate.

There are various methods available to generate these certificates:

- With openssl and many commands...
- Using CloudFlare SSL (cfssl) command and config files
- Using kubeadm

We'll use this last method as it's way simpler.

### Kubeadm

- View which certs will be generated with command:

```shell
kubeadm init phase certs
```

- All certificates will be created by default in `/etc/kubernetes/pki`
- Check available options with: `kubeadm alpha phase certs all --help`
  - You can change default destination with: `--cert-dir`
  - You can adapt the API Server exposed certificate with
    `--apiserver-advertise-address` and `--apiserver-cert-extra-sans`
  - If some certificates/keys already exist in the directory they won't be
    overwritten. So if you want to use an existing CA for your PKI, just put the
    `ca.key` and `ca.crt` in the directory and launch the command.
- Generate all certs:

```shell
kubeadm init phase certs all
```

## 2.2: Installation

We'll use `kubeadm` to spawn our cluster.

- Check what will be done by launching the command `kubeadm init --help`
- ⚠️ Don't spawn the full control plane now! We'll do it step by step.
- Check node requirements:

```shell
kubeadm init phase preflight
```

### Control plane

With kubeadm, control-plane components are spawned as static pods on the node's
kubelet.
So the first step is to start the kubelet with the right config on this node:

```shell
kubeadm init phase kubelet-start
```

You should see:

```shell
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Activating the kubelet service
```

To connect to the cluster, clients (cluster components, users, ...) must use a
kubeconfig file which holds their credentials and information on the cluster
(cluster API server endpoint, CA certificate to validate the https exposed
certificate, ...). kubeadm will generate these files for us using the already
generated certificates.

```shell
# See which kubeconfig files will be generated
kubeadm init phase kubeconfig
# Generate them all
kubeadm init phase kubeconfig all
```

We will now use kubeadm to create the static pods manifests for control-plane
components:

```shell
# See what will be generated
kubeadm init phase control-plane
# Generate them all
kubeadm init phase control-plane all
```

- Check that control-plane components containers are running:

```shell
docker container ls
```

- What is happening to the kube-apiserver container? Why?

- Deploy etcd:

```shell
kubeadm init phase etcd local
```

> Note: when deploying a HA cluster, the command
> `kubeadm join [...] --experimental-control-plane` take care of making the
> new control-plane node joining the cluster

- Check that control-plane components and etcd containers are running:

```shell
docker container ls
```

- What is happening to the kube-apiserver container? Why?
- Configure kubectl for cluster admin:

```shell
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
```

- At this point, you should be able to use `kubectl` to interact with the
  cluster.
- Check cluster state with `kubectl get nodes`
- Check cluster components pods with `kubectl get pods -n kube-system`

- Finish the control-plane set up, by
  - Marking this node as control-plane (users workloads won't be scheduled on
    it)
  - Deploying kube-proxy, which handle service to pods traffic distribution
  - Uploading kubelet and kubeadm as ConfigMaps in the cluster

```shell
kubeadm init phase mark-control-plane
kubeadm init phase addon kube-proxy
kubeadm init phase upload-config all
```

### Workers

With Kubernetes 1.12, Bootstrap TLS went GA (General Availability). Using
Bootstrap TLS, workers use a token to contact the cluster control-plane and
a client certificate is generated on the fly. This certificate will
be used by the kubelet to authenticate the node for further communications.
The token authentication information only allow to set up this certificate
request.

- Generate a bootstrap token which will be used by workers to join the
  cluster

```shell
kubeadm init phase bootstrap-token
```

- The last command generated a bootstrap token, but the following one is even
  simpler as it prints the command which we'll use to join the cluster:

```shell
kubeadm token create --print-join-command
```

- You can check the validity of tokens with: `kubeadm token list`

To deploy worker nodes, just execute the join command on each of them.
When this is done, go back to the first node and check the cluster state:

```shell
kubectl get nodes
```

## 2.3: Network solution

At this moment, the cluster isn't in a usable state. You can check it with:
`kubectl get nodes`. It still lacks a network addon to enable cross cluster
communication. There are many availables, you can check it [there](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#pod-network).
We will deploy _Weave Net_.

- Deploy a network solution on the cluster

```shell
sysctl net.bridge.bridge-nf-call-iptables=1
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
```

- Check that nodes are now ready with:

```shell
kubectl get nodes -w
```

- Create 2 pods from these descriptors `nginx-pod.yaml` and `shell-pod.yaml`
- Get the IP of nginx pod: `kubectl get pods -o wide`
- Check that the 2 pods can communicate:
  `kubectl exec shell -- wget -O- <NGINX_IP>`


## 2.4: DNS

At this point, we've enabled cross cluster pod communication.
However, our cluster doesn't have dns service discovery.
To be able to contact services by name, we should deploy a cluster DNS.
For older versions of Kubernetes, we used `kube-dns` but now you should use
`coredns`.

- Deploy CoreDNS on the cluster with kubeadm:
  `kubeadm init phase addon coredns`
- Create a service pointing to nginx with descriptor: `nginx-svc.yaml`
- Check that the service name is resolved with:
  `kubectl exec shell -- wget -O- nginx`
  

## 2.5: Storage

Setting up a full blown storage provisioner is beyond the scope of this
training.
To discover the mechanism behind a storage provisioner, we'll use
[hostpath-provisioner](https://github.com/MaZderMind/hostpath-provisioner).

- Review and deploy the storage provisioner with the provided descriptors:
  - `hostpath-provisioner-deployment.yaml`
  - `hostpath-provisioner-rbac.yaml`
  - `hostpath-provisioner-storageclass.yaml`
- Review and deploy the `test-claim.yaml` and `test-pod.yaml` to the cluster
- Check the pod to know on which node it was deployed:
  `kubectl get pods -o wide`
- Connect to the node and check that the hostpath was created:
  `ls /var/kubernetes`

<div class="pb"></div>


## Cluster up

```shell
kubeadm init phase certs all
kubeadm init phase preflight
kubeadm init phase kubelet-start
kubeadm init phase kubeconfig all
kubeadm init phase control-plane all
kubeadm init phase etcd local
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
kubeadm init phase mark-control-plane
kubeadm init phase addon kube-proxy
kubeadm init phase upload-config all
kubeadm init phase bootstrap-token
kubeadm token create --print-join-command
```
