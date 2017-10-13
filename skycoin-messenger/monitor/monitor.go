package monitor

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"encoding/json"
)

type Conn struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	SendBytes   uint64 `json:"send_bytes"`
	RecvBytes   uint64 `json:"recv_bytes"`
	LastAckTime int64  `json:"last_ack_time"`
	StartTime   int64  `json:"start_time"`
}
type NodeServices struct {
	Type        string `json:"type"`
	Apps        []App  `json:"apps"`
	SendBytes   uint64 `json:"send_bytes"`
	RecvBytes   uint64 `json:"recv_bytes"`
	LastAckTime int64  `json:"last_ack_time"`
	StartTime   int64  `json:"start_time"`
}
type App struct {
	Index      int64    `json:"index"`
	Key        string   `json:"key"`
	Attributes []string `json:"attributes"`
}

var srv *http.Server

var (
	NULL = "null"
)
var (
	BAD_REQUEST  = 400
	NOT_FOUND    = 404
	SERVER_ERROR = 500
)

type Monitor struct {
	factory *factory.MessengerFactory
	address string
}

func New(f *factory.MessengerFactory, addr string) *Monitor {
	return &Monitor{factory: f, address: addr}
}

func (m *Monitor) Close() error {
	return srv.Shutdown(nil)
}
func (m *Monitor) Start(webDir string) {
	srv = &http.Server{Addr: m.address}
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/conn/getAll", func(w http.ResponseWriter, r *http.Request) {
		cs := make([]Conn, 0)
		m.factory.ForEachAcceptedConnection(func(key cipher.PubKey, conn *factory.Connection) {
			content := Conn{
				Key:         key.Hex(),
				SendBytes:   conn.GetSentBytes(),
				RecvBytes:   conn.GetReceivedBytes(),
				StartTime:   conn.GetConnectTime(),
				LastAckTime: conn.GetLastTime()}
			if conn.IsTCP() {
				content.Type = "TCP"
			} else {
				content.Type = "UDP"
			}
			cs = append(cs, content)
		})
		js, err := json.Marshal(cs)
		if err != nil {
			http.Error(w, err.Error(), SERVER_ERROR)
			return
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(js)
	})
	http.HandleFunc("/conn/getNodeStatus", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "please use post method", BAD_REQUEST)
			return
		}
		key, err := cipher.PubKeyFromHex(r.FormValue("key"))
		if err != nil {
			http.Error(w, err.Error(), BAD_REQUEST)
			return
		}
		c, ok := m.factory.GetConnection(key)
		if !ok {
			http.Error(w, "No connection is found", NOT_FOUND)
			return
		}

		nodeService := NodeServices{
			SendBytes:   c.GetSentBytes(),
			RecvBytes:   c.GetReceivedBytes(),
			StartTime:   c.GetConnectTime(),
			LastAckTime: c.GetLastTime()}
		if c.IsTCP() {
			nodeService.Type = "TCP"
		} else {
			nodeService.Type = "UDP"
		}
		ns := c.GetServices()
		if ns != nil {
			var index int64 = 1
			for _, v := range ns.Services {
				app := App{Index: index, Key: v.Key.Hex(), Attributes: v.Attributes}
				nodeService.Apps = append(nodeService.Apps, app)
				index ++
			}
		}
		js, err := json.Marshal(nodeService)
		if err != nil {
			http.Error(w, err.Error(), SERVER_ERROR)
			return
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(js)
	})
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server: ListenAndServe() error: %s", err)
		}
	}()
	log.Debugf("http server listen on %s", m.address)
}
