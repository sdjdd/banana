# Banana API

## 搜寻路径

banana 会在工作目录寻找配置文件 `banana.yml` 和静态文件目录 `static` ，若工作目录中不存在则在程序目录下寻找。

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
403|101|当前用户没有上传权限
400|000|文件名 `filename` 为空
400|000|上一级目录不存在
400|000|上一级不是目录
400|000|文件或目录已存在
400|000|空间不足

### 创建目录

与上传文件相同，只是将 query param 中的 type 设置为 `dir`

**URL**: `http://banana-host/fs/:dirname?type=dir`

**METHOD**: POST

**ERROR**:

HTTP status|code|说明
-|-|-
403|101|当前用户没有上传权限
400|000|目录名 `dirname` 为空
400|000|上一级目录不存在
400|000|上一级不是目录
400|000|文件或目录已存在

## 静态文件

banana 提供静态文件服务。会优先在 `static` 目录下寻找与请求匹配的静态文件，供 UI 实现使用。
