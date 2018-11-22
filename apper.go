package apper_go

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"mining-pool/service/eth"
	"mining-pool/types"

	//"mining-pool/service/eth"
	//"mining-pool/types"

	"apper-go/logger"
	//"mining-pool/service/eth"
	//"mining-pool/types"
	"sync"
	"time"
)

const natsIP = "47.99.72.199:4222"
var log = logger.Log
type Apper struct {
	conn nats.Conn
}

// singleton mode
func GetApper() (*Apper, error) {
	return &Apper{}, nil
}

func (a *Apper) Connect(url string) (error) {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}else {
		a.conn = *nc
	}
	return err
}

func (a *Apper) Start(path string) ( string ,  error) {
	confs,_ :=eth.GetETHConf(eth.ETH_yaml_path)
	n_struct := types.Nats_data{*confs,"struct"}
	//序列化
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)//创建编码器
	err1 := encoder.Encode(&n_struct)//编码
	if err1!=nil{
		log.Fatal(err1)
	}
	//	fmt.Printf("序列化后：%x\n",buffer.Bytes())
	nc:=a.conn
	defer nc.Close()
	// Create a unique subject name
	uniqueReplyTo := nats.NewInbox()
	// Listen for a single response
	sub, err := nc.SubscribeSync(uniqueReplyTo)
	if err != nil {
		log.Fatal(err)
	}
	// Send the request
	if err := nc.PublishRequest("cmd", uniqueReplyTo,  buffer.Bytes()); err != nil {
		log.Fatal(err)
	}
	// Read the reply
	msg, err := sub.NextMsg(time.Minute*1)
	if err != nil {
		log.Fatal(err)
	}
	// Use the response
	//log.Printf("Reply: %s", msg.Data)
	str := string(msg.Data[:])
	// Close the connection
	nc.Close()
	return str,err
}

func (*Apper) Stop(transactionID string) error {
	return nil
}

func (*Apper) Terminate(pass string) {

}

func (a *Apper) GetVal(key, transactionID string) (interface{}, error) {
	n_struct := types.Nats_data1{key,transactionID}
	//序列化
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)//创建编码器
	err1 := encoder.Encode(&n_struct)//编码
	if err1!=nil{
		log.Fatal(err1)
	}
	nc:=a.conn
	defer nc.Close()
	// Create a unique subject name
	uniqueReplyTo := nats.NewInbox()
	// Listen for a single response
	sub, err := nc.SubscribeSync(uniqueReplyTo)
	if err != nil {
		log.Fatal(err)
	}
	// Send the request
	if err := nc.PublishRequest("cmd", uniqueReplyTo,buffer.Bytes()); err != nil {
		log.Fatal(err)
	}
	// Read the reply
	msg, err := sub.NextMsg(time.Minute*1)
	if err != nil {
		log.Fatal(err)
	}
	// Use the response
	//log.Printf("Reply: %s", msg.Data)
	byteEn:=msg.Data
	decoder := gob.NewDecoder(bytes.NewReader(byteEn)) //创建解密器
	var nterf interface{}
	err2 := decoder.Decode(&nterf)//解密
	if err2!=nil{
		log.Fatal(err2)
	}
	//fmt.Println("反序列化后：",nterf)
	nc.Close()
	return nterf,err

}

func (a *Apper) Ready(transactionID string) bool {
	var  boo bool = false
	nc:= a.conn
	defer nc.Close()
	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe
	if _, err := nc.Subscribe("res_"+transactionID, func(m *nats.Msg) {
		if string(m.Data[:]) == transactionID {
			boo = true
		}
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}
	// Wait for a message to come in
	wg.Wait()
	// Close the connection
	nc.Close()
	return boo
}
