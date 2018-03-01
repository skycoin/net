package producer

import (
	"testing"
	"github.com/skycoin/skycoin/src/util/file"
	"path/filepath"
)

func init() {
	err := Init(filepath.Join(file.UserHome(), ".skywire", "discovery", "conf.json"))
	if err != nil {
		panic(err)
	}
}
func TestSend(t *testing.T) {
	err := Send("03c773268be29f8fe48144ccb2c19045b487dceadbb8e1d7817f84d612521bc68c", "038b558e91a343c0f449b536ceab640d4055dbced85f6f5969a58ee56f2280588c", 100)
	if err != nil {
		panic(err)
		return
	}
}
