# Banana API

## 搜寻路径

banana 会在工作目录寻找配置文件 `banana.yml` 和静态文件目录 `ui` ，若工作目录中不存在则在程序目录下寻找。

## 错误

接口调用失败时会返回非 200 状态码，同时 body 为 JSON 格式的字符串。

`code` 字段为错误代码，`message` 字段为简要的错误信息。

```json
{
    "code": 123,
    "message": "some message"
}
```

## 操作

### 上传文件

**URL**: `http://banana-host/fs/:filename`

**METHOD**: POST

**ERROR**:

HTTP status|code|说明
-|-|-
500|100|内部错误
403|101|当前用户没有上传权限
400|200|文件名 `filename` 为空
400|201|上一级目录不存在
400|203|上一级不是目录
400|204|文件或目录已存在
400|205|空间不足

### 创建目录

与上传文件相同，只是将 query param 中的 type 设置为 `dir`

**URL**: `http://banana-host/fs/:dirname?type=dir`

**METHOD**: POST

**ERROR**:

HTTP status|code|说明
-|-|-
500|100|内部错误
403|101|当前用户没有上传权限
400|200|目录名 `dirname` 为空
400|201|上一级目录不存在
400|203|上一级不是目录
400|204|文件或目录已存在

### 浏览目录

浏览目录接口返回 JSON 格式的目录下文件信息，格式为 Object 数组。

Object 属性如下：

key|类型|说明
-|-|-
`name`|String|文件名
`isDir`|Boolean|是否为目录
`size`|Number|文件大小（目录的该属性为0）
`modTime`|String|RFC3339 格式的修改时间

**URL**: `http://banana-host/fs/:dirname`

**METHOD**: GET

**ERROR**:

HTTP status|code|说明
-|-|-
500|100|内部错误
400|202|目录不存在


## 静态文件

banana 提供静态文件服务。会在 `ui` 目录下寻找与请求匹配的静态文件，存放前端 UI 实现。
