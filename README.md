# configmapsync-operator

* `operator-sdk init --domain=example.com --repo=github.com/mikolajsemeniuk/configmapsync-operator`
* `operator-sdk create api --group=sync --version=v1alpha1 --kind=ConfigMapSync --resource --controller`
* update `api/v1alpha1/configmapsync_types.go` and run `make generate` and then `make manifests`
* update `internal/controller/configmapsync_controller.go`
* run `make docker-build IMG=mikolajsemeniuk/configmapsync:latest`
* run `make docker-push IMG=mikolajsemeniuk/configmapsync:latest`
* run `make install`
* run `make deploy IMG=mikolajsemeniuk/configmapsync:latest`
* check `kubectl get pods -n configmapsync-operator-system`
* check `kubectl get all -n configmapsync-operator-system`

## Test operator

* `kubectl create namespace source-ns`
* `kubectl create configmap my-config --from-literal=key1=value1 -n source-ns`
* `kubectl create namespace target-ns1`
* `kubectl create namespace target-ns2`

```yaml
# configmapsync.yaml
apiVersion: sync.example.com/v1alpha1
kind: ConfigMapSync
metadata:
  name: example-configmapsync
  namespace: source-ns
spec:
  sourceNamespace: source-ns
  sourceName: my-config
  targetNamespaces:
    - target-ns1
    - target-ns2

```

* `kubectl apply -f configmapsync.yaml`
* `kubectl get configmaps -n target-ns1`
* `kubectl get configmaps -n target-ns2`
* `kubectl create configmap my-config --from-literal=key1=value-updated -n source-ns --dry-run=client -o yaml | kubectl apply -f -`
* `kubectl get configmap my-config -n target-ns1 -o yaml`

## Debug operator

* `make install run`
* debug application

## Description

* <https://medium.com/developingnodes/mastering-kubernetes-operators-your-definitive-guide-to-starting-strong-70ff43579eb9>
* <https://suedbroecker.net/2022/03/01/debug-a-kubernetes-operator-written-in-go/>

## Getting Started

### Prerequisites

* go version v1.22.0+

* docker version 17.03+.
* kubectl version v1.11.3+.
* Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/configmapsync-operator:tag
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
make deploy IMG=<some-registry>/configmapsync-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

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
make build-installer IMG=<some-registry>/configmapsync-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/configmapsync-operator/<tag or branch>/dist/install.yaml
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
