package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"src/common"
	"src/crawler"
	"strings"
	"sync"
	"time"
)
 //routineCountTotal 线程
var(
	splitTool string           //换行符
	MyConfig crawler.CrawlerConfig
	MyLog *log.Logger
	mode=flag.Int("m",0,"mode choice,1:parse from url list,-uF Needed;" +
		"2:parse from port scan file,-pF or -xF Needed")
	urlFileName=flag.String("uF","","url file name")
	portScanFileName=flag.String("pF","","yujian port scan file")
	nmapXmlFileName=flag.String("xF","","nmap output xmlFileName")
	crawlerFlag=flag.Bool("crawler",false,"run crawler at last")
	Banner=`
___________.__                 .__      __________                
\_   _____/|  | _____    _____ |__| ____\______   \_____ __  _  __
 |    __)  |  | \__  \  /     \|  |/ __ \|     ___/\__  \\ \/ \/ /
 |     \   |  |__/ __ \|  Y Y  \  \  ___/|    |     / __ \\     / 
 \___  /   |____(____  /__|_|  /__|\___  >____|    (____  /\/\_/  
     \/              \/      \/        \/               \/        
                                                          by 半九十`
)

func init()  {
	//解析命令行，判断参数合法
	flag.Parse()
	var errConfig error
	fmt.Println(Banner)
	fmt.Println("Parse Config File...")
	MyConfig,errConfig=crawler.ParseCrawlerConfig("config.yaml")
	if errConfig!=nil{
		fmt.Println("Error occur in Parse Config file!")
		os.Exit(1)
	}
	if MyConfig.MyGetTitle.Thread<1{
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
			if (*portScanFileName=="" && *nmapXmlFileName=="") || (*portScanFileName!="" && *nmapXmlFileName!=""){
				fmt.Println("portScan filename or xml filename need one!")
				os.Exit(0)
			}
		default:
			fmt.Println("Wrong mode Number! PLZ choose 1 or 2!")
			os.Exit(0)
		}
	}
	if MyConfig.MyGetTitle.TimeOut<1{
		fmt.Println("Set timeout more than 1 second")
		os.Exit(0)
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

func getUrlFileToList(fileName string) []string {
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
	client:=&http.Client{
		Timeout:time.Duration(MyConfig.MyGetTitle.TimeOut)*time.Second,
		Transport: &http.Transport{
		//参数未知影响，目前不使用
		//TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},}             //复用client
	if MyConfig.MyGetTitle.Proxy!="none"{                   //设置代理
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(strings.TrimSpace(MyConfig.MyGetTitle.Proxy))
		}
		client=&http.Client{
			Timeout:time.Duration(MyConfig.MyGetTitle.TimeOut)*time.Second,
			Transport:&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:                  proxy,
		}}
	}
	wg:=&sync.WaitGroup{}
	target:=make(chan string)
	result:=make(chan string)
	wgScan:=&sync.WaitGroup{}
	targetScan:=make(chan string)
	resultScan:=make(chan []string)

	err:=os.MkdirAll("./Result",os.ModePerm)
	if err!=nil{
		fmt.Println("Fail to Create folder")
		os.Exit(1)
	}
	urlTitleFile,err:=os.OpenFile("./Result/urlTitle.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("Fail to open file for result")
		os.Exit(1)
	}
	defer urlTitleFile.Close()
	webToScan,err:=os.OpenFile("./Result/url200.txt",os.O_TRUNC|os.O_RDWR|os.O_CREATE,0666)
	if err!=nil{
		fmt.Println("Fail to open file for scan")
		os.Exit(1)
	}
	defer webToScan.Close()

	logFile,err:=os.OpenFile("log.txt",os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	if err!=nil{
		fmt.Println("Fail to Open file for log")
		os.Exit(1)
	}
	defer logFile.Close()
	MyLog=log.New(io.MultiWriter(logFile),"",log.Lshortfile|log.LstdFlags)

	buf:=bufio.NewWriter(urlTitleFile)
	scanBuf:=bufio.NewWriter(webToScan)
	var scanUrlSlice []string
	//接受结果，并处理判断信号
	go func() {
		for rep :=range result{
			if rep==""{            //中断信号
				close(result)
			}else {
				//文件处理传出结果
				tempList:=strings.Split(rep,"\t")
				temp:=common.GetWeb200(tempList,scanBuf,splitTool)
				scanUrlSlice=append(scanUrlSlice,temp)         //挑选扫描网址
				fmt.Fprintf(buf,"%-40s\t%s\t%-20s\t%s"+splitTool,tempList[0],tempList[1],tempList[2],tempList[3])
				buf.Flush()
			}
		}
	}()
	//根据线程分发任务
	for i:=0;i<MyConfig.MyGetTitle.Thread;i++{
		wg.Add(1)
		go common.GetOne(wg,client,MyLog,target,result)
	}
	if *crawlerFlag==true{
		go func() {               //爬虫协程
			for temp:=range resultScan{
				if temp[0]=="stop"{
					close(resultScan)
				}else {
					//处理子域名
					//fmt.Println(temp)
				}
			}
		}()
		//根据crawler线程分发任务
		if err!=nil{
			panic(err)
		}
		for i:=0;i< MyConfig.MyCrawler.CrawlerThread;i++{
			wgScan.Add(1)
			go crawler.RunCrawler(wgScan, MyConfig,targetScan,resultScan)
		}
	}

	fmt.Println("GetTitle Running...")
	//mode 1
	//接受url文件
	var reportSlice []string
	if *mode==1{
		//分发任务
		for _,baseUrl:=range getUrlFileToList(*urlFileName){
			target <-baseUrl
		}
	}else if *mode==2 && *portScanFileName!=""{ //mode 2 御剑扫描结果输入
		//分发任务
		tempSlice,reportTempSlice,err:=common.ParseYuJ(*portScanFileName,splitTool)
		if err!=nil{
			fmt.Println(err)
			os.Exit(1)
		}
		reportSlice=reportTempSlice
		if len(tempSlice)!=0{
			for _,singleUrl:=range tempSlice{
				target<-singleUrl
			}
		}
	}else if *mode==2 && *nmapXmlFileName!=""{ //mode 2 nmap扫描结果输入
		tempSlice,reportTempSlice,err:=common.ParseXml(*nmapXmlFileName,splitTool)
		if err!=nil{
			fmt.Println(err)
			os.Exit(1)
		}
		reportSlice=reportTempSlice
		if len(tempSlice)!=0{
			for _,singleUrl:=range tempSlice{
				target<-singleUrl
			}
		}
	}
	if reportSlice[0]!="0" || reportSlice[len(reportSlice)-1]!="0"{
		target<-""   //工作分发结束
	}
	wg.Wait()
	result<-""   //发出结果中断信号
	fmt.Println("GetTitle Done!")
	if *mode==2 && len(reportSlice)!=0 && *portScanFileName!=""{
		fmt.Println("Found Information:")
		fmt.Println("\tUrl:"+reportSlice[0]+"    SSH:"+reportSlice[1]+"    Telnet:"+reportSlice[2])
		fmt.Println("\tFTP:"+reportSlice[3]+"    AJP13:"+reportSlice[4]+"    Mysql:"+reportSlice[5])
		fmt.Println("\tMssql:"+reportSlice[6]+"    Redis:"+reportSlice[7]+"    MongoDB:"+reportSlice[8])
		fmt.Println("\tUnKnow:"+reportSlice[9])
	}else if *mode==2 && len(reportSlice)!=0 && *nmapXmlFileName!=""{
		fmt.Println("Found Information:")
		fmt.Println("\tUrl:"+reportSlice[0]+"    SSH:"+reportSlice[1]+"    Telnet:"+reportSlice[2])
		fmt.Println("\tFTP:"+reportSlice[3]+"    AJP13:"+reportSlice[4]+"    Mysql:"+reportSlice[5])
		fmt.Println("\tMssql:"+reportSlice[6]+"    Redis:"+reportSlice[7]+"    MongoDB:"+reportSlice[8])
		fmt.Println("\tUnKnow:"+reportSlice[9])
	}
	if *crawlerFlag==true{
		fmt.Println("Crawler Running...")
		//爬虫任务分发
		for _,scanUrl:=range scanUrlSlice{
			targetScan<-scanUrl
		}
		//爬虫协程逻辑结束
		targetScan<-""
		wgScan.Wait()
		resultScan<-[]string{"stop"}
	}
	fmt.Println("ALL DONE !")
}