#  kube-spy 功能测试用例（基于配置文件）

本测试用例目标为位于default命名空间的http-test服务

#### Namespace：

##### TestCase1：

- 输入：与服务实际所在不符的命名空间
- 预期输出：相应报错信息
- 配置文件：[namespace1.yaml](./testConfig/namespace1.yaml)

##### TestCase2：

- 输入：空

- 预期输出：相应报错信息

- 配置文件：[namespace2.yaml](./testConfig/namespace2.yaml)

  

#### VictimServices-name：

##### TestCase1：

- 输入：与服务实际名称不符的名称
- 预期输出：相应报错信息
- 配置文件：[svc_name1.yaml](./testConfig/svc_name1.yaml)

##### TestCase2：

- 输入：空
- 预期输出：相应报错信息
- 配置文件：[svc_name2.yaml](./testConfig/svc_name2.yaml)

##### TestCase3：

- 输入：两次输入相同的服务名
- 预期输出：对该服务分别进行两次故障注入
- 配置文件：[svc_name3.yaml](./testConfig/svc_name3.yaml)



#### VictimServices-ChaosList：

##### TestCase1：

- 输入：空
- 预期输出：不添加chaos，正常运行
- 配置文件：[chaoslist1.yaml](./testConfig/chaoslist1.yaml)

##### TestCase2：

- 输入：在某个服务上添加一个chaos
- 预期输出：添加chaos后正常运行
- 配置文件：[chaoslist2.yaml](./testConfig/chaoslist2.yaml)

##### TestCase3：

- 输入：在某个服务上添加多个chaos
- 预期输出：分多次测试每个chaos的添加影响
- 配置文件：[chaoslist3.yaml](./testConfig/chaoslist3.yaml)



#### VictimServices-ChaosList-replica：

##### TestCase1：

- 输入：空
- 预期输出：副本数不变，进行测试
- 配置文件：[replica1.yaml](./testConfig/replica1.yaml)

##### TestCase2：

- 输入：1（原本有多个副本，或原本就是1个副本）
- 预期输出：副本数被调整为1后进行测试
- 配置文件：[replica2.yaml](./testConfig/replica2.yaml)

##### TestCase3：

- 输入：3（原本有3个以上副本，或原本1个副本）
- 预期输出：副本数被调整为3后进行测试
- 配置文件：[replica3.yaml](./testConfig/replica3.yaml)



#### VictimServices-ChaosList-range：

##### TestCase1：

- 输入：空
- 预期输出：所有副本都没有被添加chaos
- 配置文件：[range1.yaml](./testConfig/range1.yaml)

##### TestCase2：

- 输入：2（replica为3）
- 预期输出：副本数被调整为3后，对前两个pod添加了chaos进行测试
- 配置文件：[range2.yaml](./testConfig/range2.yaml)

##### TestCase3：

- 输入：50%（replica为3）
- 预期输出：副本数被调整为3后，对第一个（取下整）pod添加了chaos进行测试
- 配置文件：[range3.yaml](./testConfig/range3.yaml)



#### VictimServices-ChaosList-ingress/egress：

##### TestCase1：

- 输入：空
- 预期输出：所有副本都没有被添加chaos
- 配置文件：[chaos1.yaml](./testConfig/chaos1.yaml)

##### TestCase2：

- 输入：ingress: "1kbps,,"
- 预期输出：服务入口流量被限速，导致服务接受请求延时，不影响发送请求
- 配置文件：[chaos2.yaml](./testConfig/chaos2.yaml)

##### TestCase3：

- 输入：egress: "1kbps,,"
- 预期输出：服务出口流量被限速，导致服务发送请求延时，不影响接受请求
- 配置文件：[chaos3.yaml](./testConfig/chaos3.yaml)

##### TestCase4：

- 输入：ingress: "1kbps,,"  egress: "1kbps,,"
- 预期输出：服务出入口流量均被限速，导致服务发送、接受请求延时
- 配置文件：[chaos4.yaml](./testConfig/chaos4.yaml)

##### TestCase5：

- 输入：ingress: ",delay,50ms" 
- 预期输出：服务入口流量被延时50ms
- 配置文件：[chaos5.yaml](./testConfig/chaos5.yaml)

##### TestCase6：

- 输入：egress: ",delay,50ms" 
- 预期输出：服务出口流量被延时50ms
- 配置文件：[chaos6.yaml](./testConfig/chaos6.yaml)

##### TestCase7：

- 输入：ingress: ",delay,50ms"  egress: ",delay,50ms" 
- 预期输出：服务入口和出口流量均被延时50ms
- 配置文件：[chaos7.yaml](./testConfig/chaos7.yaml)

