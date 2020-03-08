package main

import (
	"net"
	"encoding/hex"
	"bytes"
	"log"
	"time"
	"crypto/sha256"

	"github.com/lightningnetwork/lnd/watchtower/wtwire"
	"github.com/lightningnetwork/lnd/watchtower/blob"

	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/brontide"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)


type Client struct {
	Priv *btcec.PrivateKey
	Conn *brontide.Conn
}


// connect to watchtower server
func (c *Client) connect(toAddr *lnwire.NetAddress) {
	if c.Conn != nil {
		return
	}

	// dial(connection establishing) defined in brontide through the connection instance
	conn, err := brontide.Dial(c.Priv, toAddr, net.Dial)
	if err != nil {
		log.Printf("Unable to connect: %v", err)
		return
	}
	c.Conn = conn
}

// send Message to watchtower server
func (c *Client) sendMessage(toAddr *lnwire.NetAddress, msg wtwire.Message) {
	if c.Conn == nil { c.connect(toAddr) }

	var b bytes.Buffer
	wtwire.WriteMessage(&b, msg, 0)
	c.Conn.Write(b.Bytes())
}

// read Message from watchtower server
func (c *Client) readMessage() wtwire.Message {
	rawMsg, err := c.Conn.ReadNextMessage()
	if err != nil {
		log.Printf("Unable to read raw message: %v", err)
		return nil
	}

	msgReader := bytes.NewReader(rawMsg)
	msg, err := wtwire.ReadMessage(msgReader, 0)
	if err != nil {
		log.Printf("Unable to read message: %v", err)
		return nil
	}
	return msg
}

// prepare to send CreateSession Msg, StateUpdate Msg and DeleteSession Msg
func (c *Client) prepareToCommunicate(serverAddr *lnwire.NetAddress) {
	// define Init messege for test
	localFeatures := lnwire.NewRawFeatureVector()
	//localFeatures.Set(wtwire.AltruistSessionsOptional)
	localFeatures.Set(wtwire.AltruistSessionsRequired)
	initMsg := wtwire.NewInitMessage(localFeatures, *chaincfg.TestNet3Params.GenesisHash)
	log.Printf("client local features(InitMsg): %v", initMsg.ConnFeatures)
	log.Printf("client genesishash(InitMsg): %v", initMsg.ChainHash)

	// send Init message
	c.sendMessage(serverAddr, initMsg)

	// read Init message
	replyMsg := c.readMessage()
	initReplyMsg := replyMsg.(*wtwire.Init)
	log.Printf("server local features(InitReplyMsg): %v", initReplyMsg.ConnFeatures)
	log.Printf("server genesishash(InitReplyMsg): %v", initReplyMsg.ChainHash)
}

func (c *Client) SendCreateSessionMsg(serverAddr *lnwire.NetAddress, createSessionMsg *wtwire.CreateSession) {
	c.prepareToCommunicate(serverAddr)

	// send CreateSession message
	c.sendMessage(serverAddr, createSessionMsg)

	// read CreateSession message
	replyMsg := c.readMessage()
	createSessionReply, ok := replyMsg.(*wtwire.CreateSessionReply)
	if !ok {
		log.Printf("Unable to read CreateSession message: %T", replyMsg)
		return
	}

	log.Printf("MsgType: %v", replyMsg.MsgType())
	log.Printf("LastApplied: %v", createSessionReply.LastApplied)
	log.Printf("Code: %v", createSessionReply.Code)
	log.Printf("Data(sweepPkScript): %v", createSessionReply.Data)

	c.Conn.Close()
	c.Conn = nil
}

func (c *Client) SendStateUpdateMsg(serverAddr *lnwire.NetAddress, stateUpdate *wtwire.StateUpdate) {
	c.prepareToCommunicate(serverAddr)

	log.Printf("SeqNum: %v", stateUpdate.SeqNum)
	c.sendMessage(serverAddr, stateUpdate)

	// read message
	replyMsg := c.readMessage()
	stateUpdateReplyMsg := replyMsg.(*wtwire.StateUpdateReply)
	log.Printf("MsgType: %v", replyMsg.MsgType())
	log.Printf("Code: %v", stateUpdateReplyMsg.Code)
	log.Printf("LastApplied: %v", stateUpdateReplyMsg.LastApplied)

	c.Conn.Close()
	c.Conn = nil
}

func (c *Client) SendDeleteSessionMsg(serverAddr *lnwire.NetAddress, deleteSessionMsg *wtwire.DeleteSession) {
        c.prepareToCommunicate(serverAddr)

        // send CreateSession message
        c.sendMessage(serverAddr, deleteSessionMsg)

        // read CreateSession message
        replyMsg := c.readMessage()
        log.Printf("MsgType: %v", replyMsg.MsgType())
        deleteSessionReply, ok := replyMsg.(*wtwire.DeleteSessionReply)
        log.Printf("Code: %v", deleteSessionReply.Code)
        if !ok {
                log.Printf("Unable to read DeleteSession message: %T", replyMsg)
                return
        }

        c.Conn.Close()
        c.Conn = nil
}


func TestMsgCreateSession() *wtwire.CreateSession {
	return &wtwire.CreateSession{
		BlobType:     blob.TypeAltruistCommit,
		// returns CreateSessionCodeRejectBlobType
		// because DisableReward=true when server created for now
		//BlobType:     blob.TypeRewardCommit,

		MaxUpdates:   1000,
		RewardBase:   0,
		RewardRate:   0,
		SweepFeeRate: 10000,
	}
}

func TestMsgsStateUpdate() [](*wtwire.StateUpdate) {
        // create Justice tx, breach hint and breach encrypted blob
        jk := &blob.JusticeKit{
                SweepAddress: make([]byte, 22),

                RevocationPubKey: [33]byte{},
                CommitToLocalSig: [64]byte{},

                LocalDelayPubKey: [33]byte{},
                CSVDelay: 144,

                CommitToRemotePubKey: [33]byte{},
                CommitToRemoteSig: [64]byte{},
        }

        h := sha256.New()
        h.Write([]byte("hello world"))
        breachTxID, err := chainhash.NewHash(h.Sum(nil))
        if err != nil {
                log.Printf("Unable to create breachTxID: %v", err)
                return nil
        }
        hint, key := blob.NewBreachHintAndKeyFromHash(breachTxID)
        encBlob, err := jk.Encrypt(key, blob.TypeAltruistCommit)
        if err != nil {
                log.Printf("Unable to create encBlob: %v", err)
                return nil
        }

	// return StateUpdate Msgs
        return [](*wtwire.StateUpdate){
                &wtwire.StateUpdate{SeqNum: 1, LastApplied: 0, Hint: hint, EncryptedBlob: encBlob},
                &wtwire.StateUpdate{SeqNum: 2, LastApplied: 1, Hint: hint, EncryptedBlob: encBlob},
                &wtwire.StateUpdate{SeqNum: 3, LastApplied: 2, Hint: hint, EncryptedBlob: encBlob},
        }
}

func TestMsgDeleteSession() *wtwire.DeleteSession {
	return &wtwire.DeleteSession{}
}

func main() {
	//log.SetFlags(log.Lmicroseconds)

	// define local
	clientPriv, _ := btcec.NewPrivateKey(btcec.S256())
	client := Client{
		Priv: clientPriv,
	}

	// define server
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9911")
	hexPubkey, _ := hex.DecodeString("03fd7221f6576bea7c20ff21822081e744b492d387be2d81e95a9a38396e3de2ed")
	hexPubkey33, _ := btcec.ParsePubKey(hexPubkey, btcec.S256())
	var serverAddr = &lnwire.NetAddress{
		IdentityKey: hexPubkey33,
		Address:     addr,
	}

	log.Printf("[CreateSession]")
	client.SendCreateSessionMsg(serverAddr, TestMsgCreateSession())

	log.Printf("------")
	log.Printf("[StateUpdate]")
	for _, stateUpdate := range TestMsgsStateUpdate() {
		client.SendStateUpdateMsg(serverAddr, stateUpdate)
		time.Sleep(1 * time.Second)
		log.Printf("--")
	}

	log.Printf("------")
	log.Printf("[DeleteSession]")
	client.SendDeleteSessionMsg(serverAddr, TestMsgDeleteSession())
}
