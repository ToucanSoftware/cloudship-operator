# Cloudship Operator

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FToucanSoftware%2Fcloudship-operator.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FToucanSoftware%2Fcloudship-operator?ref=badge_shield)
![Docker Publish](https://github.com/ToucanSoftware/cloudship-operator/actions/workflows/docker-publish.yaml/badge.svg)

```console
operator-sdk init --repo=github.com/ToucanSoftware/cloudship-operator
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=Application --namespaced=true --resource --controller
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=AppService --namespaced=true --resource --controller
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=AppResource --namespaced=true --resource --controller
```
