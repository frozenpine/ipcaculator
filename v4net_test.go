package ipcalculator

import "testing"

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
