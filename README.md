# RADAR-Kubernetes [![Build Status](https://travis-ci.org/RADAR-base/RADAR-Kubernetes.svg?branch=dev)](https://travis-ci.org/RADAR-base/RADAR-Kubernetes)
The Kubernetes stack of RADAR-base platform

## About
RADAR-base is an open-source platform designed to support remote monitoring of patients by collecting continuous data from wearables and mobile applications. RADAR-Kubernetes enables installing the RADAR-base platform onto Kubernetes clusters. RADAR-base platform can be used for wide range of use-cases. Depending on the use-case, the selection of applications need to be installed can vary. Please read the [component overview and breakdown](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown) to understand the role of each component and how components work together. 

RADAR-Kubernetes setup uses [Helm](https://github.com/helm/helm) charts to package necessary Kubernetes resources for each component and [helmfile](https://github.com/roboll/helmfile) to modularize and deploy Helm charts of the platform on a Kubernetes cluster. This setup is designed to be a lightweight way to install and configure the RADAR-base components. The original images or charts may provide more and granular configurations. Please visit the `README` of respective charts to understand the configurations and visit the main repository for in depth knowledge.

## Status
RADAR-Kubernetes is one of the youngest project of RADAR-base and will be the long term supported form of deploying the platform. Even though, RADAR-Kubernetes is being used in few production environments, it is still in its early stage of development, and we are working improving the set up and documentation to enable RADAR-base community to make use of the platform.

### Disclaimer
This documentation assumes familiarity with all referenced Kubernetes concepts, utilities, and procedures and familiarity with Helm charts and helmfile. While this documentation will provide guidance for installing and configuring RADAR-base platform on a Kubernetes cluster, it is not a replacement for the official detailed documentation or tutorial of Kubernetes, Helm or Helmfile. If you are not familiar with these tools, we strongly recommend you to get familiar with these tools. Here is a [list of useful links](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2731638785/How+to+get+started+with+tools+around+RADAR-Kubernetes) to get started. 

## Prerequisites
### Infrastructure
Before installation you need to have these applications or services available.

#### Kubernetes cluster
You need to have a working Kubernetes installation and there are 3 ways to have that:
* From cloud providers such as: [AWS](https://aws.amazon.com/eks/), [GCP](https://cloud.google.com/kubernetes-engine/), [Azure](https://docs.microsoft.com/en-us/azure/aks/), etc
* Install it on your own servers with: [Rancher](https://rancher.com/), [Kubespray](https://github.com/kubernetes-sigs/kubespray), [Kubeadm](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/), [OpenShift](https://www.okd.io/), etc
* Install it on your local machine with: [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/), [K3S](https://k3s.io/), etc

If you're not using a cloud provider you need to make sure that you can [load balance](https://kubernetes.github.io/ingress-nginx/deploy/baremetal/), [expose  applications](https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/#exposing-the-service) and provide [persistent volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).

#### DNS server
Some applications are only accessible via HTTPS and it's essential to have a DNS server via providers like GoDaddy, Route53, etc

#### SMTP Server
RADAR-Base needs an SMTP server to send registration email to researchers and participants.

#### Dockerhub account
Due to introduction of [Dockerhub download rate limiter ](https://docs.docker.com/docker-hub/download-rate-limit/) and number of images used in RADAR-Kubernetes installation you might hit the limit size pretty quickly so it's advised to at least use an authenticated account or purchase a subscription. You can also use a third party container registry to cache RADAR-Base images and avoid using Dockerhub.

#### AWS KMS


#### (Optional) Object storage
An external object storage allows RADAR-Kubernetes to backup cluster data such as manifests, application configuration and data via [Velero](https://velero.io/) to a backup site.

You can also send the RADAR-Base output data to this object storage, which can provider easier management and access compared to bundled [Minio](https://min.io/) server inside RADAR-Kubernetes.

#### (Optional) Managed services
Radar-Kubernetes provides Kafka cluster, PostgreSQL server, Object storage, Log server and Monitoring system inside the cluster but you can disable these bundled  components and use an externally managed services like Confluent Cloud, Azure Database for PostgreSQL, AWS S3, etc instead.

#### (Optional) Fitbit account
Fitbit account is needed in case you want to use Fitbit wearables to collect data.


### Local machine

You need to have following tools installed in your local machine to install the stack: [Git](https://git-scm.com/downloads), [Java](https://openjdk.java.net/install/), [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/), [helm](https://github.com/helm/helm#install), [helm-diff](https://github.com/databus23/helm-diff#install), [helm-secrets](https://github.com/zendesk/helm-secrets), [sops](https://github.com/mozilla/sops), [helmfile](https://github.com/roboll/helmfile#installation). You can click on the links to go the installation page of each tools. After installing Kubectl make sure you add credentials of Kubernetes cluster to it.

## Installation
> In this guide it is assumed that you're using Linux in your local machine but you can follow these steps with little to no change in MacOS and Windows as well.

#### Prepare
First step is to clone the repository to your local machine. Following command will take care of that while also making sure confluent helm charts such as Kafka and Zookeeper are also downloaded as a submodule with `--recurse-submodules` flag. Using `--branch` will make sure to use latest development branch.
```shell
git clone --branch dev --recurse-submodules https://github.com/RADAR-base/RADAR-Kubernetes.git
```
#### Configure
Now you need to move to configure the installation. Some templates have already been provided to get to familiar with most important configuration options. If you're not sure which components you want to enable you can refer to wiki for [an overview and breakdown on RADAR-Base components and their roles](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2673967112/Component+overview+and+breakdown). For more info on configuration options for each component you can refer to their chart directory and use that documentation.

You also need to input some secrets and passwords used by RADAR-Base applications to either talk to each other or allow connections to it, it is recommended to use a random password generator to fill these secrets for more security. In order to create a password for monitoring system you need to use [this website](https://www.web2generators.com/apache-tools/htpasswd-generator) or `htpasswd` command to create an encrypted password string and put it inside `kube_prometheus_stack.nginx_auth` variable. It seems like bcrypt encryption isn't supported in current ingress-nginx so make sure that you're using MD5 encryption.

```shell
cd RADAR-Kubernetes
cp environments.yaml.tmpl environments.yaml
cp base.yaml production.yaml
cp base.yaml.gotmpl production.yaml.gotmpl
sops -e base.secrets.yaml > production.secrets.yaml
vim environments.yaml # use the files you just created
vim production.yaml  # Change setup parameters and configurations
vim production.yaml.gotmpl  # Change setup parameters that require Go templating, such as reading input files
sops production.secrets.yaml  # Change passwords and credentials
```

The final preparation step is to create a certificate for Management Portal to sign requests.
```shell
./bin/keystore-init
```

#### Install
Finally you can start installing the RADAR-Base stack. Having `--concurrency 1` will make installation slower but it is necessary because some components such as `kube-prometheus-stack` and `kafka-init` (aka `catalog-server`) should be installed in their specified order. If you've forgotten to use this flag the installation might not be successful and you should follow [Uninstallation](#uninstall) steps to clean up the Kubernetes cluster before you can try again.
```shell
helmfile sync --concurrency 1
```
Depending on you cluster this will take a few minutes to run. During the installation you can monitor the process more closely by running `kubectl get pods` and checking if new pods successfully enter Running status and are fully Ready or not.

If an application doesn't become fully ready installation will fail and you need to figure out the issue. In order to do this you can use `kubectl describe pods <podname>` and `kubcetl logs <podname>` to even more closely inspect pods status and health during the installation, fore some components such as aforementioned `kube-prometheus-stack` and `kafka-init` you might need to remove everything before trying again but for most other components you can just run installation command again.

#### Verify
After installation you can check cluster status with `kubectl get pods` command and if it has been successful you should see an output similar to this:
```
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
postgresql-postgresql-master-0                   3/3     Running   0          8m21s
postgresql-postgresql-slave-0                    1/1     Running   0          8m21s
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
radar-upload-postgresql-postgresql-master-0      3/3     Running   0          8m21s
radar-upload-postgresql-postgresql-slave-0       1/1     Running   0          8m21s
radar-upload-source-connector-c87ddc848-fxlg7    1/1     Running   0          8m21s
redis-master-0                                   2/2     Running   0          8m21s
redis-slave-0                                    2/2     Running   0          8m21s
smtp-6646fc65cd-5h2f7                            1/1     Running   0          8m21s
timescaledb-slave-0                              1/1     Running   0          8m21s
timescaledb-slave-1                              1/1     Running   0          8m21s
timescaledb-timescaledb-master-0                 3/3     Running   0          8m21s
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


If you have enabled monitoring you should also check Prometheus to see if there are any alerts firing. In next section there is a guide on how to connect to Prometheus.

Other ways to ensure that installation have been successful is to check application logs for errors and exceptions. Also to check Kafka and make sure it's functional and RADAR-Base topics are loaded in:

```bash
➜  kubectl exec -it cp-kafka-0 -c cp-kafka-broker -- kafka-topics --zookeeper cp-zookeeper-headless:2181 --list | wc -l
273
```
Number of topics can be more or less than this number depending on components that you have activated.

## Usage
### Accessing the applications
In order to access to the applications first you need to find the IP address that Nginx service is listening to and then point the domain that you've specified in `server_name` variable to this IP address via a DNS server (e.g. [Route53](https://aws.amazon.com/route53/), [Cloudflare](https://www.cloudflare.com/dns/), [Bind](https://www.isc.org/bind/)) or [`hosts` file](https://en.wikipedia.org/wiki/Hosts_(file)) in your local machine.
> For this guide we assume that you've set `server_name` to "k8s.radar-base.org" and SSL is enabled.

You can see details of Nginx service with following command:
```
➜ kubectl get service nginx-ingress-controller
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP                           PORT(S)                      AGE
nginx-ingress-controller   LoadBalancer   10.100.237.75   XXXX.eu-central-1.elb.amazonaws.com   80:31046/TCP,443:30932/TCP   1h
```
* If you're using a cloud provider you need to point the value in `EXTERNAL-IP` column (in this example `XXXX.eu-central-1.elb.amazonaws.com`) to `k8s.radar-base.org` domain in your DNS server.
* If you're not using a cloud provider you need to use a load balancer to expose `31046` and `30932` ports (will be different in your setup) to a IP address and then point `k8s.radar-base.org` domain to that IP address.
* For development and testing purposes you can run `sudo kubectl port-forward svc/nginx-ingress-controller 80:80 443:443` which will forward Nginx service ports to your local machine and you can have access to applications after adding `127.0.0.1       k8s.radar-base.org` to your `hosts` file.

**Note:** If you've enabled monitoring or loggingyou should point `*.server_name` domain to the same address as `server_name`.

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
If you want to remove the Radar-base from your cluster you and use following command to delete the applications from cluster:
```
helmfile destroy
```
After there might still be some configuration lingering inside the cluster, you can use following commands to purge them as well.
```
kubectl delete crd prometheuses.monitoring.coreos.com prometheusrules.monitoring.coreos.com servicemonitors.monitoring.coreos.com alertmanagers.monitoring.coreos.com podmonitors.monitoring.coreos.com alertmanagerconfigs.monitoring.coreos.com probes.monitoring.coreos.com thanosrulers.monitoring.coreos.com
kubectl delete psp kube-prometheus-stack-alertmanager kube-prometheus-stack-grafana kube-prometheus-stack-grafana-test kube-prometheus-stack-kube-state-metrics kube-prometheus-stack-operator kube-prometheus-stack-prometheus kube-prometheus-stack-prometheus-node-exporter
kubectl delete mutatingwebhookconfigurations prometheus-admission
kubectl delete ValidatingWebhookConfiguration prometheus-admission

kubectl delete crd certificaterequests.cert-manager.io certificates.cert-manager.io challenges.acme.cert-manager.io clusterissuers.cert-manager.io issuers.cert-manager.io orders.acme.cert-manager.io
kubectl delete pvc --all
kubectl -n cert-manager delete secrets letsencrypt-prod
kubectl -n default delete secrets radar-base-tls
kubectl -n monitoring delete secrets radar-base-tls

kubectl delete crd cephblockpools.ceph.rook.io  cephclients.ceph.rook.io cephclusters.ceph.rook.io cephfilesystems.ceph.rook.io cephnfses.ceph.rook.io cephobjectstores.ceph.rook.io cephobjectstoreusers.ceph.rook.io volumes.rook.io
kubectl delete psp 00-rook-ceph-operator
```
