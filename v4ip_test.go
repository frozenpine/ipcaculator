package ipcalculator

import "testing"

func TestIPAddress(t *testing.T) {
	ip := IPAddress{base: base{value: 1}}

	ipString := ip.ToString()

	if ipString != "0.0.0.1" {
		t.Error("ToString function failed")
	} else {
		t.Log(ipString)
	}

	ip.FromString("10.0.0.1")

	hexString := ip.ToHexString()

	if hexString != "0x0A000001" {
		t.Error("ToHex function failed")
	} else {
		t.Log(hexString)
	}

	ip.value = uint32(0x0A0A0A0A)
	ipString = ip.ToString()

	if ipString != "10.10.10.10" {
		t.Error("Fail to set value directly")
	} else if !ip.IsPrivate() {
		t.Error("IsPrivate function failed")
	} else if ip.IsLoopback() {
		t.Error("IsLoopback function failed")
	}

	ip.FromString("224.0.0.1")
	class := ip.Class()

	if class != ClassD {
		t.Error("Class function failed")
	} else {
		t.Log(ip.ToString())
		t.Log(class)
	}
}
