package main

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/listener/http"
	"github.com/axgle/mahonia" //编码转换
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"net"
	h "net/http"
	"net/url"
	"netflix-all-verify/nf"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var proxy constant.Proxy
var proxyUrl = "127.0.0.1:"
var exPath string

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

// 获取可用端口
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "127.0.0.1"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil

}

func downloadConfig() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath = filepath.Dir(ex)
	fmt.Println(exPath)
	//输入订阅链接
	fmt.Println("请输入clash订阅链接(非clash订阅请进行订阅转换)")
	var urlConfig string
	_, err = fmt.Scanln(&urlConfig)
	if err != nil {
		panic(err)
	}
	//下载配置信息
	res, err := h.Get(urlConfig)
	if err != nil {
		fmt.Println("clash的订阅链接下载失败！")
		time.Sleep(10 * time.Second)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	//创建配置文件
	f, err := os.OpenFile(exPath+"/config.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	if err != nil {
		fmt.Println("clash的订阅链接下载失败！")
		time.Sleep(10 * time.Second)
		return
	}
	_, err = io.Copy(f, res.Body)
	if err != nil {
		panic(err)
	}
}

func main() {
	downloadConfig()

	//解析配置信息
	config, err := executor.ParseWithPath(exPath + "/config.yaml")
	if err != nil {
		return
	}
	//获取端口
	port, _ := GetAvailablePort()
	proxyUrl += strconv.Itoa(port)
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
	f, err := os.OpenFile(exPath+"/netflix.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println("新建netflix.txt失败：", err)
	}

	//创建excel
	excel := excelize.NewFile()
	excel.SetCellValue("Sheet1", "A1", "节点名")
	excel.SetCellValue("Sheet1", "B1", "ip地址")
	excel.SetCellValue("Sheet1", "C1", "复用次数")
	excel.SetCellValue("Sheet1", "D1", "是否解锁")
	excel.SetCellValue("Sheet1", "E1", "详细说明")

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
		ok, out := nf.NF("http://" + proxyUrl)
		if out == "" {
			out = "完全不支持Netflix"
		}
		fmt.Println(out)
		fmt.Fprintln(f, enc.ConvertString(str+out))

		excel.SetCellValue("Sheet1", "A"+strconv.Itoa(index+1), node)
		excel.SetCellValue("Sheet1", "B"+strconv.Itoa(index+1), ip)
		if ip != "" {
			excel.SetCellFormula("Sheet1", "C"+strconv.Itoa(index+1), "= COUNTIF(B:B,B"+strconv.Itoa(index+1)+")")
		}
		excel.SetCellValue("Sheet1", "D"+strconv.Itoa(index+1), ok)
		excel.SetCellValue("Sheet1", "E"+strconv.Itoa(index+1), out)

		index++
	}

	if err := excel.SaveAs("Netflix.xlsx"); err != nil {
		fmt.Println(err)
	}
}
