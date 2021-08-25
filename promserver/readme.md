## Docs
* https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
* https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md

## Install
``` sh
# install prometheus stack
kubectl create namespace prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install -n prometheus prometheus prometheus-community/kube-prometheus-stack
# install servicemonitor
kubectl apply -n prometheus -f ./prom-servicemonitor.yaml
```

## Uninstall
``` sh
# uninstall servicemonitor
kubectl delete -n prometheus -f ./prom-servicemonitor.yaml
# uninstall prometheus stack
helm delete -n prometheus prometheus
kubectl delete -n prometheus crd alertmanagerconfigs.monitoring.coreos.com
kubectl delete -n prometheus crd alertmanagers.monitoring.coreos.com
kubectl delete -n prometheus crd podmonitors.monitoring.coreos.com
kubectl delete -n prometheus crd probes.monitoring.coreos.com
kubectl delete -n prometheus crd prometheuses.monitoring.coreos.com
kubectl delete -n prometheus crd prometheusrules.monitoring.coreos.com
kubectl delete -n prometheus crd servicemonitors.monitoring.coreos.com
kubectl delete -n prometheus crd thanosrulers.monitoring.coreos.com
```

## Check
``` sh
# prometheus
kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090:9090 -n prometheus
# grafana: admin:prom-operator
kubectl port-forward $(kubectl get pods --no-headers -o custom-columns=":metadata.name" -n prometheus | grep prometheus-grafana) 3000:3000 -n prometheus
```