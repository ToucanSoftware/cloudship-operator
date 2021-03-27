# Cloudship Operator

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FToucanSoftware%2Fcloudship-operator.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FToucanSoftware%2Fcloudship-operator?ref=badge_shield)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docker](https://github.com/ToucanSoftware/cloudship-operator/actions/workflows/docker-publish.yml/badge.svg?branch=main)](https://github.com/ToucanSoftware/cloudship-operator/actions/workflows/docker-publish.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ToucanSoftware/cloudship-operator)](https://goreportcard.com/report/github.com/ToucanSoftware/cloudship-operator)
[![GoDoc](https://godoc.org/github.com/ToucanSoftware/cloudship-operator?status.svg)](https://godoc.org/github.com/ToucanSoftware/cloudship-operator)

```console
operator-sdk init --repo=github.com/ToucanSoftware/cloudship-operator
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=Application --namespaced=true --resource --controller
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=AppService --namespaced=true --resource --controller
operator-sdk create api --group=cloudship --version=v1alpha1 --kind=AppResource --namespaced=true --resource --controller
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FToucanSoftware%2Fcloudship-operator.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FToucanSoftware%2Fcloudship-operator?ref=badge_large)
