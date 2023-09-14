# Tomcat Deploy

## Build 

### Prerequisites

- [Install Golang](https://go.dev/doc/install)

### Build
1. Init project
```shell
go mod init tomcat_deploy
```
2. Get requirements
```shell
go mod tidy
```
3. Build project
```shell
go build ./tomcat_deploy
```

## Overview

### Flags
```shell
./tomcat_deploy.exe -help
```

- `-h` - host to deploy. Usage: `-h rc.erp.sperasoft.com`. If not set, deploy will be performed to k8s tomcat pod

- `-b` - build before deploy. Default: `true`. Usage: `-b=false`

### Process
1. Build a maven project in the dir where `tomcat_deploy.exe` executed.
2. Find `*.war` file in the `/target` dir
3. Rename `*.war` file to service name. e.g. `wa-event2email-1.0-SNAPSHOT.war` will be renamed to `event2email.war`
4. Copy `*.war` file to the certain stand/k8s pod
- k8s pod will be defined in `getPodName` func
- deploy to host uses scp and ssh. Requires to input password twice, for scp copying and for move file as sudo.