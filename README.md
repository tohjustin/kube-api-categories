# kube-api-categories

[![build](https://github.com/tohjustin/kube-api-categories/actions/workflows/build.yaml/badge.svg)](https://github.com/tohjustin/kube-api-categories/actions/workflows/build.yaml)
[![kubernetes compatibility](https://aegisbadges.appspot.com/static?subject=k8s%20compatibility&status=v1.19%2B&color=318FE0)](https://endoflife.date/kubernetes)
[![license](https://aegisbadges.appspot.com/static?subject=license&status=Apache-2.0&color=318FE0)](./LICENSE.md)

A CLI tool to display all API categories & their respective resources in a Kubernetes cluster.

## Usage

```shell
$ kube-api-categories
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
$ kube-api-categories api-extensions
apiservices.apiregistration.k8s.io
customresourcedefinitions.apiextensions.k8s.io
mutatingwebhookconfigurations.admissionregistration.k8s.io
validatingwebhookconfigurations.admissionregistration.k8s.io

$ kube-api-categories prometheus-operator
alertmanagerconfigs.monitoring.coreos.com
alertmanagers.monitoring.coreos.com
podmonitors.monitoring.coreos.com
probes.monitoring.coreos.com
prometheuses.monitoring.coreos.com
prometheusrules.monitoring.coreos.com
servicemonitors.monitoring.coreos.com
thanosrulers.monitoring.coreos.com
```

## Installation

### Install from Source

```shell
$ git clone git@github.com:tohjustin/kube-api-categories.git && cd kube-api-categories
$ make install

$ kube-api-categories --version
```
