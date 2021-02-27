grafana-embed-proxy
===================

웹페이지에서 iframe embed를 통해 Grafana의 패널을 삽입할 수 있다.
그러나 Grafana는 iframe을 사용했을 때, 인증 등의 프로세스는 제공하지 않기 때문에
`grafana-embed-proxy`는 이러한 과정을 제공해주는 것이 목표다.


## How to test

### Prepare

- Docker
- Kind
- Helm

### Install test environment

```bash
$ cd ./test

$ make create-cluster
$ make install-ingress
$ make install-grafana
```
