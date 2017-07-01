# Skycoin Networking Framework



## Protocol


                      +--+--------+--------+--------------------+
    msg protocol      |  |        |        |                    |
                      +-++-------++-------++---------+----------+
                        |        |        |          |
                        v        |        v          v
                      msg type   |     msg len    msg body
                       1 byte    v     4 bytes
                              msg seq
                              4 bytes
    
    
    
                      +-----------+--------+---------+----------+
    normal msg        |01|  seq   |  len   |  pubkey |  body    |
                      +-----------+--------+---------+----------+
    
    
                      +-----------+--------+---------+
    reg msg           |02|  seq   |  len   |  pubkey |
                      +-----------+--------+---------+
    
    
                      +-----------+
    ack msg           |80|  seq   |
                      +-----------+
    
    
                      +--+
    ping msg          |81|
                      +--+
    
    
                      +--+
    pong msg          |82|
                      +--+


## Flow Chart

    node                                    server                                    node
    
    +--+                                     +--+                                     +--+
    |  |                                     |  |                                     |  |
    |  | +------+register+pubkey+----------> |  | <------+register+pubkey+----------+ |  |
    |  |                                     |  |                                     |  |
    |  | <-------------+ack+---------------+ |  | +-------------+ack+---------------> |  |
    |  |                                     |  |                                     |  |
    |  | +------+send+msg+to+pubkey+-------> |  | +------+forward+msg+to+pubkey+----> |  |
    |  |                                     |  |                                     |  |
    |  | <-------------+ack+---------------+ |  | <-------------+ack+---------------+ |  |
    |  |                                     |  |                                     |  |
    |  | <------+forward+resp+msg+---------+ |  | <------+resp+msg+to+pubkey+-------+ |  |
    |  |                                     |  |                                     |  |
    |  | +-------------+ack+---------------> |  | +-------------+ack+---------------> |  |
    |  |                                     |  |                                     |  |
    |  |                                     |  |                                     |  |
    |  |                                     |  |                                     |  |
    |  |                                     |  |                                     |  |
    +--+                                     +--+                                     +--+


## TCP/UDP Client Example

### Client 0xf1
```
factory := client.NewClientConnectionFactory()
factory.Connect("tcp", ":8080", cipher.PubKey([33]byte{0xf1}))
conn := factory.Dial(cipher.PubKey([33]byte{0xf2}))
conn.Out <- []byte("Hello 0xf2")

for {
    select {
    case m, ok := <-conn.In:
        if !ok {
            log.Println("conn closed")
            return
        }
        log.Printf("msg In %x", m)
    }
}
```

### Client 0xf2
```
factory := client.NewClientConnectionFactory()
factory.SetIncomingCallback(func(conn *client.ClientConnection, data []byte) bool {
    log.Printf("msg from %s In %s", conn.Key.Hex(), data)

    go func() {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("recovered err %v", err)
            }
        }()
        for {
            select {
            case m, ok := <-conn.In:
                if !ok {
                    return
                }
                log.Printf("received msg %s", m)
            }
        }
    }()


    // return true for save this conn in factory so can use conn.Out for resp something
    // otherwise conn.Out can not be used, because no receiver goroutine exists
    return true
})
factory.Connect("udp", ":8081", cipher.PubKey([33]byte{0xf2}))
```

## Server Example

```
s := server.New(":8080", ":8081")
go func() {
    log.Println("listening udp")
    if err := s.ListenUDP(); err != nil {
        panic(err)
    }
}()
log.Println("listening tcp")
if err := s.ListenTCP(); err != nil {
    panic(err)
}
```