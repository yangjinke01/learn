# Delete Namespace Terminating

## Export namespace json

```shell
kubectl get namespace <YOUR_NAMESPACE> -o json > <YOUR_NAMESPACE>.json
```

## Edit Json

remove kubernetes from finalizers array which is under spec

## Replace namespace

```shell
kubectl replace --raw "/api/v1/namespaces/<YOUR_NAMESPACE>/finalize" -f ./<YOUR_NAMESPACE>.json
```
