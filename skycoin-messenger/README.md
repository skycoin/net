# Skycoin Messenger

## TCP/UDP Client Example

```
client := client.New()
//create a tcp client
err := c.Connect("tcp", ":8080")
//create a udp client
//err := c.Connect("udp", ":8081")

if err != nil {
   panic(err)
}

//read and write from socket
go c.Loop()

//send a reg msg about pubkey to server
c.Reg(cipher.PubKey([33]byte{0xf1}))

//require a conn by pubkey
factory := client.NewClientConnectionFactory(c)
conn := factory.GetConn(cipher.PubKey([33]byte{0xf2}))

//write msg to pubkey 0xf2
conn.Write([]byte("Hello 0xf2"))

//msg in
for {
   select {
   case m, ok := <-c.In:
      if !ok {
         log.Println("conn closed")
         return
      }
      log.Printf("msg In %x", m)
   }
}
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