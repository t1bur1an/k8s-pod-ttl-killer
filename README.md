# k8s-pod-ttl-killer

## Description

Project designed to kill pods in kubernetes cluster filtered by annotation and diff between annotation ttl in seconds and pod ReadyTime.

## Environment variables

All default environment variables can be found in `.env` file. For now it is:
```
CHECK_INTERVAL_SECONDS=10
TTL_ANNOTATION=k8s-pod-ttl-killer
HTTP_LISTEN_ADDRESS=0.0.0.0
HTTP_LISTEN_PORT=8080
```

So service would check pods every 10 seconds.

Check pod meta->annotations->k8s-pod-ttl-killer which is integer of seconds allowed to pod to live.

If ReadyTime+annotation seconds >= now() service would kill the pod by sending `delete` command to kubernetes api.

## Example

To play around with service there is `example` folder contains manifest file with example pod with annotation. Service would find it easy and kill without any problems.

Service can be also runned out of PC without and use personal `.kube/config` file by default.

In kubernetes it uses kubernetes API with Service Account.

## Helm

Service can be started in cluster by helm chart. It will start one instance of service and create all needed resources like ClusterRole and ClusterRoleBinding to project Service Account.
