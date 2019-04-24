# RADAR-Kubernetes
Kubernetes deployment of RADAR-base


```shell
git clone https://github.com/confluentinc/cp-helm-charts
helm install --name zookeeper -f zookeeper.yaml cp-helm-charts/charts/cp-zookeeper
helm install --name kafka -f kafka.yaml cp-helm-charts/charts/cp-kafka
helm install --name schema-registry -f schema-registry.yaml cp-helm-charts/charts/cp-schema-registry
helm install --name rest-proxy -f rest-proxy.yaml cp-helm-charts/charts/cp-kafka-rest
helm install --name postgres -f postgres.yaml stable/postgresql
helm install --name kafka-manager -f kafka-manager.yaml stable/kafka-manager
```
