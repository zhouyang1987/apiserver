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
* **[应用接口](#10)**
  - [查询应用](#10.1)
  - [创建应用](#10.2)
  - [启动应用](#10.3)
  - [停止应用](#10.4)
  - [删除应用](#10.5)
  - [重新部署](#10.6)
* **[服务接口](#11)**
  - [查询服务](#11.1)
  - [创建服务](#11.2)
  - [删除服务](#11.3)
  - [弹性伸缩](#11.4)
  - [灰度升级](#11.5)
  - [重新部署](#11.6)
  - [动态扩容](#11.7)
* **[容器接口](#12)**
  - [查询容器](#12.1)
  - [重新部署容器](#12.2)
* **[日志接口](#13)**
  - [获取服务的事件](#13.1)
  - [获取应用pod的cpu/memory实用情况](#13.2)
  - [获取容器的事件](#13.3)
* **[镜像接口](#14)**
  - [获取镜像列表](#14.1)
* **[构建接口](#15)**
  - [ 构建应用](#15.1)
* **[服务配置接口](#16)**
  - [添加服务配置](#16.1)
  - [添加服务配置子文件](#16.2)
  - [删除服务配置](#16.3)
  - [删除服务配置子文件](#16.4)
  - [查询服务配置](#16.5)

# <span id="2">协议</span>

## <span id="3">版本</span>
---

**v1**

## <span id="4">提示</span>
---

本文为markdown格式文本，可使用beyond compare或类似工具比较版本间的修改。
改动时请拉取最新代码进行改动(推荐)，不要用空格缩进，而应该用tab缩进。

## <span id="5">更改</span>
---
- 2017/06/9, 黄佳, 1.1
  * 文档更新


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

> 严格符合REST风格；

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
  "apiversion": "v1",
  "status": "200",
  "data": {
    "apps": [],
    "total": 0
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

## <span id="8.5">请求地址</span>

> RootURI: http://192.168.1.101:9090/

> ApiURI: RootURI/api/v1/

## <span id="9">协议</span>
---
统一采用http协议

## <span id="10">应用接口</span>
---

#### <span id="10.1">查询应用</span>


URI: ApiURI/api/v1/{namespace}/apps?pageCnt=10&pageNum=0&name={appName}

Method: GET

**参数说明**

- namespace: 应用所属租户   必须字段
- pageCnt: 分页查询每页大小  必须字段
- pageNum: 分页查询页码  必须字段
- name: 应用名称  可选字段，当不传时，默认查询当前命名空间下的所有的应用，传入值时，以该值模糊查询结果


**请求**

- ApiURI/api/v1/huangjia/apps?pageCnt=10&pageNum=0&name=nginx


**响应**


```
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "apps": [
      {
        "id": 1,
        "createAt": "2017-06-09T12:27:00+08:00",
        "nmae": "nginx",
        "nameSpace": "huangjia",
        "description": "this is a test nginx",
        "serviceCount": 1,
        "external": "http://10.39.1.45:30327",
        "services": [
          {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "name": "nginx-test",
            "image": "nginx",
            "instanceCount": 1,
            "external": "http://10.39.1.45:30327",
            "loadbalanceIp": "10.39.1.45",
            "config": {
              "id": 1,
              "createAt": "2017-06-09T12:27:00+08:00",
              "base": {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "cpu": "12Mi",
                "memory": "12m",
                "ServiceConfigId": 1
              },
              "config": {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "name": "nginx-test",
                "content": "{\"appName\":\"nginx\"}",
                "containerPath": "/opt",
                "ServiceConfigId": 1
              },
              "super": {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "envs": [
                  {
                    "id": 1,
                    "createAt": "2017-06-09T12:27:00+08:00",
                    "key": "test",
                    "val": "1",
                    "SuperConfigId": 1
                  }
                ],
                "ports": [
                  {
                    "id": 1,
                    "createAt": "2017-06-09T12:27:00+08:00",
                    "containerPort": 8080,
                    "servicePort": 8080,
                    "protocol": "TCP",
                    "SuperConfigId": 1
                  }
                ],
                "ServiceConfigId": 1
              },
              "ServiceId": 1
            },
            "containers": [
              {
                "id": 1,
                "createAt": "2017-06-09T12:27:01+08:00",
                "name": "nginx-test-2991595585-kgxd9",
                "image": "nginx",
                "ServiceId": 1
              }
            ],
            "appId": 1
          }
        ]
      }
    ],
    "total": 1
  }
}
```

#### <span id="10.2">创建应用</span>


URI: ApiURI/api/v1/{namespace}/apps

Method: POST

**参数说明**

**请求**

- ApiURI/api/v1/huangjia/apps


```
{
  "nmae": "nginx",
  "nameSpace": "huangjia",
  "description": "this is a test nginx",
  "serviceCount": 1,
  "services": [
    {
      "name": "nginx-test",
      "image": "nginx",
      "instanceCount": 1,
      "Config": {
        "base": {
          "cpu": "12Mi",
          "memory": "12m",
          "type": 0
        },
        "config": {
          "name": "nginx-test",
          "content": "{\"appName\":\"nginx\"}",
          "containerPath": "/opt"
        },
        "super": {
          "envs": [
            {
              "key": "test",
              "val": "1"
            }
          ],
          "ports": [
            {
              "containerPort": 8080,
              "servicePort": 8080,
              "protocol": "TCP"
            }
          ]
        }
      }
    }
  ]
}
```

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}
```

#### <span id="10.3">启动应用</span>


URI: ApiURI/api/v1/{namespace}/apps/{id}/{verb}

Method: PATCH

**参数说明**
- namespace：应用所属租户
- id：应用id
- verb：操作类型，目前支持三种：start，stop，redeploy

**请求**

- ApiURI/api/v1/huangjia/apps/1/start

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}

```

#### <span id="10.4">停止应用</span>


URI: ApiURI/api/v1/{namespace}/apps{id}/{verb}

Method: PATCH

**参数说明**

**参数说明**
- namespace：应用所属租户
- id：应用id
- verb：操作类型，目前支持三种：start，stop，redeploy


**请求**

- ApiURI/api/v1/huangjia/apps/1/stop

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}

```


#### <span id="10.5">删除应用</span>


URI: ApiURI/api/v1/{namespace}/apps/{id}

Method: DELETE

**参数说明**
- namespace：应用所属租户
- id：应用id



**请求**

- ApiURI/api/v1/huangjia/apps

**响应**

```
{
  "apiversion": "v1",
  "status": "204",
  "data": "ok"
}

```


#### <span id="10.6">重新部署应用</span>


URI: ApiURI/api/v1/{namespace}/apps/{id}/{verb}

Method: PATCH

**参数说明**

**请求**

- ApiURI/api/v1/huangjia/apps

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}
```


## <span id="11">服务接口</span>
---

#### <span id="11.1">查询服务</span>


URI: ApiURI/api/v1/{namespace}/services?pageCnt=10&pageNum=0&name={serviceName}

Method: GET

**参数说明**

- namespace: 服务所属租户   必须字段
- pageCnt: 分页查询每页大小  必须字段
- pageNum: 分页查询页码  必须字段
- name: 应用名称  可选字段，当不传时，默认查询当前命名空间下的所有的应用，传入值时，以该值模糊查询结果


**请求**

- ApiURI/api/v1/huangjia/services?pageCnt=10&pageNum=0&name=nginx


**响应**


```
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "services": [
      {
        "id": 1,
        "createAt": "2017-06-09T12:27:00+08:00",
        "name": "nginx-test",
        "image": "nginx",
        "instanceCount": 1,
        "external": "http://10.39.1.45:30327",
        "loadbalanceIp": "10.39.1.45",
        "config": {
          "id": 1,
          "createAt": "2017-06-09T12:27:00+08:00",
          "base": {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "cpu": "12Mi",
            "memory": "12m",
            "ServiceConfigId": 1
          },
          "config": {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "name": "nginx-test",
            "content": "{\"appName\":\"nginx\"}",
            "containerPath": "/opt",
            "ServiceConfigId": 1
          },
          "super": {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "envs": [
              {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "key": "test",
                "val": "1",
                "SuperConfigId": 1
              }
            ],
            "ports": [
              {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "containerPort": 8080,
                "servicePort": 8080,
                "protocol": "TCP",
                "SuperConfigId": 1
              }
            ],
            "ServiceConfigId": 1
          },
          "ServiceId": 1
        },
        "containers": [
          {
            "id": 1,
            "createAt": "2017-06-09T12:27:01+08:00",
            "name": "nginx-test-2991595585-kgxd9",
            "image": "nginx",
            "ServiceId": 1
          }
        ],
        "appId": 1
      }
    ],
    "total": 0
  }
}
```

#### <span id="11.2">创建服务</span>


URI: ApiURI/api/v1/{namespace}/services

Method: POST

**参数说明**
- namespace: 服务所属租户

**请求**

- ApiURI/api/v1/huangjia/services

```
{   
  "appId":1,
  "name": "nginx-test-1",
  "image": "nginx",
  "instanceCount": 1,
  "loadbalanceIp": "10.39.1.45",
  "config": {
    "base": {
      "cpu": "12Mi",
      "memory": "12m",
      "type": 0
    },
    "config": {
      "name": "nginx-test-1",
      "content": "{\"appName\":\"nginx\"}",
      "containerPath": "/opt"
    },
    "super": {
      "envs": [
        {
          "key": "test",
          "val": "1"
        }
      ],
      "ports": [
        {
          "containerPort": 8080,
          "servicePort": 8080,
          "protocol": "TCP"
        }
      ]
    }
  }
}
```
**参数说明**

- appId：服务所属应用id
- config->base：服务的基本配置，包括cpu，memory，还有服务类型，服务类型两种：有状态  1 无状态 0 ，默认是无状态
- config->config：服务的配置文件挂载
- config->super：包括服务的端口映射和环境变量的定义


**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}
```

#### <span id="11.3">删除服务</span>


URI: ApiURI/api/v1/{namespace}/services/{id}

Method: DELETE

**参数说明**
- namespace: 服务所属租户

**请求**

- ApiURI/api/v1/huangjia/services/2

**响应**

```
{
  "apiversion": "v1",
  "status": "204",
  "data": "ok"
}

```

#### <span id="11.4">弹性伸缩服务</span>


URI: ApiURI/api/v1/{namespace}/services/{id}

Method: PUT

**参数说明**
- namespace: 服务所属租户

**请求**

- ApiURI/api/v1/huangjia/services/2

```
{
  "serviceInstanceCnt":2
}
```

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}

```

#### <span id="11.5">灰度升级服务</span>


URI: ApiURI/api/v1/{namespace}/services/{id}

Method: PUT

**参数说明**
- namespace: 服务所属租户

**请求**

- ApiURI/api/v1/huangjia/services/2

```
{
    "image":"nginx:latest",
    "config":{
          "name": "nginx-test",
          "content": "{\"name\":\"huangjia\"}",
          "containerPath": "/opt"
        }
}
```
**说明**

- image：升级发布的镜像
- config：该服务挂载的配置文件（目前没支持多个配置文件挂载，看需求，如果需要挂载多个配置文件，我会提供支持）
- config->name: 这个name一定要和对应的service的name一致
- config->content: 配置文件内容
- config->containerPath: 挂载到容器中的目录


**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}

```

#### <span id="11.6">重新部署</span>


URI: ApiURI/api/v1/{namespace}/services/{id}

Method: PATCH

**参数说明**

- namespace: 服务所属租户

**请求**

- ApiURI/api/v1/huangjia/services/2

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}

```


#### <span id="11.7">动态扩容</span>


URI: ApiURI/api/v1/{namespace}/services/{id}

Method: PUT

**参数说明**

- namespace: 服务所属租户
- id: 服务id 

**请求**

- ApiURI/api/v1/huangjia/services/2

**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}

```



## <span id="12">容器接口</span>
---

#### <span id="12.1">查询容器</span>


URI: ApiURI/api/v1/{namespace}/containers?pageCnt=10&pageNum=0&name=nginx

Method: GET

**参数说明**

- namespace: 容器所属租户   必须字段
- pageCnt: 分页查询每页大小  必须字段
- pageNum: 分页查询页码  必须字段
- name: 应用名称  可选字段，当不传时，默认查询当前命名空间下的所有的应用，传入值时，以该值模糊查询结果


**请求**

- ApiURI/api/v1/huangjia/containers?pageCnt=10&pageNum=0&name=nginx


**响应**


```
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "containers": [
      {
        "id": 1,
        "createAt": "2017-06-09T12:27:01+08:00",
        "name": "nginx-test-2991595585-kgxd9",
        "image": "nginx",
        "config": {
          "id": 1,
          "createAt": "2017-06-09T12:27:01+08:00",
          "base": {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "cpu": "12Mi",
            "memory": "12m",
            "ServiceConfigId": 1
          },
          "config": {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "name": "nginx-test",
            "content": "{\"appName\":\"nginx\"}",
            "containerPath": "/opt",
            "ServiceConfigId": 1
          },
          "super": {
            "id": 1,
            "createAt": "2017-06-09T12:27:00+08:00",
            "envs": [
              {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "key": "test",
                "val": "1",
                "SuperConfigId": 1
              }
            ],
            "ports": [
              {
                "id": 1,
                "createAt": "2017-06-09T12:27:00+08:00",
                "containerPort": 8080,
                "servicePort": 8080,
                "protocol": "TCP",
                "SuperConfigId": 1
              }
            ],
            "ServiceConfigId": 1
          },
          "ContainerId": 1
        },
        "ServiceId": 1
      }
    ],
    "total": 0
  }
}
```

#### <span id="12.2">重新部署容器</span>


URI: ApiURI/api/v1/{namespace}/containers/{id}
Method: PATCH

**参数说明**

- namespace: 容器所属租户   必须字段
- id: 容器id



**请求**

- ApiURI/api/v1/huangjia/containers/2


**响应**

```
{
  "apiversion": "v1",
  "status": "201",
  "data": "ok"
}
```

## <span id="13">日志接口</span>
---

#### <span id="13.1">获取服务的事件</span>


URI: ApiURI/api/v1/{namespace}/services/{name}/events

Method: GET

**参数说明**

- namespace: 镜像所属租户   必须字段
- name: 服务名称 必须字段

**请求**

- ApiURI//api/v1/huangjia/services/nginx-test/events


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "events": [
            {
                "reason": "Scheduled",
                "message": "Successfully assigned nginx-test-1891245937-t2j17 to slave3",
                "lastTimestamp": "2017-06-11T11:15:45Z",
                "type": "Normal"
            },
            {
                "reason": "Pulled",
                "message": "Container image \"nginx\" already present on machine",
                "lastTimestamp": "2017-06-11T11:15:46Z",
                "type": "Normal"
            },
            {
                "reason": "Created",
                "message": "Created container with id 1ca3818f1e607b7c2ac1a429ef264c3143cec84a2e257e4d0792027dca920751",
                "lastTimestamp": "2017-06-11T11:15:46Z",
                "type": "Normal"
            },
            {
                "reason": "Started",
                "message": "Started container with id 1ca3818f1e607b7c2ac1a429ef264c3143cec84a2e257e4d0792027dca920751",
                "lastTimestamp": "2017-06-11T11:15:47Z",
                "type": "Normal"
            },
            {
                "reason": "SuccessfulCreate",
                "message": "Created pod: nginx-test-1891245937-t2j17",
                "lastTimestamp": "2017-06-11T11:15:45Z",
                "type": "Normal"
            },
            {
                "reason": "ScalingReplicaSet",
                "message": "Scaled up replica set nginx-test-1891245937 to 1",
                "lastTimestamp": "2017-06-11T11:15:45Z",
                "type": "Normal"
            }
        ]
    }
}
```


**metricName参考：**

```
[
  "network/tx",
  "network/tx_errors_rate",
  "memory/working_set",
  "network/tx_errors",
  "cpu/limit",
  "memory/major_page_faults",
  "memory/page_faults_rate",
  "cpu/request",
  "network/rx_rate",
  "cpu/usage_rate",
  "memory/limit",
  "memory/usage",
  "memory/cache",
  "network/rx_errors",
  "network/rx_errors_rate",
  "network/tx_rate",
  "memory/major_page_faults_rate",
  "cpu/usage",
  "network/rx",
  "memory/rss",
  "memory/page_faults",
  "memory/request",
  "uptime"
 ]

```

#### <span id="13.2">获取容器cpu实时使用情况</span>


URI: ApiURI/api/v1/{namespace}/metrics/{name}/{metric}/{type}

Method: GET

**参数说明**

- namespace: 镜像所属租户   必须字段
- name: 容器名称 必须字段
- metric: 粒度名称   {metric}/{type} 组成metricName
- type: 操作类型 {metric}/{type} 组成metricName

**请求**

- ApiURI/api/v1/kube-system/metrics/calico-node-r116q/memory/usage


**响应**

```
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "metrics": {
      "latestTimestamp": "2017-06-09T08:19:00Z",
      "metrics": [
        {
          "timestamp": "2017-06-09T08:05:00Z",
          "value": 69935104
        },
        {
          "timestamp": "2017-06-09T08:06:00Z",
          "value": 69951488
        },
        {
          "timestamp": "2017-06-09T08:07:00Z",
          "value": 69963776
        },
        {
          "timestamp": "2017-06-09T08:08:00Z",
          "value": 69980160
        },
        {
          "timestamp": "2017-06-09T08:09:00Z",
          "value": 69992448
        },
        {
          "timestamp": "2017-06-09T08:10:00Z",
          "value": 70004736
        },
        {
          "timestamp": "2017-06-09T08:11:00Z",
          "value": 70021120
        },
        {
          "timestamp": "2017-06-09T08:12:00Z",
          "value": 70017024
        },
        {
          "timestamp": "2017-06-09T08:13:00Z",
          "value": 70926336
        },
        {
          "timestamp": "2017-06-09T08:14:00Z",
          "value": 70934528
        },
        {
          "timestamp": "2017-06-09T08:15:00Z",
          "value": 70959104
        },
        {
          "timestamp": "2017-06-09T08:16:00Z",
          "value": 70975488
        },
        {
          "timestamp": "2017-06-09T08:17:00Z",
          "value": 70987776
        },
        {
          "timestamp": "2017-06-09T08:18:00Z",
          "value": 71131136
        },
        {
          "timestamp": "2017-06-09T08:19:00Z",
          "value": 71016448
        }
      ]
    }
  }
}
```

**请求**

- ApiURI/api/v1/kube-system/metrics/calico-node-r116q/cpu/usage


**响应**

```
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "metrics": {
      "latestTimestamp": "2017-06-09T08:18:00Z",
      "metrics": [
        {
          "timestamp": "2017-06-09T08:16:00Z",
          "value": 9210518
        },
        {
          "timestamp": "2017-06-09T08:17:00Z",
          "value": 9210518
        },
        {
          "timestamp": "2017-06-09T08:18:00Z",
          "value": 9210518
        }
      ]
    }
  }
}
```

#### <span id="13.3">获取容器的事件</span>


URI: ApiURI/api/v1/{namespace}/containers/{name}/events

Method: GET

**参数说明**

- namespace: 镜像所属租户   必须字段
- name: 容器名称 必须字段

**请求**

- ApiURI//api/v1/huangjia/containers/nginx-test-1891245937-t2j17/events


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "events": [
            {
                "reason": "Scheduled",
                "message": "Successfully assigned nginx-test-1891245937-t2j17 to slave3",
                "lastTimestamp": "2017-06-11T11:15:45Z",
                "type": "Normal"
            },
            {
                "reason": "Pulled",
                "message": "Container image \"nginx\" already present on machine",
                "lastTimestamp": "2017-06-11T11:15:46Z",
                "type": "Normal"
            },
            {
                "reason": "Created",
                "message": "Created container with id 1ca3818f1e607b7c2ac1a429ef264c3143cec84a2e257e4d0792027dca920751",
                "lastTimestamp": "2017-06-11T11:15:46Z",
                "type": "Normal"
            },
            {
                "reason": "Started",
                "message": "Started container with id 1ca3818f1e607b7c2ac1a429ef264c3143cec84a2e257e4d0792027dca920751",
                "lastTimestamp": "2017-06-11T11:15:47Z",
                "type": "Normal"
            }
        ]
    }
}
```


## <span id="14">镜像接口</span>
---


#### <span id="14.1">镜像接口</span>


URI: ApiURI/api/v1/{namespace}/images?pageCnt=10&pageNum=0&name=kube

Method: GET

**参数说明**

- namespace: 镜像所属租户   必须字段
- pageCnt: 分页查询每页大小  必须字段
- pageNum: 分页查询页码  必须字段
- name: 应用名称  可选字段，当不传时，默认查询当前命名空间下的所有的应用，传入值时，以该值模糊查询结果


**请求**

- ApiURI/api/v1/huangjia/services?pageCnt=10&pageNum=0&name=kube


**响应**


```
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "images": [
      {
        "name": "google_containers/kube-proxy-amd64",
        "tagLen": 1,
        "manifest": [
          {
            "namespace": "huangjia",
            "name": "google_containers/kube-proxy-amd64",
            "tag": "v1.6.3",
            "architecture": "amd64",
            "os": "linux",
            "author": "",
            "id": "58e5a78aa714d219fd42cd658ff1738a64b717bf2ae257b896b508caa6a141ed",
            "parent": "e9dd45dc6fc23a2963b4c2bcfbc032350ceef996fd3da568b05ee9754871d7ab",
            "created": "2017-05-10T15:58:19.155908842Z",
            "docker_version": "1.12.6",
            "pull": "docker pull http://10.39.1.48/google_containers/kube-proxy-amd64:v1.6.3"
          }
        ]
      }
    ],
    "total": 1
  }
}
```



## <span id="16">服务配置接口</span>
---


#### <span id="16.1">添加服务配置</span>


URI: ApiURI/api/v1/{namespace}/configs

Method: POST

**参数说明**

- namespace: 镜像所属租户   必须字段



**请求**

- ApiURI/api/v1/huangjia/configs

```
{
  "name":"test"
}
```

**响应**

```
{
    "apiversion": "v1",
    "status": "201",
    "data": "ok"
}
```


#### <span id="16.2">添加服务配置子文件</span>


URI: ApiURI/api/v1/{namespace}/configs/2/items

Method: POST

**参数说明**

- namespace: 镜像所属租户   必须字段



**请求**

- ApiURI/api/v1/huangjia/configs

```
{
  "name":"config",
  "content":"{\"test\":\"huangjia\"}"
}
```

**响应**

```
{
    "apiversion": "v1",
    "status": "201",
    "data": "ok"
}
```


#### <span id="16.3">删除服务配置</span>


URI: ApiURI/api/v1/{namespace}/configs/2/items

Method: DELETE

**参数说明**

- namespace: 镜像所属租户   必须字段



**请求**

- ApiURI/api/v1/huangjia/configs/2

```
{
  "name":"config",
  "content":"{\"test\":\"huangjia\"}"
}
```

**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": "ok"
}
```


#### <span id="16.4">删除服务配置子文件</span>


URI: ApiURI/api/v1/{namespace}/configs/2/items/1

Method: DELETE

**参数说明**

- namespace: 镜像所属租户   必须字段



**请求**

- ApiURI/api/v1/huangjia/configs



**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": "ok"
}
```


#### <span id="16.5">查询服务配置</span>


URI: ApiURI/api/v1/{namespace}/configs?pageNum=0&pageCnt=10

Method: GET

**参数说明**

- namespace: 应用所属租户   必须字段
- pageCnt: 分页查询每页大小  必须字段
- pageNum: 分页查询页码  必须字段
- name: 应用名称  可选字段，当不传时，默认查询当前命名空间下的所有的应用，传入值时，以该值模糊查询结果



**请求**

- ApiURI/api/v1/huangjia/configs?pageNum=0&pageCnt=10


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "configs": [
            {
                "id": 2,
                "createAt": "2017-06-14T01:39:50+08:00",
                "namespace": "huangjia",
                "name": "test",
                "items": [
                    {
                        "id": 3,
                        "createAt": "2017-06-14T01:40:27+08:00",
                        "name": "config",
                        "content": "{\"test\":\"huangjia\"}",
                        "ServiceConfigId": 0,
                        "ConfigId": 2
                    }
                ]
            },
            {
                "id": 1,
                "createAt": "2017-06-14T01:35:02+08:00",
                "namespace": "huangjia",
                "name": "nignx",
                "items": [
                    {
                        "id": 2,
                        "createAt": "2017-06-14T01:36:00+08:00",
                        "name": "config",
                        "content": "{\"test\":\"huangjia\"}",
                        "ServiceConfigId": 0,
                        "ConfigId": 1
                    }
                ]
            }
        ],
        "total": 2
    }
}
```