# Skycoin Messenger

## TCP/UDP Client Example

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

```
factory := client.NewClientConnectionFactory()
factory.SetIncomingCallback(func(conn *client.ClientConnection, data []byte) bool {
   log.Printf("msg from %s In %s", conn.Key.Hex(), data)
   
   // return true for save this conn in factory so can use conn.Out for resp something
   // otherwise conn.Out can not be used, because no receiver goroutine exists
   return true
})
factory.Connect("udp", ":8081", cipher.PubKey([33]byte{0xf2}))
```

## Server Example

```
func main() {
   log.SetFlags(log.LstdFlags|log.Lshortfile)
   s := server.New()
   server.DefaultConnectionFactory.TCPClientHandler = handleTCP
   server.DefaultConnectionFactory.UDPClientHandler = handleUDP
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
}

func handleTCP(connection *conn.TCPConn) {
   for {
      select {
      case m, ok := <-connection.In:
         if !ok {
            log.Println("conn closed")
            return
         }
         log.Printf("msg in %x", m)
         key := cipher.NewPubKey(m[:33])
         c := server.DefaultConnectionFactory.GetConn(key.Hex())
         if c == nil {
            log.Printf("pubkey not found in factory %x", m)
            continue
         }
         publicKey := connection.GetPublicKey()
         copy(m[:33], publicKey[:])
         c.Write(m)
      }
   }
}

func handleUDP(connection *conn.UDPConn) {
   for {
      select {
      case m, ok := <-connection.In:
         if !ok {
            log.Println("udp conn closed")
            return
         }
         log.Printf("msg in %x", m)
         key := cipher.NewPubKey(m[:33])
         c := server.DefaultConnectionFactory.GetConn(key.Hex())
         if c == nil {
            log.Printf("pubkey not found in factory %x", m)
            continue
         }
         publicKey := connection.GetPublicKey()
         copy(m[:33], publicKey[:])
         c.Write(m)
      }
   }
}
```