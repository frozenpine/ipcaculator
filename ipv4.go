package ipcalculator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

	maxMaskBit = 30
	minMaskBit = 1
)

var (
	ipPattern   = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))`)
	maskPattern = regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}`)

	// PrivateA A类私网地址
	PrivateA, _ = CreateNetwork("10.0.0.0/8")
	// PrivateB B类私网地址
	PrivateB, _ = CreateNetwork("172.16.0.0/12")
	// PrivateC C类私网地址
	PrivateC, _ = CreateNetwork("192.168.0.0/16")
	// Loopback 回环口网段
	Loopback, _ = CreateNetwork("127.0.0.0/8")
)

// Errors
var (
	errIPPattern   = errors.New("invalid IPv4 pattern")
	errMaskPattern = errors.New("invalid MASKv4 pattern")
)

type base struct {
	value uint32
}

// IPAddress IPv4地址数据结构
type IPAddress struct {
	base
}

// MaskAddress IPv4掩码数据结构
type MaskAddress struct {
	base
	isWild bool
}

// Network IPv4网络定义, 包含IPAddress + MaskAddress
type Network struct {
	address IPAddress
	mask    MaskAddress
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

// CreateNetwork Network工厂方法
func CreateNetwork(value string) (Network, error) {
	values := strings.Split(value, "/")

	ip := IPAddress{}
	err := ip.FromString(values[0])

	if err != nil {
		return Network{}, err
	}

	mask := MaskAddress{}
	maskBit, err := strconv.Atoi(values[1])

	if err != nil {
		if err = mask.FromString(values[1]); err != nil {
			return Network{}, err
		}
	} else if err = mask.FromMaskBit(maskBit); err != nil {
		return Network{}, err
	}

	return Network{address: ip, mask: mask}, nil
}

// ToString 显示网络
func (net Network) ToString() string {
	return fmt.Sprintf("%s/%d", net.NetID().ToString(), net.mask.MaskBit())
}

// Contains 测试IP地址是否在Network指定范围内
func (net Network) Contains(ip IPAddress) bool {
	netID := net.NetID()
	broad := net.Broadcast()

	if ip.value >= netID.value && ip.value <= broad.value {
		return true
	}

	return false
}

// Includes 测试是否为Network的子网
func (net Network) Includes(network Network) bool {
	outNet := net.NetID()
	inNet := network.NetID()

	if outNet.value&inNet.value == outNet.value {
		return true
	}

	return false
}

// FirstIP 获取网段第一个可用的IP地址
func (net Network) FirstIP() IPAddress {
	first := IPAddress{}

	first.value = net.NetID().value + 1

	return first
}

// LastIP 获取网段最后一个可用的IP地址
func (net Network) LastIP() IPAddress {
	last := IPAddress{}

	last.value = net.Broadcast().value - 1

	return last
}

// IPRange 获取网段的可用地址范围
func (net Network) IPRange() (IPAddress, IPAddress) {
	return net.FirstIP(), net.LastIP()
}

// NetID 获取网段的网络号
func (net Network) NetID() IPAddress {
	netID := IPAddress{}

	netID.value = net.address.value & net.mask.value

	return netID
}

// Broadcast 获取网段的广播地址
func (net Network) Broadcast() IPAddress {
	broadcast := IPAddress{}

	broadcast.value = net.address.value | net.mask.WildMask().value

	return broadcast
}

// MaxSubCount 获取最大子网数
func (net Network) MaxSubCount() int {
	return 1 << uint(maxMaskBit-net.mask.MaskBit())
}

// Split 将网络分割成多个子网
func (net Network) Split(count int) ([]Network, error) {
	countBit := 0

	// 值为0则count为2的幂
	if count&(count-1) == 0 {
		countBit = -1
	}

	for count > 0 {
		count >>= 1
		countBit++
	}

	originMaskBit := net.mask.MaskBit()

	maxSubBit := maxMaskBit - originMaskBit

	if maxSubBit <= countBit {
		return []Network{}, errors.New("Exceed max subnets count")
	}

	subCount := 1 << uint(countBit)
	splited := make([]Network, subCount, subCount)

	for i := 0; i < subCount; i++ {
		sub := Network{}
		sub.address.value = net.NetID().value + uint32(i*1<<uint(32-originMaskBit-countBit))
		sub.mask.FromMaskBit(originMaskBit + countBit)

		splited[i] = sub
	}

	return splited, nil
}
