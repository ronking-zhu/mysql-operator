# mysql-crd
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites
- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
```bash
[root@localhost mysql-crd]# go version
go version go1.24.4 (Red Hat 1.24.4-1.module+el8.10.0+23323+67916f33) linux/amd64
[root@localhost mysql-crd]# minikube version
minikube version: v1.37.0
commit: 65318f4cfff9c12cc87ec9eb8f4cdd57b25047f3
[root@localhost mysql-crd]# kubectl version
Client Version: v1.34.3
Kustomize Version: v5.7.1
Server Version: v1.28.0
Warning: version difference between client (1.34) and server (1.28) exceeds the supported minor version skew of +/-1
[root@localhost mysql-crd]#
```
### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/mysql-crd:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/mysql-crd:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

***run demo***
```bash
``` controller
[root@localhost mysql-crd]# make install
/root/mysql-crd/bin/controller-gen-v0.14.0 rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/root/mysql-crd/bin/kustomize-v5.3.0 build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/mysqlclusters.db.my.domain created
[root@localhost mysql-crd]# kubectl get crd
NAME                         CREATED AT
mysqlclusters.db.my.domain   2025-12-10T15:10:56Z
[root@localhost mysql-crd]# make run
/root/mysql-crd/bin/controller-gen-v0.14.0 rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/root/mysql-crd/bin/controller-gen-v0.14.0 object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
internal/controller/mysqlcluster_controller.go
go vet ./...
go run ./cmd/main.go
2025-12-10T10:11:11-05:00       INFO    setup   starting manager
2025-12-10T10:11:11-05:00       INFO    starting server {"kind": "health probe", "addr": "[::]:8081"}
2025-12-10T10:11:11-05:00       INFO    Starting EventSource    {"controller": "mysqlcluster", "controllerGroup": "db.my.domain", "controllerKind": "MySQLCluster", "source": "kind source: *v1.MySQLCluster"}
2025-12-10T10:11:11-05:00       INFO    Starting Controller     {"controller": "mysqlcluster", "controllerGroup": "db.my.domain", "controllerKind": "MySQLCluster"}
2025-12-10T10:11:11-05:00       INFO    Starting workers        {"controller": "mysqlcluster", "controllerGroup": "db.my.domain", "controllerKind": "MySQLCluster", "worker count": 1}
2025-12-10T10:11:33-05:00       INFO    Creating a new StatefulSet      {"controller": "mysqlcluster", "controllerGroup": "db.my.domain", "controllerKind": "MySQLCluster", "MySQLCluster": {"name":"mysqlcluster-sample","namespace":"default"}, "namespace": "default", "name": "mysqlcluster-sample", "reconcileID": "59b9a59a-aaf9-4bbd-9681-d176d946a679", "Namespace": "default", "Name": "mysqlcluster-sample"}

```
```bash
--- teminal

[root@localhost mysql-crd]# cat config/samples/db_v1_mysqlcluster.yaml
apiVersion: db.my.domain/v1
kind: MySQLCluster
metadata:
  labels:
    app.kubernetes.io/name: mysql-crd
    app.kubernetes.io/managed-by: kustomize
  name: mysqlcluster-sample
spec:
  # TODO(user): Add fields here
        replicas: 1
[root@localhost mysql-crd]# kubectl apply -f config/samples/db_v1_mysqlcluster.yaml
mysqlcluster.db.my.domain/mysqlcluster-sample created
[root@localhost mysql-crd]# kubectl get pods
NAME                    READY   STATUS              RESTARTS   AGE
mysqlcluster-sample-0   0/1     ContainerCreating   0          7s
[root@localhost mysql-crd]# kubectl get pods
NAME                    READY   STATUS    RESTARTS   AGE
mysqlcluster-sample-0   1/1     Running   0          44s
[root@localhost mysql-crd]# kubectl get pods
NAME                    READY   STATUS    RESTARTS   AGE
mysqlcluster-sample-0   1/1     Running   0          45s
```
### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/mysql-crd:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/mysql-crd/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

