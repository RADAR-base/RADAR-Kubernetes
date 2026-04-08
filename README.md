# RADAR-Kubernetes

[![GitHub release](https://img.shields.io/github/v/release/radar-base/radar-kubernetes)](https://github.com/RADAR-base/RADAR-Kubernetes/releases/latest)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/radar-base/radar-kubernetes/push.yaml)
[![Apache Licensed](https://img.shields.io/github/license/radar-base/radar-kubernetes)](LICENSE)
![GitHub Discussions](https://img.shields.io/github/discussions/radar-base/radar-kubernetes) ![Maintenance](https://img.shields.io/maintenance/yes/2023)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/radar-base/radar-kubernetes/dev)
[![Join our community Slack](https://img.shields.io/badge/slack-radarbase-success.svg?logo=slack)](https://docs.google.com/forms/d/e/1FAIpQLScKNZ-QonmxNkekDMLLbP-b_IrNHyDRuQValBy1BAsLOjEFpg/viewform)

The Kubernetes stack of RADAR-base platform.

## Table of contents

<!-- TOC -->
* [RADAR-Kubernetes](#radar-kubernetes)
  * [Table of contents](#table-of-contents)
  * [About](#about)
  * [Status](#status)
  * [Prerequisites](#prerequisites)
    * [Knowledge requirements](#knowledge-requirements)
    * [Software Compatibility](#software-compatibility)
    * [Hosting](#hosting)
    * [Third party services](#third-party-services)
    * [Local machine](#local-machine)
  * [Installation](#installation)
    * [Prepare](#prepare)
    * [Project Structure](#project-structure)
    * [Configure](#configure)
    * [Install](#install)
      * [Install RADAR-Kubernetes on your cluster.](#install-radar-kubernetes-on-your-cluster)
      * [Monitor and verify the installation process.](#monitor-and-verify-the-installation-process)
        * [Note](#note)
      * [Ensure Kafka cluster is functional and RADAR-base topics are loaded](#ensure-kafka-cluster-is-functional-and-radar-base-topics-are-loaded)
  * [Usage and accessing the applications](#usage-and-accessing-the-applications)
  * [Service-specific documentation](#service-specific-documentation)
  * [Troubleshooting](#troubleshooting)
  * [Volume expansion](#volume-expansion)
  * [Uninstall](#uninstall)
  * [Update charts](#update-charts)
  * [Developer documentation](#developer-documentation)
  * [Feedback and Contributions](#feedback-and-contributions)
<!-- TOC -->

## About

RADAR-base is an open-source platform designed to support remote clinical trials by collecting continuous data from
wearables and mobile applications. RADAR-Kubernetes enables installing the RADAR-base platform onto Kubernetes clusters.
RADAR-base platform can be used for wide range of use-cases. Depending on the use-case, the selection of applications
need to be installed can vary. Please read
the [component overview and breakdown](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown)
to understand the role of each component and how components work together.

RADAR-Kubernetes setup uses [Helm](https://github.com/helm/helm) charts to package necessary Kubernetes resources for
each component and [helmfile](https://github.com/roboll/helmfile) to modularize and deploy Helm charts of the platform
on a Kubernetes cluster. This setup is designed to be a lightweight way to install and configure the RADAR-base
components. The original images or charts may provide more and granular configurations. Please visit the `README` of
respective charts in [radar-helm-charts](https://github.com/RADAR-base/radar-helm-charts) to understand the
configurations and visit the main repository for in depth knowledge.

## Status

RADAR-Kubernetes is one of the youngest project of RADAR-base and will be the **long term supported form of deploying
the platform**. Even though, RADAR-Kubernetes is being used in few production environments, it is still in its early
stage of development. We are working on improving the setup and documentation to enable RADAR-base community to make
use of the platform.

## Prerequisites

### Knowledge requirements

This documentation assumes familiarity with all referenced Kubernetes concepts, utilities, and procedures and
familiarity with Helm charts and helmfile, depending on your environment you might need to have knowledge of other
hosting infrastructure such as DNS and mail servers as well. While this documentation will provide guidance for
installing and configuring RADAR-base platform on a Kubernetes cluster and tries to make is as simple and possible, it
is not a replacement for the detailed knowledge of the tools that have been used. If you are not familiar with these
tools, we strongly recommend you to get familiar with these tools. Here is
a [list of useful links](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2731638785/How+to+get+started+with+tools+around+RADAR-Kubernetes)
to get started.

### Software Compatibility

Currently RADAR-Kubernetes is tested and supported on following component versions:
| Component | Version |
| ---- | ------- |
| Kubernetes | v1.30, v1.31, v1.32 and v1.33 |
| K3s | v1.30.6+k3s1, v1.31.10+k3s1, v1.32.6+k3s1 and v1.33.2+k3s1 |
| Kubectl | v1.30, v1.31, v1.32 and v1.33 |
| Helm | v3.16.3 |
| Helm diff | 3.9.12 |
| Helmfile | v0.169.1 |
| YQ | v4.44.3 |

It's possible to install RADAR-Kubernetes on different version of tools as well, but you might encounter compatibility
issues. Make sure that `kubectl` version matches or exceeds the Kubernetes (or K3s) version that you're using. For
other tools, such as Git and Java, as long as the versions are not very old, it's not very impactful.

### Hosting

Kubernetes can be installed on wide variety of platforms and in turn you can install RADAR-Base on most places that
Kubernetes can run. However, your infrastructure needs to meet the requirements listed below:

| Component                              | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | Required |
|----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| Kubernetes cluster                     | An infrastructure with working installation of Kubernetes services. Read [this article](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2744942595?draftShareId=e09429e8-38c8-4b71-955d-5df8de94b694) for available options. Minimum requirements for a single node: 8 vCPU's, 32 GB memory, 200 GB storage. Minimum requirements for a cluster: 3 nodes with 3 vCPUs, 16 GB memory, 100 GB storage each and 200 GB shared storage.                                                                                                                   | Required |
| DNS Server                             | Some applications are only accessible via HTTPS and it's essential to have a DNS server via providers like GoDaddy, Route53, etc                                                                                                                                                                                                                                                                                                                                                                                                                             | Required |
| SMTP Server                            | RADAR-Base needs an SMTP server to send registration email to researchers and participants.                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | Required |
| Whitelisted access to ports 80 and 443 | We use Let's Encrypt to create SSL certificates and in the default configuration we use HTTP challenge. This means that the RADAR-Base installation needs to be visible to Let's Encrypt servers for the verification, so make sure these ports are white listed in your firewall. If you want to have a private installation you should change Let's Encrypt configuration to use DNS challenge.                                                                                                                                                            | Required |
| Object storage                         | An external object storage allows RADAR-Kubernetes to backup cluster data such as manifests, application configuration and data via Velero to a backup site. You can also send the RADAR-Base output data to this object storage, which can provider easier management and access compared to bundled Minio server inside RADAR-Kubernetes.                                                                                                                                                                                                                  | Optional |
| Managed services                       | RADAR-Kubernetes includes all necessary components to run the platform as a standalone application. However, you can also opt to use managed services such as with the platform, e.g. Confluent cloud for Kafka and schema registry, Postgres DB for storage, Azure blob storage or AWS S3 instead of Minio. If you're using a managed object storage that you have to pay per request (such as AWS S3), it's recommended that to install the Minio just for the `radar-intermediate-storage` since the applications send a lot of API calls to that bucket. | Optional |

If you want to deploy a production-ready cluster on AWS you can
use [RADAR-K8s-Infrastructure](https://github.com/RADAR-base/RADAR-K8s-Infrastructure) which hosts the Terraform scripts
required to create an AWS EKS cluster and other managed services required for RADAR-Base platform.

Alternatively, in order to have a simple single node Kubernetes server you can run these commands on a Linux machine:

```shell
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.33.2+k3s1" K3S_KUBECONFIG_MODE="644" INSTALL_K3S_SYMLINK="skip" sh -s - --disable traefik --disable-helm-controller
```

### Third party services

Depending on which components you've enabled you might need credentials for Fitbit, REDCap, Google Firebase, etc. You
need to provide them in order for the respective component to work correctly.

### Local machine

The following tools should be installed in your local machine to install the RADAR-Kubernetes on your Kubernetes
cluster.

| Component                                                          | Description                                                                                                                                                                                                      |
|--------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Git](https://git-scm.com/downloads)                               | RADAR-Kubernetes uses Git-submodules to use some third party Helm charts. Thus Git is required to properly download and sync correct versions of this repository and its dependent repositories                  |
| [Java](https://openjdk.java.net/install/)                          | The installation setup uses Java Keytools to create Keystore files + for signing access tokens.                                                                                                                  |
| [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) | Kubernetes command-line tool, kubectl, allows you to run commands against Kubernetes clusters                                                                                                                    |
| [helm 3](https://github.com/helm/helm#install)                     | Helm Charts are used to package Kubernetes resources for each component                                                                                                                                          |
| [helmfile](https://github.com/helmfile/helmfile#installation)      | RADAR-Kubernetes uses helmfiles to deploy Helm charts.                                                                                                                                                           |
| [helm-diff](https://github.com/databus23/helm-diff#install)        | A dependency for Helmfile.                                                                                                                                                                                       |
| [yq](https://github.com/mikefarah/yq#install)                      | Used to run `init`, `generate-secrets` and `chart-updates` scripts.                                                                                                                                              |
| openssl                                                            | Used in `init` and `generate-secrets` scripts to generate secret for Prometheus Nginx authentication. This binary is in `openssl` package for Ubuntu, it's also easily available on other distributions as well. |

**Once you have a working installation of a Kubernetes cluster,
please [configure Kubectl with the appropriate Kubeconfig](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#verify-kubectl-configuration)
to enable Kubectl to find and access your cluster. Then proceed to the installation section.**

## Installation

> The following instructions on this guide are for local machines which runs on Linux operating systems. You can still
> use the same instructions with small to no changes on a MacOS device as well.

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

It is recommended make a private clone of this repository, if you want to version control your configurations and/or
share with other people.

**You must keep `etc/secrets.yaml` secure and confidential once you have started installing the platform** and the best
practice to share your platform configurations is by **sharing the encrypted version of `etc/secrets.yaml`, this can be
done via the combination of [sops](https://github.com/getsops/sops)
and [helm-secrets](https://github.com/jkroepke/helm-secrets) but it's outside the scope of this document**.

### Project Structure

- `bin/`: Contains initialization scripts.
- `etc/`: Contains configurations for the Helm charts.
- `helmfile.d/`: Contains Helmfiles for modular deployment of the platform.
- `environments.yaml`: Defines current environment files in order to be used by helmfile and where to find the
  configuration files.
- `etc/production.yaml`: Production helmfile template to configure and install RADAR-base components. Inspect the file
  to enable, disable and configure components required for your use case. The default helmfile enables all core
  components that are needed to run RADAR-base platform with pRMT and aRMT apps. If you're not sure which components you
  want to enable you can refer to wiki
  for [an overview and breakdown on RADAR-Base components and their roles](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown).
- `etc/production.yaml.gotmpl`: Some helm charts need an external file during installation, you should put those files
  in the specified path and uncomment the respective lines.
- `etc/secrets.yaml`: Passwords and client secrets used by the installation.
- `mods/`: contains helmfile value files that configure groups of services for specific purposes. Mods are applied in
  `environments.yaml` file, often via setting _master_ options in `etc/production.yaml`.

To deploy RADAR-base platform, you need the following:

- A Kubernetes cluster
- Helm 3.16.3+ ([setup instructions](#installing-helm))
- Helmfile 0.169.1 ([setup instructions](#installing-helmfile))
- `kubectl` access to your cluster
- `yq` (for generating secrets)
- Optionally, a domain name with access to its DNS settings

### Installing Helm

#### Option 1: Using Homebrew (macOS)
```bash
brew install helm@3.16.3
```

#### Option 2: Manual installation

##### macOS
```bash
# Download Helm
curl -Lo helm.tar.gz https://get.helm.sh/helm-v3.16.3-darwin-amd64.tar.gz

# Extract the archive
tar -zxvf helm.tar.gz

# Move the helm binary to a directory in your PATH
sudo mv darwin-amd64/helm /usr/local/bin/

# Verify installation
helm version
# Should output: version.BuildInfo{Version:"v3.16.3", ...}

# Clean up
rm -rf helm.tar.gz darwin-amd64
```

##### Linux
```bash
# Download Helm
curl -Lo helm.tar.gz https://get.helm.sh/helm-v3.16.3-linux-amd64.tar.gz

# Extract the archive
tar -zxvf helm.tar.gz

# Move the helm binary to a directory in your PATH
sudo mv linux-amd64/helm /usr/local/bin/

# Verify installation
helm version
# Should output: version.BuildInfo{Version:"v3.16.3", ...}

# Clean up
rm -rf helm.tar.gz linux-amd64
```

### Installing Helmfile

⚠️ **IMPORTANT**: Helmfile version 0.169.1 is required. Newer versions like v1.0.0 may cause template processing issues with RADAR-Kubernetes.

#### Option 1: Using Homebrew (macOS)
```bash
# Remove existing version if present
brew uninstall helmfile

# Tap the helmfile repository
brew tap helmfile/tap

# Install the specific version
brew install helmfile/tap/helmfile@0.169.1
```

#### Option 2: Manual installation

##### macOS
```bash
# Download Helmfile
curl -Lo helmfile.tar.gz https://github.com/helmfile/helmfile/releases/download/v0.169.1/helmfile_0.169.1_darwin_amd64.tar.gz

# Extract the archive
tar -zxvf helmfile.tar.gz

# Make it executable and move it to your PATH
chmod +x helmfile
sudo mv helmfile /usr/local/bin/

# Clean up
rm -f helmfile.tar.gz LICENSE README.md

# Verify installation
helmfile --version
# Should output: helmfile version v0.169.1
```

##### Linux
```bash
# Download Helmfile
curl -Lo helmfile.tar.gz https://github.com/helmfile/helmfile/releases/download/v0.169.1/helmfile_0.169.1_linux_amd64.tar.gz

# Extract the archive
tar -zxvf helmfile.tar.gz

# Make it executable and move it to your PATH
chmod +x helmfile
sudo mv helmfile /usr/local/bin/

# Clean up
rm -f helmfile.tar.gz LICENSE README.md

# Verify installation
helmfile --version
# Should output: helmfile version v0.169.1
```

### Installing Required Helm Plugins

Helmfile requires the helm-diff plugin to function properly:

```bash
helm plugin install https://github.com/databus23/helm-diff --version 3.9.12
```

### Configure

1. Configure the `etc/production.yaml`. Make sure to read the comments in the file and change the values that are
   relevant to your installation. You at least want to change the `server_name` and `management_portal.smtp`
   configuration. Optionally, you can also enable or disable other components that are configured otherwise by default.

   ```shell
   nano etc/production.yaml  # Change setup parameters and configurations
   ```

2. In `etc/production.yaml.gotmpl` file, change setup parameters for charts that are reading input files. You most
   likely just want to put the file in the default location specified in the file and uncomment the respective lines.
   Make sure to remove both `#` and `{{/*` from the line in order to uncomment it.

   ```shell
   nano etc/production.yaml.gotmpl
   ```

3. (Optional) If you are installing `radar-appserver`, it needs to be authorized with the Google Firebase also used by
   the aRMT / Questionnaire app. In Firebase, go to _Project settings_ -> _Service accounts_ and download a Firebase
   Admin SDK private key. Store the generated key as `etc/radar-appserver/firebase-adminsdk.json` and uncomment the
   respective section in `etc/production.yaml.gotmpl`.

4. In `etc/secrets.yaml` file, add any passwords, client secrets or API credentials that are determined by services
   outside RADAR-base (for instance Fitbit or Garmin Connect). After the installation you can find login credentials to
   the components in this file. Be sure to keep it private.

   ```shell
   nano etc/secrets.yaml
   ```

### Install

Once all configuration files are ready, the RADAR-Kubernetes can be deployed on a Kubernetes cluster.

#### Install RADAR-Kubernetes on your cluster.

```shell
helmfile sync
```

The `helmfile sync` will synchronize all the Kubernetes resources defined in the helmfiles with your Kubernetes cluster.
Depending on your cluster specification, this may take around 30 minutes when installed for the first time. Note that
during first installation, a large portion of the time is spent downloading docker images. This may result in timeout errors
that
may be ignored by repeatedly running the `helmfile sync` command.

| :exclamation: Note |
|:-------------------|

By default, RADAR-Kubernetes deploys in _atomic_ mode where any error results in the automatic rollback of all changes.
For analysis of potential problem it is advised to disable this feature by setting the `atomicInstall: false` in
`etc/production.yaml`.

Graylog and fluent-bit services in the `graylog` namespace will not immediately be operational, first it needs an input
source defined. Log into `graylog.<server name>` with the Graylog credentials. Then navigate to _System_ -> _Inputs_,
select GELF TCP in the dropdown and _Launch new input_. Set it as a global input on port 12222.

#### Monitor and verify the installation process.

Once the installation is done or in progress, you can check the status using `kubectl get pods --all-namespaces`.

If the installation has been successful, you should see an output similar to the list below when using the services
installed by default. However, depending on which components that you've enabled for installation this list and be
longer or shorter.

```shell
➜ kubectl get pods --all-namespaces
NAMESPACE      NAME                                                   READY   STATUS      RESTARTS   AGE
cert-manager   cert-manager-77c5f7bf75-f4kjp                          1/1     Running     0          15m
cert-manager   cert-manager-cainjector-669d85f6cf-hqz4b               1/1     Running     0          15m
cert-manager   cert-manager-webhook-585b8b6bfc-2lz42                  1/1     Running     0          15m
default        catalog-server-85896fdc49-d96kj                        1/1     Running     2          10m
default        cloudnativepg-operator-cloudnative-pg-6c95784c56-ggqbf 1/1     Running     0          15m
default        cp-kafka-0                                             1/1     Running     0          13m
default        cp-kafka-1                                             1/1     Running     0          13m
default        cp-kafka-2                                             1/1     Running     0          12m
default        cp-schema-registry-5bf975844f-v2jcn                    1/1     Running     0          11m
default        cp-zookeeper-0                                         1/1     Running     0          15m
default        cp-zookeeper-1                                         1/1     Running     0          15m
default        cp-zookeeper-2                                         1/1     Running     0          14m
default        ingress-nginx-controller-699dc9cf5c-ld65d              1/1     Running     0          15m
default        management-portal-57d8d8978b-slxqz                     1/1     Running     0          9m42s
default        minio-0                                                1/1     Running     0          15m
default        minio-1                                                1/1     Running     0          15m
default        minio-2                                                1/1     Running     0          15m
default        minio-3                                                1/1     Running     0          15m
default        minio-provisioning-gn8db                               0/1     Completed   0          15m
default        radar-cloudnative-postgresql-cluster-1                 1/1     Running     0          13m
default        radar-cloudnative-postgresql-cluster-2                 1/1     Running     0          13m
default        radar-gateway-849bd7f76-bklk7                          1/1     Running     0          10m
default        radar-home-8655d45964-w65bt                            1/1     Running     0          15m
default        radar-output-868b578858-rnpz8                          1/1     Running     0          7m5s
default        radar-redis-replication-0                              1/1     Running     0          15m
default        radar-redis-replication-1                              1/1     Running     0          15m
default        radar-redis-replication-2                              1/1     Running     0          15m
default        radar-s3-connector-6d8898775c-cf77v                    1/1     Running     1          7m39s
default        redis-operator-569485b5d5-6vb9k                        1/1     Running     0          15m
kube-system    coredns-ccb96694c-jdpsl                                1/1     Running     0          17m
kube-system    local-path-provisioner-5cf85fd84d-wzd2j                1/1     Running     0          17m
kube-system    metrics-server-5985cbc9d7-hmm92                        1/1     Running     0          17m
kube-system    svclb-ingress-nginx-controller-8d7925e1-kq7qv          2/2     Running     0          15m
```

If you have enabled monitoring+logging (`enable_logging_monitoring: true` in `etc/production.yaml`) and HTTPS (
`disable_tls: false` in `etc/production.yaml`) you should see an extended list of pods in the `monitoring` and `graylog`
namespaces (as well as increased companion containers in pods for each service):

```shell
➜ kubectl -n monitoring get pods
NAMESPACE      NAME                                                   READY   STATUS      RESTARTS   AGE
monitoring     alertmanager-kube-prometheus-stack-alertmanager-0      2/2     Running     0          8m12s
monitoring     kube-prometheus-stack-grafana-64c8964dd4-66l7w         3/3     Running     0          8m14s
monitoring     kube-prometheus-stack-kube-state-metrics-6d57649657-46fvq 1/1  Running     0          8m14s
monitoring     kube-prometheus-stack-operator-5554dcc847-78h6s        1/1     Running     0          8m14s
monitoring     kube-prometheus-stack-prometheus-node-exporter-98pqw   1/1     Running     0          8m14s
monitoring     prometheus-kube-prometheus-stack-prometheus-0          2/2     Running     0          8m12s

➜ kubectl -n graylog get pods
NAMESPACE      NAME                                                   READY   STATUS      RESTARTS   AGE
graylog        elasticsearch-master-0                                 1/1     Running     0          7m50s
graylog        elasticsearch-master-1                                 1/1     Running     0          7m50s
graylog        elasticsearch-master-2                                 1/1     Running     0          7m50s
graylog        fluent-bit-zlnst                                       1/1     Running     4          2m20s
graylog        graylog-0                                              1/1     Running     0          6m9s
graylog        mongodb-0                                              2/2     Running     0          7m47s
graylog        mongodb-1                                              2/2     Running     0          7m22s
graylog        mongodb-arbiter-0                                      1/1     Running     1          7m47s
```

In most cases seeing `1/1` or `2/2` in `READY` column and `Running` in `STATUS` column indicates that the application is
running and healthy. Other ways to ensure that installation have been successful is to check application logs for errors
and exceptions.

##### Note

The first time the cluster is installed, a considerable time is consumed by the download of docker images. This may
result in timeout errors that may be ignored by repeatedly running the `helmfile sync` command.

#### Ensure Kafka cluster is functional and RADAR-base topics are loaded

```shell
➜  `kubectl exec -it cp-kafka-0 -c cp-kafka-broker -- kafka-topics --bootstrap-server localhost:9092 --list | wc -l`
425
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

In order to access to the applications first you need to find the IP address that Nginx service is listening to and then
point the domain that you've specified in `server_name` variable to this IP address via a DNS server (
e.g. [Route53](https://aws.amazon.com/route53/), [Cloudflare](https://www.cloudflare.com/dns/), [Bind](https://www.isc.org/bind/))
or [`hosts` file](<https://en.wikipedia.org/wiki/Hosts_(file)>) in your local machine.

> For this guide we assume that you've set `server_name` to "k8s.radar-base.org" and SSL is enabled. Please replace it
> with a DNS domain under your control.

You can see details of Nginx service with following command:

```shell
➜ kubectl get service nginx-ingress-controller
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP                           PORT(S)                      AGE
nginx-ingress-controller   LoadBalancer   10.100.237.75   XXXX.eu-central-1.elb.amazonaws.com   80:31046/TCP,443:30932/TCP   1h
```

- If you're using a cloud provider you need to point the value in `EXTERNAL-IP` column (in this example
  `XXXX.eu-central-1.elb.amazonaws.com`) to `k8s.radar-base.org` domain in your DNS server.
- Some of the RADAR-base applications are accessible through sub-domains and you need to configure the DNS server to
  allow access to those applications. The easy way to do this is to create two wildcard CNAME records:
  ```
  *.k8s.radar-base.org              IN  CNAME  k8s.radar-base.org
  *.*.k8s.radar-base.org            IN  CNAME  k8s.radar-base.org
  ```
- If you're not using a cloud provider you need to use a load balancer to expose `31046` and `30932` ports (will be
  different in your setup) to a IP address and then point `k8s.radar-base.org` domain to that IP address.
- For development and testing purposes you can run
  `sudo kubectl port-forward svc/nginx-ingress-controller 80:80 443:443` which will forward Nginx service ports to your
  local machine and you can have access to applications after adding `127.0.0.1       k8s.radar-base.org` to your
  `hosts` file.

Now when you go to this IP address you should see a home page with a few links to applications that are installed in the
cluster:

```
https://k8s.radar-base.org
```

**Note:** If you have enabled the SSL you might see invalid certificate error when you try to access to the websites, in
this case wait a couple of minutes until `cert-manager` issues those certificates.

Now you can head over to
the [Management Portal](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/49512484/Management+Portal) guide for
next steps.

## Service-specific documentation

- [Configuration of external PostgreSQL and TimescaleDB databases](docs/external_postgresql_databases.md)
- [Data Dashboard Backend data transformation](docs/ksql-server_for_data-dashboard-backend.md)

## Troubleshooting

If an application doesn't become fully ready, installation will not be successful. In this case, you should investigate
the root cause by investigating the relevant component. It's suggested to run the following command when `helmfile sync`
command is running so you can keep an eye on the installation:

```shell
# on linux
watch kubectl get pods

# on other platforms
kubectl get pods --watch
```

This can help you identify potential issues faster.

It is suggested to change value of `atomicInstall` to `false` in `etc/production.yaml` file during the installation.
This will help troubleshooting potential installation issues easier since it will leave the broken components in place
for further inspection, be sure to enable this flag after the installation to prevent broken components causing
disruption in case of a faulty update.

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

For more information on how `kubectl` can be used to manage a Kubernetes application, please
visit [Kubectl documentation](https://kubernetes.io/docs/reference/kubectl/cheatsheet/).
| :exclamation: Note |
|--------------------|
| For most of the components, you can reinstall them without additional actions. However, for some components such as
`kube-prometheus-stack` and `kafka-init`, you may need to remove everything before trying again.|

Once you've solved the issue, you need to run the `helmfile sync` command again.

If you have enabled monitoring you should also check **Prometheus** to see if there are any alerts. In next section
there is a guide on how to connect to Prometheus.

## Volume expansion

If want to resize a volumes after its initialization you need to make sure that it's supported by its underlying volume
plugin:
https://kubernetes.io/docs/concepts/storage/persistent-volumes/#expanding-persistent-volumes-claims

If it's supported then it should be an easy process like this:
https://www.jeffgeerling.com/blog/2019/expanding-k8s-pvs-eks-on-aws

## Uninstall

If you can spin up a new Kubernetes cluster in a few mintues it's generally suggested to recreate the cluster since the
installation creates various components that might need to be manually removed. If that's not an option you can run
following commands to delete the applications from cluster:

```shell
kubectl get redis -oname | xargs kubectl delete
kubectl get redisreplications -oname | xargs kubectl delete
kubectl get redisclusters -oname |  | xargs kubectl delete
helmfile destroy
```

Some configurations can still linger inside the cluster. Try using following commands to purge them as well.

```shell
kubectl get sts -oname | xargs kubectl delete
kubectl get jobs -oname | xargs kubectl delete
PATTERN="monitoring.coreos.com|redis.opstreelabs.in|postgresql.cnpg.io|cert-manager.io"
kubectl get crds -oname | grep -E $PATTERN | xargs kubectl delete
kubectl delete pvc --all
kubectl delete pv --all
kubectl -n cert-manager delete secrets --all
kubectl -n default delete secrets --all
kubectl -n monitoring delete secrets --all
```

## Update charts

To find any updates to the Helm charts that are listed in the repository, run

```shell
bin/chart-updates
```

## Developer documentation

See [Development guide](docs/development_guide.md)

## Feedback and Contributions

Enabling RADAR-base community to use RADAR-Kubernetes is important for us. If you have troubles setting up the platform
using provided instructions, you can create an dicussion with exact details to reproduce and the expected behavior.
You can also reach out to the RADAR-base community via RADAR-base Slack on *
*[radar-kubernetes channel](https://radardevelopment.slack.com/archives/C021AGGESC9)**. The RADAR-base developers
support the community on a voluntary basis and will pick up your requests as time permits.
If you'd like to contribute to this project, please
checkout [CONTRIBUTING.md](https://github.com/RADAR-base/RADAR-Kubernetes/blob/main/CONTRIBUTING.md) file.
