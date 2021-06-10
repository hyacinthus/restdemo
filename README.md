# 用 Golang 快速开发 RESTful API

这是我 2018 年一个演讲上用的 demo，用 golang 实现一个五脏俱全的 API 项目。

计划大概写一篇博客做说明，但是还没抽出空来写，可以看当时的演讲 PPT ，大概只是个
提纲。项目根目录有个 PDF 文件就是。

2021年6月10日：

我看到偶尔还有同学通过PPT或者视频来到这个repo，这是三年前的代码，虽然各种原理和写法可以对着演讲参考，
但是代码本身和它的依赖都已经过时了。大家看个意思就好，这些年来，我也有很多新的经验，它都被反映在了我新的项目中。

这里有一个开源项目是我一直维护着的，代码也比较简洁，大家可以去看这个项目学习：

[https://github.com/hack-fan/skadi](https://github.com/hack-fan/skadi)

## 使用方法

只要安装了 docker ，并开启了 Swarm 模式，在项目目录执行:

```bash
docker stack deploy -c docker-compose.yml demo
```

如果安装了 make ， 也可以直接 `make up`，会自动帮你执行以上语句。

然后访问 your-host:1328/swagger/index.html 可以看到所有的 api

如果没有开放 swarm 模式，也可以自行改改 compose 文件用 docker-compose 启动.

## 关于项目组织

开发微服务的时候，建议不用建太多文件夹来横向分层。尽量把同一个实体的模型和业务逻
辑全放一个文件里其实是最利于代码维护的。因为微服务项目每个服务维护的实体其实很少
，所以并不会很乱。

现在文件看起来比较乱是因为有太多的公用模块，业务逻辑只有 `note.go` 一个文件。但
是实际上，真的写微服务的时候，最好做一个自己的公用 package，剩下的文件大部分就都
是业务文件了。可以参考我的 [x](https://github.com/hyacinthus/x) 和
[ske](https://github.com/hyacinthus/ske) 这两个 demo

## 一些更新

还没有体现在这个 demo 里

- uuid 库我现在已经弃用，更喜欢用 [xid](https://github.com/rs/xid)
- 日志库 zap 比 logrus 更好
- swag 其实可以把展示统一维护在项目之外，在 ci 的时候生成 json 文件，业务代码里
  只有注释，不用提供 swagger 的 endpoint ，这样是侵入最小的。
