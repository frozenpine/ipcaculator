package ipcalculator

import (
	"errors"
	"fmt"
	"regexp"
)

// IPClass class type in IPv4
type IPClass uint8

// ToString get class string
func (c IPClass) ToString() string {
	switch c {
	case ClassA:
		return "A"
	case ClassB:
		return "B"
	case ClassC:
		return "C"
	case ClassD:
		return "D"
	case ClassE:
		return "E"
	default:
		panic("invalid IP class")
	}
}

const (
	// ClassA A class 0******* 1~127
	ClassA IPClass = 0x7F
	// ClassB B class 10****** 128~191
	ClassB IPClass = 0xBF
	// ClassC C class 110***** 192~223
	ClassC IPClass = 0xDF
	// ClassD D Class 1110**** 224~239
	ClassD IPClass = 0xEF
	// ClassE E Class 1111**** 240~256 reserved
	ClassE IPClass = 0xFF
)

var (
	ipPattern = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)

	// Errors
	errIPPattern = errors.New("invalid IPv4 pattern")
)

// IPAddress IPv4地址数据结构
type IPAddress struct {
	base
}

// FromString 以点分十进制格式为IPAddress赋值
func (ip *IPAddress) FromString(value string) error {
	if !ipPattern.MatchString(value) {
		return errIPPattern
	}

	return ip.setValueByString(value)
}

// Class 获取IPv4地址的分类
func (ip IPAddress) Class() IPClass {
	sec := ip.Sec1()

	if sec <= uint8(ClassA) {
		return ClassA
	}

	if sec <= uint8(ClassB) {
		return ClassB
	}

	if sec <= uint8(ClassC) {
		return ClassC
	}

	if sec <= uint8(ClassD) {
		return ClassD
	}

	if sec <= uint8(ClassE) {
		return ClassE
	}

	panic(fmt.Errorf("can't get class for this address: %s", ip.ToString()))
}

// IsPrivate 测试地址是否为私有地址
func (ip IPAddress) IsPrivate() bool {
	if PrivateA.Contains(ip) {
		return true
	}

	if PrivateB.Contains(ip) {
		return true
	}

	if PrivateC.Contains(ip) {
		return true
	}

	return false
}

// IsLoopback 测试地址是否为回环地址
func (ip IPAddress) IsLoopback() bool {
	if Loopback.Contains(ip) {
		return true
	}

	return false
}
