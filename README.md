# xenia

A discord project solution.

## toolchain

- Compute Platform
  - [Azure Kubernetes Service](https://azure.microsoft.com/en-us/products/kubernetes-service) w/ [Kubernetes](https://kubernetes.io/)
  - [Azure Virtual Machines](https://azure.microsoft.com/en-us/products/virtual-machines) w/ [Debian](https://www.debian.org/)
- Package Management
  - [Helm](https://helm.sh)
  - [podman](https://podman.io/)
  - [Azure Container Registry](https://azure.microsoft.com/en-us/products/container-registry)
- Ingress Controller
  - [nginx](https://www.nginx.com/)
  - [haproxy](https://www.haproxy.org/)
  - [???](https://linkerd.io/2.12/tasks/using-ingress/)
- Secrets Management
  - [akv2k8s](https://akv2k8s.io/)
  - [sops](https://github.com/mozilla/sops)
- Service Mesh
  - [linkerd](https://linkerd.io/)
- Observability
  - [Azure Monitor managed service for Prometheus](https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/prometheus-metrics-overview)
- Pipelines
  - [github actions](https://github.com/features/actions)
  - [tekton](https://tekton.dev/)
  - [argocd](https://argo-cd.readthedocs.io/en/stable/)
- IaC
  - [cdktf](https://developer.hashicorp.com/terraform/cdktf) w/ [go](https://go.dev/)

## references

- [terraform cdk go examples](https://github.com/hashicorp/terraform-cdk/tree/main/examples/go)
