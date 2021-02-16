package common

import (
	"bufio"
	"fmt"
	"strconv"
)

func GetWeb200(urlSlice []string,scanWriter *bufio.Writer,sp string) string{
	statuCode,_:=strconv.Atoi(urlSlice[1])
	if statuCode >=200 && statuCode <300{
		fmt.Fprintf(scanWriter,urlSlice[0]+sp)
		scanWriter.Flush()
	}
	return urlSlice[0]
}

