![](./images/k8sArch.png)

- the controller makes sure that the real state matches the desired state

Deployment
—Replicate
——Pod
————Container

- Actually, a Deployment doesn’t manage replicas directly: instead, it automatically creates an associated object called a ReplicaSet
- A *Pod* is the Kubernetes object that represents a group of one or more containers
- Why doesn’t a Deployment just manage an individual container directly? The answer is that sometimes a set of containers needs to be scheduled together, running on the same node, and communicating locally, perhaps sharing storage. 

```shell
# kubectl create deployment NAME --image=image [--dry-run=server|client|none] [options]
kubectl create deployment mongo-deploy --image=mongo
kubectl describe deployments.apps mongo-deploy
kubectl logs pod/mongo-deploy-76c8d88f55-4x5fp
kubectl edit deployments.apps mongo-deploy
kubectl delete deployments.apps mongo-deploy
kubectl get endpoints
```

```shell
yum install kubectl-1.19.16-0 kubelet-1.19.16-0 kubeadm-1.19.16-0 --disableexcludes=kubernetes
kubeadm upgrade plan
kubeadm upgrade apply v1.19.16 --certificate-renewal=false 
```