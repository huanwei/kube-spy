# kube-spy

## Config格式
> "**\***"注明的选项为可选项

#### `namespace` 命名空间
```
# Your services' namespace in kubernetes
Namespace: "default"
```
namespace用于指定被测试应用在集群中所在的namespace，目前一次测试只支持在一个命名空间内的应用。

#### `VictimServices` 服务列表

```
# Your service list, the first one will be considered the API server
VictimServices:
- name: "http-test-service1"
  chaosList:
  - replica: 1
    ingress: ",delay,100ms"
    egress: ",delay,100ms"
- name: "http-test-service2"
  chaosList:
  - ingress: ",delay,100ms"
    egress: ",delay,100ms"
- name: "http-test-service3"
  chaosList:
  - ingress: ",delay,100ms"
    egress: ",delay,100ms"
```
本列表用于指定将被用于测试的服务，测试按照列表顺序进行，每个服务分为有`name`和*`chaosList`两项参数，`name`指定该服务在Kubernetes集群中的service名，而\*`chaosList`指定在该服务上注入的故障列表。

*`chaosList`上可以有任意项，每一项有三个参数：

* *`replica`指定这个服务所对应的deployment的副本数，可以用这个参数来进行副本数伸缩测试，0或不填代表不进行副本数控制；

* *`ingress`和\*`egress`分别指定该服务所对应Pod上的入境流量和出境流量的故障参数。

#### *`APIServerAddr` 服务接口地址
```
# This can override the address of the API server to enable out-of-cluster test
APIServerAddr: "httpbin.org"
```
本参数指定应用对外的API接口地址，如果不指定，则默认为`VictimServices`中第一个服务的地址。后续的API测试都将对这个地址发起请求。

#### *`APISetting` API全局设置
```
# Every request will carry them, and they can be overridden
APISetting:
#  authToken: "BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F"
#  headers:
#    Content-Type: "application/json"
```
在这里可以为所有API测试用例进行设置，比如设置权限认证，或者请求头等，本参数格式与单个的测试用例相同。

如果在随后的测试用例中有与本设置冲突的参数，则优先使用测试用例中的设置。

#### `TestCases` 测试用例列表
```
# These test case will be tested in every loop
TestCases:
# Set request json body and headers
  - method: "GET"
    url: "/headers"
    headers:
      Content-Type: "application/json"
    body: "{ languages: [ 'Ruby', 'Perl', 'Python', 'c' ] }"

# Set form and files
  - method: "POST"
    url: "/post"
    form:
      first_name: "Jeevanandam"
      last_name:  "M"
    files:
      1: "/tmp/spy.INFO"
      2: "/tmp/spy.INFO"


# Set query params and bearer auth
  - method: "Get"
    url: "/get"
    params:
      data: "a"
    headers:
    authToken: "C6A79608-782F-4ED0-A11D-BD82FAD829CD"

# Set multi value form and params
  - method: "Post"
    url: "/post"
    multiValueForm:
      data:
        - "a"
        - "b"
      text:
        - "c"
    multiValueParams:
      data:
        - "a"
        - "b"
      text:
        - "c"

# Set path params
  - method: "Get"
    url: "/base64/{value}"
    pathParams:
      value: "a3ViZS1zcHk="

# Set basic auth
- method: "get"
  url: "/"
    basicAuth:
      username: "root"
      password: "123456"
```
#### *`ClientSetting` 客户端设置

```
# Client retry settings
ClientSetting:
  retryCount: 0
  retryWait: 1000
  retryMaxWait: 1000
  timeout: 3000

```

## 流程图
![](img/execProcess.png)
