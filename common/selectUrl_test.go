package common

import (
	"bufio"
	"os"
	"testing"
)

func TestSelectUrl(t *testing.T)  {
	var(
		a=[]string{
			"http://221.226.253.50:7001	404	Apache-Coyote/1.1	NULL_title!",
			"http://221.226.253.50:8083	200	nginx/1.17.6	NULL_title!",
			"http://221.226.253.50:85	404	NULL_server!	Error 404--Not Found"}
	)
	f,_:=os.OpenFile("scanFile.txt",os.O_CREATE|os.O_RDWR|os.O_TRUNC,0666)
	buf:=bufio.NewWriter(f)
	for _,i:=range a{
		GetWeb200(i,buf,"\r\n")
	}
	f.Close()
	t.Log("success")
}
