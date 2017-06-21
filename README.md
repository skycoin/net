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

