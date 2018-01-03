package conn

import "testing"

func TestFec(t *testing.T) {
	ds := 4
	ps := 2
	decoder := newFECDecoder(ds, ps, 0)
	encoder := newFECEncoder(ds, ps, 0)

	datas := [][]byte{
		{0x1},
		{0x2, 0x3},
		{0x4},
		{0x5},
	}

	var pp [][]byte
	for _, d := range datas {
		p, err := encoder.encode(d)
		if err != nil {
			t.Error(err)
		}
		t.Log(p)
		if p != nil {
			pp = p
		}
	}

	datas[1] = nil

	for i, d := range datas {
		g, err := decoder.decode(uint32(i+1), d)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v", g)
	}

	g, err := decoder.decode(5, pp[0])
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", g)

	g, err = decoder.decode(6, pp[1])
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", g)
}
