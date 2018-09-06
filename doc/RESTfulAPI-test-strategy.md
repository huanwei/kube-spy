## GET
幂等

**参数**：

---

* URL
* 请求Body(可选)（URL）

---

* 返回结果(StatusCode+Body)

---

**测试内容**：

1. 正确返回结果

## POST
非幂等

**参数**：

---

* URL
* 正确请求（Body）
* 错误请求（Body）
* 对应GET（URL+Body）

---

* 对应GET的资源创建正确返回结果（StatusCode+Body）
* 对应GET的资源创建失败返回结果（StatusCode+Body）

---

**测试内容**：

1. 参数正确则返回状态码：200
2. 使用对应GET确认资源已创建/更新
3. 参数错误则返回错误或资源创建/更新失败

## PUT
幂等

**参数**：

---

* URL
* 正确请求（Body）
* 错误请求（Body）
* 对应GET（URL+Body）

---

* 参数正确的状态码（StatusCode）
* 对应GET的资源创建正确返回结果（StatusCode+Body）
* 对应GET的资源创建失败返回结果（StatusCode+Body）

---

**测试内容**：

1. 参数正确则返回状态码：2xx
2. 使用对应GET确认资源已创建/更新
3. 多次相同测试执行结果相同
4. 参数错误则返回错误并且资源没有创建/更新

## PATCH
非幂等

**参数**：

---

* URL
* 请求体(Body)
* 对应GET（URL+Body）

---

* 参数正确的状态码（StatusCode）
* 对应GET的资源创建正确返回结果（StatusCode+Body）
* 对应GET的资源创建失败返回结果（StatusCode+Body）

---

**测试内容**：

1. 参数正确则返回状态码：2xx
2. 使用对应GET确认资源已更新
3. 参数错误则返回错误并且资源没有创建/更新

## DELETE
幂等

**参数**：

---

* URL
* 对应POST（URL+Body）
* 对应GET （URL+Body）

---

* 对应POST的资源创建正确返回结果（StatusCode+Body）
* 对应GET的资源创建正确返回结果（StatusCode+Body）
* DELETE的正确返回结果（StatusCode+Body）
* 对应GET的资源不存在的返回结果（StatusCode+Body）

---

**测试内容**：

1. 使用POST创建资源
2. 使用GET确认资源已创建
3. 删除POST创建的资源（**若资源地址在POST返回值中，则需要取得POST返回值**）
4. 使用GET后返回404

## HEAD
幂等

**参数**：

---

* URL

---

* 预期的headers列表
* 预期的状态码

---

**测试内容**：

* 使用多种query parameter测试API是否正常

## OPTIONS
幂等

**参数**：

---

* URL

---

* 预期的headers列表
* 预期的状态码

---

**测试内容**：

* 测试API是否能正常响应不支持的请求



