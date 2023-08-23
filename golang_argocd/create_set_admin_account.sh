#!/bin/bash

set -euo pipefail

echo create argocd namespace
kubectl create ns argocd  -o=yaml --dry-run=client | kubectl apply --server-side -f -

echo install argocd
kubectl apply --server-side -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml > /dev/null

# wait for argocd pods to be ready
echo wait for argocd-server pods to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=argocd-server -n argocd --timeout=300s

echo patch argocd-cm configmap
kubectl patch configmaps -n argocd argocd-cm --type merge --patch-file argocd-cm.json

echo patch argocd-rbac-cm configmap
kubectl patch configmap argocd-rbac-cm -n argocd --type merge --patch-file argocd-rbac-cm.json

echo set namespace to argocd
kubectl config set-context --current --namespace=argocd

echo argocd login
argocd login --core

echo get token for foo account
echo "------------------------------------------"
echo ""
echo "Token"
argocd account generate-token --account foo --core
echo ""
echo "------------------------------------------"

echo create guestbook argocd application
argocd app create guestbook --repo https://github.com/argoproj/argocd-example-apps.git --path guestbook --dest-server https://kubernetes.default.svc --dest-namespace default

echo sync guestbook application
argocd app sync argocd/guestbook --assumeYes > /dev/null


echo Get argocd initial password
echo "------------------------------------------"
echo ""
argocd admin initial-password | head -n 1
echo ""
echo "------------------------------------------"
