grafana-embed-proxy
===================

웹페이지에서 iframe embed를 통해 Grafana의 패널을 삽입할 수 있다.
그러나 Grafana는 iframe을 사용했을 때, 인증 등의 프로세스는 제공하지 않기 때문에
`grafana-embed-proxy`는 이러한 과정을 제공해주는 것이 목표다.


## Concept

- 인증용 API 키를 2개로 나눈다.
  - 외부용 API: 클라이언트와 Proxy 간의 인증키
    - iframe을 위해 form으로 POST를 쏠 것이기 때문에, form-data 형식이여야 한다.
  - 내부용 API: Proxy와 Grafana 간의 인증키
    - Grafana의 인증 기능을 그대로 사용한다.
- HTTP Proxy는 Go의 내장된 Proxy 모듈을 사용한다.
- 인증 등을 위한 필터를 간편히 쓰기 위해 Go-Gin 라이브러리를 사용한다.
- 현재는 개념을 검증하기 위한 코드로, 인증에 대한 세부 사항을 제공하지 않는다.

### Flow

1. 웹페이지에서는 외부용 API 키와 함께 POST로 Proxy에 전송하고, target은 iframe으로 지정한다.
    - 현재 코드는 `testtest` 값을 보낸다.
2. Proxy는 POST 요청을 받아서 인증 및 쿠키 값을 설정한다.
    - Grafana JS가 추가로 호출하는 요청에 대해 인증을 유지하기 위해
3. Proxy는 내부용 Grafana API 키를 헤더에 담아서, Grafana에 첫 HTML을 요청하고 클라이언트로 토스한다.
4. 웹페이지는 토스받은 HTML에서 추가적으로 CSS 및 JS를, 인증키가 담긴 쿠키와 함께 다시 Proxy에게 요청한다.
    - 토스받은 HTML 내의 link 주소가 상대 경로이기 때문에 Grafana가 아닌 Proxy로 요청된다.
5. Proxy는 인증키가 담긴 쿠키를 확인하여, 인증키가 유효하다면 CSS 및 JS를 클라이언트로 토스한다.
6. 모든 Grafana 패널에 대한 리소스가 Proxy 완료되었다.


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

- Grafana URL: http://localhost:10080
- Grafana의 기본 Username 및 Password는 admin/admin 이다.  

### Create Grafana API Key

Grafana의 Configuration > API Keys 메뉴에서 `Add API Key` 버튼을 클릭한다.

Example:

```
eyJrIjoiUElXZHFtNUtHUmFveWFxTWFLa3NlMVJzWUVmRjh0SzUiLCJuIjoiZ3JhZmFuYS1lbWJlZC1wcm94eSIsImlkIjoxfQ==
```

### Start Proxy

```bash
$ go run . \
-grafana-url=http://localhost:10080 \
-grafana-api-key=${생성한 Grafana API 키}
```

### View HTML

`./test/iframe.html` 파일을 웹브라우저에서 오픈하여, Grafana 패널이 잘 나오는지 확인한다.
