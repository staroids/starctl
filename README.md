# starctl

CLI for staroid.com, cloud platform based on Kubernetes that funds open-source developers.

## Download

https://github.com/staroids/starctl/releases

## Usage

### Cluster
```
export STAROID_ACCESS_TOKEN=xxxxxxxxxx

# list all clusters
starctl cluster list
```

### Namespace

```
# list all namespaces in the clusuter
starctl namespace -org <org> -cluster <cluster> list

# create a namespace
starctl namespace -org <org> -cluster <cluster> -wait create <alias>

# delete a namespace
starctl namespace -org <org> -cluster <cluster> -wait delete <alias>

# stop all deployments/pods/jobs in namespace (but keep configmaps, secrets)
starctl namespace -org <org> -cluster <cluster> -wait stop <alias>

# bring all deployment/pod/job back online 
starctl namespace -org <org> -cluster <cluster> -wait start <alias>
```

### Shell

```
# Start a shell service in the namespace
starctl shell -org <org> -cluster <cluster> start <alias>

# Stop a shell service in the namespace
starctl shell -org <org> -cluster <cluster> stop <alias>
```

### Tunnel

```
# Open a tunnel to Kubernetets api proxy.
# Local port 8001 (can change using --kube-proxy-port) will be connected to the Kubernetes API of the cluster.
# (Tunnel 'kubectl --server localhost:8001' will )

starctl tunnel --kube-proxy


# Open tunnels to services running clustter.
# (traffic to Local port 7000 will be forwarded to 'my-service1:8000' in the cluster, in following example)

starctl tunnel 7000:my-service1:8000 1234:my-service2:5678 ...

```

## Environment variables

| Variable name | Optional | Description |
| --------- | -------- | --------- |
| STAROID_ACCESS_TOKEN | Required | Access token string. (e.g. `v0hsolmc6vu1tpnp4vtv8c8solvgt0`) Get from [Access Tokens menu](https://staroid.com/settings/accesstokens). |

## Build

```
$ make
```
