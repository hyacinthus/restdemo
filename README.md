# 用 Golang 快速开发 RESTful API

这是我 2018 年一个演讲上用的 demo，用 golang 实现一个五脏俱全的 API 项目。

计划大概写一篇博客做说明，但是还没抽出空来写，可以看当时的演讲 PPT ，大概只是个
提纲。项目根目录有个 PDF 文件就是。

## 使用方法

只要安装了 docker ，在项目目录执行:

```bash
docker stack deploy -c docker-compose.yml demo
```

如果安装了 make ， 也可以直接 `make up`，会自动帮你执行以上语句。

然后访问 your-host:1328/swagger/index.html 可以看到所有的 api

我也在线部署了一份，[点击这里查看](https://demo.crandom.com/swagger/index.html)

## 关于项目组织

开发微服务的时候，建议不用弄太多文件夹来分层。尽量把同一个实体的模型和业务逻辑全
放一个文件里其实是比较好维护的。

现在文件看起来比较乱是因为有太多的公用模块，业务逻辑只有 `note.go` 一个文件。但
是实际上，真的写微服务的时候，最好做一个自己的公用 package，剩下的文件大部分就都
是业务文件了。可以参考我的 [x](https://github.com/hyacinthus/x) 和
[ske](https://github.com/hyacinthus/ske) 这两个 demo

## 一些更新

还没有体现在这个 demo 里

- uuid 库我现在已经弃用，更喜欢用 [xid](https://github.com/rs/xid)
- swag 已经弃用，更为成熟的是
  [go-swagger](https://github.com/go-swagger/go-swagger)
