# Sample code to check watchtower protocol of lnd

# environment

* Golang
    * v1.13
* boot lnd with watchtower server on testnet
    * v0.9.1
    * watchtower setup reference: https://github.com/lightningnetwork/lnd/blob/master/docs/watchtower.md

# setup & run

```
$ git clone https://github.com/takaya-imai/sample_code_for_watchtower_protocol_of_lnd.git $GOPATH/src/github.com/takaya-imai/sample_code_for_watchtower_protocol_of_lnd
$ cd $GOPATH/src/github.com/takaya-imai/sample_code_for_watchtower_protocol_of_lnd
$ go get -u github.com/golang/dep/cmd/dep
$ dep ensure
$ go run sample_code.go
```

# example of code outputs

```
2020/03/08 14:49:38 [CreateSession]
2020/03/08 14:49:38 client local features(InitMsg): &{map[0:true]}
2020/03/08 14:49:38 client genesishash(InitMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:38 server local features(InitReplyMsg): &{map[1:true]}
2020/03/08 14:49:38 server genesishash(InitReplyMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:38 MsgType: MsgCreateSessionReply
2020/03/08 14:49:38 LastApplied: 0
2020/03/08 14:49:38 Code: CodeOK
2020/03/08 14:49:38 Data(sweepPkScript): []
2020/03/08 14:49:38 ------
2020/03/08 14:49:38 [StateUpdate]
2020/03/08 14:49:38 client local features(InitMsg): &{map[0:true]}
2020/03/08 14:49:38 client genesishash(InitMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:38 server local features(InitReplyMsg): &{map[1:true]}
2020/03/08 14:49:38 server genesishash(InitReplyMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:38 SeqNum: 1
2020/03/08 14:49:38 MsgType: MsgStateUpdateReply
2020/03/08 14:49:38 Code: CodeOK
2020/03/08 14:49:38 LastApplied: 1
2020/03/08 14:49:39 --
2020/03/08 14:49:39 client local features(InitMsg): &{map[0:true]}
2020/03/08 14:49:39 client genesishash(InitMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:39 server local features(InitReplyMsg): &{map[1:true]}
2020/03/08 14:49:39 server genesishash(InitReplyMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:39 SeqNum: 2
2020/03/08 14:49:39 MsgType: MsgStateUpdateReply
2020/03/08 14:49:39 Code: CodeOK
2020/03/08 14:49:39 LastApplied: 2
2020/03/08 14:49:40 --
2020/03/08 14:49:40 client local features(InitMsg): &{map[0:true]}
2020/03/08 14:49:40 client genesishash(InitMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:40 server local features(InitReplyMsg): &{map[1:true]}
2020/03/08 14:49:40 server genesishash(InitReplyMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:40 SeqNum: 3
2020/03/08 14:49:40 MsgType: MsgStateUpdateReply
2020/03/08 14:49:40 Code: CodeOK
2020/03/08 14:49:40 LastApplied: 3
2020/03/08 14:49:41 --
2020/03/08 14:49:41 ------
2020/03/08 14:49:41 [DeleteSession]
2020/03/08 14:49:41 client local features(InitMsg): &{map[0:true]}
2020/03/08 14:49:41 client genesishash(InitMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:41 server local features(InitReplyMsg): &{map[1:true]}
2020/03/08 14:49:41 server genesishash(InitReplyMsg): 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943
2020/03/08 14:49:41 MsgType: MsgDeleteSessionReply
2020/03/08 14:49:41 Code: CodeOK
```
