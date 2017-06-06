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
  - [获取应用的pod](#10.10)
  - [获取应用pod的事件](#10.11)
  - [获取应用pod的cpu实用情况](#10.12)
  - [获取应用pod的内存实用情况](#10.13)
* **[docker-build组件模块](#11)**
  - [ 构建应用](#11.1)
* **[registry组件模块](#12)**
  - [获取镜像列表](#12.1)

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

URI: ApiURI/{namespace}/apps

Method: GET

**请求**

- ApiURI/{namespace}/apps


**说明**：lessee 表示租户的意思，查询应用其实就是查询当前租户下的所有的应用，租户对应到k8s中的namespace

**响应**

- HTTP Status: 201;
- JSON:

```text
{
  "apiversion": "v1",
  "status": "200",
  "data": [
    {
      "id": 2,
      "createAt": "2017-06-06T01:01:15+08:00",
      "nmae": "nginx",
      "nameSpace": "huangjia",
      "description": "this is a test nginx",
      "serviceCount": 1,
      "intanceCount": 1,
      "external": "http://192.168.99.109:30976",
      "services": [
        {
          "id": 2,
          "createAt": "2017-06-06T01:01:15+08:00",
          "name": "nginx-test",
          "image": "nginx",
          "instanceCount": 1,
          "external": "http://192.168.99.109:30976",
          "loadbalanceIp": "192.168.99.109",
          "Config": null,
          "containers": [
            {
              "id": 2,
              "createAt": "2017-06-06T01:01:15+08:00",
              "name": "test-1-123213",
              "image": "nginx",
              "Config": {
                "id": 2,
                "createAt": "2017-06-06T01:01:15+08:00",
                "base": {
                  "id": 1,
                  "createAt": "2017-06-06T00:38:05+08:00",
                  "cpu": "12Mi",
                  "memory": "12m",
                  "type": 1,
                  "ServiceConfigId": 2
                },
                "config": {
                  "id": 1,
                  "createAt": "2017-06-06T00:38:05+08:00",
                  "name": "nginx-test",
                  "content": "{\"appName\":\"nginx\"}",
                  "containerPath": "/opt",
                  "ServiceConfigId": 2
                },
                "super": {
                  "id": 1,
                  "createAt": "2017-06-06T00:38:05+08:00",
                  "ServiceConfigId": 2
                },
                "ContainerId": 2
              },
              "ServiceId": 2
            }
          ],
          "appId": 2
        }
      ]
    }
  ]
}
```

### <span id="10.2">创建应用</span>

部署应用。

URI: ApiURI/{namespace}/apps

Method: POST

**请求**

- ApiURI/apps 
- JSON
```text
{
         "nmae": "nginx",
         "nameSpace": "huangjia",
         "description": "this is a test nginx",
         "serviceCount": 1,
         "intanceCount": 1,
         "external": "",
         "services": [
             {
                 "name": "nginx-test",
                 "image": "nginx",
                 "instanceCount": 1,
                 "loadbalanceIp": "192.168.99.109",
                 "Config": {
                     "base": {
                         "cpu": "12Mi",
                         "memory": "12m",
                         "type": 1,
                         "volumes": null,
                         "ServiceConfigId": 0
                     },
                     "config": {
                         "name": "nginx-test",
                         "content": "{\"appName\":\"nginx\"}",
                         "containerPath": "/opt",
                         "ServiceConfigId": 0
                     },
                     "super": {
                         "envs": [
                             {
                                 "key": "test",
                                 "val": "1",
                                 "SuperConfigId": 0
                             }
                         ],
                         "ports": [
                             {
                                 "containerPort": 8080,
                                 "servicePort": 8080,
                                 "protocol": "TCP",
                                 "SuperConfigId": 0
                             }
                         ]
                     }
                 }
             }
         ]
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
  "data":"ok"
}
```

### <span id="10.3">启动应用</span>

启动应用

URI: ApiURI/{namespace}/apps/{id}/start

Method: PATCH

**请求**

- ApiURI/{namespace}/apps/{id}/start

```
- 说明：


**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": "200",
  "data":"ok"
}
```

### <span id="10.4">停止应用</span>

停止应用。

URI: ApiURI/{namespace}/apps/{id}/stop

Method: PATCH

**请求**

- ApiURI/{namespace}/apps/{id}/stop

```
- 说明：


**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": "200",
  "data":"ok"
}
```

### <span id="10.5">删除应用</span>

删除应用。

- ApiURI/{namespace}/apps/{id}

Method: DELETE

**请求**



**响应**

- HTTP Status: 204;
- JSON:

```text
{
  "apiversion": "v1",
  "status": "204",
  "data": "ok"
}
```

### <span id="10.6">弹性伸缩</span>

弹性伸缩。

URI: ApiURI/{namespace}/apps/{id}/scale

Method: PUT

**请求**

{
         "nmae": "nginx",
         "nameSpace": "huangjia",
         "description": "this is a test nginx",
         "serviceCount": 1,
         "intanceCount": 2,
         "external": "",
         "services": [
             {
                 "name": "nginx-test",
                 "image": "nginx",
                 "instanceCount": 2,
                 "loadbalanceIp": "192.168.99.109",
                 "Config": {
                     "base": {
                         "cpu": "12Mi",
                         "memory": "12m",
                         "type": 1,
                         "volumes": null,
                         "ServiceConfigId": 0
                     },
                     "config": {
                         "name": "nginx-test",
                         "content": "{\"appName\":\"nginx\"}",
                         "containerPath": "/opt",
                         "ServiceConfigId": 0
                     },
                     "super": {
                         "envs": [
                             {
                                 "key": "test",
                                 "val": "1",
                                 "SuperConfigId": 0
                             }
                         ],
                         "ports": [
                             {
                                 "containerPort": 8080,
                                 "servicePort": 8080,
                                 "protocol": "TCP",
                                 "SuperConfigId": 0
                             }
                         ]
                     }
                 }
             }
         ]
     }

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": "200",
  "data":"ok"
}
```

### <span id="10.7">灰度升级</span>

灰度升级。

URI: ApiURI/{namespace}/apps/{id}/roll

Method: PUT

**请求**

{
         "nmae": "nginx",
         "nameSpace": "huangjia",
         "description": "this is a test nginx",
         "serviceCount": 1,
         "intanceCount": 2,
         "external": "",
         "services": [
             {
                 "name": "nginx-test",
                 "image": "nginx:1.8",
                 "instanceCount": 2,
                 "loadbalanceIp": "192.168.99.109",
                 "Config": {
                     "base": {
                         "cpu": "12Mi",
                         "memory": "12m",
                         "type": 1,
                         "volumes": null,
                         "ServiceConfigId": 0
                     },
                     "config": {
                         "name": "nginx-test",
                         "content": "{\"appName\":\"nginx\"}",
                         "containerPath": "/opt",
                         "ServiceConfigId": 0
                     },
                     "super": {
                         "envs": [
                             {
                                 "key": "test",
                                 "val": "1",
                                 "SuperConfigId": 0
                             }
                         ],
                         "ports": [
                             {
                                 "containerPort": 8080,
                                 "servicePort": 8080,
                                 "protocol": "TCP",
                                 "SuperConfigId": 0
                             }
                         ]
                     }
                 }
             }
         ]
     }

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": "200",
  "data":"ok"
}

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
- 说明：暂时不打算提供该接口


### <span id="10.9">动态扩容</span>

动态扩容。

URI: ApiURI/{namespace}/apps/{id}/roll

Method: PUT

**请求**

{
         "nmae": "nginx",
         "nameSpace": "huangjia",
         "description": "this is a test nginx",
         "serviceCount": 1,
         "intanceCount": 2,
         "external": "",
         "services": [
             {
                 "name": "nginx-test",
                 "image": "nginx:1.8",
                 "instanceCount": 2,
                 "loadbalanceIp": "192.168.99.109",
                 "Config": {
                     "base": {
                         "cpu": "24Mi",
                         "memory": "24m",
                         "type": 1,
                         "volumes": null,
                         "ServiceConfigId": 0
                     },
                     "config": {
                         "name": "nginx-test",
                         "content": "{\"appName\":\"nginx\"}",
                         "containerPath": "/opt",
                         "ServiceConfigId": 0
                     },
                     "super": {
                         "envs": [
                             {
                                 "key": "test",
                                 "val": "1",
                                 "SuperConfigId": 0
                             }
                         ],
                         "ports": [
                             {
                                 "containerPort": 8080,
                                 "servicePort": 8080,
                                 "protocol": "TCP",
                                 "SuperConfigId": 0
                             }
                         ]
                     }
                 }
             }
         ]
     }

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": "200",
  "data":"ok"
}
```

### <span id="10.10">获取应用的pod</span>

获取应用的pod。

URI: ApiURI/apps/deployment/pods?appName=test&namespace=huangjia

**参数说明**：
- appName：应用的名称
- namespace：应用的namespace

Method: GET

**请求**
- ApiURI/apps/deployment/pods?appName=test&namespace=huangjia

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "apiversion": "v1",
  "status": "200",
  "data": [
    {
      "metadata": {
        "name": "test-4240827775-7l8g1",
        "generateName": "test-4240827775-",
        "namespace": "huangjia",
        "selfLink": "/api/v1/namespaces/huangjia/pods/test-4240827775-7l8g1",
        "uid": "e5b0387f-449f-11e7-bc89-0800278ff542",
        "resourceVersion": "27368",
        "creationTimestamp": "2017-05-29T18:52:00Z",
        "labels": {
          "name": "test",
          "pod-template-hash": "4240827775"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"huangjia\",\"name\":\"test-4240827775\",\"uid\":\"e5aa72c8-449f-11e7-bc89-0800278ff542\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"27322\"}}\n",
          "name": "test"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "test-4240827775",
            "uid": "e5aa72c8-449f-11e7-bc89-0800278ff542",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-jn556",
            "secret": {
              "secretName": "default-token-jn556",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "test",
            "image": "registry:latest",
            "resources": {
              "limits": {
                "cpu": "128m",
                "memory": "128Mi"
              },
              "requests": {
                "cpu": "128m",
                "memory": "128Mi"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-jn556",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "slave2",
        "securityContext": {},
        "schedulerName": "default-scheduler",
        "tolerations": [
          {
            "key": "node.alpha.kubernetes.io/notReady",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          },
          {
            "key": "node.alpha.kubernetes.io/unreachable",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          }
        ]
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-05-29T18:33:23Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-05-29T18:33:38Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-05-29T18:52:00Z"
          }
        ],
        "hostIP": "192.168.99.112",
        "podIP": "10.244.2.33",
        "startTime": "2017-05-29T18:33:23Z",
        "containerStatuses": [
          {
            "name": "test",
            "state": {
              "running": {
                "startedAt": "2017-05-29T18:33:37Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "docker.io/registry:latest",
            "imageID": "docker-pullable://docker.io/registry@sha256:a3551c422521617e86927c3ff57e05edf086f1648f4d8524633216ca363d06c2",
            "containerID": "docker://d83c025584ff1b4c3b30ae1d77465d040ed63bb14f748e27d5017118bd03ac20"
          }
        ],
        "qosClass": "Guaranteed"
      }
    },
    {
      "metadata": {
        "name": "test-4240827775-z8jpf",
        "generateName": "test-4240827775-",
        "namespace": "huangjia",
        "selfLink": "/api/v1/namespaces/huangjia/pods/test-4240827775-z8jpf",
        "uid": "e5b08d54-449f-11e7-bc89-0800278ff542",
        "resourceVersion": "27364",
        "creationTimestamp": "2017-05-29T18:52:00Z",
        "labels": {
          "name": "test",
          "pod-template-hash": "4240827775"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"huangjia\",\"name\":\"test-4240827775\",\"uid\":\"e5aa72c8-449f-11e7-bc89-0800278ff542\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"27322\"}}\n",
          "name": "test"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "test-4240827775",
            "uid": "e5aa72c8-449f-11e7-bc89-0800278ff542",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-jn556",
            "secret": {
              "secretName": "default-token-jn556",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "test",
            "image": "registry:latest",
            "resources": {
              "limits": {
                "cpu": "128m",
                "memory": "128Mi"
              },
              "requests": {
                "cpu": "128m",
                "memory": "128Mi"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-jn556",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "slave2",
        "securityContext": {},
        "schedulerName": "default-scheduler",
        "tolerations": [
          {
            "key": "node.alpha.kubernetes.io/notReady",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          },
          {
            "key": "node.alpha.kubernetes.io/unreachable",
            "operator": "Exists",
            "effect": "NoExecute",
            "tolerationSeconds": 300
          }
        ]
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-05-29T18:33:23Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-05-29T18:33:37Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-05-29T18:52:00Z"
          }
        ],
        "hostIP": "192.168.99.112",
        "podIP": "10.244.2.32",
        "startTime": "2017-05-29T18:33:23Z",
        "containerStatuses": [
          {
            "name": "test",
            "state": {
              "running": {
                "startedAt": "2017-05-29T18:33:37Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "docker.io/registry:latest",
            "imageID": "docker-pullable://docker.io/registry@sha256:a3551c422521617e86927c3ff57e05edf086f1648f4d8524633216ca363d06c2",
            "containerID": "docker://6ba5de0602e1774211dfe44f9894fc5089b2f99e95c7a87a807396b4b13ba221"
          }
        ],
        "qosClass": "Guaranteed"
      }
    }
  ]
}
```

### <span id="10.11">获取应用的pod事件</span>

获取应用的pod事件。

URI: ApiURI/apps/pods/events?podName=test-4240827775-7l8g1&namespace=huangjia

**参数说明**：
- podName：pod的name
- namespace：pod的namespace

Method: GET

**请求**
- ApiURI/apps/pods/events?podName=test-4240827775-7l8g1&namespace=huangjia

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "apiversion": "v1",
  "status": "200",
  "data": [
    {
      "metadata": {
        "name": "test-4240827775-7l8g1.14c32763893b979e",
        "namespace": "huangjia",
        "selfLink": "/api/v1/namespaces/huangjia/events/test-4240827775-7l8g1.14c32763893b979e",
        "uid": "e9bb6aa7-449f-11e7-bc89-0800278ff542",
        "resourceVersion": "27350",
        "creationTimestamp": "2017-05-29T18:52:06Z"
      },
      "involvedObject": {
        "kind": "Pod",
        "namespace": "huangjia",
        "name": "test-4240827775-7l8g1",
        "uid": "e5b0387f-449f-11e7-bc89-0800278ff542",
        "apiVersion": "v1",
        "resourceVersion": "27329",
        "fieldPath": "spec.containers{test}"
      },
      "reason": "Pulled",
      "message": "Container image \"registry:latest\" already present on machine",
      "source": {
        "component": "kubelet",
        "host": "slave2"
      },
      "firstTimestamp": "2017-05-29T18:33:29Z",
      "lastTimestamp": "2017-05-29T18:33:29Z",
      "count": 1,
      "type": "Normal"
    },
    {
      "metadata": {
        "name": "test-4240827775-7l8g1.14c327639280b342",
        "namespace": "huangjia",
        "selfLink": "/api/v1/namespaces/huangjia/events/test-4240827775-7l8g1.14c327639280b342",
        "uid": "e9d316a5-449f-11e7-bc89-0800278ff542",
        "resourceVersion": "27352",
        "creationTimestamp": "2017-05-29T18:52:07Z"
      },
      "involvedObject": {
        "kind": "Pod",
        "namespace": "huangjia",
        "name": "test-4240827775-7l8g1",
        "uid": "e5b0387f-449f-11e7-bc89-0800278ff542",
        "apiVersion": "v1",
        "resourceVersion": "27329",
        "fieldPath": "spec.containers{test}"
      },
      "reason": "Created",
      "message": "Created container with id d83c025584ff1b4c3b30ae1d77465d040ed63bb14f748e27d5017118bd03ac20",
      "source": {
        "component": "kubelet",
        "host": "slave2"
      },
      "firstTimestamp": "2017-05-29T18:33:29Z",
      "lastTimestamp": "2017-05-29T18:33:29Z",
      "count": 1,
      "type": "Normal"
    },
    {
      "metadata": {
        "name": "test-4240827775-7l8g1.14c32765602e94d4",
        "namespace": "huangjia",
        "selfLink": "/api/v1/namespaces/huangjia/events/test-4240827775-7l8g1.14c32765602e94d4",
        "uid": "ee710667-449f-11e7-bc89-0800278ff542",
        "resourceVersion": "27363",
        "creationTimestamp": "2017-05-29T18:52:14Z"
      },
      "involvedObject": {
        "kind": "Pod",
        "namespace": "huangjia",
        "name": "test-4240827775-7l8g1",
        "uid": "e5b0387f-449f-11e7-bc89-0800278ff542",
        "apiVersion": "v1",
        "resourceVersion": "27329",
        "fieldPath": "spec.containers{test}"
      },
      "reason": "Started",
      "message": "Started container with id d83c025584ff1b4c3b30ae1d77465d040ed63bb14f748e27d5017118bd03ac20",
      "source": {
        "component": "kubelet",
        "host": "slave2"
      },
      "firstTimestamp": "2017-05-29T18:33:37Z",
      "lastTimestamp": "2017-05-29T18:33:37Z",
      "count": 1,
      "type": "Normal"
    },
    {
      "metadata": {
        "name": "test-4240827775-7l8g1.14c328662114d2ab",
        "namespace": "huangjia",
        "selfLink": "/api/v1/namespaces/huangjia/events/test-4240827775-7l8g1.14c328662114d2ab",
        "uid": "e5c72e7e-449f-11e7-bc89-0800278ff542",
        "resourceVersion": "27339",
        "creationTimestamp": "2017-05-29T18:52:00Z"
      },
      "involvedObject": {
        "kind": "Pod",
        "namespace": "huangjia",
        "name": "test-4240827775-7l8g1",
        "uid": "e5b0387f-449f-11e7-bc89-0800278ff542",
        "apiVersion": "v1",
        "resourceVersion": "27327"
      },
      "reason": "Scheduled",
      "message": "Successfully assigned test-4240827775-7l8g1 to slave2",
      "source": {
        "component": "default-scheduler"
      },
      "firstTimestamp": "2017-05-29T18:52:00Z",
      "lastTimestamp": "2017-05-29T18:52:00Z",
      "count": 1,
      "type": "Normal"
    }
  ]
}

```
pod 指标：

```
[
  "memory/major_page_faults_rate",
  "network/rx_errors_rate",
  "cpu/usage",
  "network/rx_rate",
  "memory/cache",
  "network/tx_errors",
  "uptime",
  "cpu/limit",
  "network/rx_errors",
  "memory/request",
  "memory/page_faults_rate",
  "cpu/request",
  "memory/major_page_faults",
  "cpu/usage_rate",
  "network/tx_errors_rate",
  "memory/usage",
  "memory/limit",
  "network/rx",
  "network/tx",
  "memory/working_set",
  "memory/page_faults",
  "network/tx_rate",
  "memory/rss"
]
```


### <span id="10.12">获取应用的pod的cpu使用情况</span>

获取应用的pod的cpu使用情况

URI: ApiURI/api/v1/apps/pods/metrics?podName=test-4240827775-1phdj&namespace=huangjia&metricName=cpu/usage

**参数说明**：
- podName：pod的name
- namespace：pod的namespace
- metricName：指标名称

Method: GET

**请求**
- ApiURI/api/v1/apps/pods/metrics?podName=test-4240827775-1phdj&namespace=huangjia&metricName=cpu/usage

```

**响应**

- HTTP Status: 200;
- JSON:

```text

{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "latestTimestamp": "2017-05-29T23:12:00Z",
    "metrics": [
      {
        "timestamp": "2017-05-29T23:10:00Z",
        "value": 0
      },
      {
        "timestamp": "2017-05-29T23:11:00Z",
        "value": 0
      },
      {
        "timestamp": "2017-05-29T23:12:00Z",
        "value": 0
      }
    ]
  }
}
```

### <span id="10.13">获取应用的pod的内存使用情况</span>

获取应用的pod的内存使用情况

URI: ApiURI/api/v1/apps/pods/metrics?podName=test-4240827775-1phdj&namespace=huangjia&metricName=memory/usage

**参数说明**：
- podName：pod的name
- namespace：pod的namespace
- metricName：指标名称

Method: GET

**请求**
- ApiURI/api/v1/apps/pods/metrics?podName=test-4240827775-1phdj&namespace=huangjia&metricName=memory/usage

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "latestTimestamp": "2017-05-29T23:13:00Z",
    "metrics": [
      {
        "timestamp": "2017-05-29T22:59:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:00:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:01:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:02:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:03:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:04:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:05:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:06:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:07:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:08:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:09:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:10:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:11:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:12:00Z",
        "value": 16371712
      },
      {
        "timestamp": "2017-05-29T23:13:00Z",
        "value": 16371712
      }
    ]
  }
}
```

### <span id="10.14">获取应用的pod的网络使用情况</span>

获取应用的pod的网络使用情况

URI: ApiURI/api/v1/apps/pods/metrics?podName=test-4240827775-1phdj&namespace=huangjia&metricName=network/rx_rate

**参数说明**：
- podName：pod的name
- namespace：pod的namespace
- metricName：指标名称

Method: GET

**请求**
- ApiURI/api/v1/apps/pods/metrics?podName=test-4240827775-1phdj&namespace=huangjia&metricName=network/rx_rate

```

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "apiversion": "v1",
  "status": "200",
  "data": {
    "latestTimestamp": "2017-05-29T23:18:00Z",
    "metrics": [
      {
        "timestamp": "2017-05-29T23:16:00Z",
        "value": 0
      },
      {
        "timestamp": "2017-05-29T23:17:00Z",
        "value": 0
      },
      {
        "timestamp": "2017-05-29T23:18:00Z",
        "value": 0
      }
    ]
  }
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


- registry组件api

## <span id="12">registry组件模块</span>
---

### <span id="12.1">获取镜像列表</span>

获取镜像列表。

URI: ApiURI/images

Method: GET

**请求**

- ApiURI/images?name=imageName&&pageNum=0&&pageCnt=10

```
**参数说明**：
- name:镜像名称，没有传入默认查询全部的镜像
- pageNum:分页页码，
- pageCnt:分页每页行数

**响应**

- HTTP Status: 200;
- JSON:

```text
{
     "apiversion": "v1",
     "status": "200",
     "date": {
         "images": [
             {
                 "name": "busybox",
                 "tagLen": 1,
                 "manifest": [
                     {
                         "name": "busybox",
                         "tag": "latest",
                         "architecture": "amd64",
                         "os": "linux",
                         "author": "",
                         "id": "21bd05c98a33998aba2cea975e0fcdc4c8b051070b70ed36f28c0bc55bcdacb6",
                         "parent": "86549330fef190e649817430dfaba05934d46b450fe2004cc1e2afc36587054c",
                         "created": "2017-03-09T18:28:04.586987216Z",
                         "docker_version": "1.12.6",
                         "pull": "docker pull http://10.4.94.98:5000/busybox:latest"
                     }
                 ]
             },
             {
                 "name": "huangjia/pause-amd64",
                 "tagLen": 1,
                 "manifest": [
                     {
                         "name": "huangjia/pause-amd64",
                         "tag": "3.0",
                         "architecture": "amd64",
                         "os": "linux",
                         "author": "",
                         "id": "bebc58b30ecc163fe8a56301e2fff15de40d225663c95134ff5f242ebb8a516e",
                         "parent": "ce598428d6bb655a3c88cf9d60d4e728bffb94d82f578fbfa30236bea68c92a0",
                         "created": "2016-05-04T06:26:41.522308365Z",
                         "docker_version": "1.9.1",
                         "pull": "docker pull http://10.4.94.98:5000/huangjia/pause-amd64:3.0"
                     }
                 ]
             },
             {
                 "name": "pause-amd64",
                 "tagLen": 2,
                 "manifest": [
                     {
                         "name": "pause-amd64",
                         "tag": "3.1",
                         "architecture": "amd64",
                         "os": "linux",
                         "author": "",
                         "id": "bebc58b30ecc163fe8a56301e2fff15de40d225663c95134ff5f242ebb8a516e",
                         "parent": "ce598428d6bb655a3c88cf9d60d4e728bffb94d82f578fbfa30236bea68c92a0",
                         "created": "2016-05-04T06:26:41.522308365Z",
                         "docker_version": "1.9.1",
                         "pull": "docker pull http://10.4.94.98:5000/pause-amd64:3.1"
                     },
                     {
                         "name": "pause-amd64",
                         "tag": "3.0",
                         "architecture": "amd64",
                         "os": "linux",
                         "author": "",
                         "id": "bebc58b30ecc163fe8a56301e2fff15de40d225663c95134ff5f242ebb8a516e",
                         "parent": "ce598428d6bb655a3c88cf9d60d4e728bffb94d82f578fbfa30236bea68c92a0",
                         "created": "2016-05-04T06:26:41.522308365Z",
                         "docker_version": "1.9.1",
                         "pull": "docker pull http://10.4.94.98:5000/pause-amd64:3.0"
                     }
                 ]
             }
         ],
         "total": 3
     }
 }
```

-- 说明：
images:镜像列表

{
   "name": "busybox",
   "tagLen": 1,
   "manifest": [
       {
           "name": "busybox",
           "tag": "latest",
           "architecture": "amd64",
           "os": "linux",
           "author": "",
           "id": "21bd05c98a33998aba2cea975e0fcdc4c8b051070b70ed36f28c0bc55bcdacb6",
           "parent": "86549330fef190e649817430dfaba05934d46b450fe2004cc1e2afc36587054c",
           "created": "2017-03-09T18:28:04.586987216Z",
           "docker_version": "1.12.6",
           "pull": "docker pull http://10.4.94.98:5000/busybox:latest"
       }
   ]
}

-- nmae:镜像名称
-- tagLen:镜像版本个数
-- manifest:镜像的详细信息



