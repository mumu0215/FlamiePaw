package crawler

import (
	"encoding/json"
	"github.com/json-iterator/go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type CrawlerJson struct {
	AllReqList []interface{} `json:"all_req_list"`
	ReqList []interface{} `json:"req_list"`
	AllDomainList []string `json:"all_domain_list"`
	SubDomainList []string `json:"sub_domain_list"`
}

//yaml解析
type CustomHeader struct {
	UserAgent string `json:"User-Agent"`
}
type Chrome struct {
	Path string `yaml:"Path"`
	MaxTapCount string `yaml:"MaxTapCount"`
}
type Crawler struct {
	CrawlerThread int `yaml:"CrawlerThread"`
	FilterMode string `yaml:"FilterMode"`
	Proxy string `yaml:"Proxy"`
	PushToPoolMax string `yaml:"PushToPoolMax"`
}
type CrawlerConfig struct {
	MyChrome Chrome `yaml:"Chrome"`
	MyCrawler Crawler `yaml:"Crawler"`
}
var UserAgent=[]string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:28.0) Gecko/20100101 Firefox/31.0",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; WOW64; en-US; rv:2.0.4) Gecko/20120718 AskTbAVR-IDW/3.12.5.17700 Firefox/14.0.1",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36 Mozilla/5.0 (iPad; U; CPU OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B334b Safari/531.21.10",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0.1 Safari/602.2.14",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36",
}
//.\crawlergo.exe -c "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe"
//-t 5 -f smart --fuzz-path --output-mode json http://hffengyun.com/

func getRandomUA() string {
	rand.Seed(time.Now().Unix())
	temp:=rand.Intn(6)   // [0,100)的随机值，返回值为int
	return UserAgent[temp]
}
func getHeader() string {
	var header CustomHeader
	header.UserAgent=getRandomUA()
	jsonHeader,err:=json.Marshal(header)
	if err!=nil{
		panic(err)
	}
	return string(jsonHeader)
}

func ParseCrawlerConfig(fileName string) (CrawlerConfig,error) {
	var config CrawlerConfig
	configFile,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		return CrawlerConfig{},err
	}
	err=yaml.Unmarshal(configFile,&config)
	if err!=nil{
		return CrawlerConfig{},err
	}
	return config,nil
}
func getOutputDirectly(name string, args ...string) (output []byte) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output() // 等到命令执行完, 一次性获取输出
	if err != nil {
		panic(err)
	}
	output, err = simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	if err != nil {
		panic(err)
	}
	return
}

func getOutputContinually(name string, args ...string) []string{
	cmd := exec.Command(name, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdoutPipe.Close()
	//go func() {
	//	scanner := bufio.NewScanner(stdoutPipe)
	//	for scanner.Scan() { // 命令在执行的过程中, 实时地获取其输出
	//		data, err := simplifiedchinese.GB18030.NewDecoder().Bytes(scanner.Bytes()) // 防止乱码
	//		if err != nil {
	//			fmt.Println("transfer error with bytes:", scanner.Bytes())
	//			continue
	//		}
	//
	//		fmt.Printf("%s\n", string(data))
	//	}
	//}()
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	var crawlerTemp CrawlerJson
	if opBytes,err:=ioutil.ReadAll(stdoutPipe);err!=nil{
		panic(err)
	}else {
		temp:=strings.Split(string(opBytes),"--[Mission Complete]--")
		tempj:=strings.TrimSpace(temp[1])
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		json.Unmarshal([]byte(tempj),&crawlerTemp)
	}
	cmd.Wait()
	return crawlerTemp.SubDomainList
}
func RunCrawler(group *sync.WaitGroup,myConfig CrawlerConfig,baseUrl chan string,domains chan []string){
	for url:=range baseUrl{
		if url==""{
			close(baseUrl)
		}else {
			subDomains:=getOutputContinually("cmd","/C",".\\crawlergo.exe","-c", myConfig.MyChrome.Path,
				"-t", myConfig.MyChrome.MaxTapCount,"-f", myConfig.MyCrawler.FilterMode,"--fuzz-path",
				"--output-mode","json","--custom-headers",getHeader(),url)
			domains <-subDomains
		}
	}
	group.Done()
}

//func main1() {
//	wg:=&sync.WaitGroup{}
//	target:=make(chan string)
//	result:=make(chan []string)
//	myConfig,err:=ParseCrawlerConfig("config.yaml")
//	if err!=nil{
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	go func() {
//		for temp:=range result{
//			if temp[0]=="stop"{
//				close(result)
//			}else {
//				fmt.Println(temp)
//			}
//		}
//	}()
//	for i:=0;i<myConfig.MyCrawler.CrawlerThread;i++{
//		wg.Add(1)
//		go RunCrawler(wg,myConfig,target,result)
//	}
//	temp:=[]string{"http://testphp.vulnweb.com/","http://hffengyun.com/","http://www.hw008.com/"}
//	for _,url:=range temp{
//		target<-url
//	}
//	target<-""
//	wg.Wait()
//	result<-[]string{"stop"}
//}