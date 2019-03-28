package ipcalculator

import (
	"fmt"
	"strconv"
	"strings"
)

type base struct {
	value uint32
}

// ToString 返回IPv4地址/掩码的点分十进制字符串
func (v base) ToString() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.Sec1(), v.Sec2(), v.Sec3(), v.Sec4())
}

// Sec1 first section in DEC format
func (v base) Sec1() uint8 {
	return uint8(v.value >> 24)
}

// Sec2 second section in DEC format
func (v base) Sec2() uint8 {
	return uint8((v.value & 0x00FF0000) >> 16)
}

// Sec3 third section in DEC format
func (v base) Sec3() uint8 {
	return uint8((v.value & 0x0000FF00) >> 8)
}

// Sec4 forth section in DEC format
func (v base) Sec4() uint8 {
	return uint8(v.value & 0x000000FF)
}

// ToHexString 返回IPv4地址/掩码的十六进制字符串
func (v base) ToHexString() string {
	hexString := fmt.Sprintf("%X", v.value)
	return fmt.Sprintf("0x%08s", hexString)
}

func (v *base) setValueByString(value string) error {
	value = strings.TrimSpace(value)
	secStrings := strings.Split(value, ".")
	var secValues [4]uint32

	for idx, secString := range secStrings {
		sec, err := strconv.ParseUint(secString, 10, 32)

		if err != nil {
			return err
		}

		secValues[idx] = uint32(sec)
	}

	v.value = (secValues[0] << 24) | (secValues[1] << 16) | (secValues[2] << 8) | secValues[3]

	return nil
}
