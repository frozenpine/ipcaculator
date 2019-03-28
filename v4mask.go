package ipcalculator

import (
	"errors"
	"regexp"
)

const (
	maxMaskBit = 30
	minMaskBit = 1
)

var (
	maskPattern = regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)

	// errors
	errMaskPattern = errors.New("invalid MASKv4 pattern")
)

// MaskAddress IPv4掩码数据结构
type MaskAddress struct {
	base
	isWild bool
}

// FromString 以点分十进制格式为NetMask赋值
func (mask *MaskAddress) FromString(value string) error {
	if !maskPattern.MatchString(value) {
		return errMaskPattern
	}

	return mask.setValueByString(value)
}

// MaskBit 获取IPvMask的掩码位
func (mask MaskAddress) MaskBit() int {
	var maskValue uint32

	if mask.isWild {
		maskValue = mask.value ^ 0xFFFFFFFF
	} else {
		maskValue = mask.value
	}

	bit := 0

	for maskValue > 0 {
		maskValue = maskValue << 1
		bit++
	}

	return bit
}

// FromMaskBit 使用Bit位设置NetMask
func (mask *MaskAddress) FromMaskBit(bit int) error {
	if bit < minMaskBit || bit > maxMaskBit {
		return errors.New("Mask bit out of range")
	}

	mask.value = 0

	for i := 0; i < bit; i++ {
		mask.value = mask.value>>1 | 0x80000000
	}

	return nil
}

// WildMask 获取当前掩码的反掩码
func (mask MaskAddress) WildMask() MaskAddress {
	if mask.isWild {
		return mask
	}

	wildMask := MaskAddress{}

	wildMask.value = mask.value ^ 0xFFFFFFFF
	wildMask.isWild = true

	return wildMask
}
