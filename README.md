# Cloudship Operator

```console
operator-sdk init --repo=github.com/ToucanSoftware/cloudship-operator
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=Application --namespaced=true --resource --controller
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=AppService --namespaced=true --resource --controller
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=AppResource --namespaced=true --resource --controller
```
