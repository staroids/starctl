# StarCtl

CLI for staroid.com, cloud platform based on Kubernetes that funds open-source developers.

## Download

TBD

## Usage

```
star <create|delete|get> ske <name>
star <create|delete|get|stop|start> namespace --alias=<name> --ske=<name> --project <org/repo:branch#commit>
star <create|delete|get> proxy --alias=<name> --ske=<name>
star <create|delete> proxy --alias=<name> --ske=<name> <local port> [<local-host>:<local-port>:<remote-host>:<remote-port>] [...]
```
