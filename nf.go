package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var proxyPort = flag.String("proxy", "", "clash的代理端口")
var controlPort = flag.String("control", "", "clash的控制端口")
var netflix = flag.String("netflix", "1", "是否只显示能解锁的节点")

const Netflix = "https://www.netflix.com/title/"

func RequestIP(requrl string, ip string) string {
	if ip == "" {
		return "Error"
	}
	urlValue, err := url.Parse(requrl)
	if err != nil {
		return "Error"
	}
	host := urlValue.Host
	if ip == "" {
		ip = host
	}

	proxyUrl := "http://127.0.0.1:" + *proxyPort
	proxy, _ := url.Parse(proxyUrl)
	netTransport := &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: host},
		Proxy:           http.ProxyURL(proxy),
	}

	newrequrl := strings.Replace(requrl, host, ip, 1)
	client := &http.Client{
		Transport:     netTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       5 * time.Second,
	}
	req, err := http.NewRequest("GET", newrequrl, nil)
	if err != nil {
		//return errors.New(strings.ReplaceAll(err.Error(), newrequrl, requrl))
		return "Error"
	}
	req.Host = host
	req.Header.Set("USER-AGENT", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		//return errors.New(strings.ReplaceAll(err.Error(), newrequrl, requrl))
		return "Error"
	}
	defer resp.Body.Close()

	Header := resp.Header

	if Header["X-Robots-Tag"] != nil {
		if Header["X-Robots-Tag"][0] == "index" {
			return "us"
		}
	}

	if Header["Location"] == nil {
		return "Ban"
	} else {
		return strings.Split(Header["Location"][0], "/")[3]
	}
}

func ParseIP(s string) int {
	ip := net.ParseIP(s)
	if ip == nil {
		return 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return 4
		case ':':
			return 6
		}
	}
	return 0
}

func UnblockTest(MoiveID int, ip string) bool {
	testURL := Netflix + strconv.Itoa(MoiveID)
	reCode := RequestIP(testURL, ip)
	if strings.Contains(reCode, "Ban") {
		return false
	} else {
		return true
	}
}

func FindCountry(Code string) string {
	countryName := []string{"美国", "阿富汗", "奥兰群岛", "阿尔巴尼亚", "阿尔及利亚", "美属萨摩亚", "安道尔", "安哥拉", "安圭拉", "南极洲", "安提瓜和巴布达", "阿根廷", "亚美尼亚", "阿鲁巴", "澳大利亚", "奥地利", "阿塞拜疆", "巴哈马", "巴林", "孟加拉国", "巴巴多斯", "白俄罗斯", "比利时", "伯利兹", "贝宁", "百慕大", "不丹", "玻利维亚", "波黑", "博茨瓦纳", "布维岛", "巴西", "英属印度洋领地", "文莱", "保加利亚", "布基纳法索", "布隆迪", "柬埔寨", "喀麦隆", "加拿大", "佛得角", "开曼群岛", "中非", "乍得", "智利", "中国", "圣诞岛", "科科斯（基林）群岛", "哥伦比亚", "科摩罗", "刚果（布）", "刚果（金）", "库克群岛", "哥斯达黎加", "科特迪瓦", "克罗地亚", "古巴", "塞浦路斯", "捷克", "丹麦", "吉布提", "多米尼克", "多米尼加", "厄瓜多尔", "埃及", "萨尔瓦多", "赤道几内亚", "厄立特里亚", "爱沙尼亚", "埃塞俄比亚", "福克兰群岛（马尔维纳斯）", "法罗群岛", "斐济", "芬兰", "法国", "法属圭亚那", "法属波利尼西亚", "法属南部领地", "加蓬", "冈比亚", "格鲁吉亚", "德国", "加纳", "直布罗陀", "希腊", "格陵兰", "格林纳达", "瓜德罗普", "关岛", "危地马拉", "格恩西岛", "几内亚", "几内亚比绍", "圭亚那", "海地", "赫德岛和麦克唐纳岛", "梵蒂冈", "洪都拉斯", "香港", "匈牙利", "冰岛", "印度", "印度尼西亚", "伊朗", "伊拉克", "爱尔兰", "英国属地曼岛", "以色列", "意大利", "牙买加", "日本", "泽西岛", "约旦", "哈萨克斯坦", "肯尼亚", "基里巴斯", "朝鲜", "韩国", "科威特", "吉尔吉斯斯坦", "老挝", "拉脱维亚", "黎巴嫩", "莱索托", "利比里亚", "利比亚", "列支敦士登", "立陶宛", "卢森堡", "澳门", "前南马其顿", "马达加斯加", "马拉维", "马来西亚", "马尔代夫", "马里", "马耳他", "马绍尔群岛", "马提尼克", "毛利塔尼亚", "毛里求斯", "马约特", "墨西哥", "密克罗尼西亚联邦", "摩尔多瓦", "摩纳哥", "蒙古", "黑山", "蒙特塞拉特", "摩洛哥", "莫桑比克", "缅甸", "纳米比亚", "瑙鲁", "尼泊尔", "荷兰", "荷属安的列斯", "新喀里多尼亚", "新西兰", "尼加拉瓜", "尼日尔", "尼日利亚", "纽埃", "诺福克岛", "北马里亚纳", "挪威", "阿曼", "巴基斯坦", "帕劳", "巴勒斯坦", "巴拿马", "巴布亚新几内亚", "巴拉圭", "秘鲁", "菲律宾", "皮特凯恩", "波兰", "葡萄牙", "波多黎各", "卡塔尔", "留尼汪", "罗马尼亚", "俄罗斯联邦", "卢旺达", "圣赫勒拿", "圣基茨和尼维斯", "圣卢西亚", "圣皮埃尔和密克隆", "圣文森特和格林纳丁斯", "萨摩亚", "圣马力诺", "圣多美和普林西比", "沙特阿拉伯", "塞内加尔", "塞尔维亚", "塞舌尔", "塞拉利昂", "新加坡", "斯洛伐克", "斯洛文尼亚", "所罗门群岛", "索马里", "南非", "南乔治亚岛和南桑德韦奇岛", "西班牙", "斯里兰卡", "苏丹", "苏里南", "斯瓦尔巴岛和扬马延岛", "斯威士兰", "瑞典", "瑞士", "叙利亚", "台湾", "塔吉克斯坦", "坦桑尼亚", "泰国", "东帝汶", "多哥", "托克劳", "汤加", "特立尼达和多巴哥", "突尼斯", "土耳其", "土库曼斯坦", "特克斯和凯科斯群岛", "图瓦卢", "乌干达", "乌克兰", "阿联酋", "英国", "美国本土外小岛屿", "乌拉圭", "乌兹别克斯坦", "瓦努阿图", "委内瑞拉", "越南", "英属维尔京群岛", "美属维尔京群岛", "瓦利斯和富图纳", "西撒哈拉", "也门", "赞比亚", "津巴布韦"}
	countryCode := []string{"us", "af", "ax", "al", "dz", "as", "ad", "ao", "ai", "aq", "ag", "ar", "am", "aw", "au", "at", "az", "bs", "bh", "bd", "bb", "by", "be", "bz", "bj", "bm", "bt", "bo", "ba", "bw", "bv", "br", "io", "bn", "bg", "bf", "bi", "kh", "cm", "ca", "cv", "ky", "cf", "td", "cl", "cn", "cx", "cc", "co", "km", "cg", "cd", "ck", "cr", "ci", "hr", "cu", "cy", "cz", "dk", "dj", "dm", "do", "ec", "eg", "sv", "gq", "er", "ee", "et", "fk", "fo", "fj", "fi", "fr", "gf", "pf", "tf", "ga", "gm", "ge", "de", "gh", "gi", "gr", "gl", "gd", "gp", "gu", "gt", "gg", "gn", "gw", "gy", "ht", "hm", "va", "hn", "hk", "hu", "is", "in", "id", "ir", "iq", "ie", "im", "il", "it", "jm", "jp", "je", "jo", "kz", "ke", "ki", "kp", "kr", "kw", "kg", "la", "lv", "lb", "ls", "lr", "ly", "li", "lt", "lu", "mo", "mk", "mg", "mw", "my", "mv", "ml", "mt", "mh", "mq", "mr", "mu", "yt", "mx", "fm", "md", "mc", "mn", "me", "ms", "ma", "mz", "mm", "na", "nr", "np", "nl", "an", "nc", "nz", "ni", "ne", "ng", "nu", "nf", "mp", "no", "om", "pk", "pw", "ps", "pa", "pg", "py", "pe", "ph", "pn", "pl", "pt", "pr", "qa", "re", "ro", "ru", "rw", "sh", "kn", "lc", "pm", "vc", "ws", "sm", "st", "sa", "sn", "rs", "sc", "sl", "sg", "sk", "si", "sb", "so", "za", "gs", "es", "lk", "sd", "sr", "sj", "sz", "se", "ch", "sy", "tw", "tj", "tz", "th", "tl", "tg", "tk", "to", "tt", "tn", "tr", "tm", "tc", "tv", "ug", "ua", "ae", "gb", "um", "uy", "uz", "vu", "ve", "vn", "vg", "vi", "wf", "eh", "ye", "zm", "zw"}
	for i, v := range countryCode {
		if strings.Contains(Code, v) {
			return countryName[i]
		}
	}
	return Code
}

func nf() (bool, string) {
	var ipv4 string

	var areaAvailableID = 80018499
	var SelfMadeAvailableID = 80197526
	var NonSelfMadeAvailableID = 70143836

	dns := "www.netflix.com"

	flag.Parse()

	// 解析ip地址
	ns, err := net.LookupHost(dns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s", err.Error())
		return false, ""
	}

	switch {
	case len(ns) != 0:
		for _, n := range ns {
			if ParseIP(n) == 4 {
				ipv4 = n
			}
		}

	}
	// 拼接非自制剧的URL
	testURL := Netflix + strconv.Itoa(NonSelfMadeAvailableID)
	ipv4CountryCode := RequestIP(testURL, ipv4)

	/***
	 * 检查CountryCode返回值:
	 * Error 代表该网络访问失败
	 * Ban 代表无法解锁这个ID种类的影片
	 * 此处如果显示值不为Error则都应该继续检测
	***/
	if !strings.Contains(ipv4CountryCode, "Error") {
		//开启换行信号,在IPv4检测完毕后换行
		//检测是否为自定义测试模式
		//如果反馈为Ban，那么进一步检测是否支持Netflix地区解锁
		if strings.Contains(ipv4CountryCode, "Ban") {
			//检测该IP所在的地区是否支持NF
			if UnblockTest(areaAvailableID, ipv4) {
				//所在地区支持NF
				//检测是否支持自制
				if UnblockTest(SelfMadeAvailableID, ipv4) {
					testURL2 := Netflix + strconv.Itoa(SelfMadeAvailableID)
					ipv4CountryCode2 := RequestIP(testURL2, ipv4)
					ip := "NF库识别的IP地域信息：" + FindCountry(ipv4CountryCode2) + "区(" + strings.ToUpper(strings.Split(ipv4CountryCode2, "-")[0]) + ") NetFlix 非原生IP"
					return false, "您的出口IP不能解锁Netflix，仅支持自制剧的观看\n" + ip
					//支持自制剧
				} else {
					//不支持自制剧
					return false, "不支持解锁带有强版权的自制剧"
				}
			} else {
				//所在地区不支持NF
				return false, ""
			}
		} else {
			//如果支持非自制剧的解锁，则直接跳过自制剧的解锁
			ip := "原生IP地域解锁信息：" + FindCountry(ipv4CountryCode) + "区(" + strings.ToUpper(strings.Split(ipv4CountryCode, "-")[0]) + ") NetFlix 原生IP"
			return true, "您的出口IP完整解锁Netflix，支持非自制剧的观看\n" + ip
		}
	}
	return false, ""
}

func setGlobal() bool {
	type clashMode struct {
		Mode string `json:"mode"`
	}
	jsonStr := clashMode{
		Mode: "Global",
	}

	dataJson, _ := json.Marshal(jsonStr)
	urlGlobal := "http://127.0.0.1:" + *controlPort + "/configs"

	req, _ := http.NewRequest(http.MethodPatch, urlGlobal, bytes.NewBuffer(dataJson))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	if resp.StatusCode != 204 {
		return false
	}
	return true
}
func getNodes() []string {
	type clashProxies struct {
		Proxies struct {
			Global struct {
				All []string `json:"all"`
			} `json:"GLOBAL"`
		}
	}
	urlGet := "http://127.0.0.1:" + *controlPort + "/proxies"
	resp, err := http.Get(urlGet)
	if err != nil {
		return nil
	}

	if resp.StatusCode != 200 {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("clash 节点获取失败")
			time.Sleep(10 * time.Second)
			os.Exit(3)
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	var nodes clashProxies
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		return nil
	}
	return nodes.Proxies.Global.All
}
func setNode(node string) {
	type nodePost struct {
		Name string `json:"name"`
	}
	jsonStr := nodePost{
		Name: node,
	}
	client := &http.Client{}
	dataJson, _ := json.Marshal(jsonStr)
	urlSet := "http://127.0.0.1:" + *controlPort + "/proxies/GLOBAL"
	req, _ := http.NewRequest(http.MethodPut, urlSet, bytes.NewBuffer(dataJson))
	req.Header.Set("Content-Type", "application/json")
	_, err := client.Do(req)
	if err != nil {
		fmt.Println("clash 设置节点失败")
		time.Sleep(10 * time.Second)
		os.Exit(3)
	}
}
func main() {
	flag.Parse()
	fmt.Println("检测开始")
	//设置为全局
	status := setGlobal()
	if !status {
		fmt.Println("端口错误或者clash软件没有打开")
		time.Sleep(10 * time.Second)
		os.Exit(3)
	}
	//获取节点
	nodes := getNodes()
	if nodes == nil {
		fmt.Println("端口错误或者clash软件没有打开")
		time.Sleep(10 * time.Second)
		os.Exit(3)
	}

	index := 1
	f, err := os.OpenFile("netflix.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, node := range nodes {
		setNode(node)
		isNetflix, out := nf()
		if *netflix == "1" {
			if isNetflix == false {
				continue
			}
		}
		if out != "" {
			fmt.Printf("%d-%d   节点名: %s\n", index, len(nodes), node)
			fmt.Println(out)
			fmt.Fprintf(f, "%d-%d   节点名: %s\n", index, len(nodes), node)
			fmt.Fprintln(f, out)
		} else {
			fmt.Printf("%d-%d   节点名: %s\n", index, len(nodes), node)
			fmt.Println("完全不支持Netflix")
			fmt.Fprintf(f, "%d-%d   节点名: %s\n", index, len(nodes), node)
			fmt.Fprintln(f, "完全不支持Netflix")
		}
		index++
	}
	fmt.Println("检测结束")
}
