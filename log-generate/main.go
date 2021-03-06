package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// 资源结构体
type resource struct {
	url    string   // 访问的url
	target string   // 字符串替换的目标，便于字符串替换
	start  int      // url中的开始id
	end    int      // url中的结束id
}

// 生成页面 url
func ruleResource() []resource {
	var res []resource // 结构体切片

	// 首页
	r1 := resource{
		url:    "http://localhost:8888/",
		target: "",
		start:  0,
		end:    0,
	}

	// 列表页
	r2 := resource{
		url:    "http://localhost:8888/list/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    21,
	}

	// 详情页
	r3 := resource{
		url:    "http://localhost:8888/movie/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    12924,
	}
	// 追加，形成url资源切片
	res = append(append(append(res, r1), r2), r3)
	return res
}

// 构建具体的 url 列表
func buildUrl(res []resource) []string {
	var list []string
	for _, resItem := range res {
		if len(resItem.target) == 0 {
			list = append(list, resItem.url)
		} else {
			for i := resItem.start; i <= resItem.end; i++ {
				// 字符串替换
				// strconv.Itoa(i)   int => string
				urlStr := strings.Replace(resItem.url, resItem.target, strconv.Itoa(i), -1)
				list = append(list, urlStr)
			}
		}
	}
	return list
}

// 模拟日志
func makeLog(current, refer, ua string) string {
	// 利用 net/url 库生成 key-value 键值对（type Values map[string][]string）
	u := url.Values{}
	// 时间
	u.Set("time", "1")
	// 当前 url
	u.Set("url", current)
	// refer 跳转前的url
	u.Set("refer", refer)
	// user-agent
	u.Set("ua", ua)

	// URL encoded
	paramsStr := u.Encode()

	logTemplate := "172.20.10.4 - - [23/Jul/2019:12:28:48 +0800] \"GET /dig?{$paramsStr}  HTTP/1.1\" 200 43 \"-\" \"{$ua}\" \"-\""
	log := strings.Replace(logTemplate, "{$paramsStr}", paramsStr, -1)
	log = strings.Replace(log, "{$ua}", ua, -1)
	return log
}

// 生成随机数
func randInt(min, max int) int {
	// 利用随机值生成随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if min > max {
		return max
	}
	// 返回[0, max-min) + min，即[min, max)之间的随机数
	return r.Intn(max-min) + min
}

// 常用的user-agent
var userAgentList = []string {

	// Android平台原生浏览器
	"Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19",
	"Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (Linux; U; Android 2.2; en-gb; GT-P1000 Build/FROYO) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",

	// Firefox火狐
	"Mozilla/5.0 (Android; Mobile; rv:14.0) Gecko/14.0 Firefox/14.0",
	"Mozilla/5.0 (Android; Tablet; rv:14.0) Gecko/14.0 Firefox/14.0",
	"Mozilla/5.0 (Windows NT 6.2; WOW64; rv:21.0) Gecko/20100101 Firefox/21.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:21.0) Gecko/20130331 Firefox/21.0",

	// Google chrome
	"Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19",
	"Mozilla/5.0 (Linux; Android 4.1.2; Nexus 7 Build/JZ054K) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166 Safari/535.19",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Ubuntu/11.10 Chromium/27.0.1453.93 Chrome/27.0.1453.93 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) CriOS/27.0.1453.10 Mobile/10B350 Safari/8536.25",

	// Internet Explore
	"Mozilla/5.0 (compatible; WOW64; MSIE 10.0; Windows NT 6.2)", //IE10
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)", //IE9
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)", //IE8
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)", //IE7
	"Mozilla/4.0 (Windows; MSIE 6.0; Windows NT 5.2)", //IE6

	// Opera
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.9.168 Version/11.52", //Mac
	"Opera/9.80 (Windows NT 6.1; WOW64; U; en) Presto/2.10.229 Version/11.62", //Windows

	// Safari
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_6; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27", //Mac
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27", //windows
	"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3", //iPad
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3", //iPhone

	// iOS
	"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3", //iPad
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.", //iPhone
	"Mozilla/5.0 (iPod; U; CPU like Mac OS X; en) AppleWebKit/420.1 (KHTML, like Gecko) Version/3.0 Mobile/3A101a Safari/419.3", //iPod

	// Windows Phone
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows Phone OS 7.0; Trident/3.1; IEMobile/7.0; LG; GW910)", //windows Phone 7
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; SAMSUNG; SGH-i917)",// Windows Phone7.5
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows Phone 8.0; Trident/6.0; IEMobile/10.0; ARM; Touch; NOKIA; Lumia 920)", //windows phone 8
}

func main() {
	// 命令行参数， 推荐使用 pflag
	// total, filePath 都是指针
	total := flag.Int("total", 100, "how many rows by created")
	filePath := flag.String("filePath", "D:\\GoWorkspace\\traffic-analysis\\generate-log\\dig.log", "log file path")
	// 解析命令行参数，必须调用
	flag.Parse()
	// fmt.Println(*total, *filePath)

	// 构造出真实网站的 url 集合
	res := ruleResource()
	// fmt.Println("res: ", res)
	// url 列表
	list := buildUrl(res)
	// fmt.Println("list:", list)

	// 生成 total 行日志内容
	logStr := ""
	for i := 0; i <= *total; i++ {
		// list 中随机选取一个 url 作为当前 url
		currentUrl := list[randInt(0, len(list))]
		// list 中随机选取一个 url 作为跳转前 url
		referUrl := list[randInt(0, len(list))]
		// userAgentList 中随机选取一个 user-agent
		ua := userAgentList[randInt(0, len(userAgentList))]

		// 日志行拼接
		logStr = logStr + makeLog(currentUrl, referUrl, ua) + "\n"
		fmt.Println("logStr: ", logStr)
		// ioutil.WriteFile(*filePath, []byte(logStr), 0644)
	}
	// 一次性写入日志文件
	fd, err := os.OpenFile(*filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open log file error: %s\n", err.Error())
	}
	defer fd.Close()
	// fmt.Println("logStr: ", logStr)

	_, err = fd.Write([]byte(logStr))
	if err != nil {
		fmt.Printf("write log file error: %s\n", err.Error())
	}

	fmt.Println("done.")
}