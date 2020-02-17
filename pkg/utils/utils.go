package utils

import "bytes"

// ancient magic =)
// 534d/5350 swastika
// 2744 snowflake
// 2620 Skull And Crossbones
const swasti = "\u2744"
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
