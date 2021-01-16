package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"src/common"
	"strconv"
	"strings"
	"sync"
)
 //routineCountTotal 线程
var(
	userAgent ="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	splitTool string           //换行符

	mode=flag.Int("m",0,"mode choice,1:parse from url list,-uF Needed;" +
		"2:parse from port scan file,-pF or -xF Needed")
	routineCountTotal=flag.Int("t",15,"thread")
	myProxy=flag.String("p","","proxy")
	urlFileName=flag.String("uF","","url file name")
	portScanFileName=flag.String("pF","","yujian port scan file")
	nmapXmlFileName=flag.String("xF","","nmap output xmlFileName")
)

func init()  {
	//解析命令行，判断参数合法
	flag.Parse()
	if *routineCountTotal<1{
		fmt.Println("Thread must more than one!")
		os.Exit(0)
	}
	if *mode==0{
		fmt.Println("Mode choice Needed!")
		os.Exit(0)
	}else{
		switch *mode {
		case 1:
			if *urlFileName==""{
				fmt.Println("url FileName input Needed!")
				os.Exit(0)
			}
		case 2:
			if *portScanFileName=="" && *nmapXmlFileName==""{
				fmt.Println("portScan filename or xml filename at least need one!")
				os.Exit(0)
			}
		default:
			fmt.Println("Wrong mode Number! PLZ choose 1 or 2!")
			os.Exit(0)
		}
	}

	//根据系统设置换行符
	switch runtime.GOOS {
	case "windows":
		splitTool="\r\n"
	case "linux":
		splitTool="\n"
	case "darwin":
		splitTool="\r"
	default:
		splitTool="\r\n"
	}
}

func getOne(group *sync.WaitGroup,client *http.Client,baseUrl chan string,rep chan string) {  //处理每个请求
	for url :=range baseUrl{
		if url==""{
			close(baseUrl)             //关闭target channel
		}else {
			temp,err:=oneWorker(client,url)
			if err!=nil{
				fmt.Println(url+": "+temp)
				//rep <-err.Error()
			}else {
				rep<-temp
			}
		}
	}
	group.Done()
}
func oneWorker(client *http.Client,baseUrl string) (string,error) {
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err!=nil{
		return  "Fail to create request!",err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", baseUrl)
	req.Header.Set("Accept-Encoding","")
	res, err := client.Do(req)
	if err!=nil{
		return "Error occur at Sending Request!",err
	}
	defer res.Body.Close()
	//fmt.Println(res.StatusCode)
	//fmt.Println(res.Header["Server"][0])
	server:=""
	if len(res.Header["Server"])!=0{
		server+=res.Header["Server"][0]
	}else {
		server+="NULL_server!"
	}
	tempBody,_:=common.DetermineDecoding(res.Body)
	doc,err:=goquery.NewDocumentFromReader(tempBody)
	if err!=nil{
		return "Fail to parse response!",err
	}
	title:=strings.TrimSpace(doc.Find("title").First().Text())
	if title==""{
		title+="NULL_title!"
	}
	report:=baseUrl+"\t"+strconv.Itoa(res.StatusCode)+"\t"+server+"\t"+title
	return report,nil
}
func getTxtToList(fileName string) []string {
	dataByte,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		fmt.Println("fail to open file "+fileName)
		os.Exit(0)
	}
    data:=strings.TrimSpace(string(dataByte))
    return strings.Split(data,splitTool)
}
func main() {
	//flag.Parse()
	client:=&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}             //复用client
	if *myProxy!=""{                   //设置代理
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(strings.TrimSpace(*myProxy))
		}
		client=&http.Client{Transport:&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:                  proxy,
		}}
	}
	wg:=&sync.WaitGroup{}
	target:=make(chan string)
	result:=make(chan string)

	fileResult,err:=os.OpenFile("urlTitle.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("Fail to open file for result")
		os.Exit(0)
	}
	defer fileResult.Close()
	buf:=bufio.NewWriter(fileResult)
	//接受结果，并处理判断信号
	go func() {
		for rep :=range result{
			if rep==""{            //中断信号
				close(result)
			}else {
				//文件处理传出结果
				//fileResult.WriteString(rep+"\r\n")
				tempList:=strings.Split(rep,"\t")
				fmt.Fprintf(buf,"%-60s\t%s\t%-20s\t%s"+splitTool,tempList[0],tempList[1],tempList[2],tempList[3])
				buf.Flush()
			}
		}
	}()

	for i:=0;i<*routineCountTotal;i++{
		wg.Add(1)
		go getOne(wg,client,target,result)
	}

	//分发任务
	for _,baseUrl:=range getTxtToList(*urlFileName){
		target <-baseUrl
	}
	target<-""   //工作分发结束
	wg.Wait()
	result<-""   //发出结果中断信号
}