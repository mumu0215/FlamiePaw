package common

import (
	"reflect"
	"testing"
)

func TestParseXml(t *testing.T) {
	var(
		in = "test.xml"
		out=30
		outSlice=[]string{"0","0","30"}
	)
	urlSlice,countSlice,_:=ParseXml(in,"\r\n")
	if out!=len(urlSlice) || !reflect.DeepEqual(outSlice,countSlice){
		t.Error("fail to pass test")
	}
}