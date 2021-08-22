# netflix-all-verify
NetFlix批量检测,golang编写

Netflix的大规模封锁,网上的Netflix检测脚本只有本地检测的功能,
因而诞生了本项目

本项目基于[netflix-verify](https://github.com/sjlleo/netflix-verify) 修改,增加了对批量检测的功能

## 原理

基于clash for windows,clashx pro等clash客户端软件的external-controller控制端口的控制,全局模式遍历节点并切换,从而达到
检测节点是否解锁Netflix的目的

## 使用方法
1. proxy: clash的代理端口

2. control: clash的控制端口/api port

3. netflix: 是否只显示能解锁的节点, 1-只显示能解锁的节点  0-全显示
端口获取方式

- clash for windows: 通过点击主页的 打开目录或Home Directory, 目录中的config.yaml文件中可得代理端口为mixed-port,
控制端口为external-controller
- clashx pro页面的帮助-端口中可得代理端口为http port, 控制端口为api port

打开clash软件情况下,访问127.0.0.1:control端口,可检验该端口是否正确

```bash
# go运行
go run nf.go -proxy 7890 -control 49506 -netflix 0
# go编译
go build -v -o nf .
./nf -proxy 7890 -control 49506 -netflix 0
```

## 未来工作
由于仓促,该项目高度依赖clash客户端,未来将独立(有空时)

- [ ] 集成clash，不依靠客户端
- [ ] 内置订阅转换,便于不同订阅地址的使用
- [ ] web界面
- [ ] 结果导出图片化
- [ ] 落地机ip检测和测速

## 感谢

1. 感谢 [netflix-verify](https://github.com/sjlleo/netflix-verify)
