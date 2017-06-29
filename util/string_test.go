package util

import "testing"

func TestString2ByteSlice(t *testing.T) {
	s := "hello"
	bs := String2ByteSlice(s)
	if string(bs) != s {
		t.Fail()
	}
	t.Logf("%s", bs)
}
