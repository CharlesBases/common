#!/bin/zsh

set -e

DOCKER_VERSION=master
DASHBOARD_PORT=8000

gopath=$GOPATH
gopath=${gopath%:*}

mkdir -p $gopath/src/github.com
cd $gopath/src/github.com
git clone https://github.com/AliyunContainerService/k8s-for-docker-desktop.git kubernetes
cd kubernetes
git checkout ${DOCKER_VERSION}
./load_images.sh
rm -rf ../kubernetes

# restart docker
killall Docker && open /Applications/Docker.app

# docker-for-desktop
kubectl config use-context docker-for-desktop

# Dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml

# external access
# kubectl port-forward kubernetes-dashboard-7798c48646-wkgk4 ${DASHBOARD_PORT}:8443 --namespace=kube-system &

# local access
kubectl proxy &
# http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/

# Token
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep kubernetes-dashboard-token | awk '{print $1}')
