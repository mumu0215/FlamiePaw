package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//var fileName=flag.String("f","","filename")
func doWork(input []string,sp string) (string,[]string) {
	url:=""
	ssh:=""
	telnet:=""
	ftp:=""
	mysql:=""
	mssql:=""
	ajp13:=""
	//计数
	countUrl:=0
	countSsh:=0
	countTelnet:=0
	countFtp:=0
	countMysql:=0
	countMssql:=0
	countAjp13:=0
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
		default:
		}
	}
	err:=os.MkdirAll("./result",os.ModePerm)
	if err!=nil{
		fmt.Println("Fail to Create folder")
		os.Exit(1)
	}
	fileUrl,err:=os.OpenFile("./result/url.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open url file")
		os.Exit(0)
	}
	fileUrl.WriteString(url)
	fileUrl.Close()

	fileAjp13,err:=os.OpenFile("./result/ajp13.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open ajp file")
		os.Exit(0)
	}
	fileAjp13.WriteString(ajp13)
	fileAjp13.Close()

	fileTelnet,err:=os.OpenFile("./result/telnet.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open telnet file")
		os.Exit(0)
	}
	fileTelnet.WriteString(telnet)
	fileTelnet.Close()

	fileFtp,err:=os.OpenFile("./result/ftp.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open ftp file")
		os.Exit(0)
	}
	fileFtp.WriteString(ftp)
	fileFtp.Close()

	fileSSH,err:=os.OpenFile("./result/ssh.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open ssh file")
		os.Exit(0)
	}
	fileSSH.WriteString(ssh)
	fileSSH.Close()

	fileMysql,err:=os.OpenFile("./result/mysql.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open mysql file")
		os.Exit(0)
	}
	fileMysql.WriteString(mysql)
	fileMysql.Close()

	fileMssql,err:=os.OpenFile("./result/mssql.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open mssql file")
		os.Exit(0)
	}
	fileMssql.WriteString(mssql)
	fileMssql.Close()

	return strings.TrimSpace(url),[]string{strconv.Itoa(countUrl),strconv.Itoa(countSsh),
		strconv.Itoa(countTelnet),strconv.Itoa(countFtp),strconv.Itoa(countAjp13),
		strconv.Itoa(countMysql),strconv.Itoa(countMssql)}
}
func ParseYuJ(portScanFile string,sp string) (a []string,b []string) {
	dataByte,err:=ioutil.ReadFile(portScanFile)
	if err!=nil{
		fmt.Println("Fail to open portscanFile!")
		os.Exit(0)
	}
	dataSlice:=strings.Split(strings.TrimSpace(string(dataByte)),sp)
	tempStr,tempCountSlice:=doWork(dataSlice,sp)
	tempSlice:=strings.Split(tempStr,sp)
	return tempSlice,tempCountSlice
}