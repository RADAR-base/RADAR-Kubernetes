# RADAR-Kubernetes
Kubernetes deployment of RADAR-base.

**Note:**
This repository is still in **alpha** stage and it's not ready for production use.


#### Installation
Install [helm](https://github.com/helm/helm#install), [helm-diff](https://github.com/databus23/helm-diff#install) and [helmfile](https://github.com/roboll/helmfile#installation) and then run following commands:
```shell
git clone https://github.com/RADAR-base/cp-helm-charts.git
kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.7/deploy/manifests/00-crds.yaml
kubectl create namespace cert-manager
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
./bin/keystore-init
cp env.template .env
vim .env  # Change setup parameters and configurations
source .env
helmfile sync --concurrency 1
```

Having `--concurrency 1` is necessary because some components such as `prometheus-operator` and `kafka-init` (aka `catalog-server`) should be installed in their specified order, if you've forgotten to use this flag the installation will not be successful and you need to go some cleaning before you can try to install again:

###### prometheus-operator
The Prometheus-operator will define a `ServiceMonitor` CRD that other services with monitoring enabled will use, so please make sure that Prometheus chart installs successfully before preceding. By default it's configured to wait for Prometheus deployment to be finished in 10 minutes if this time isn't enough for your environment change it accordingly. If the deployment has been failed for the first time then you should delete it first and then try installing the stack again:
```
helm del --purge prometheus
kubectl delete crd prometheuses.monitoring.coreos.com prometheusrules.monitoring.coreos.com servicemonitors.monitoring.coreos.com alertmanagers.monitoring.coreos.com
```

###### kafka-init
This service will create schemes in Kafka and it should do it before other services access data from it. If order hasn't been preserved you should remove Kafka and initialize again:
```
helm del --purge cp-kafka cp-zookeeper
kubectl delete pvc datadir-0-cp-kafka-{0,1,2} datadir-cp-zookeeper-{0,1,2} datalogdir-cp-zookeeper-{0,1,2}
```

#### Volume expansion

If want to resize a volumes after its initialization you need to make sure that it's supported by its underlying volume plugin:
https://kubernetes.io/docs/concepts/storage/persistent-volumes/#expanding-persistent-volumes-claims

If it's supported then it should be an easy process like this:
https://www.jeffgeerling.com/blog/2019/expanding-k8s-pvs-eks-on-aws

#### Authentication for monitoring
If you want to have authentication you need to generate it somehow and there are couple of ways to do that. You can use these websites:

http://www.htaccesstools.com/htpasswd-generator/

https://www.web2generators.com/apache-tools/htpasswd-generator

Or following command:
```
sudo docker run --rm httpd:2.4-alpine htpasswd -nbB <username> <password>
```
And put the output to the `MONITORING_NGINX_AUTH` variable in `.env` file.
