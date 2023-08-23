#!/bin/bash

set -euo pipefail

kubectl create ns argocd  -o=yaml --dry-run=client | kubectl apply --server-side -f -

kubectl apply --server-side -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml > /dev/null

# wait for argocd pods to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=argocd-server -n argocd --timeout=300s

kubectl patch configmaps -n argocd argocd-cm --type merge --patch-file argocd-cm.json

kubectl patch configmap argocd-rbac-cm -n argocd --type merge --patch-file argocd-rbac-cm.json

kubectl config set-context --current --namespace=argocd

kubectl port-forward -n argocd svc/argocd-server -n argocd 8080:443 &

argocd login 127.0.0.1:8080 --insecure --username=admin --password="$(argocd admin initial-password | head -n 1)"

echo "------------------------------------------"
echo ""
echo "Token"
argocd account generate-token --account foo --core
echo ""
echo "------------------------------------------"


argocd app create guestbook --repo https://github.com/argoproj/argocd-example-apps.git --path guestbook --dest-server https://kubernetes.default.svc --dest-namespace default

argocd app sync argocd/guestbook --assumeYes > /dev/null
