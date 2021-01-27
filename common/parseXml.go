package common

import (
	"github.com/n0ncetonic/nmapxml"
	"os"
	"strconv"
	"strings"
)
//解析nmap输出的xml文件
//error统一返回上层处理
func ParseXml(xmlFileName string,sp string) ([]string,[]string,error){
	scanData,err:=nmapxml.Readfile(xmlFileName)
	if err!=nil{
		//fmt.Println("Read XML FileName Failed!")
		return []string{},[]string{},err
	}
	tempUrl,countSlice,err:=dealWithRun(scanData,sp)
	if err!=nil{
		return []string{},[]string{},err
	}
	return strings.Split(tempUrl,sp),countSlice,nil
}
func dealWithRun(r nmapxml.Run,sp string) (string,[]string,error) {
	url:=""
	ssh:=""
	//telnet:=""
	//ftp:=""
	//mysql:=""
	//mssql:=""
	//ajp13:=""

	//计数
	countUrl:=0
	unknown:=0
	countSsh:=0
	//countTelnet:=0
	//countFtp:=0
	//countMysql:=0
	//countMssql:=0
	//countAjp13:=0
	hostSlice:=r.Host
	for _,host:=range hostSlice{
		ipAddr:=host.Address.Addr
		for _,portInfo:=range *host.Ports.Port{
			if portInfo.State.State=="open"{
				portID:=portInfo.PortID
				service:=portInfo.Service.Name
				switch service {
				case "http":
					url+="http://"+ipAddr+":"+portID+sp
					countUrl+=1
				case "https":
					url+="https://"+ipAddr+":"+portID+sp
					countUrl+=1
				//这边因为服务名可能有字符串出入，暂时不判断
				case "ssh":
					ssh+=ipAddr+":"+portID+sp
					countSsh+=1
				//case "telnet":
				//case "ftp":
				//case "mysql":
				//case "mssql":
				//case "ajp":
				default:      //未分类全部送去web检测
					url+="http://"+ipAddr+":"+portID+sp
					unknown+=1
				}
			}
		}
	}
	err:=os.MkdirAll("./XmlResult",os.ModePerm)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileUrl,err:=os.OpenFile("./XmlResult/url.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileUrl.WriteString(url)
	fileUrl.Close()

	if countSsh>0{
		fileSsh,err:=os.OpenFile("./XmlResult/ssh.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileSsh.WriteString(ssh)
		fileSsh.Close()
	}

	return strings.TrimSpace(url),[]string{strconv.Itoa(countUrl),strconv.Itoa(countSsh),strconv.Itoa(unknown)},nil
}


//func main1() {
//	scanData, err := nmapxml.Readfile("output.xml")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	jsonData, err := json.Marshal(scanData)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Printf("%+v", string(jsonData))
//}
