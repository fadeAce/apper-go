package apper_go

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"time"
)

type Apper struct {
	conn *nats.Conn
}

// singleton mode
func GetApper() (*Apper, error) {
	return &Apper{}, nil
}

func (a *Apper) Connect(url string) error {
	nc, err := nats.Connect(url)
	if err != nil {
		Log.Error(err)
	} else {
		a.conn = nc
	}
	return err
}

func (a *Apper) Start(path string) (string, error) {
	var conf Conf
	f, err := ioutil.ReadFile(path)
	err = yaml.Unmarshal(f, &conf)
	nStruct := Nats_data{conf, "struct"}
	//序列化
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) //创建编码器
	err1 := encoder.Encode(&nStruct)   //编码
	if err1 != nil {
		Log.Error(err1)
	}
	//	fmt.Printf("序列化后：%x\n",buffer.Bytes())
	nc := a.conn
	defer nc.Close()
	// Send the request
	msg := &nats.Msg{}
	if msg, err = nc.Request("cmd", buffer.Bytes(), time.Minute*1); err != nil {
		Log.Error(err)
	}
	str := string(msg.Data[:])
	// Close the connection
	nc.Close()
	return str, err
}

func (*Apper) Stop(transactionID string) error {
	return nil
}

func (*Apper) Terminate(pass string) {

}

func (a *Apper) GetVal(key, transactionID string) (interface{}, error) {
	nStruct := Nats_data1{key, transactionID}
	//序列化
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) //创建编码器
	err1 := encoder.Encode(&nStruct)   //编码
	if err1 != nil {
		Log.Error(err1)
	}
	nc := a.conn
	defer nc.Close()
	// Create a unique subject name
	uniqueReplyTo := nats.NewInbox()
	// Listen for a single response
	sub, err := nc.SubscribeSync(uniqueReplyTo)
	if err != nil {
		Log.Error(err)
	}
	// Send the request
	if err := nc.PublishRequest("cmd", uniqueReplyTo, buffer.Bytes()); err != nil {
		Log.Error(err)
	}
	// Read the reply
	msg, err := sub.NextMsg(time.Minute * 1)
	if err != nil {
		Log.Error(err)
	}
	// Use the response
	//Log.Printf("Reply: %s", msg.Data)
	byteEn := msg.Data
	decoder := gob.NewDecoder(bytes.NewReader(byteEn)) //创建解密器
	var nterf interface{}
	err2 := decoder.Decode(&nterf) //解密
	if err2 != nil {
		Log.Error(err2)
	}
	//fmt.Println("反序列化后：",nterf)
	nc.Close()
	return nterf, err

}

func (a *Apper) Ready(transactionID string) bool {
	var boo bool = false
	nc := a.conn
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
		Log.Error(err)
	}
	// Wait for a message to come in
	wg.Wait()
	// Close the connection
	nc.Close()
	return boo
}
