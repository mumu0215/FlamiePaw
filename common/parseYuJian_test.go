package common

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseYuJ(t *testing.T) {
	var(
		in ="test.txt"
		outCount=[]string{"22","0","0","0","0","1","0","11"}
		urlcount=33
	)
	tempSlice,c,_:=ParseYuJ(in,"\r\n")
	if len(tempSlice)!=urlcount || !reflect.DeepEqual(c,outCount){
		fmt.Println(len(tempSlice))
		fmt.Println(c)
		t.Error("fail to pass test")
	}
}
