<p align="center">

  <a href="https://hub.docker.com/r/banzaicloud/kafka-operator/">
    <img src="https://img.shields.io/docker/cloud/automated/banzaicloud/kafka-operator.svg" alt="Docker Automated build">
  </a>

  <a href="https://circleci.com/gh/banzaicloud/kafka-operator/tree/master">
    <img src="https://circleci.com/gh/banzaicloud/kafka-operator/tree/master.svg?style=shield" alt="CircleCI">
  </a>


  <a href="https://goreportcard.com/report/github.com/banzaicloud/kafka-operator">
    <img src="https://goreportcard.com/badge/github.com/banzaicloud/kafka-operator" alt="Go Report Card">
  </a>

  <a href="https://github.com/banzaicloud/kafka-operator/">
    <img src="https://img.shields.io/badge/license-Apache%20v2-orange.svg" alt="license">
  </a>

</p>

# Kafka-Operator

Kafka-operator is a Kubernetes operator to deploy and manage [Apache Kafka](https://kafka.apache.org).

## Overview

[Kafka](https://kafka.apache.org) is an open-source distributed streaming platform. It has three capabilities:
- Publish and subscribe to streams of records, similar to a message queue or enterprise messaging system.
- Store streams of records in a fault-tolerant durable way.
- Process streams of records as they occur.

The goal of the **Kafka-operator** is to automate and simplify installing and managing Kafka on Kubernetes. 

![Kafka Operator](/docs/img/operator.png)

### Motivation

At [Banzai Cloud](https://banzaicloud.com) we are building a Kubernetes distribution, [PKE](https://github.com/banzaicloud/pke) and platform, [Pipeline](https://github.com/banzaicloud/pipeline) and operate managed Kafka clusters for our customers.

There is a huge interest in the Kafka community for a solution which enables Kafka on Kubernetes.
There were several initiatives to simplify Kafka on Kubernetes:
- [Helm chart](https://github.com/confluentinc/cp-helm-charts/tree/master/charts/cp-kafka)
- [Yaml files](https://github.com/Yolean/kubernetes-kafka)
- [Strimzi Kafka Operator](https://github.com/strimzi/strimzi-kafka-operator)

however, none of these gave a full solution to **automate** the Kafka experience and make it consumable for the wider audience. Our motivation is to build an open source solution and a community which drives the innovation and features of this operator.

If you are willing to kickstart your Kafka experience using Pipeline, check out the free developer beta:
<p align="center">
  <a href="https://beta.banzaicloud.io">
  <img src="https://camo.githubusercontent.com/a487fb3128bcd1ef9fc1bf97ead8d6d6a442049a/68747470733a2f2f62616e7a6169636c6f75642e636f6d2f696d672f7472795f706970656c696e655f627574746f6e2e737667">
  </a>
</p>

## Installation

The operator installs the 2.1.0 version of Kafka, and can run on Minikube v0.33.1+ and Kubernetes 1.12.0+.

As a pre-requisite it needs a Kubernetes cluster (you can create one using [Pipeline](https://github.com/banzaicloud/pipeline)).
Also, Kafka requires Zookeeper so you need to first start a Zookeeper server if you don't already have one.

> Banzai Cloud's Kafka operator is a pure Kafka Operator without Zookeeper support.

##### Install Zookeeper

To install Zookeeper we recommend using [Pravega's Zookeeper Operator](https://github.com/pravega/zookeeper-operator).
You can deploy Zookeeper by using a Helm chart.

```bash
helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com/
helm install --name zookeeper -namespace=zookeeper banzaicloud-stable/zookeeper
```

1. Set `KUBECONFIG` pointing towards your cluster 
2. Run `make deploy` (deploys the operator in the `kafka` namespace to the cluster)
3. Set your Kafka configurations in a Kubernetes custom resource (sample: `config/samples/banzaicloud_v1alpha1_kafkacluster.yaml`) and run this command to deploy the Kafka components:

```bash
# Add your zookeeper svc name to the configuration
kubectl create -n kafka -f config/samples/banzaicloud_v1alpha1_kafkacluster.yaml
```

##### Alert

TODO

### Installation with Helm

Alternatively, if you are using Helm, you can deploy the operator using a Helm chart [Helm chart](https://github.com/banzaicloud/kafka-operator/tree/master/charts):

```bash
helm repo add banzaicloud-stable https://kubernetes-charts.banzaicloud.com/
helm install --name=kafka --namespace=kafka banzaicloud-stable/kafka-operator
```

## Development

Check out the [developer docs](docs/developer.md).

## Features

Check out the [supported features](docs/feature.md).

## Issues, feature requests and roadmap

Please note that the Kafka operator is constantly under development and new releases might introduce breaking changes. We are striving to keep backward compatibility as much as possible while adding new features at a fast pace. Issues, new features or bugs are tracked on the projects [GitHub page](https://github.com/banzaicloud/kafka-operator/issues) - please feel free to add yours!

To track some of the significant features and future items from the roadmap please visit the [roadmap doc](docs/roadmap.md).

## Contributing

If you find this project useful here's how you can help:

- Send a pull request with your new features and bug fixes
- Help new users with issues they may encounter
- Support the development of this project and star this repo!

## License

Copyright (c) 2017-2019 [Banzai Cloud, Inc.](https://banzaicloud.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
