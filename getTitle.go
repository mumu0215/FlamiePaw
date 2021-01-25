package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"src/common"
	"strings"
	"sync"
	"time"
)
 //routineCountTotal 线程
var(
	splitTool string           //换行符

	mode=flag.Int("m",0,"mode choice,1:parse from url list,-uF Needed;" +
		"2:parse from port scan file,-pF or -xF Needed")
	routineCountTotal=flag.Int("t",15,"thread")
	myProxy=flag.String("p","","proxy")
	urlFileName=flag.String("uF","","url file name")
	portScanFileName=flag.String("pF","","yujian port scan file")
	nmapXmlFileName=flag.String("xF","","nmap output xmlFileName")
	timeOut=flag.Int("T",2,"request timeout seconds")
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
			if (*portScanFileName=="" && *nmapXmlFileName=="") || (*portScanFileName!="" && *nmapXmlFileName!=""){
				fmt.Println("portScan filename or xml filename need one!")
				os.Exit(0)
			}
		default:
			fmt.Println("Wrong mode Number! PLZ choose 1 or 2!")
			os.Exit(0)
		}
	}
	if *timeOut<1{
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
		Timeout:time.Duration(*timeOut)*time.Second,
		Transport: &http.Transport{
		//参数未知影响，目前不使用
		//TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},}             //复用client
	if *myProxy!=""{                   //设置代理
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(strings.TrimSpace(*myProxy))
		}
		client=&http.Client{
			Timeout:time.Duration(*timeOut)*time.Second,
			Transport:&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:                  proxy,
		}}
	}
	wg:=&sync.WaitGroup{}
	target:=make(chan string)
	result:=make(chan string)

	urlTitleFile,err:=os.OpenFile("urlTitle.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("Fail to open file for result")
		os.Exit(0)
	}
	defer urlTitleFile.Close()
	buf:=bufio.NewWriter(urlTitleFile)
	//接受结果，并处理判断信号
	go func() {
		for rep :=range result{
			if rep==""{            //中断信号
				close(result)
			}else {
				//文件处理传出结果
				tempList:=strings.Split(rep,"\t")
				fmt.Fprintf(buf,"%-60s\t%s\t%-20s\t%s"+splitTool,tempList[0],tempList[1],tempList[2],tempList[3])
				buf.Flush()
			}
		}
	}()

	for i:=0;i<*routineCountTotal;i++{
		wg.Add(1)
		go common.GetOne(wg,client,target,result)
	}
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
		for _,singleUrl:=range tempSlice{
			target<-singleUrl
		}
	}else if *mode==2 && *nmapXmlFileName!=""{ //mode 2 nmap扫描结果输入
		tempSlice,reportTempSlice,err:=common.ParseXml(*nmapXmlFileName,splitTool)
		if err!=nil{
			fmt.Println(err)
			os.Exit(1)
		}
		reportSlice=reportTempSlice
		for _,singleUrl:=range tempSlice{
			target<-singleUrl
		}
	}
	target<-""   //工作分发结束
	wg.Wait()
	result<-""   //发出结果中断信号
	if *mode==2 && len(reportSlice)!=0 && *portScanFileName!=""{
		fmt.Println("Found Information:")
		fmt.Println("\tUrl:"+reportSlice[0]+"    SSH:"+reportSlice[1]+"    Telnet:"+reportSlice[2])
		fmt.Println("\tFTP:"+reportSlice[3]+"    AJP13:"+reportSlice[4]+"    Mysql:"+reportSlice[5])
		fmt.Println("\tMssql:"+reportSlice[6]+"    UnKnow:"+reportSlice[7])
	}else if *mode==2 && len(reportSlice)!=0 && *nmapXmlFileName!=""{
		fmt.Println("Found Information:")
		fmt.Println("\tUrl:"+reportSlice[0]+"    SSH:"+reportSlice[1]+"    UnKnow:"+reportSlice[2])
	}
}