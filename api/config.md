# Banana 配置说明

## 完整配置样例

```yaml
listen: 127.0.0.1:8080  # 监听地址
root: /var/share/banana # 共享目录
size: 100M              # 空间大小
users: # 用户列表
  anonymous: # 匿名用户
    expire: 2019-12-31 23:59:59 # 过期时间
    privilege: [download]
  sdjdd: # 用户名
    password: secret # 密码
    privilege:       # 权限
      - upload
      - download
      - delete
```

## 字段说明

### listen

banana 监听的地址，默认为 `0.0.0.0:8080` 。

### root

共享目录，必填。

### size

共享空间大小，支持浮点数，单位为字节。可添加 K/M/G 后缀将单位设置为千字节、兆字节和吉字节。

### users

用户列表，键为用户名。可选属性如下：

- expire: 过期时间，格式为 `YYYY-mm-dd HH:MM:ss`。未设置即为永不过期。
- password: 密码。
- privilege: 用户权限，可设置 `download`/`upload`/`delete` ，表示下载、上传和删除权限。

设置名为 `anonymous` 的用户即启用匿名用户，可使用 `anonymous` 或空用户名登录，banana 不对匿名用户的密码进行验证。
