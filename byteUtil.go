package gospider

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type ByteHandler interface {
	Handle(s []byte) ([]byte, error)
}

type GBKByteHandler struct {
}

func (h *GBKByteHandler) Handle(s []byte) ([]byte, error) {
	return DecodeGBK(s)
}

//gbk转utf8
func DecodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
