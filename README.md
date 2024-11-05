# RADAR-Kubernetes

[![GitHub release](https://img.shields.io/github/v/release/radar-base/radar-kubernetes)](https://github.com/RADAR-base/RADAR-Kubernetes/releases/latest)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/radar-base/radar-kubernetes/push.yaml)
[![Apache Licensed](https://img.shields.io/github/license/radar-base/radar-kubernetes)](LICENSE)
![GitHub Discussions](https://img.shields.io/github/discussions/radar-base/radar-kubernetes) ![Maintenance](https://img.shields.io/maintenance/yes/2023)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/radar-base/radar-kubernetes/dev)
[![Join our community Slack](https://img.shields.io/badge/slack-radarbase-success.svg?logo=slack)](https://docs.google.com/forms/d/e/1FAIpQLScKNZ-QonmxNkekDMLLbP-b_IrNHyDRuQValBy1BAsLOjEFpg/viewform)

The Kubernetes stack of RADAR-base platform.

## Table of contents

<!-- TOC start (generated with https://github.com/derlin/bitdowntoc) -->

- [About](#about)
- [Status](#status)
- [Prerequisites](#prerequisites)
   * [Knowledge requirements ](#knowledge-requirements)
   * [Software Compatibility](#software-compatibility)
   * [Hosting ](#hosting)
   * [Third party services](#third-party-services)
   * [Local machine](#local-machine)
- [Installation](#installation)
   * [Prepare](#prepare)
   * [Project Structure](#project-structure)
   * [Configure](#configure)
   * [Install](#install)
- [Usage and accessing the applications](#usage-and-accessing-the-applications)
- [Service-specific documentation](#service-specific-documentation)
- [Troubleshooting](#troubleshooting)
- [Volume expansion](#volume-expansion)
- [Uninstall](#uninstall)
- [Update charts](#update-charts)
- [Development automation](#development-automation)
- [Feedback and Contributions](#feedback-and-contributions)

<!-- TOC end -->

## About

RADAR-base is an open-source platform designed to support remote clinical trials by collecting continuous data from wearables and mobile applications. RADAR-Kubernetes enables installing the RADAR-base platform onto Kubernetes clusters. RADAR-base platform can be used for wide range of use-cases. Depending on the use-case, the selection of applications need to be installed can vary. Please read the [component overview and breakdown](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown) to understand the role of each component and how components work together.

RADAR-Kubernetes setup uses [Helm](https://github.com/helm/helm) charts to package necessary Kubernetes resources for each component and [helmfile](https://github.com/roboll/helmfile) to modularize and deploy Helm charts of the platform on a Kubernetes cluster. This setup is designed to be a lightweight way to install and configure the RADAR-base components. The original images or charts may provide more and granular configurations. Please visit the `README` of respective charts in [radar-helm-charts](https://github.com/RADAR-base/radar-helm-charts) to understand the configurations and visit the main repository for in depth knowledge.

## Status

RADAR-Kubernetes is one of the youngest project of RADAR-base and will be the **long term supported form of deploying the platform**. Even though, RADAR-Kubernetes is being used in few production environments, it is still in its early stage of development. We are working on improving the set up and documentation to enable RADAR-base community to make use of the platform.

## Prerequisites

### Knowledge requirements 

This documentation assumes familiarity with all referenced Kubernetes concepts, utilities, and procedures and familiarity with Helm charts and helmfile, depending on your environment you might need to have knowledege of other hosting infrastructure such as DNS and mail servers as well.
While this documentation will provide guidance for installing and configuring RADAR-base platform on a Kubernetes cluster and tries to make is as simple and possible, it is not a replacement for the detailed knowledege of the tools that have been used. If you are not familiar with these tools, we strongly recommend you to get familiar with these tools. Here is a [list of useful links](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2731638785/How+to+get+started+with+tools+around+RADAR-Kubernetes) to get started.

### Software Compatibility

Currently RADAR-Kubernetes is tested and supported on following component versions:
| Component | Version |
| ---- | ------- |
| Kubernetes | v1.27 to v1.30 |
| K3s | v1.27.14+k3s1 to v1.30.1+k3s1 |
| Kubectl | v1.27 to v1.30 |
| Helm | v3.11.3 |
| Helm diff | v3.6.0 |
| Helmfile | v0.152.0 |
| YQ | v4.33.3 |

It's possible to install RADAR-Kubernetes on different version of tools as well, but you might encounter compatibility issues. Make sure that `kubectl` version matches or it's higher than Kubernetes or K3s version that you're using. For other tools such as Git or Java the version, as long as it's not very old, it's not very impactful.

### Hosting 
Kubernetes can be installed on wide varaity of platforms and in turn you can install RADAR-Base on most places that Kuberentes can run. However your infrastructure needs to have a set up requirements listed below:

| Component          | Description                                                                                                                                                                                                                                                                                                                                                                                                                                | Required |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| Kubernetes cluster | An infrastructure with working installation of Kubernetes services. Read [this article](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2744942595?draftShareId=e09429e8-38c8-4b71-955d-5df8de94b694) for available options. Minimum requirements for a single node: 8 vCPU's, 32 GB memory, 200 GB storage. Minimum requirements for a cluster: 3 nodes with 3 vCPUs, 16 GB memory, 100 GB storage each and 200 GB shared storage. | Required |
| DNS Server         | Some applications are only accessible via HTTPS and it's essential to have a DNS server via providers like GoDaddy, Route53, etc                                                                                                                                                                                                                                                                                                           | Required |
| SMTP Server        | RADAR-Base needs an SMTP server to send registration email to researchers and participants.                                                                                                                                                                                                                                                                                                                                                | Required |
| Whitelisted acces to ports 80 and 443 | We use Let's Encrypt to create SSL certificates and in the default configuration we use HTTP challenge. This means that the RADAR-Base installation needs to be visible to Let's Encrypt servers for the verification, so make sure these ports are white listed in your firewall. If you want to have a private installation you should change Let's Encrypt configuration to use DNS challenge.  | Required |
| Object storage     | An external object storage allows RADAR-Kubernetes to backup cluster data such as manifests, application configuration and data via Velero to a backup site. You can also send the RADAR-Base output data to this object storage, which can provider easier management and access compared to bundled Minio server inside RADAR-Kubernetes.                                                                                                | Optional |
| Managed services   | RADAR-Kubernetes includes all necessary components to run the platform as a standalone application. However, you can also opt to use managed services such as with the platform, e.g. Confluent cloud for Kafka and schema registry, Postgres DB for storage, Azure blob storage or AWS S3 instead of Minio. If you're using a managed object storage that you have to pay per request (such as AWS S3), it's recommeneded that to install the Minio just for the `radar-intermediate-storage` since the applications send a lot of API calls to that bucket. | Optional |

In order to have a simple single node Kubernetes server you can run these commands on a Linux server:
```shell
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.26.3+k3s1" K3S_KUBECONFIG_MODE="644" INSTALL_K3S_SYMLINK="skip" sh -s - --disable traefik --disable-helm-controller
```


### Third party services
Depending on which components you've enabled you might need credentials for Fitbit, REDCap, Google Firebase, etc. You need to provide them in order for the respective component to work properly.

### Local machine

The following tools should be installed in your local machine to install the RADAR-Kubernetes on your Kubernetes cluster.

| Component                                                          | Description                                                                                                                                                                                                      |
| ------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Git](https://git-scm.com/downloads)                               | RADAR-Kubernetes uses Git-submodules to use some third party Helm charts. Thus Git is required to properly download and sync correct versions of this repository and its dependent repositories                  |
| [Java](https://openjdk.java.net/install/)                          | The installation setup uses Java Keytools to create Keystore files necessary for signing access tokens.                                                                                                          |
| [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) | Kubernetes command-line tool, kubectl, allows you to run commands against Kubernetes clusters                                                                                                                    |
| [helm 3](https://github.com/helm/helm#install)                     | Helm Charts are used to package Kubernetes resources for each component                                                                                                                                          |
| [helmfile](https://github.com/helmfile/helmfile#installation)      | RADAR-Kubernetes uses helmfiles to deploy Helm charts.                                                                                                                                                           |
| [helm-diff](https://github.com/databus23/helm-diff#install)        | A dependency for Helmfile.                                                                                                                                                                                       |
| [yq](https://github.com/mikefarah/yq#install)                      | Used to run `init`, `generate-secrets` and `chart-updates` scripts.                                                                                                                                              |
| openssl                                                            | Used in `init` and `generate-secrets` scripts to generate secret for Prometheus Nginx authentication. This binary is in `openssl` package for Ubuntu, it's also easily available on other distributions as well. |

**Once you have a working installation of a Kubernetes cluster, please [configure Kubectl with the appropriate Kubeconfig](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#verify-kubectl-configuration) to enable Kubectl to find and access your cluster. Then proceed to the installation section.**

## Installation

> The following instructions on this guide are for local machines which runs on Linux operating systems. You can still use the same instructions with small to no changes on a MacOS device as well.

### Prepare

1. Clone the repository to your local machine by using following command.

   ```shell
   git clone https://github.com/RADAR-base/RADAR-Kubernetes.git
   ```

2. Run the initialization script to create basic configuration files.

   ```shell
   cd RADAR-Kubernetes
   bin/init
   ```

It is recommended make a private clone of this repository, if you want to version control your configurations and/or share with other people.

**You must keep `etc/secrets.yaml` secure and confidential once you have started installing the platform** and the best practice to share your platform configurations is by **sharing the encrypted version of `etc/secrets.yaml`, this can be done via the combination of [sops](https://github.com/getsops/sops) and [helm-secrets](https://github.com/jkroepke/helm-secrets) but it's outside the scope of this document**.

### Project Structure

- `bin/`: Contains initialization scripts.
- `etc/`: Contains configurations for the Helm charts.
- `helmfile.d/`: Contains Helmfiles for modular deployment of the platform.
- `environments.yaml/`: Defines current environment files in order to be used by helmfile and where to find the configuration files.
- `etc/production.yaml`: Production helmfile template to configure and install RADAR-base components. Inspect the file to enable, disable and configure components required for your use case. The default helmfile enables all core components that are needed to run RADAR-base platform with pRMT and aRMT apps. If you're not sure which components you want to enable you can refer to wiki for [an overview and breakdown on RADAR-Base components and their roles](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown).
- `etc/production.yaml.gotmpl`: Some helm charts need an external file during installation, you should put those files in the specifed path and uncomment the respective lines.
- `etc/secrets.yaml`: Passwords and client secrets used by the installation.

### Configure

1. Configure the `etc/production.yaml`. Make sure to read the comments in the file and change the values that are relevant to your installation. You at least want to change the `server_name` and `management_portal.smtp` configuration. Optionally, you can also enable or disable other components that are configured otherwise by default.

   ```shell
   nano etc/production.yaml  # Change setup parameters and configurations
   ```

   When doing a clean install, you are advised to change the `postgresql`, `radar_appserver_postgresql` `radar_upload_postgresql` image tags to the latest PostgreSQL version. Likewise, the timescaledb image tag should use the latest timescaledb version. PostgreSQL passwords and major versions cannot easily be updated after installation.

3. In `etc/production.yaml.gotmpl` file, change setup parameters for charts that are reading input files. You most likely just want to put the file in the default location specified in the file and uncomment the respective lines. Make sure to remove both `#` and `{{/*` from the line in order to uncomment it.

   ```shell
   nano etc/production.yaml.gotmpl
   ```

4. (Optional) If you are installing `radar-appserver`, it needs to be authorized with the Google Firebase also used by the aRMT / Questionnaire app. In Firebase, go to _Project settings_ -> _Service accounts_ and download a Firebase Admin SDK private key. Store the generated key as `etc/radar-appserver/firebase-adminsdk.json` and uncomment the respective section in `etc/production.yaml.gotmpl`.

5. In `etc/secrets.yaml` file, change any passwords, client secrets or API credentials like for Fitbit or Garmin Connect. After the installation you can find login credentials to the components in this file. Be sure to keep it private.
   ```shell
   nano etc/secrets.yaml
   ```


### Install

Once all configuration files are ready, the RADAR-Kubernetes can be deployed on a Kubernetes cluster.

#### Install RADAR-Kubernetes on your cluster.

```shell
helmfile sync --concurrency 1
```

The `helmfile sync` will synchronize all the Kubernetes resources defined in the helmfiles with your Kubernetes cluster. Having `--concurrency 1` will make sure components are installed in required order. Depending on your cluster specification, this may take around 30 minutes when installed for the first time.

| :exclamation: Note |
| :----------------- |

| Installing the stack with `--concurrency 1` may make the installation slower. However, it is necessary because some components such as `kube-prometheus-stack` and `catalog-server`, should be installed in their specified order. If you've forgotten to use this flag, then the installation may not be successful. To continue, follow [Uninstallation](#uninstall) steps to clean up the Kubernetes cluster before you can try again.

Graylog and fluent-bit services in the `graylog` namespace will not immediately be operational, first it needs an input source defined. Log into `graylog.<server name>` with the Graylog credentials. Then navigate to _System_ -> _Inputs_, select GELF TCP in the dropdown and _Launch new input_. Set it as a global input on port 12222.

#### Monitor and verify the installation process.

Once the installation is done or in progress, you can check the status using `kubectl get pods`.

If the installation has been successful, you should see an output similar to the list below. However depending on which components that you've enabled for installation this list and be longer or shorter.

```shell
➜ kubectl get pods
NAME                                             READY   STATUS    RESTARTS   AGE
catalog-server-5c6767cbd8-dc6wc                  1/1     Running   0          8m21s
cp-kafka-0                                       2/2     Running   0          8m21s
cp-kafka-1                                       2/2     Running   0          8m21s
cp-kafka-2                                       2/2     Running   0          8m21s
cp-kafka-rest-5c654995d4-vtzbb                   2/2     Running   0          8m21s
cp-schema-registry-7968b7c554-mznv7              2/2     Running   0          8m21s
cp-schema-registry-7968b7c554-sz6w2              2/2     Running   0          8m21s
cp-zookeeper-0                                   2/2     Running   0          8m21s
cp-zookeeper-1                                   2/2     Running   0          8m21s
cp-zookeeper-2                                   2/2     Running   0          8m21s
management-portal-56cd7f88c6-vmqfk               1/1     Running   0          8m21s
minio-0                                          1/1     Running   0          8m21s
minio-1                                          1/1     Running   0          8m21s
minio-2                                          1/1     Running   0          8m21s
minio-3                                          1/1     Running   0          8m21s
nginx-ingress-controller-748f5b5b88-9j882        1/1     Running   0          8m21s
postgresql-0                                     3/3     Running   0          8m21s
radar-fitbit-connector-594d8b668c-h8m4d          2/2     Running   0          8m21s
radar-gateway-5c4b8c8645-c8zrh                   2/2     Running   0          8m21s
radar-grafana-75b698d68-8gq94                    1/1     Running   0          8m21s
radar-integration-75b76c785c-j8cl9               1/1     Running   0          8m21s
radar-jdbc-connector-677d9dd8c7-txchq            1/1     Running   0          8m21s
radar-output-5d58db5bff-t96vx                    1/1     Running   0          8m21s
radar-rest-sources-authorizer-848cfcbcdf-sdxjg   1/1     Running   0          8m21s
radar-rest-sources-backend-f5895cfd5-d8cdm       1/1     Running   0          8m21s
radar-s3-connector-864b69bd5d-68dvk              1/1     Running   0          8m21s
radar-upload-connect-backend-6d446f885d-29jnc    1/1     Running   0          8m21s
radar-upload-connect-frontend-547cb69595-9ld5s   1/1     Running   0          8m21s
radar-upload-postgresql-postgresql-0             3/3     Running   0          8m21s
radar-upload-source-connector-c87ddc848-fxlg7    1/1     Running   0          8m21s
redis-master-0                                   2/2     Running   0          8m21s
timescaledb-postgresql-0                         3/3     Running   0          8m21s
```

If you have enabled monitoring, logging and HTTPS you should see these as well:

```shell
➜ kubectl -n monitoring get pods
NAME                                                       READY   STATUS    RESTARTS   AGE
kube-prometheus-stack-grafana-674bb6887f-2pgxh             2/2     Running   0          8m29s
kube-prometheus-stack-kube-state-metrics-bbf56d7f5-tm2kg   1/1     Running   0          8m29s
kube-prometheus-stack-operator-7d456878d7-bwrsx            1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-84n2m       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-h5kgc       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-p6mkb       1/1     Running   0          8m29s
prometheus-kube-prometheus-stack-prometheus-0              2/2     Running   1          8m21s
prometheus-kube-prometheus-stack-prometheus-1              2/2     Running   1          8m21s
prometheus-kube-prometheus-stack-prometheus-2              2/2     Running   1          8m21s

➜ kubectl -n graylog get pods
NAME                           READY   STATUS      RESTARTS   AGE
elasticsearch-master-0         1/1     Running     0          8m21s
elasticsearch-master-1         1/1     Running     0          8m21s
elasticsearch-master-2         1/1     Running     0          8m21s
fluentd-6jmhn                  1/1     Running     0          8m21s
fluentd-9lc2g                  1/1     Running     0          8m21s
fluentd-cfzqv                  1/1     Running     0          8m21s
graylog-0                      1/1     Running     0          8m21s
graylog-1                      1/1     Running     0          8m21s
mongodb-mongodb-replicaset-0   2/2     Running     0          8m21s
mongodb-mongodb-replicaset-1   2/2     Running     0          8m21s
mongodb-mongodb-replicaset-2   2/2     Running     0          8m21s

➜ kubectl -n cert-manager get pods
NAME                            READY   STATUS    RESTARTS   AGE
cert-manager-77bbfd565c-rf7wh              1/1     Running   0          8m21s
cert-manager-cainjector-75b6bc7b8b-dv2js   1/1     Running   0          8m21s
cert-manager-webhook-8444c4bc77-jhzgb      1/1     Running   0          8m21s
```

In most cases seeing `1/1` or `2/2` in `READY` column and `Running` in `STATUS` column indicates that the application is running and healthy. Other ways to ensure that installation have been successful is to check application logs for errors and exceptions.

#### Ensure Kafka cluster is functional and RADAR-base topics are loaded

```shell
➜  kubectl exec -it cp-kafka-0 -c cp-kafka-broker -- kafka-topics --bootstrap-server localhost:9092 --list | wc -l
273
```

This output means there are 273 topics loaded in the Kafka cluster.
In your setup, the number of topics can be more or less, depending on components that you have activated.

Other useful Kafka commands can be found by running

```shell
kubectl exec -it cp-kafka-0 -c cp-kafka-broker -- sh -c "ls /usr/bin/kafka*"
```

Use the `--help` flag with each tool to see its purpose.

View the data in a kafka topic by running:

```shell
topic=... # kafka topic to read from
# add any arguments to kafka-avro-console-consumer, e.g. --from-beginning or --max-messages 100
args="--property print.key=true --bootstrap-server cp-kafka-headless:9092"
command="unset JMX_PORT; kafka-avro-console-consumer"
pod=$(kubectl get pods --selector=app=cp-schema-registry -o jsonpath="{.items[0].metadata.name}")
kubectl exec -it $pod -c cp-schema-registry-server -- sh -c "$command --topic $topic $args"
```

## Usage and accessing the applications

In order to access to the applications first you need to find the IP address that Nginx service is listening to and then point the domain that you've specified in `server_name` variable to this IP address via a DNS server (e.g. [Route53](https://aws.amazon.com/route53/), [Cloudflare](https://www.cloudflare.com/dns/), [Bind](https://www.isc.org/bind/)) or [`hosts` file](<https://en.wikipedia.org/wiki/Hosts_(file)>) in your local machine.

> For this guide we assume that you've set `server_name` to "k8s.radar-base.org" and SSL is enabled. Please replace it with a DNS domain under your control.

You can see details of Nginx service with following command:

```shell
➜ kubectl get service nginx-ingress-controller
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP                           PORT(S)                      AGE
nginx-ingress-controller   LoadBalancer   10.100.237.75   XXXX.eu-central-1.elb.amazonaws.com   80:31046/TCP,443:30932/TCP   1h
```

- If you're using a cloud provider you need to point the value in `EXTERNAL-IP` column (in this example `XXXX.eu-central-1.elb.amazonaws.com`) to `k8s.radar-base.org` domain in your DNS server.
- Some of the RADAR-base applications are accesible through sub-domains and you need to configure the DNS server to allow access to those applications. The easy way to do this is to create two wildcard CNAME records:
  ```
  *.k8s.radar-base.org              IN  CNAME  k8s.radar-base.org
  *.*.k8s.radar-base.org            IN  CNAME  k8s.radar-base.org
  ```
- If you're not using a cloud provider you need to use a load balancer to expose `31046` and `30932` ports (will be different in your setup) to a IP address and then point `k8s.radar-base.org` domain to that IP address.
- For development and testing purposes you can run `sudo kubectl port-forward svc/nginx-ingress-controller 80:80 443:443` which will forward Nginx service ports to your local machine and you can have access to applications after adding `127.0.0.1       k8s.radar-base.org` to your `hosts` file.

Now when you go to this IP address you should see a home page with a few links to applications that are installed in the cluster:

```
https://k8s.radar-base.org
```

**Note:** If you have enabled the SSL you might see invalid certificate error when you try to access to the websites, in this case wait a couple of minutes until `cert-manager` issues those certificates.

Now you can head over to the [Management Portal](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/49512484/Management+Portal) guide for next steps.


## Service-specific documentation

- [Data Dashboard Backend data transformation](docs/ksql-server_for_data-dashboard-backend.md)


## Troubleshooting

If an application doesn't become fully ready, installation will not be successful. In this case, you should investigate the root cause by investigating the relevant component. It's suggested to run the following command when `helmfile sync` command is running so you can keep an eye on the installation:
```shell
# on linux
watch kubectl get pods

# on other platforms
kubectl get pods --watch
```
This can help you to findout potential issues faster.

It is suggested to change value of `atomicInstall` to `false` in `etc/production.yaml` file during the installation. This will help troubleshooting potential installation issues easier since it will leave the broken components in place for further inspection, be sure to enable this flag after the installation to prevent broken components causing distruption in case of a faulty update.

Some useful commands for troubleshooting a component are mentioned below.

- Describe a pod to understand current status:

```shell
kubectl describe pods <podname>
```

- Investigate the logs of the pod:

```shell
kubectl logs <podname>
```

- To check last few lines:

```shell
kubectl logs --tail 100 <podname>
```

- To continuously monitor the logs:

```shell
kubectl logs -f <podname>
```

For more information on how `kubectl` can be used to manage a Kubernetes application, please visit [Kubectl documentation](https://kubernetes.io/docs/reference/kubectl/cheatsheet/).
| :exclamation: Note |
|--------------------|
| For most of the components, you can reinstall them without additional actions. However, for some components such as `kube-prometheus-stack` and `kafka-init`, you may need to remove everything before trying again.|

Once you've solved the issue, you need to run the `helmfile sync` command again.

If you have enabled monitoring you should also check **Prometheus** to see if there are any alerts. In next section there is a guide on how to connect to Prometheus.


## Volume expansion

If want to resize a volumes after its initialization you need to make sure that it's supported by its underlying volume plugin:
https://kubernetes.io/docs/concepts/storage/persistent-volumes/#expanding-persistent-volumes-claims

If it's supported then it should be an easy process like this:
https://www.jeffgeerling.com/blog/2019/expanding-k8s-pvs-eks-on-aws

## Uninstall

If you can spin up a new Kubernetes cluster in a few mintues it's generally suggested to recreate the cluster since the installation creates various components that might need to be manually removed. If that's not an option you can run following commands to delete the applications from cluster:

```shell
helmfile destroy
```

Some configurations can still linger inside the cluster. Try using following commands to purge them as well.

```shell
kubectl delete crd prometheuses.monitoring.coreos.com prometheusrules.monitoring.coreos.com servicemonitors.monitoring.coreos.com alertmanagers.monitoring.coreos.com podmonitors.monitoring.coreos.com alertmanagerconfigs.monitoring.coreos.com probes.monitoring.coreos.com thanosrulers.monitoring.coreos.com
kubectl delete psp kube-prometheus-stack-alertmanager kube-prometheus-stack-grafana kube-prometheus-stack-grafana-test kube-prometheus-stack-kube-state-metrics kube-prometheus-stack-operator kube-prometheus-stack-prometheus kube-prometheus-stack-prometheus-node-exporter
kubectl delete mutatingwebhookconfigurations prometheus-admission
kubectl delete ValidatingWebhookConfiguration prometheus-admission

kubectl delete crd certificaterequests.cert-manager.io certificates.cert-manager.io challenges.acme.cert-manager.io clusterissuers.cert-manager.io issuers.cert-manager.io orders.acme.cert-manager.io
kubectl delete pvc --all
kubectl delete pv --all
kubectl -n cert-manager delete secrets letsencrypt-prod
kubectl -n default delete secrets radar-base-tls
kubectl -n monitoring delete secrets radar-base-tls
```

## Update charts

To find any updates to the Helm charts that are listed in the repository, run

```shell
bin/chart-updates
```

## Development automation

This repository can be used for development automation for instance on a k3s or k3d (dockerized k3s) cluster. The example below shows how to deploy on a k3d cluster.

1. Install k3d (see [here](https://github.com/k3d-io/k3d#get))
2. Create a k3d cluster that is configured to run RADAR-base

```shell
k3d cluster create my-test-cluster --port '80:80@loadbalancer' --config=.github/ci_config/k3d-config.yaml
```

This example creates a cluster named `my-test-cluster` with a load balancer that forwards local port 80 to the cluster. The
configuration file `.github/ci_config/k3d-config.yaml` is used to configure the cluster. This cluster will be accessible
in _kubectl_ with context name _k3d-my-test-cluster_.

3. Initialize the RADAR-Kubernetes deployment. Run:

```shell
./bin/init
```

4. In file _etc/production.yaml_:

- set _kubeContext_ to _k3d-my-test-cluster_
- set _dev_deployment_ to _true_
- (optional) enable/disable components as needed with the __install_ fields

5. Install RADAR-Kubernetes on the k3d cluster:

```shell
helmfile sync
```

When installation is complete, you can access the applications at `http://localhost`.



## Feedback and Contributions

Enabling RADAR-base community to use RADAR-Kubernetes is important for us. If you have troubles setting up the platform using provided instructions, you can create an dicussion with exact details to reproduce and the expected behavior.
You can also reach out to the RADAR-base community via RADAR-base Slack on **[radar-kubernetes channel](https://radardevelopment.slack.com/archives/C021AGGESC9)**. The RADAR-base developers support the community on a voluntary basis and will pick up your requests as time permits.
If you'd like to contribute to this project, please checkout [CONTRIBUTING.md](https://github.com/RADAR-base/RADAR-Kubernetes/blob/main/CONTRIBUTING.md) file.
