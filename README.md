# kube-api-categories

[![build](https://github.com/tohjustin/kube-api-categories/actions/workflows/build.yaml/badge.svg)](https://github.com/tohjustin/kube-api-categories/actions/workflows/build.yaml)
[![release](https://aegisbadges.appspot.com/static?subject=release&status=v0.1.0&color=318FE0)](https://github.com/tohjustin/kube-api-categories/releases)
[![kubernetes compatibility](https://aegisbadges.appspot.com/static?subject=k8s%20compatibility&status=v1.19%2B&color=318FE0)](https://endoflife.date/kubernetes)
[![license](https://aegisbadges.appspot.com/static?subject=license&status=Apache-2.0&color=318FE0)](./LICENSE.md)

A CLI tool to display all API categories & resources in a Kubernetes cluster.

## Usage

```shell
$ kube-api-categories
RESOURCE                          APIGROUP                       NAMESPACED   CATEGORIES
bindings                                                         true         []
componentstatuses                                                false        []
configmaps                                                       true         []
endpoints                                                        true         []
events                                                           true         []
limitranges                                                      true         []
namespaces                                                       false        []
nodes                                                            false        []
persistentvolumeclaims                                           true         []
persistentvolumes                                                false        []
pods                                                             true         [all]
podtemplates                                                     true         []
replicationcontrollers                                           true         [all]
resourcequotas                                                   true         []
secrets                                                          true         []
serviceaccounts                                                  true         []
services                                                         true         [all]
mutatingwebhookconfigurations     admissionregistration.k8s.io   false        [api-extensions]
validatingwebhookconfigurations   admissionregistration.k8s.io   false        [api-extensions]
customresourcedefinitions         apiextensions.k8s.io           false        [api-extensions]
...
```

View list of available API categories

```shell
$ kube-api-categories --output=category
all
api-extensions
apim-istio-io
cert-manager
cert-manager-acme
constraint
constraints
networking-istio-io
policy-istio-io
prometheus-operator
rbac-istio-io
security-istio-io
telemetry-istio-io
```

View resources in a specific API category

```shell
$ kube-api-categories --output=resource --categories=all
cronjobs.batch
daemonsets.apps
deployments.apps
horizontalpodautoscalers.autoscaling
jobs.batch
pods
replicasets.apps
replicationcontrollers
services
statefulsets.apps

$ kube-api-categories --output=resource --categories=prometheus-operator
alertmanagerconfigs.monitoring.coreos.com
alertmanagers.monitoring.coreos.com
podmonitors.monitoring.coreos.com
probes.monitoring.coreos.com
prometheuses.monitoring.coreos.com
prometheusrules.monitoring.coreos.com
servicemonitors.monitoring.coreos.com
thanosrulers.monitoring.coreos.com
```

### Flags

Flags for configuring relationship discovery parameters

| Flag | Description |
| ---- | ----------- |
| `--api-group=''`     | Limit to resources in the specified API group |
| `--cached=false`     | Use the cached list of resources if available |
| `--categories=[]`    | Limit to resources that belong to the specified categories |
| `--namespaced=true`  | If false, non-namespaced resources will be returned, <br/> otherwise returning namespaced resources by default. |
| `--no-headers=false` | When using the default output format, don't print headers (default print headers) |
| `--output=''`, `-o`  | Output format. One of: category\|resource. |

Use the following commands to view the full list of supported flags

```shell
$ kube-api-categories --help
```

## Installation

### Install via [krew](https://krew.sigs.k8s.io/)

The tool is available via the [tohjustin/kubectl-plugins](https://github.com/tohjustin/kubectl-plugins) custom plugin index

```shell
$ kubectl krew index add tohjustin https://github.com/tohjustin/kubectl-plugins.git
$ kubectl krew install tohjustin/api-categories

$ kubectl api-categories --version
```

### Install from Source

```shell
$ git clone git@github.com:tohjustin/kube-api-categories.git && cd kube-api-categories
$ make install

$ kube-api-categories --version
```
