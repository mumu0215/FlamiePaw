package crawler

import (
	"encoding/json"
	"fmt"
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
			fmt.Println("Crawler Scan: "+url+"...")
			subDomains:=getOutputContinually("cmd","/C",".\\crawlergo.exe","-c", myConfig.MyChrome.Path,
				"-t", myConfig.MyChrome.MaxTapCount,"-f", myConfig.MyCrawler.FilterMode,"--fuzz-path",
				"--output-mode","json","--custom-headers",getHeader(),"--push-to-proxy",myConfig.MyCrawler.Proxy,
				"--push-pool-max",myConfig.MyCrawler.PushToPoolMax,url)
			domains <-subDomains
			fmt.Println("Finish Scan: "+url+".")
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