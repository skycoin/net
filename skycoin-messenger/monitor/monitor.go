package monitor

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"encoding/json"
)

type Conns struct {
	Conns []Conn `json:"conns"`
}
type Conn struct {
	Key string `json:"key"`
	Type string `json:"type"`
	SendBytes uint64 `json:"send_bytes"`
	ReceivedBytes uint64 `json:"received_bytes"`
	LastAckTime	int64 `json:"last_ack_time"`
	StartTime int64 `json:"start_time"`
}

var srv *http.Server

type Monitor struct {
	factory *factory.MessengerFactory
	address string
}

func New(f *factory.MessengerFactory, addr string) *Monitor {
	return &Monitor{factory:f,address:addr}
}

func (m *Monitor) Close() error {
	return srv.Shutdown(nil)
}
func (m *Monitor) Start(webDir string) {
	srv = &http.Server{Addr: m.address}
	http.Handle("/",http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/conn/getAll", func(w http.ResponseWriter, r *http.Request) {
		cs := make([]Conn,0)
		m.factory.ForEachAcceptedConnection(func(key cipher.PubKey, conn *factory.Connection) {
			c := conn.Connection.Connection
			content := Conn{Key:key.Hex(),
				SendBytes: c.GetSentBytes(),
				ReceivedBytes: c.GetReceivedBytes(),
				StartTime:conn.GetConnectTime(),
				LastAckTime: c.GetLastTime()}
			if conn.Connection.Connection.IsTCP() {
				content.Type = "tcp"
			}else {
				content.Type = "udp"
			}
			cs = append(cs,content)
		})
		js,err := json.Marshal(Conns{Conns:cs})
		if err != nil {
			http.Error(w,err.Error(),500)
			return
		}
		w.Header().Add("Access-Control-Allow-Origin","*")
		w.Write(js)
	})
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server: ListenAndServe() error: %s", err)
		}
	}()
	log.Debugf("http server listen on %s",m.address)
}