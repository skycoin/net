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

func TestByteSlice2String(t *testing.T) {
	bs := []byte("world")
	s := ByteSlice2String(bs)
	if s != string(bs) {
		t.Fail()
	}
	t.Logf("%s", s)
}