package utils

import "bytes"

// ancient magic =)
const swasti = "\u5350"
const space = ` `

//StringSplitter split \n & \r from byte slice & return string
func StringSplitter(b bytes.Buffer) string {
	var prepStr []byte
	for _, sm := range b.Bytes() {
		if string(sm) == "\n" || string(sm) == "\r" {
			prepStr = append(prepStr, space...)
		} else {
			prepStr = append(prepStr, sm)
		}
	}
	return string(prepStr)
}

//StringDecorator join string with prerfix & postfix decorators
func StringDecorator(s string) string {
	return swasti + space + s + space + swasti
}
