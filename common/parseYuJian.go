package common

import (
	"errors"
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//var fileName=flag.String("f","","filename")
func doWork(input []string,sp string) (string,[]string,error) {
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
	telnet:=""
	ftp:=""
	mysql:=""
	mssql:=""
	ajp13:=""
	redis:=""
	mongodb:=""
	oracle:=""
	//计数
	countMongoDb:=0
	countRedis:=0
	countUrl:=0
	countSsh:=0
	countTelnet:=0
	countFtp:=0
	countMysql:=0
	countMssql:=0
	countAjp13:=0
	countOracle:=0
	countUnknow:=0
	for _,temp:=range input{
		sList:=strings.Split(strings.TrimSpace(temp),"\t")
		//if strings.Split(strings.TrimSpace(sList[1]),":")[]
		switch sList[2] {
		case "http":
			url+="http://"+sList[1]+sp
			countUrl+=1
		case "ssl":
			url+="https://"+sList[1]+sp
			countUrl+=1
		case "ssl/http":
			url+="https://"+sList[1]+sp
			countUrl+=1
		case "-":
			url+="http://"+sList[1]+sp
			url+="https://"+sList[1]+sp
			countUrl+=2
		//case "unknow":
		//	url+="http://"+sList[1]+sp
		//	countUrl+=1
		case "ssh":
			ssh+=sList[1]+sp
			sshSlice=append(sshSlice,sList[1])
			countSsh+=1
		case "https":
			url+="https://"+sList[1]+sp
			countUrl+=1
		case "telnet":
			telnet+=sList[1]+sp
			telnetSlice=append(telnetSlice,sList[1])
			countTelnet+=1
		case "ftp":
			ftp+=sList[1]+sp
			ftpSlice=append(ftpSlice,sList[1])
			countFtp+=1
		case "mysql":
			mysql+=sList[1]+sp
			mysqlSlice=append(mysqlSlice,sList[1])
			countMysql+=1
		case "ms-sql-s":
			mssql+=sList[1]+sp
			mssqlSlice=append(mssqlSlice,sList[1])
			countMssql+=1
		case "ajp13":
			ajp13+=sList[1]+sp
			ajp13Slice=append(ajp13Slice,sList[1])
			countAjp13+=1
		case "redis":
			redis+=sList[1]+sp
			redisSlice=append(redisSlice,sList[1])
			countRedis+=1
		case "mongodb":
			mongodb+=sList[1]+sp
			mongodbSlice=append(mongodbSlice,sList[1])
			countMongoDb+=1
		case "oracle-tns":
			oracle+=sList[1]+sp
			oracleSlice=append(oracleSlice,sList[1])
			countOracle+=1
		default:
			url+="http://"+sList[1]+sp //未分类送去web检测
			url+="https://"+sList[1]+sp
			countUnknow+=2
		}
	}
	err:=os.MkdirAll("./Result",os.ModePerm)
	if err!=nil{
		//fmt.Println("Fail to Create folder")
		//os.Exit(1)
		return strings.TrimSpace(url),[]string{},err
	}

	fileUrl,err:=os.OpenFile("./Result/YuJianUrl.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileUrl.WriteString(url)
	fileUrl.Close()

	fileService,err:=os.OpenFile("./Result/YuJianService.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	if countAjp13>0{
		fileService.WriteString("ajp1.3:"+sp)
		fileService.WriteString(ajp13)
		tempSlice=append(tempSlice,Service{
			Service:    "ajp13",
			IpPortList: ajp13Slice,
		})
	}
	if countMongoDb>0{
		fileService.WriteString("mongoDB:"+sp)
		fileService.WriteString(mongodb)
		tempSlice=append(tempSlice,Service{
			Service:    "mongoDB",
			IpPortList: mongodbSlice,
		})
	}
	if countRedis>0{
		fileService.WriteString("redis:"+sp)
		fileService.WriteString(redis)
		tempSlice=append(tempSlice,Service{
			Service:    "redis",
			IpPortList: redisSlice,
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
	if countFtp>0{
		fileService.WriteString("ftp:"+sp)
		fileService.WriteString(ftp)
		tempSlice=append(tempSlice,Service{
			Service:    "ftp",
			IpPortList: ftpSlice,
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
	if countMssql>0{
		fileService.WriteString("mssql:"+sp)
		fileService.WriteString(mssql)
		tempSlice=append(tempSlice,Service{
			Service:    "mssql",
			IpPortList: mssqlSlice,
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
	fileJson,err:=os.OpenFile("./Result/YuJian.json",os.O_RDWR|os.O_TRUNC|os.O_CREATE,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileJson.WriteString(jsonData)
	fileJson.Close()

	return strings.TrimSpace(url),[]string{strconv.Itoa(countUrl),strconv.Itoa(countSsh),
		strconv.Itoa(countTelnet),strconv.Itoa(countFtp),strconv.Itoa(countAjp13),
		strconv.Itoa(countMysql),strconv.Itoa(countMssql),strconv.Itoa(countRedis),
		strconv.Itoa(countMongoDb),strconv.Itoa(countOracle),strconv.Itoa(countUnknow)},nil
}
func ParseYuJ(portScanFile string,sp string) (a []string,b []string,err1 error) {
	if temp:=strings.Split(portScanFile,".");temp[len(temp)-1]!="txt"{
		return []string{},[]string{},errors.New("file type error,need portScan file")
	}
	dataByte,err:=ioutil.ReadFile(portScanFile)
	if err!=nil{
		//fmt.Println("Fail to open portscanFile!")
		//os.Exit(0)
		return []string{},[]string{},err
	}
	dataSlice:=strings.Split(strings.TrimSpace(string(dataByte)),sp)
	tempStr,tempCountSlice,err:=doWork(dataSlice,sp)
	if err!=nil{
		return []string{},[]string{},err
	}
	tempSlice:=strings.Split(tempStr,sp)
	return tempSlice,tempCountSlice,nil
}