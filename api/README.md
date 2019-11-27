# Banana API

## 认证

Banana 使用 HTTP Basic authentication 进行用户认证，关于用户配置请查阅[配置文件说明](config.md)。

## 错误

接口调用失败时会返回非 200 状态码，同时 body 为 JSON 格式的字符串。

`code` 字段为错误代码，`message` 字段为简要的错误信息。

部分可通过 HTTP status 明确表达信息的接口错误值不包含 `code` 字段。

```json
{
    "code": 123,
    "message": "something wrong"
}
```

## 操作

### 上传文件/创建目录

根据 query string 中的 `dir` 判断是否为创建目录操作。

**URL** : `http://banana-host/fs/:name?dir=(true|false)`

**METHOD** : `POST`

**ERROR** :

HTTP status|code|说明
-|-|-
403|-|当前用户没有上传权限
400|1|`name` 为空
400|2|上一级目录不存在
400|3|上一级不是目录
400|4|文件或目录已存在
400|5|空间不足

### 浏览目录/下载文件

链接指向目录时，返回 Object 数组，属性如下：

key|类型|说明
-|-|-
`name`|String|文件名
`isDir`|Boolean|是否为目录
`size`|Number|文件大小（目录的该属性为0）
`modTime`|String|RFC3339 格式的修改时间

链接指向文件时返回文件内容。

**URL** : `http://banana-host/fs/:dirname`

**METHOD** : `GET`

**ERROR** :

HTTP status|说明
-|-|-
403|当前用户没有下载权限
404|文件或目录不存在

### 删除文件/目录

**URL** : `http://banana-host/fs/:name`

**METHOD** : `DELETE`

**ERROR** :

HTTP status|code|说明
-|-|-
403|-|当前用户没有删除权限
404|-|文件或目录不存在
400|101|`name` 为空

### 移动文件

需提供文件原路径 `from` 和目的路径 `to`，支持 form-data、x-www-form-urlencoded 和 json 格式的请求数据。

用户需要同时具备上传和删除权限才能移动文件。

**URL** : `http://banana-host/mv`

**METHOD** : POST

**ERROR** :

HTTP status|code|说明
-|-|-
403|-|当前用户没有上传和删除权限
404|-|源文件不存在
400|104|目的文件已存在
