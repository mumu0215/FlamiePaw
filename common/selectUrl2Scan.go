package common

import (
	"bufio"
	"fmt"
	"strconv"
)

func GetWeb200(urlSlice []string,scanWriter *bufio.Writer,sp string) string{
	statuCode,_:=strconv.Atoi(urlSlice[1])
	if statuCode >=200 && statuCode <300{
		fmt.Fprintf(scanWriter,"%-40s\t%s\t%-20s\t%s"+sp,urlSlice[0],urlSlice[1],urlSlice[2],urlSlice[3])
		scanWriter.Flush()
	}
	return urlSlice[0]
}

