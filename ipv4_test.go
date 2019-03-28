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

func TestMaskAddress(t *testing.T) {
	mask := MaskAddress{}
	mask.FromString("255.255.255.0")

	if mask.value != uint32(0xFFFFFF00) {
		t.Error("FromString function failed")
	} else {
		t.Log(mask.ToString())
	}

	maskBit := mask.MaskBit()

	if maskBit != 24 {
		t.Error("MaskBit function failed")
	} else {
		t.Log(maskBit)
	}

	mask2 := MaskAddress{}
	mask2.FromMaskBit(30)
	mask2String := mask2.ToString()

	if mask2String != "255.255.255.252" {
		t.Error("FromMaskBit function failed")
	} else {
		t.Log(mask2String)
	}

	wildMask := mask.WildMask()
	maskString := wildMask.ToString()

	if !wildMask.isWild || maskString != "0.0.0.255" {
		t.Error("WildMask function failed")
	} else {
		t.Log(maskString)
		t.Log(wildMask.MaskBit())
		t.Log(wildMask.ToHexString())
	}
}

func TestNetwork(t *testing.T) {
	n, err := CreateNetwork("191.168.4.6/23")

	if err != nil {
		t.Error(err)
	} else {
		netID := n.NetID()
		broad := n.Broadcast()
		first, last := n.IPRange()

		if netID.value != 0xBFA80400 {
			t.Error("NetID function failed")
		} else {
			t.Log(netID.ToString())
		}

		if broad.value != 0xBFA805FF {
			t.Error("Broadcast function failed")
		} else {
			t.Log(broad.ToString())
		}

		if first.value != 0xBFA80401 {
			t.Error("FirstIP function failed")
		} else {
			t.Log(first.ToString())
		}

		if last.value != 0xBFA805FE {
			t.Error("LastIP function failed")
		} else {
			t.Log(last.ToString())
		}

		ip := IPAddress{}
		ip.FromString("192.168.6.4")

		if n.Contains(ip) {
			t.Error("Contains function failed")
		}

		inner, _ := CreateNetwork("191.168.5.0/24")

		if !n.Includes(inner) {
			t.Error("Includes function failed")
		}
	}

	n, err = CreateNetwork("10.0.0.1/255.255.255.0")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(n.ToString())
	}
}

func TestNetworkSplit(t *testing.T) {
	n, err := CreateNetwork("191.168.4.6/23")

	splited, _ := n.Split(2)

	if len(splited) != 2 {
		t.Error("Split function failed")
	}

	splited, _ = n.Split(6)
	if len(splited) != 8 {
		t.Error("Split function failed")
	} else {
		for _, n := range splited {
			t.Log(n.ToString())
		}
	}

	splited, err = n.Split(129)

	if err == nil {
		t.Error("Split function failed")
	} else {
		t.Log(n.MaxSubCount())
	}
}
