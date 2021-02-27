create-cluster:
	kind create cluster --name grafana-embed-proxy \
	--config ./cluster-config.yaml \
	--kubeconfig ./kubeconfig

delete-cluster:
	rm -rf ./kubeconfig
	kind delete clusters grafana-embed-proxy

install-ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

install-grafana:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	helm install grafana grafana/grafana -f grafana_values.yaml

upgrade-grafana:
	helm upgrade grafana grafana/grafana -f grafana_values.yaml
