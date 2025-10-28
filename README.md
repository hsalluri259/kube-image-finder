# kube-image-finder 🕵️‍♂️

A lightweight Go CLI tool to **search for Docker images across all Kubernetes deployments in all namespaces and all kube contexts**.  
Useful for cluster audits, image version tracking, and migration planning.

---

## 🚀 Features

- Searches for a specific image (or partial match) in **all namespaces**  
- Lists matching **Kube Context**, **Namespaces**, and **Pod names**  
- Fast and efficient — uses the Kubernetes Go client 

---

## Installation

### Prerequisites
- Go 1.22+
- Access to a valid kubeconfig (~/.kube/config)

### Build
```bash
git clone https://github.com/<your-username>/kube-image-finder.git
cd kube-image-finder
go build -o kube-image-finder
./kube-image-finder --image "docker.io/busybox" --kubeconfig /path/to/.kube/config
```

## 🧰 Example Usage

```bash
# This example shows without build
# Search for all deployments using "docker.io/busybox" images
go run main.go --image "docker.io/busybox" --kubeconfig "/Users/harish/.kube/config"
```