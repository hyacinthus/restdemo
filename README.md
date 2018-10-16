# 用 Golang 快速开发 RESTful API

## 使用方法

只要安装了 docker ，在项目目录执行:

```bash
docker stack deploy -c docker-compose.yml demo
```

如果安装了make ， 也可以直接 `make up`，会自动帮你执行以上语句。

然后访问 your-host:1328/swagger/index.html
可以看到所有的 api

我也在线部署了一份，[点击这里查看](https://demo.crandom.com/swagger/index.html)
