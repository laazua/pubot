### pubot

- 创建管理员
```bash
# 首次运行新建用户(构建前将代码中AuthMw中间件先去除,新建用户后再加上)
curl -XPOST http://127.0.0.1:7777/api/user \
  -H "Content-Type: application/json" \
  -d '{"username":"admin", "password":"123456", "role":"admin"}'
```

- 任务模板示例
```yaml
# YAML 示例
name: demo1
build:
  - if [ ! -d pubot-web ];then git clone git@github.com:laazua/pubot-web.git;fi
  - cd pubot-web
  - npm install && npm run build
deploy:
  platform: linux
  run:
    - echo run1 && sleep 4
    - echo run2 && sleep 9
```