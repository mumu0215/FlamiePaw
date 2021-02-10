package common

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

//baseUrl+"\t"+strconv.Itoa(res.StatusCode)+"\t"+server+"\t"+title
func GetWeb200(urlString string,scanWriter *bufio.Writer,sp string)  {
	tempSlice:=strings.Split(urlString,"\t")
	statuCode,_:=strconv.Atoi(tempSlice[1])
	if statuCode >=200 && statuCode <300{
		fmt.Fprintf(scanWriter,tempSlice[0]+sp)
		scanWriter.Flush()
	}
}