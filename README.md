# RADAR-Kubernetes [![Build Status](https://travis-ci.org/RADAR-base/RADAR-Kubernetes.svg?branch=master)](https://travis-ci.org/RADAR-base/RADAR-Kubernetes)
Kubernetes deployment of RADAR-base.


## Installation
You need to have a working Kubernetes installation and there are 3 ways to have that:
* From cloud providers suchs as:
  * [AWS](https://aws.amazon.com/eks/)
  * [GCP](https://cloud.google.com/kubernetes-engine/)
  * [Azure](https://docs.microsoft.com/en-us/azure/aks/)
  * And many more...
* Install it on your own servers with:
  * [Rancher](https://rancher.com/)
  * [Kubespray](https://github.com/kubernetes-sigs/kubespray)
  * [Kubeadm](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)
  * [OpenShift](https://www.okd.io/)
  * And many more...
* Install it on your local machine with:
  * [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/)
  * [K3S](https://k3s.io/)
  * [Docker Desktop](https://www.docker.com/products/docker-desktop)

**Note 1:** This setup is currently only tested on AWS EKS, OpenStack Magnum and Azure ASK however because of cloud agnostic approach of Kubernetes you should be able install this stack on any Kubernetes installation. Also the idea behind using `helm` and `helmfile` has been allowing more complex and customized setups without too much change in the original repository, so if current approach isn't working in your environment you can easily change components to your needs.

**Note 2:** If you're not using a cloud provider you need to make sure that you can [load balance](https://kubernetes.github.io/ingress-nginx/deploy/baremetal/) and [expose  applications](https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/#exposing-the-service) and provide [persistent volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).

You need to have following tools installed in your machine to install the stack:
* Git
* Java (Used by `keystore-init` script to generate keys)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [helm](https://github.com/helm/helm#install) (Server side component of Helm should be installed in Kubernetes cluster as well)
* [helm-diff](https://github.com/databus23/helm-diff#install)
* [helmfile](https://github.com/roboll/helmfile#installation)

After installing them run following commands:
```shell
git clone --recurse-submodules https://github.com/RADAR-base/RADAR-Kubernetes.git
cd RADAR-Kubernetes
cp environments.yaml.tmpl environments.yaml
cp base.yaml production.yaml
cp base.yaml.gotmpl production.yaml.gotmpl
vim environments.yaml # use the files you just created
vim production.yaml  # Change setup parameters and configurations
vim production.yaml.gotmpl  # Change setup parameters that require Go templating, such as reading input files
./bin/keystore-init
helmfile sync --concurrency 1
```

Having `--concurrency 1` will make installation slower but it is necessary because some components such as `kube-prometheus-stack` and `kafka-init` (aka `catalog-server`) should be installed in their specified order, if you've forgotten to use this flag the installation will not be successful and you need to go some cleaning before you can try to install again:

### kube-prometheus-stack
The kube-prometheus-stack will define a `ServiceMonitor` CRD that other services with monitoring enabled will use, so please make sure that Prometheus chart installs successfully before preceding. By default it's configured to wait for Prometheus deployment to be finished in 10 minutes if this time isn't enough for your environment change it accordingly. If the deployment has been failed for the first time then you should delete it first and then try installing the stack again:
```
helm del --purge kube-prometheus-stack
kubectl delete crd prometheuses.monitoring.coreos.com prometheusrules.monitoring.coreos.com servicemonitors.monitoring.coreos.com alertmanagers.monitoring.coreos.com podmonitors.monitoring.coreos.com alertmanagerconfigs.monitoring.coreos.com probes.monitoring.coreos.com thanosrulers.monitoring.coreos.com
kubectl delete psp kube-prometheus-stack-alertmanager kube-prometheus-stack-grafana kube-prometheus-stack-grafana-test kube-prometheus-stack-kube-state-metrics kube-prometheus-stack-operator kube-prometheus-stack-prometheus kube-prometheus-stack-prometheus-node-exporter
kubectl delete mutatingwebhookconfigurations prometheus-admission
kubectl delete ValidatingWebhookConfiguration prometheus-admission
```

### kafka-init
This service will create schemes in Kafka and it should do it before other services access data from it. If order hasn't been preserved you should remove Kafka and initialize again:
```
helm del --purge cp-zookeeper cp-kafka
kubectl delete pvc datadir-0-cp-kafka-{0,1,2} datadir-cp-zookeeper-{0,1,2} datalogdir-cp-zookeeper-{0,1,2}
```

### Uninstall
If you want to remove the Radar-base from your cluster you need set all of your `RADAR_INSTALL_*` variables in `.env` file to `false` and then run the `helmfile sync --concurrency 1` command to delete the charts after that you need to run following commands to remove all of the traces of the installation:
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

## Volume expansion

If want to resize a volumes after its initialization you need to make sure that it's supported by its underlying volume plugin:
https://kubernetes.io/docs/concepts/storage/persistent-volumes/#expanding-persistent-volumes-claims

If it's supported then it should be an easy process like this:
https://www.jeffgeerling.com/blog/2019/expanding-k8s-pvs-eks-on-aws

## Authentication for monitoring
If you want to have authentication you need to generate it somehow and there are couple of ways to do that. You can use these websites:

http://www.htaccesstools.com/htpasswd-generator/

https://www.web2generators.com/apache-tools/htpasswd-generator

Or following command:
```
sudo docker run --rm httpd:2.4-alpine htpasswd -nb <username> <password>
```
And put the output to the `MONITORING_NGINX_AUTH` variable in `.env` file.
It seems like bcrypt encryption isn't supported in current ingress-nginx so make sure that you're using MD5 encryption.

## Usage
After installation you can check cluster status with `kubectl get pods` command and if it has been successful you should see an output similar to this:
```
NAME                                             READY   STATUS    RESTARTS   AGE
catalog-server-65d7cd5fd5-df7mm                  1/1     Running   0          1h
cp-kafka-0                                       2/2     Running   0          1h
cp-kafka-1                                       2/2     Running   0          1h
cp-kafka-2                                       2/2     Running   0          1h
cp-kafka-rest-7f8fbc6d9-ktlwd                    2/2     Running   0          1h
cp-schema-registry-6bf9fb9f89-5qdzc              2/2     Running   0          1h
cp-zookeeper-0                                   2/2     Running   0          1h
cp-zookeeper-1                                   2/2     Running   0          1h
cp-zookeeper-2                                   2/2     Running   0          1h
hdfs-datanode-0                                  1/1     Running   0          1h
hdfs-datanode-1                                  1/1     Running   0          1h
hdfs-datanode-2                                  1/1     Running   0          1h
hdfs-namenode-0                                  1/1     Running   0          1h
kafka-manager-6c9c75676c-qdckg                   1/1     Running   0          1h
management-portal-748655c69-ghd7z                1/1     Running   0          1h
nginx-ingress-controller-85d89d6776-tgfwv        1/1     Running   0          1h
nginx-ingress-default-backend-6694789b87-jstq6   1/1     Running   0          1h
postgresql-postgresql-master-0                   2/2     Running   0          1h
postgresql-postgresql-slave-0                    1/1     Running   0          1h
radar-backend-monitor-b468ff5d9-2gtsm            1/1     Running   0          1h
radar-backend-stream-bd8c768cb-4rf8v             1/1     Running   0          1h
radar-connect-hdfs-sink-5b74664d84-n2vfv         1/1     Running   0          1h
radar-dashboard-6fbb7b646c-rmmf7                 1/1     Running   0          1h
radar-gateway-74496dd958-6vz5v                   2/2     Running   0          1h
radar-hotstorage-84777f4fc-q7mhd                 1/1     Running   0          1h
radar-mongodb-connector-5b6b9b77d4-4g7c6         1/1     Running   0          1h
radar-output-7579864f58-c9bzd                    2/2     Running   0          1h
radar-restapi-bfd4c87-dvl8v                      1/1     Running   0          1h
smtp-57fff69b4f-gvrqv                            1/1     Running   0          1h
```

If you have enabled monitoring and SSL you should see these as well:
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
➜ kubectl -n cert-manager get pods
NAME                            READY   STATUS    RESTARTS   AGE
cert-manager-776cd4f499-688bm   1/1     Running   0          1h
```

### Accessing the applications
In order to access to the applications first you need to find the IP address that Nginx service is listening to and then point the domain that you've specified in `SERVER_NAME` variable to this IP address via a DNS server (e.g. [Route53](https://aws.amazon.com/route53/), [Cloudflare](https://www.cloudflare.com/dns/), [Bind](https://www.isc.org/bind/)) or [`hosts` file](https://en.wikipedia.org/wiki/Hosts_(file)) in your local machine.
> For this guide we assume that you've set `SERVER_NAME` to "k8s.radar-base.org" and SSL is enabled.

You can see details of Nginx service with following command:
```
➜ kubectl get service nginx-ingress-controller
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP                           PORT(S)                      AGE
nginx-ingress-controller   LoadBalancer   10.100.237.75   XXXX.eu-central-1.elb.amazonaws.com   80:31046/TCP,443:30932/TCP   1h
```
* If you're using a cloud provider you need to point the value in `EXTERNAL-IP` column (in this example `XXXX.eu-central-1.elb.amazonaws.com`) to `k8s.radar-base.org` domain in your DNS server.
* If you're not using a cloud provider you need to use a load balancer to expose `31046` and `30932` ports (will be different in your setup) to a IP address and then point `k8s.radar-base.org` domain to that IP address.
* For development and testing purposes you can run `sudo kubectl port-forward svc/nginx-ingress-controller 80:80 443:443` which will forward Nginx service ports to your local machine and you can have access to applications after adding `127.0.0.1       k8s.radar-base.org` to your `hosts` file.

**Note:** If you've enabled monitoring you should point `*.SERVER_NAME` domain to the same address as `SERVER_NAME`.

Now depending on your setup you should have access to following URLs:
```
https://alertmanager.k8s.radar-base.org
https://grafana.k8s.radar-base.org
https://k8s.radar-base.org/api
https://k8s.radar-base.org/dashboard
https://k8s.radar-base.org/kafka
https://k8s.radar-base.org/kafkamanager
https://k8s.radar-base.org/managementportal
https://k8s.radar-base.org/schema
https://prometheus.k8s.radar-base.org
```

**Note:** If you have enabled the SSL you might see invalid certificate error when you try to access to the websites, in this case wait a couple of minutes until `cert-manager` issues those certificates.

#### Radar output
If `RADAR_INSTALL_HDFS` is set to `true` you can have access to Radar output via SFTP protocol. Default username is `dl` and password login is disabled and you can only connect to it via a ssh key pair (you should put the public key in `RADAR_OUTPUT_SFTP_PUBLIC_KEY` variable).\
In order to get `host` you should run this command:
```
➜ kubectl get service radar-output
NAME           TYPE           CLUSTER-IP      EXTERNAL-IP                           PORT(S)        AGE
radar-output   LoadBalancer   10.100.43.218   XXXX.eu-central-1.elb.amazonaws.com   22:32606/TCP   1h
```
And now you can use value in `EXTERNAL-IP` column (in this example `XXXX.eu-central-1.elb.amazonaws.com`) as `host` to connect to SFTP server.\
Alternatively you can forward SSH port to your local machine and connect locally via this command:
```
kubectl port-forward svc/radar-output 2222:22
```
Now you can use "127.0.0.1" as `host` and "2222" as the `port` to connect to SFTP server.
