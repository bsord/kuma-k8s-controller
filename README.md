# kuma-k8s-operator
Kubernetes operator that posts ingresses to an uptime monitor.
## Install
```sh
helm repo add ingress-nginx https://cfcr.io/bsord/helm-charts
helm repo update
helm install kuma-k8s-operator bsord/kuma-k8s-operator --set monitorUrl="https://yourkumauptimehost/api"
```

## Setting up development environment
It is assumed you have golang installed.
[Get started with Golang](https://linuxize.com/post/how-to-install-go-on-ubuntu-20-04/)

1. Install a local cluster using minikube/microk8s/kind
This example will use a local kind cluster:
```sh
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/bin/kind
```

2. Install helm
Helm is used for simplicity of deployment of other useful resources and is a dependency of Skaffold
```sh
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

3. Install Skaffold
Then enables real time and local development of this project in a kubernetes environment.
```sh
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
sudo install skaffold /usr/local/bin/
```

4. Deploy an ingress controller
An ingress controller is required in the cluster for this project to be useful, nginx is used here but others should work.
```sh
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx
```

5. Deploy this project in development mode
Start the development environment
```sh
skaffold dev
```

6. Create a deployment with an ingress to trigger an event(Optional):
This is useful for testing.
```sh
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx
```

Switch back to your other kube contexts if you have any, as installing a local cluster may have changed your default context.
```
cat ~/.kube/config # Find your context here
kubectl config use-context YOUR_OTHER_CONTEXT # switch to your desired context
```

## TODO:
- [x] Initial implementation
- [ ] Use watch and kubeinformers + cache to catch events real time (ideally behind a Cobra config/flag)
- [ ] Restore proper sig handling in Cobra
- [ ] Implement config map, secrets, arg passing from helm to Cobra
- [ ] Define basic models and attributes needed to build a monitor definitions
- [ ] Implement http post method to Slack on ingress events to test models and posting implementation
- [ ] Implement reconciliation/sync for list modes to avoid making unnecessary http posts
- [ ] Cache go modules/packages during docker build process to avoid lengthy build times
- [ ] Add github actions workflow to auto version bump, build/publish docker, and publish helm chart
- [ ] :allthethings: