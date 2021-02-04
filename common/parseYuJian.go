package common

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//var fileName=flag.String("f","","filename")
func doWork(input []string,sp string) (string,[]string,error) {
	url:=""
	ssh:=""
	telnet:=""
	ftp:=""
	mysql:=""
	mssql:=""
	ajp13:=""
	redis:=""
	mongodb:=""
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
			countUrl+=1
		//case "unknow":
		//	url+="http://"+sList[1]+sp
		//	countUrl+=1
		case "ssh":
			ssh+=sList[1]+sp
			countSsh+=1
		case "https":
			url+="https://"+sList[1]+sp
			countUrl+=1
		case "telnet":
			telnet+=sList[1]+sp
			countTelnet+=1
		case "ftp":
			ftp+=sList[1]+sp
			countFtp+=1
		case "mysql":
			mysql+=sList[1]+sp
			countMysql+=1
		case "ms-sql-s":
			mssql+=sList[1]+sp
			countMssql+=1
		case "ajp13":
			ajp13+=sList[1]+sp
			countAjp13+=1
		case "redis":
			redis+=sList[1]+sp
			countRedis+=1
		case "mongodb":
			mongodb+=sList[1]+sp
			countMongoDb+=1
		default:
			url+="http://"+sList[1]+sp //未分类送去web检测
			countUnknow+=1
		}
	}
	err:=os.MkdirAll("./Result",os.ModePerm)
	if err!=nil{
		//fmt.Println("Fail to Create folder")
		//os.Exit(1)
		return strings.TrimSpace(url),[]string{},err
	}

	fileUrl,err:=os.OpenFile("./Result/url.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		return strings.TrimSpace(url),[]string{},err
	}
	fileUrl.WriteString(url)
	fileUrl.Close()

	if countAjp13>0{
		fileAjp13,err:=os.OpenFile("./Result/ajp13.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileAjp13.WriteString(ajp13)
		fileAjp13.Close()
	}
	if countMongoDb>0{
		fileMongoDb,err:=os.OpenFile("./Result/mongodb.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileMongoDb.WriteString(mongodb)
		fileMongoDb.Close()
	}
	if countRedis>0{
		fileRedis,err:=os.OpenFile("./Result/redis.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileRedis.WriteString(redis)
		fileRedis.Close()
	}
	if countTelnet>0{
		fileTelnet,err:=os.OpenFile("./Result/telnet.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileTelnet.WriteString(telnet)
		fileTelnet.Close()
	}

	if countFtp>0{
		fileFtp,err:=os.OpenFile("./Result/ftp.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileFtp.WriteString(ftp)
		fileFtp.Close()
	}

	if countSsh>0{
		fileSSH,err:=os.OpenFile("./Result/ssh.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileSSH.WriteString(ssh)
		fileSSH.Close()
	}

	if countMysql>0{
		fileMysql,err:=os.OpenFile("./Result/mysql.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileMysql.WriteString(mysql)
		fileMysql.Close()
	}

	if countMssql>0{
		fileMssql,err:=os.OpenFile("./Result/mssql.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
		if err!=nil{
			return strings.TrimSpace(url),[]string{},err
		}
		fileMssql.WriteString(mssql)
		fileMssql.Close()
	}

	return strings.TrimSpace(url),[]string{strconv.Itoa(countUrl),strconv.Itoa(countSsh),
		strconv.Itoa(countTelnet),strconv.Itoa(countFtp),strconv.Itoa(countAjp13),
		strconv.Itoa(countMysql),strconv.Itoa(countMssql),strconv.Itoa(countRedis),
		strconv.Itoa(countMongoDb),strconv.Itoa(countUnknow)},nil
}
func ParseYuJ(portScanFile string,sp string) (a []string,b []string,err1 error) {
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