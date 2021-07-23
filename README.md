![Node-Check Logo](https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png)

[![Build Status](https://drone.support.tools/api/badges/mattmattox/node-check/status.svg)](https://drone.support.tools/mattmattox/node-check)
![Docker Pulls](https://img.shields.io/docker/pulls/supporttools/node-check.svg)

# What is Node-Check
Node Check is a small pod that runs on all nodes in a Kubernetes cluster. This pod listens on port 8888 and responds with the current node health status. Node Check was primarily designed for unauthenticated services like an external load balancer to make simple GET requests to get the current node status. This is needed for an external load balancer to safely drain connections for cordoned or currently listed as not ready nodes, even if they are still responding to port 80.

## How does Node-Check work?
Node Check is a GO webserver that connects to the Kubernetes API to capture the node conditions for the node where the pod is running.

NOTE: All access is Read-Only and is as minimum as possible.

## Installation steps

### kubectl apply
```
kubectl apply -f https://raw.githubusercontent.com/mattmattox/node-check/main/deployment.yaml
```

### Helm chart
```
helm repo add support-tools https://charts.support.tools
helm repo update
helm upgrade --install \
node-check support-tools/node-check \
--namespace node-check
```