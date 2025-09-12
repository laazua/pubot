### pubot

- 创建管理员
```bash
# 首次运行新建用户(构建前将代码中AuthMw中间件先去除,新建用户后再加上)
curl -XPOST http://127.0.0.1:7777/api/user \
  -H "Content-Type: application/json" \
  -d '{"username":"admin", "password":"123456", "role":"admin"}'
```