package ipcalculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// PrivateA A类私网地址
	PrivateA, _ = CreateNetwork("10.0.0.0/8")
	// PrivateB B类私网地址
	PrivateB, _ = CreateNetwork("172.16.0.0/12")
	// PrivateC C类私网地址
	PrivateC, _ = CreateNetwork("192.168.0.0/16")
	// Loopback 回环口网段
	Loopback, _ = CreateNetwork("127.0.0.0/8")
)

// Network IPv4网络定义, 包含IPAddress + MaskAddress
type Network struct {
	address IPAddress
	mask    MaskAddress
}

// ToString 显示网络
func (net *Network) ToString() string {
	return fmt.Sprintf("%s/%d", net.NetID().ToString(), net.mask.MaskBit())
}

// Contains 测试IP地址是否在Network指定范围内
func (net *Network) Contains(ip IPAddress) bool {
	netID := net.NetID()
	broad := net.Broadcast()

	if ip.value >= netID.value && ip.value <= broad.value {
		return true
	}

	return false
}

// Includes 测试是否为Network的子网
func (net *Network) Includes(network *Network) bool {
	outNet := net.NetID()
	inNet := network.NetID()

	if outNet.value&inNet.value == outNet.value {
		return true
	}

	return false
}

// FirstIP 获取网段第一个可用的IP地址
func (net *Network) FirstIP() IPAddress {
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
func (net *Network) NetID() IPAddress {
	netID := IPAddress{}

	netID.value = net.address.value & net.mask.value

	return netID
}

// Broadcast 获取网段的广播地址
func (net *Network) Broadcast() IPAddress {
	broadcast := IPAddress{}

	broadcast.value = net.address.value | net.mask.WildMask().value

	return broadcast
}

// MaxSubCount 获取最大子网数
func (net *Network) MaxSubCount() int {
	return 1 << uint(maxMaskBit-net.mask.MaskBit())
}

// Split 将网络分割成多个子网
func (net *Network) Split(count int) ([]*Network, error) {
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
		return nil, errors.New("Exceed max subnets count")
	}

	subCount := 1 << uint(countBit)
	splited := make([]*Network, subCount, subCount)

	for i := 0; i < subCount; i++ {
		sub := Network{}
		sub.address.value = net.NetID().value + uint32(i*1<<uint(32-originMaskBit-countBit))
		sub.mask.FromMaskBit(originMaskBit + countBit)

		splited[i] = &sub
	}

	return splited, nil
}

// CreateNetwork Network工厂方法
func CreateNetwork(value string) (*Network, error) {
	values := strings.Split(value, "/")

	ip := IPAddress{}
	err := ip.FromString(values[0])

	if err != nil {
		return nil, err
	}

	mask := MaskAddress{}
	maskBit, err := strconv.Atoi(values[1])

	if err != nil {
		if err = mask.FromString(values[1]); err != nil {
			return nil, err
		}
	} else if err = mask.FromMaskBit(maskBit); err != nil {
		return nil, err
	}

	return &Network{address: ip, mask: mask}, nil
}
