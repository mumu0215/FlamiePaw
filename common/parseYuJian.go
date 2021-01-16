package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//var fileName=flag.String("f","","filename")
func doWork(input []string,sp string)  {
	url:=""
	ssh:=""
	telnet:=""
	ftp:=""
	mysql:=""
	mssql:=""
	ajp13:=""
	for _,temp:=range input{
		sList:=strings.Split(strings.TrimSpace(temp),"\t")
		//if strings.Split(strings.TrimSpace(sList[1]),":")[]
		switch sList[2] {
		case "http":
			url+="http://"+sList[1]+sp
		case "ssl":
			url+="https://"+sList[1]+sp
		case "ssl/http":
			url+="https://"+sList[1]+sp
		case "ssh":
			ssh+=sList[1]+sp
		case "https":
			url+="https://"+sList[1]+sp
		case "telnet":
			telnet+=sList[1]+sp
		case "ftp":
			ftp+=sList[1]+sp
		case "mysql":
			mysql+=sList[1]+sp
		case "ms-sql-s":
			mssql+=sList[1]+sp
		case "ajp13":
			ajp13+=sList[1]+sp
		default:
		}
	}
	fileUrl,err:=os.OpenFile("url.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open url file")
		os.Exit(0)
	}
	fileUrl.WriteString(url)
	fileAjp,err:=os.OpenFile("ajp13.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open ajp file")
		os.Exit(0)
	}
	fileAjp.WriteString(ajp13)
	fileTelnet,err:=os.OpenFile("telnet.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open telnet file")
		os.Exit(0)
	}
	fileTelnet.WriteString(telnet)

	fileFtp,err:=os.OpenFile("ftp.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open ftp file")
		os.Exit(0)
	}
	fileFtp.WriteString(ftp)

	fileSSH,err:=os.OpenFile("ssh.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open ssh file")
		os.Exit(0)
	}
	fileSSH.WriteString(ssh)

	fileMysql,err:=os.OpenFile("mysql.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open mysql file")
		os.Exit(0)
	}
	fileMysql.WriteString(mysql)

	fileMssql,err:=os.OpenFile("mssql.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open mssql file")
		os.Exit(0)
	}
	fileMssql.WriteString(mssql)
}
func ParseYuJ(portScanFile string,sp string)  {
	dataByte,err:=ioutil.ReadFile(portScanFile)
	if err!=nil{
		fmt.Println("Fail to open portscanFile!")
		os.Exit(0)
	}
	dataSlice:=strings.Split(strings.TrimSpace(string(dataByte)),sp)
	doWork(dataSlice,sp)
}