package common

import (
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	//"golang.org/x/text/encoding/simplifiedchinese"
	//"golang.org/x/text/transform"
)

func DetermineDecoding(rep io.Reader) (io.Reader,error) {     //处理网页编码问题
	OldReader := bufio.NewReader(rep)
	bytes, err := OldReader.Peek(1024)
	if err != nil {
		return rep,err
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	reader := transform.NewReader(OldReader, e.NewDecoder())
	return reader,nil
}
