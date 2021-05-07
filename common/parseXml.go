package common

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/n0ncetonic/nmapxml"
	"os"
	"strconv"
	"strings"
)
//解析nmap输出的xml文件
//error统一返回上层处理
func ParseXml(xmlFileName string,sp string) ([]string,[]string,error){
	if temp:=strings.Split(xmlFileName,".");temp[len(temp)-1]!="xml"{
		return []string{},[]string{},errors.New("file type error,need xml file")
	}
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
	var (
		sshSlice []string
		telnetSlice []string
		ftpSlice []string
		mysqlSlice []string
		mssqlSlice []string
		ajp13Slice []string
		redisSlice []string
		mongodbSlice []string
		oracleSlice []string
		serviceList ServiceList
		tempSlice []Service
	)

	ssh:=""
	redis:=""
	telnet:=""
	ftp:=""
	mysql:=""
	mssql:=""
	ajp13:=""
	mongodb:=""
	oracle:=""
	//计数
	countRedis:=0
	countUrl:=0
	unknown:=0
	countSsh:=0
	countTelnet:=0
	countFtp:=0
	countMysql:=0
	countMssql:=0
	countAjp13:=0
	countMongoDB:=0
	countOracle:=0
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
					sshSlice=append(sshSlice,ipAddr+":"+portID)
					countSsh+=1
				case "redis":
					redis+=ipAddr+":"+portID+sp
					redisSlice=append(redisSlice,ipAddr+":"+portID)
					countRedis+=1
				case "telnet":
					telnet+=ipAddr+":"+portID+sp
					telnetSlice=append(telnetSlice,ipAddr+":"+portID)
					countTelnet+=1
				case "ftp":
					ftp+=ipAddr+":"+portID+sp
					ftpSlice=append(ftpSlice,ipAddr+":"+portID)
					countFtp+=1
				case "mysql":
					mysql+=ipAddr+":"+portID+sp
					mysqlSlice=append(mysqlSlice,ipAddr+":"+portID)
					countMysql+=1
				case "ms-sql-s":
					mssql+=ipAddr+":"+portID+sp
					mssqlSlice=append(mssqlSlice,ipAddr+":"+portID)
					countMssql+=1
				case "ajp13":
					ajp13+=ipAddr+":"+portID+sp
					ajp13Slice=append(ajp13Slice,ipAddr+":"+portID)
					countAjp13+=1
				case "mongodb":
					mongodb+=ipAddr+":"+portID+sp
					mongodbSlice=append(mongodbSlice,ipAddr+":"+portID)
					countMongoDB+=1
				case "oracle-tns":
					oracle+=ipAddr+":"+portID+sp
					oracleSlice=append(oracleSlice,ipAddr+":"+portID)
					countOracle+=1
				default:      //未分类全部送去web检测
					url+="http://"+ipAddr+":"+portID+sp
					url+="https://"+ipAddr+":"+portID+sp
					unknown+=2
				}
			}
		}
	}
	err:=os.MkdirAll("./Result",os.ModePerm)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileUrl,err:=os.OpenFile("./Result/NmapUrl.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileUrl.WriteString(url)
	fileUrl.Close()

	fileService,err:=os.OpenFile("./Result/NmapService.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}

	if countRedis>0{
		fileService.WriteString("redis:"+sp)
		fileService.WriteString(redis)
		tempSlice=append(tempSlice,Service{
			Service:    "redis",
			IpPortList: redisSlice,
		})
	}
	if countMssql>0{
		fileService.WriteString("mssql:"+sp)
		fileService.WriteString(mssql)
		tempSlice=append(tempSlice,Service{
			Service:    "mssql",
			IpPortList: mssqlSlice,
		})
	}
	if countAjp13>0{
		fileService.WriteString("ajp1.3:"+sp)
		fileService.WriteString(ajp13)
		tempSlice=append(tempSlice,Service{
			Service:    "ajp13",
			IpPortList: ajp13Slice,
		})
	}
	if countFtp>0{
		fileService.WriteString("ftp:"+sp)
		fileService.WriteString(ftp)
		tempSlice=append(tempSlice,Service{
			Service:    "ftp",
			IpPortList: ftpSlice,
		})
	}
	if countTelnet>0{
		fileService.WriteString("telnet:"+sp)
		fileService.WriteString(telnet)
		tempSlice=append(tempSlice,Service{
			Service:    "telnet",
			IpPortList: telnetSlice,
		})
	}

	if countSsh>0{
		fileService.WriteString("ssh:"+sp)
		fileService.WriteString(ssh)
		tempSlice=append(tempSlice,Service{
			Service:    "ssh",
			IpPortList: sshSlice,
		})
	}
	if countMysql>0{
		fileService.WriteString("mysql:"+sp)
		fileService.WriteString(mysql)
		tempSlice=append(tempSlice,Service{
			Service:    "mysql",
			IpPortList: mysqlSlice,
		})
	}
	if countMongoDB>0{
		fileService.WriteString("mongodb:"+sp)
		fileService.WriteString(mongodb)
		tempSlice=append(tempSlice,Service{
			Service:    "mongodb",
			IpPortList: mongodbSlice,
		})
	}
	if countOracle>0{
		fileService.WriteString("oracle:"+sp)
		fileService.WriteString(oracle)
		tempSlice=append(tempSlice,Service{
			Service:    "oracle-tns",
			IpPortList: oracleSlice,
		})
	}
	fileService.Close()

	serviceList.ServiceList=tempSlice
	var json=jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData,err:=json.MarshalToString(serviceList)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileJson,err:=os.OpenFile("./Result/Nmap.json",os.O_RDWR|os.O_TRUNC|os.O_CREATE,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileJson.WriteString(jsonData)
	fileJson.Close()

	return strings.TrimSpace(url),[]string{strconv.Itoa(countUrl),strconv.Itoa(countSsh),strconv.Itoa(countTelnet),
		strconv.Itoa(countFtp),strconv.Itoa(countAjp13),strconv.Itoa(countMysql),strconv.Itoa(countMssql),
		strconv.Itoa(countRedis),strconv.Itoa(countMongoDB),strconv.Itoa(countOracle),strconv.Itoa(unknown)},nil
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
