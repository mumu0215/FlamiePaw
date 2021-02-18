package common

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)
var userAgent ="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
func GetOne(group *sync.WaitGroup,client *http.Client,myLog *log.Logger,baseUrl chan string,rep chan string) {  //处理每个请求
	for url :=range baseUrl{
		if url==""{
			close(baseUrl)             //关闭target channel
		}else {
			temp,err:=oneWorker(client,url)
			if err!=nil{
				myLog.Println(err)
				//fmt.Println(url+": "+temp)
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
	tempBody,_:=DetermineDecoding(res.Body)
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
