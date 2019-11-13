# Lab 3

## 3.1: EFK Deployment

- Check logs manually from various sources (kubelet, container runtime, ...)
- Deploy EFK with provided descriptors: `kubectl apply -f .`

- Check EFK deployment state: `kubectl get pods -n kube-system`
- Label nodes with:

```shell
kubectl label node node-{0,1,2,3} beta.kubernetes.io/fluentd-ds-ready=true
```

- Find Kibana exposed NodePort: `kubectl get svc kibana-logging -n kube-system`
- Connect to Kibana check logs

- Deploy a container producing logs
- Check the container logs on kibana

## 3.2: Prometheus Deployment

Check exposed metrics (Is there a endpoint accessible with kube metrics?)

Deploy prometheus

## 3.3: Backup/Restore

Check etcd state with CLI

Perform a backup

Check backup

## 3.4: Capacity planning

### 3.4.1: Cluster metrics

Test kubectl top pods/kubectl top nodes
Check if it comes from heapster or metrics-server

### 3.4.2: Limits/Requests

Define some Limits/Requests
Check Pod behaviour
Check assinged QoS

### 3.4.3: PodDisruptionBudget/PriorityClass (Optional)

Configure a PodDisruptionBudget and a PriorityClass

## 3.5: Garbage collection

Configure 

## 3.6: Ingress

Deploy an Ingress Controller
Create an Ingress
Check that it's working

## 3.8: Affinity

Configure a Pod with Affinity
Configure a Pod with Anti-affinity

## 3.9: Taints & Tolerations

Configure a Pod with Tolerations
Assign a Taint to a Node


## 3.10: Schedulers

Configure a custom scheduler
Configure a Pod using this custom scheduler

## 3.11: Troubleshooting

Break a Node
Check what happens (kubectl get nodes, logs, ...)

<div class="pb"></div>
