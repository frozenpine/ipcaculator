package ipcalculator

import "testing"

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
