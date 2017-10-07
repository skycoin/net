package data

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/skycoin/net/skycoin-messenger/factory"
)

var (
	keysPath  string
	keys      = make(map[string]*factory.SeedConfig)
	keysMutex = &sync.Mutex{}
)

func walkFunc(path string, info os.FileInfo, err error) (e error) {
	if err != nil || info.IsDir() {
		return
	}
	sc, e := factory.ReadSeedConfig(path)
	if e != nil {
		return e
	}
	keys[sc.PublicKey] = sc
	return
}

func InitData(path string) (err error) {
	keysMutex.Lock()
	err = filepath.Walk(path, walkFunc)
	keysPath = path
	keysMutex.Unlock()
	return
}

func GetData() map[string]*factory.SeedConfig {
	if len(keysPath) < 1 {
		return nil
	}

	keysMutex.Lock()
	k := keys
	err := filepath.Walk(keysPath, walkFunc)
	keysMutex.Unlock()
	if err != nil {
		return k
	}
	return keys
}
