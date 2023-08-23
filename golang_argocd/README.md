Inspired by this great [article](https://nemre.medium.com/manage-argocd-resources-programmatically-with-golang-5fa825f1f36e)

You will need a .env file with the following variables:
```
ARGOCD_TOKEN=your_argocd_token
```

Before running the app you need to access the argocd server. You can do this by port forwarding:
```
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

## Creating, setting admin account and getting a token

```
bash create_set_admin_account.sh
```
