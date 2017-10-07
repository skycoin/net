package monitor

import (
	"net/http"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
)

type Conns struct {
	Conns []Conn `json:"conns"`
}
type Conn struct {
	Key string `json:"key"`
	Type string `json:"type"`
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

func (m *Monitor) Start() {
	srv = &http.Server{Addr: m.address}
	http.Handle("/",http.FileServer(http.Dir("./github.com/skycoin/net/skycoin-messenger/monitor/web/dist")))
	http.HandleFunc("/conn/getAll", func(w http.ResponseWriter, r *http.Request) {
		conns := m.factory.GetRegConns()
		cs := make([]Conn,0,len(conns))
		for k, v := range conns {
			c := Conn{Key:k.Hex()}
			if v.Connection.Connection.IsTCP() {
				c.Type = "tcp"
			}else {
				c.Type = "udp"
			}
			cs = append(cs,c)
		}
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