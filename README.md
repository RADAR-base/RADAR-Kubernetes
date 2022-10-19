# RADAR-Kubernetes [![GitHub release](https://img.shields.io/github/v/release/radar-base/radar-kubernetes)](https://github.com/RADAR-base/RADAR-Kubernetes/releases/latest)

The Kubernetes stack of RADAR-base platform.

## About

RADAR-base is an open-source platform designed to support remote clinical trials by collecting continuous data from wearables and mobile applications. RADAR-Kubernetes enables installing the RADAR-base platform onto Kubernetes clusters. RADAR-base platform can be used for wide range of use-cases. Depending on the use-case, the selection of applications need to be installed can vary. Please read the [component overview and breakdown](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown) to understand the role of each component and how components work together. 

RADAR-Kubernetes setup uses [Helm 3](https://github.com/helm/helm) charts to package necessary Kubernetes resources for each component and [helmfile](https://github.com/roboll/helmfile) to modularize and deploy Helm charts of the platform on a Kubernetes cluster. This setup is designed to be a lightweight way to install and configure the RADAR-base components. The original images or charts may provide more and granular configurations. Please visit the `README` of respective charts in [radar-helm-charts](https://github.com/RADAR-base/radar-helm-charts) to understand the configurations and visit the main repository for in depth knowledge.

## Status
RADAR-Kubernetes is one of the youngest project of RADAR-base and will be the **long term supported form of deploying the platform**. Even though, RADAR-Kubernetes is being used in few production environments, it is still in its early stage of development. We are working on improving the set up and documentation to enable RADAR-base community to make use of the platform.

### Disclaimer
This documentation assumes familiarity with all referenced Kubernetes concepts, utilities, and procedures and familiarity with Helm charts and helmfile. While this documentation will provide guidance for installing and configuring RADAR-base platform on a Kubernetes cluster, it is not a replacement for the official detailed documentation or tutorial of Kubernetes, Helm or Helmfile. If you are not familiar with these tools, we strongly recommend you to get familiar with these tools. Here is a [list of useful links](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2731638785/How+to+get+started+with+tools+around+RADAR-Kubernetes) to get started. 

## Prerequisites
### Hosting Infrastructure
| Component | Description | Required |
|-----|------|------|
| Kubernetes cluster | An infrastructure with working installation of Kubernetes services. Read [this article](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2744942595?draftShareId=e09429e8-38c8-4b71-955d-5df8de94b694) for available options. Minimum requirements for a single node: 8 vCPU's, 32 GB memory, 200 GB storage. Minimum requirements for a cluster: 3 nodes with 3 vCPUs, 16 GB memory, 100 GB storage each and 200 GB shared storage. | Required |
| DNS Server | Some applications are only accessible via HTTPS and it's essential to have a DNS server via providers like GoDaddy, Route53, etc| Required |
| SMTP Server | RADAR-Base needs an SMTP server to send registration email to researchers and participants. | Required |
| Object storage | An external object storage allows RADAR-Kubernetes to backup cluster data such as manifests, application configuration and data via Velero to a backup site. You can also send the RADAR-Base output data to this object storage, which can provider easier management and access compared to bundled Minio server inside RADAR-Kubernetes. | Optional |
| Managed services | RADAR-Kubernetes includes all necessary components to run the platform as a standalone application. However, you can also opt to use managed services such as with the platform, e.g. Confluent cloud for Kafka and schema registry, Postgres DB for storage, Azure blob storage or AWS S3 instead of minio. | Optional |

### Local machine

The following tools should be installed in your local machine to install the RADAR-Kubernetes on your Kubernetes cluster.

| Component | Description |
|-----|------|
| [Git](https://git-scm.com/downloads) | RADAR-Kubernetes uses Git-submodules to use some third party Helm charts. Thus Git is required to properly download and sync correct versions of this repository and its dependent repositories |
| [Java](https://openjdk.java.net/install/)| The installation setup uses Java Keytools to create Keystore files necessary for signing access tokens.|
| [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)| Kubernetes command-line tool, kubectl, allows you to run commands against Kubernetes clusters|
| [helm 3](https://github.com/helm/helm#install)| Helm Charts are used to package Kubernetes resources for each component|
| [helmfile](https://github.com/roboll/helmfile#installation)| RADAR-Kubernetes uses helmfiles to deploy Helm charts.|
| [helm-diff](https://github.com/databus23/helm-diff#install)| A dependency for Helmfile| 

**Once you have a working installation of a Kubernetes cluster, please [configure Kubectl with the appropriate Kubeconfig](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#verify-kubectl-configuration) to enable Kubectl to find and access your cluster. Then proceed to the installation section.** 

## Installation
> The following instructions on this guide are for local machines which runs on Linux operating systems. You can still use the same instructions with small to no changes on a MacOS device as well.

### Prepare
1. Clone the repository to your local machine by using following command.

    ```shell
    git clone https://github.com/RADAR-base/RADAR-Kubernetes.git
    ```
 
2. Create basic config files using template files.

    ```shell
    cd RADAR-Kubernetes
    cp environments.yaml.tmpl environments.yaml
    cp etc/base.yaml etc/production.yaml
    cp etc/base.yaml.gotmpl etc/production.yaml.gotmpl
    ```
   
    It is recommended make a private clone of this repository, if you want to version control your configurations and/or share with other people.
      
    The best practise to share your platform configurations is by **sharing the encrypted version of `production.yaml`**.

### Configure

#### Project Structure
- [/bin](bin): Contains initialization scripts
- [/etc](etc): Contains configurations for some Helm charts.
- [/secrets](secrets): Contains secrets configuration for helm charts.
- [/helmfile.d](helmfile.d): Contains Helmfiles for modular deployment of the platform
- [environments.yaml](environments.yaml): Defines current environment files in order to be used by helmfile. Read more about `bases` [here](https://github.com/roboll/helmfile/blob/master/docs/writing-helmfile.md).
- `etc/production.yaml`: Production helmfile template to configure and install RADAR-base components. Inspect the file to enable, disable and configure components required for your use case. The default helmfile enables all core components that are needed to run RADAR-base platform with pRMT and aRMT apps. If you're not sure which components you want to enable you can refer to wiki for [an overview and breakdown on RADAR-Base components and their roles](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown). 
- `etc/production.yaml.gotmpl`: Change setup parameters that require Go templating, such as reading input files

1. Configure the [environments.yaml](environments.yaml) to use the files that you have created by copying the template files.
    ```shell
    vim environments.yaml # use the files you just created
    ```
2. Configure the `etc/production.yaml`. In this file you are required to fill in secrets and passwords used by RADAR-base applications. It is strongly recommended to use random password generator to fill these secrets. **You must keep this file secure and confidential once you have started installing the platform.**
    
     To create an encrypted password string for monitoring system you need to use [this website](https://www.web2generators.com/apache-tools/htpasswd-generator) or `htpasswd` command to create an encrypted password string and put it inside `kube_prometheus_stack.nginx_auth` variable. It seems like bcrypt encryption isn't supported in current ingress-nginx so make sure that you're using MD5 encryption.
    
    Optionally, you can also enable or disable other components that are configured otherwise by default.
  
    ```shell
    vim etc/production.yaml  # Change setup parameters and configurations
    ```

    When doing a clean install, you are advised to change the `postgresql`, `radar_appserver_postgresql` `radar_upload_postgresql` image tags to the latest PostgreSQL version. Likewise, the timescaledb image tag should use the latest timescaledb version. PostgreSQL passwords and major versions cannot easily be updated after installation.
3. In `etc/production.yaml.gotmpl` file, change setup parameters that require Go templating, such as reading input files and selecting an option for the `keystore.p12`
    ```shell
    vim etc/production.yaml.gotmpl 
    ```
4. Run `bin/keystore-init` to create the Keystore file which used to sign JWT access tokens by [Management Portal](https://github.com/RADAR-base/radar-helm-charts/blob/main/charts/management-portal/README.md)

    ```shell
    bin/keystore-init
    ```

    To prevent the tool from querying variables interactively, please provide a DNAME in the following format, replacing each of the placeholders `<...>` with their proper value:

    ```shell
    DNAME="CN=<name>,O=<organization>,L=<city>,C=<2 letter country code>" bin/keystore-init
    ```

    Consult the full [X.500 Distinguished Name syntax](https://docs.oracle.com/javase/8/docs/technotes/tools/windows/keytool.html#CHDHBFGJ) for more information.

### Install
Once you are done with all configurations, the RADAR-Kubernetes can be deployed on a Kubernetes cluster. 

#### Install RADAR-Kubernetes on your cluster. 

```shell
helmfile sync --concurrency 1
```

The `helmfile sync` will synchronize all the Kubernetes resources defined in the helmfiles with your Kubernetes cluster. Having `--concurrency 1` will make sure components are installed in required order. Depending on your cluster specification, this may take a few minutes when installed for the first time. 

| :exclamation:  Note|
|:----------------------------------------|
| Installing the stack with `--concurrency 1` may make the installation slower. However, it is necessary because some components such as `kube-prometheus-stack` and `kafka-init` (aka `catalog-server`) should be installed in their specified order. If you've forgotten to use this flag, then the installation may not be successful. To continue, follow [Uninstallation](#uninstall) steps to clean up the Kubernetes cluster before you can try again.

Graylog and fluent-bit services in the `graylog` namespace will not immediately be operational, first it needs an input source defined. Log into  `graylog.<server name>` with the Graylog credentials. Then navigate to _System_ -> _Inputs_, select GELF TCP in the dropdown and _Launch new input_. Set it as a global input on port 12222.

#### Monitor and verify the installation process.
Once the installation is done or in progress, you can check the status using `kubectl get pods`.

If the installation has been successful, you should see an output similar to the list below.

```
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
kafka-manager-6858986866-qhphj                   1/1     Running   0          8m21s
management-portal-56cd7f88c6-vmqfk               1/1     Running   0          8m21s
minio-0                                          1/1     Running   0          8m21s
minio-1                                          1/1     Running   0          8m21s
minio-2                                          1/1     Running   0          8m21s
minio-3                                          1/1     Running   0          8m21s
nginx-ingress-controller-748f5b5b88-9j882        1/1     Running   0          8m21s
nginx-ingress-default-backend-659bd647bd-kk922   1/1     Running   0          8m21s
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
```
➜ kubectl -n monitoring get pods                                                            
NAME                                                       READY   STATUS    RESTARTS   AGE
kube-prometheus-stack-grafana-674bb6887f-2pgxh             2/2     Running   0          8m29s
kube-prometheus-stack-kube-state-metrics-bbf56d7f5-tm2kg   1/1     Running   0          8m29s
kube-prometheus-stack-operator-7d456878d7-bwrsx            1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-84n2m       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-h5kgc       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-p6mkb       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-tmsk7       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-vvk6d       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-wp2t7       1/1     Running   0          8m29s
kube-prometheus-stack-prometheus-node-exporter-zsls7       1/1     Running   0          8m29s
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
fluentd-g88cr                  1/1     Running     0          8m21s
fluentd-ks5zx                  1/1     Running     0          8m21s
fluentd-mdg8p                  1/1     Running     0          8m21s
fluentd-qnn8b                  1/1     Running     0          8m21s
fluentd-x4vjd                  1/1     Running     0          8m21s
fluentd-zwzfw                  1/1     Running     0          8m21s
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

Other ways to ensure that installation have been successful is to check application logs for errors and exceptions. 

#### Ensure Kafka cluster is functional and RADAR-base topics are loaded

```bash
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

#### Troubleshoot
If an application doesn't become fully ready installation will not be successful. In this case, you should investigate the root cause by investigating the relevant component. 

Some useful commands for troubleshooting a component are mentioned below.

1. Describe a pod to understand current status
```shell
kubectl describe pods <podname>
```
2. Investigate the logs of the pod
```shell
kubectl logs <podname>
```
To check last few lines
```shell
kubectl logs --tail 100 <podname>
```
To continue monitoring the logs
```shell
kubectl logs -f <podname>
```
For more information on how `kubectl` can be used to manage a Kubernetes application, please visit [Kubectl documentation](https://kubernetes.io/docs/reference/kubectl/cheatsheet/). 
| :exclamation: Note |
|--------------------|
| For most of the components, you can reinstall them without additional actions. However, for some components such as `kube-prometheus-stack` and `kafka-init`, you may need to remove everything before trying again.|
 
If you have enabled monitoring you should also check **Prometheus** to see if there are any alerts. In next section there is a guide on how to connect to Prometheus.

#### Optional
If you are installing `radar-appserver`, it needs to be authorized with the Google Firebase also used by the aRMT / Questionnaire app. In Firebase, go to _Project settings_ -> _Service accounts_ and download a Firebase Admin SDK private key. Store the generated key as `etc/radar-appserver/firebase-adminsdk.json`.

## Upgrade instructions

Run the following instructions to upgrade an existing RADAR-Kubernetes cluster.

| :exclamation: Note |
|--------------------|
| Upgrading the major version of a PostgreSQL image is not supported. If necessary, we propose to use a `pg_dump` to dump the current data and a `pg_restore` to restore that data on a newer version. Please find instructions for this elsewhere. |

### Upgrade to RADAR-Kubernetes version 1.0.0

Before running the upgrade, compare `etc/base.yaml` and `etc/base.yaml.gotmpl` with their `production.yaml` counterparts. Please ensure that all properties in `etc/base.yaml` are overridden in your `production.yaml` or that the `base.yaml` default value is fine, in which case no value needs to be provided in `production.yaml`.

To upgrade the initial services, run
```
kubectl delete -n monitoring deployments kube-prometheus-stack-kube-state-metrics kafka-manager
helm -n graylog uninstall mongodb
kubectl delete -n graylog pvc datadir-mongodb-0 datadir-mongodb-1
```
Note that this will remove your graylog settings but not your actual logs. This step is unfortunately needed to enable credentials on the Graylog database hosted by the mongodb chart. You will need to recreate the GELF TCP input source as during install.

Then run
```
helmfile -f helmfile.d/00-init.yaml --selector name=cert-manager apply
helmfile -f helmfile.d/00-init.yaml apply --concurrency 1
```

To update the Kafka stack, run:
```
helmfile -f helmfile.d/10-base.yaml apply --concurrency 1
```
After this has succeeded, edit your `production.yaml` and change the `cp_kafka.customEnv.KAFKA_INTER_BROKER_PROTOCOL_VERSION` to the corresponding version documented in the [Confluent upgrade instructions](https://docs.confluent.io/platform/current/installation/upgrade.html) of your Kafka installation. Find the currently installed version of Kafka with `kubectl exec cp-kafka-0 -c cp-kafka-broker -- kafka-topics --version`.
When the `cp_kafka.customEnv.KAFKA_INTER_BROKER_PROTOCOL_VERSION` is updated, again run
```
helmfile -f helmfile.d/10-base.yaml apply
```

To upgrade to the latest PostgreSQL helm chart, in `production.yaml`, uncomment the line `postgresql.primary.persistence.existingClaim: "data-postgresql-postgresql-0"` to use the same data storage as previously. Then run
```
kubectl delete secrets postgresql
kubectl delete statefulsets postgresql-postgresql
helmfile -f helmfile.d/10-managementportal.yaml apply
```

If installed, `radar-appserver-postgresql`, uncomment the `production.yaml` line `radar_appserver_postgresql.primary.existingClaim: "data-radar-appserver-postgresql-postgresql-0"`. Then run
```
kubectl delete secrets radar-appserver-postgresql
kubectl delete statefulsets radar-appserver-postgresql-postgresql
helmfile -f helmfile.d/20-appserver.yaml apply
```

If installed, to upgrade `timescaledb`, uncomment the `production.yaml` line `timescaledb.primary.existingClaim: "data-timescaledb-postgresql-0"`. Then run
```
kubectl delete secrets timescaledb-postgresql
kubectl delete statefulsets timescaledb-postgresql
helmfile -f helmfile.d/20-grafana.yaml apply
```

If installed, to upgrade `radar-upload-postgresql`, uncomment the `production.yaml` line `radar_upload_postgresql.primary.existingClaim: "data-radar-upload-postgresql-postgresql-0"`. Then run
```
kubectl delete secrets radar-upload-postgresql
kubectl delete statefulsets radar-upload-postgresql-postgresql
helmfile -f helmfile.d/20-upload.yaml apply
```

Delete the redis stateful set (this will not delete the data on the volume) 
```
kubectl delete statefulset redis-master
helmfile -f helmfile.d/20-s3.yaml sync --concurrency 1 
```


## Usage
### Accessing the applications
In order to access to the applications first you need to find the IP address that Nginx service is listening to and then point the domain that you've specified in `server_name` variable to this IP address via a DNS server (e.g. [Route53](https://aws.amazon.com/route53/), [Cloudflare](https://www.cloudflare.com/dns/), [Bind](https://www.isc.org/bind/)) or [`hosts` file](https://en.wikipedia.org/wiki/Hosts_(file)) in your local machine.
> For this guide we assume that you've set `server_name` to "k8s.radar-base.org" and SSL is enabled. Please replace it with a DNS domain under your control.

You can see details of Nginx service with following command:
```
➜ kubectl get service nginx-ingress-controller
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP                           PORT(S)                      AGE
nginx-ingress-controller   LoadBalancer   10.100.237.75   XXXX.eu-central-1.elb.amazonaws.com   80:31046/TCP,443:30932/TCP   1h
```
* If you're using a cloud provider you need to point the value in `EXTERNAL-IP` column (in this example `XXXX.eu-central-1.elb.amazonaws.com`) to `k8s.radar-base.org` domain in your DNS server.
* RADAR-base uses a few subdomains, which need to be present in DNS for the certificate manager to get an SSL certificate. The easy way to do this is to create a wildcard CNAME record:
     ```
    *.k8s.radar-base.org            IN  CNAME  k8s.radar-base.org
    ```
    Alternatively, create each DNS entry manually:
    ```
    dashboard.k8s.radar-base.org    IN  CNAME  k8s.radar-base.org
    graylog.k8s.radar-base.org      IN  CNAME  k8s.radar-base.org
    alertmanager.k8s.radar-base.org IN  CNAME  k8s.radar-base.org
    s3.k8s.radar-base.org           IN  CNAME  k8s.radar-base.org
    prometheus.k8s.radar-base.org   IN  CNAME  k8s.radar-base.org
    grafana.k8s.radar-base.org      IN  CNAME  k8s.radar-base.org
    ```
    With the latter method, when a new subdomain is needed for a new service, also a new CNAME entry needs to be created.
* If you're not using a cloud provider you need to use a load balancer to expose `31046` and `30932` ports (will be different in your setup) to a IP address and then point `k8s.radar-base.org` domain to that IP address.
* For development and testing purposes you can run `sudo kubectl port-forward svc/nginx-ingress-controller 80:80 443:443` which will forward Nginx service ports to your local machine and you can have access to applications after adding `127.0.0.1       k8s.radar-base.org` to your `hosts` file.

Now depending on your setup you should have access to following URLs:
```
https://k8s.radar-base.org/managementportal
https://s3.k8s.radar-base.org
https://k8s.radar-base.org/upload
https://k8s.radar-base.org/rest-sources/authorizer
https://k8s.radar-base.org/kafkamanager/
https://graylog.k8s.radar-base.org   # Log management
https://prometheus.k8s.radar-base.org   # Monitoring stack
https://alertmanager.k8s.radar-base.org
https://grafana.k8s.radar-base.org
```

**Note:** If you have enabled the SSL you might see invalid certificate error when you try to access to the websites, in this case wait a couple of minutes until `cert-manager` issues those certificates.

## Volume expansion

If want to resize a volumes after its initialization you need to make sure that it's supported by its underlying volume plugin:
https://kubernetes.io/docs/concepts/storage/persistent-volumes/#expanding-persistent-volumes-claims

If it's supported then it should be an easy process like this:
https://www.jeffgeerling.com/blog/2019/expanding-k8s-pvs-eks-on-aws

## Uninstall
If you want to remove the RADAR-base from your cluster you and use following command to delete the applications from cluster:
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
kubectl -n cert-manager delete secrets letsencrypt-prod
kubectl -n default delete secrets radar-base-tls
kubectl -n monitoring delete secrets radar-base-tls
```
## Update charts

To find any updates to the Helm charts that are listed in the repository, run
```
bin/chart-updates
```

## Feedback and Contributions
Enabling RADAR-base community to use RADAR-Kubernetes is important for us. If you have troubles setting up the platform using provided instructions, you can create an issue with exact details to reproduce and the expected behaviour.
You can also reach out to the RADAR-base community via RADAR-base Slack on **[radar-kubernetes channel](https://radardevelopment.slack.com/archives/C021AGGESC9)**. The RADAR-base developers support the community on a voluntary basis and will pick up your requests as time permits.
