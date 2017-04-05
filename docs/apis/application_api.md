# <span id="1">目录</span>

* **[目录](#1)**
* **[协议](#2)**
* **[版本](#3)**
* **[提示](#4)**
* **[更改](#5)**
* **[传输](#6)**
* **[格式](#7)**
* **[规范](#8)**
    - [规范](#8.1)
    - [请求](#8.2)
    - [响应](#8.3)
    - [错误码](#8.4)
* **[协议](#9)**
* **[apiserver组件模块](#10)**
	- [查询应用](#10.1)
  - [创建应用](#10.2)
  - [启动应用](#10.3)
  - [停止应用](#10.4)
  - [删除应用](#10.5)
  - [弹性伸缩](#10.6)
  - [灰度升级](#10.7)
  - [重新部署](#10.8)
  - [动态扩容](#10.9)
* **[docker-build组件模块](#11)**
  - [ 构建应用](#11.1)

# <span id="2">协议</span>

## <span id="3">版本</span>
---

**alpha**

## <span id="4">提示</span>
---

本文为markdown格式文本，可使用beyond compare或类似工具比较版本间的修改。
改动时请拉取最新代码进行改动(推荐)，或者在git.asts365.com对应文件下进行编辑。
不要用空格缩进，而应该用tab缩进。

## <span id="5">更改</span>
---
- 2017/03/21, 黄佳, 1.0
  * 协议模板创建


## <span id="6">传输</span>
---

> 使用HTTP作为传输层; 

> 使用UTF-8编码; 

## <span id="7">格式</span>
---

> 请求使用原始的HTTP格式；

> 响应使用JSON封装，详情见下面响应说明；

> 时间格式采用如下形式：yyyy:mm:dd hh:mm:ss;

消息格式为Json,
参考：http://www.json.org/json-zh.html

## <span id="8">规范</span>
---

### <span id="8.1">规范</span>

> 大体上符合REST风格，非严格的restful,类restful风格；

> URL都采用单数，复数的情况使用路径文件夹形式，例如POST BaseURI/app/, 注意最后的'/'表示文件夹；

> 命名采用小写开头，驼峰格式，例如appId;

### <span id="8.2">请求</span>

> GET: 用于读取信息，参数在query中，成功返回200；幂等；

> POST: 主要用于创建，也可以用于更改，参数在body中，成功返回201；非幂等；

> PUT: 用于更改已有资源，参数和POST一样，成功返回201；非幂等；

> DELETE：用于删除资源，成功返回204；非幂等；

### <span id="8.3">响应</span>

- 格式如下：

```text
{
    "apiversion": "alpha",
    "code": 200,
    "err": 0,
    "msg":"",
    "data": {
        "totalSize": 200
    }
}
```

- 空数组：
 "data": []

- 空对象：
 "data": {}

### <span id="8.4">错误码</span>

- 200 OK - [GET]：服务器成功返回用户请求的数据，该操作是幂等的（Idempotent）。
- 201 CREATED - [POST/PUT/PATCH]：用户新建或修改数据成功。
- 204 NO CONTENT - [DELETE]：用户删除数据成功。
- 400 INVALID REQUEST - [POST/PUT/PATCH]：用户发出的请求有错误，服务器没有进行新建或修改数据的操作，该操作是幂等的。。
- 404 NOT FOUND - [*]：用户发出的请求针对的是不存在的记录，服务器没有进行操作，该操作是幂等的。
- 500 INTERNAL SERVER ERROR - [*]：服务器发生错误，用户将无法判断发出的请求是否成功。

### <span id="8.5">请求地址</span>

> RootURI: http://192.168.1.101:9090/

> ApiURI: RootURI/api/v1/

## <span id="9">协议</span>
---

- apiserver组件api

## <span id="10">apiserver组件模块</span>
---

### <span id="10.1">查询应用</span>

查询应用。

URI: ApiURI/apps

Method: GET

**请求**

- JSON:

```text
{
    "lessee":"jxcf"
}
```
**说明**：lessee 表示租户的意思，查询应用其实就是查询当前租户下的所有的应用，租户对应到k8s中的namespace

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "apiversion": "alpha",
  "code": 200,
  "err": 0,
  "msg": "",
  "data": [
    {
      "id": 0,
      "name": "",
      "region": "",
      "memory": "",
      "cpu": "",
      "instanceCount": 0,
      "envs": "",
      "ports": "",
      "image": "",
      "userName": "",
      "remark": ""
    }
  ]
}
```

### <span id="10.2">创建应用</span>

部署应用。

URI: ApiURI/apps

Method: POST

**请求**

- ApiURI/apps 
- JSON
```text
{
  "name": "test",
  "region": "north",
  "memory": "128Mi",
  "cpu": "128m",
  "instanceCount": 1,
  "image": "index.tenxcloud.com/carrotzpc/docker-2048:latest",
  "userName": "huangjia",
  "remark": "this is a test web application"
}

```
- 说明:

**响应**

- HTTP Status: 201
- JSON:

```text
{
  "api": "1.0",
  "status": "201",
  "err": "OK",
  "msg": "create app successed"
}
```

### <span id="10.3">启动应用</span>

上传文件。

URI: ApiURI/apps

Method: PUT

**请求**

- ApiURI/apps

```text
{
  "app_name":"test",
  "ns":"huangjia"
}
```
- 说明：


**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": "200",
  "err": "OK",
  "msg": "start application named test successed"
}
```

### <span id="10.4">停止应用</span>

停止应用。

URI: ApiURI/apps

Method: PATCH

**请求**

- ApiURI/apps

```text
{
	"app_name":"test",
  "ns":"huangjia"
}
```

- 说明：


**响应**

```text
{
  "api": "1.0",
  "status": "200",
  "err": "OK",
  "msg": "stop application named test successed"
}
```

### <span id="10.5">删除应用</span>

删除应用。

URI: ApiURI/apps

Method: DELETE

**请求**

```text
{
  "app_name":"test",
  "ns":"huangjia"
}
```

**响应**

- HTTP Status: 204;
- JSON:

```text
{
  "api": "1.0",
  "status": "204",
  "err": "OK",
  "msg": "delete app successed"
}
```

### <span id="10.6">弹性伸缩</span>

弹性伸缩。

URI: ApiURI/apps/scale

Method: PATCH

**请求**

```text
{
  "app_name":"test",
  "ns":"huangjia",
  "app_cnt":3
}
```

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "api": "1.0",
  "status": "201",
  "err": "OK",
  "msg": "scale application named test successed"
}
```

### <span id="10.7">灰度升级</span>

灰度升级。

URI: ApiURI/apps/rollupdate

Method: POST

**请求**

```text
{
  "app_name":"test",
  "ns":"huangjia",
  "image":"index.tenxcloud.com/carrotzpc/docker-2048:20150730"
}
```

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "api": "1.0",
  "status": "201",
  "err": "OK",
  "msg": "rolling update application named test successed"
}
```

### <span id="10.8">重新部署</span>

重新部署。（该接口暂时还未提供）

URI: ApiURI/apps

Method: UPDATE

**请求**

```text
<!-- {
  "id":1
} -->
```

**响应**

- HTTP Status: 201;
- JSON:

```text
<!-- {
  "apiversion": "alpha",
  "code": 201,
  "err": 0,
  "msg": "rolling update app successed",
} -->
```

### <span id="10.9">动态扩容</span>

动态扩容。

URI: ApiURI/apps/expansion

Method: PUT

**请求**

```text

{
  "app_name":"test",
  "ns":"huangjia",
	"cpu":"512Mi",
	"memory":"512m"
}
```

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "api": "1.0",
  "status": "201",
  "err": "OK",
  "msg": "Expansion application named test successed"
}
```

- docker-build组件api

## <span id="11">docker-build组件模块</span>
---

### <span id="11.1">在线构建应用</span>

查询应用。

URI: ApiURI/builds

Method: POST

**请求**

- JSON:

```text
 {
    "app_name": "", 
    "version": "", 
    "remark": "", 
    "registry": "", 
    "repository": "", 
    "branch": ""
}

```
**参数说明**：
- app_name：构建应用的名称，该名称会用作生成镜像名称,例如：my/xx/app_name:vaersion
- version：构建应用的名称，该名称会用作生成镜像名称的tag,例如：my/xx/app_name:version
- remark：构建应用的描述信息
- registry： 应用上传的镜像仓库地址
- repository：应用的项目代码地址
- branch：应用的项目代码的分支

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "api": "1.0",
  "status": "201",
  "err": "OK",
  "msg": "build application successed"
}
```

### <span id="11.2">离线构建应用</span>

查询应用。

URI: ApiURI/builds

Method: PUT

**请求**

- JSON:

```text
{
    "app_name": "", 
    "version": "", 
    "remark": "", 
    "registry": "",
    "baseImage": "", 
    "tarball": ""
}

**参数说明**：
- app_name：构建应用的名称，该名称会用作生成镜像名称,例如：my/xx/app_name:vaersion
- version：构建应用的名称，该名称会用作生成镜像名称的tag,例如：my/xx/app_name:version
- remark：构建应用的描述信息
- registry： 应用上传的镜像仓库地址
- baseImage：构建应用的基础镜像
- tarball：应用的压缩包文件

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "api": "1.0",
  "status": "201",
  "err": "OK",
  "msg": "build application successed"
}
```