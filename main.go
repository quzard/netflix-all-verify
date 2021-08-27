package main

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/listener/http"
	"github.com/axgle/mahonia" //编码转换
	"io"
	"io/ioutil"
	"net"
	h "net/http"
	"net/url"
	"netflix-all-verify/nf"
	"os"
	"path/filepath"
	"time"
)

var proxy constant.Proxy
var proxyUrl = "127.0.0.1:10000"

func getIP() string {
	proxy, _ := url.Parse("http://" + proxyUrl)
	client := h.Client{
		Timeout: 5 * time.Second,
		Transport: &h.Transport{
			// 设置代理
			Proxy: h.ProxyURL(proxy),
		},
	}
	resp, err := client.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func relay(l, r net.Conn) {
	go io.Copy(l, r)
	io.Copy(r, l)
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	//输入订阅链接
	fmt.Println("请输入clash订阅链接(非clash订阅请进行订阅转换)")
	var urlConfig string
	fmt.Scanln(&urlConfig)
	//下载配置信息
	res, err := h.Get(urlConfig)
	if err != nil {
		fmt.Println("clash的订阅链接下载失败！")
		time.Sleep(10 * time.Second)
		return
	}
	defer res.Body.Close()
	//创建配置文件
	f, err := os.OpenFile(exPath+"/config.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("clash的订阅链接下载失败！")
		time.Sleep(10 * time.Second)
		return
	}
	io.Copy(f, res.Body)
	//解析配置信息
	config, err := executor.ParseWithPath(exPath + "/config.yaml")
	if err != nil {
		return
	}
	//开启代理
	in := make(chan constant.ConnContext, 100)
	defer close(in)
	l, err := http.New(proxyUrl, in)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	println("listen at:", l.Address())

	//设置编码
	enc := mahonia.NewDecoder("utf8")

	//监听代理
	go func() {
		for c := range in {
			conn := c
			metadata := conn.Metadata()
			go func() {
				remote, err := proxy.DialContext(context.Background(), metadata)
				if err != nil {
					conn.Conn().Close()
					return
				}
				relay(remote, conn.Conn())
			}()
		}
	}()

	//创建netflix.txt
	f, err = os.OpenFile(exPath+"/netflix.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println("新建netflix.txt失败：", err)
	}
	index := 1
	nodes := config.Proxies

	for node, server := range nodes {
		if server.Type() != constant.Shadowsocks && server.Type() != constant.ShadowsocksR && server.Type() != constant.Snell && server.Type() != constant.Socks5 && server.Type() != constant.Http && server.Type() != constant.Vmess && server.Type() != constant.Trojan {
			continue
		}
		proxy = server
		//落地机IP
		ip := getIP()
		str := fmt.Sprintf("%d   节点名: %s ip地址:%s\n", index, node, ip)
		fmt.Print(str)
		//Netflix检测
		_, out := nf.NF("http://" + proxyUrl)
		if out == "" {
			out = "完全不支持Netflix"
		}
		fmt.Println(out)
		fmt.Fprintln(f, enc.ConvertString(str+out))
		index++
	}
}
