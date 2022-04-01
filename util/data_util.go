package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func CombineString(values ...string) string {
	var buffer bytes.Buffer

	for _, v := range values {
		buffer.WriteString(v)
	}

	return buffer.String()
}

func SubString(str string, begin, lenght int) (substr string) {

	//fmt.Println("Substring =", str)
	rs := []rune(str)
	lth := len(rs)
	//fmt.Printf("begin=%d, end=%d, lth=%d\n", begin, lenght, lth)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + lenght

	if end > lth {
		end = lth
	}
	//fmt.Printf("begin=%d, end=%d, lth=%d\n", begin, lenght, lth)
	return string(rs[begin:end])
}

func ConverToJsonString(datas interface{}) []byte {
	jsonStr, err := json.Marshal(datas)

	if err != nil {
		fmt.Println("json.Marshalerror:", err)
	}
	return jsonStr
}

func HexToBinaryStr(hexStr string, count int) string {

	base, err := strconv.ParseInt(hexStr, 16, 32)

	if err != nil {
		fmt.Printf("%s \n", err.Error())
		return ""
	}

	return fmt.Sprintf("%0*s", count, strconv.FormatInt(base, 2))
}
