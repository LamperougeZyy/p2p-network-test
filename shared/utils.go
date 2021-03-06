package shared

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func MessageIn(c Conn, b []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(b, m) // 对json编码进行解码，保存到Message结构体中

	//本来这里要对得到的文件进行是否加密的判断，暂时先省略，消息都以明文进行传输
	if err != nil {
		log.Print(err)
		return m, err
	}

	return m, nil
}

func MessageOut(c Conn, m *Message) ([]byte, error) {
	//获得message的json编码
	b, err := json.Marshal(m)
	if err != nil {
		return b, err
	}

	//判断message是否需要加密，需要加密的则要先获得Secret
	/*
		if m.Encrypt {
			var s [32]byte
			s, err = c.GetSecret()
			if err != nil {
				return b, fmt.Errorf("cannot encrypt with an empty secret")
			}

			b, err = crypto.Encrypt(b, s)
			if err != nil {
				return b, err
			}
		}*/

	return b, nil
}

func route(client Client, cs Conns, c Conn, m *Message) {
	switch m.Type {
	case "greeting":
		return greetingHandler(client, c, m)
	case "register":
		return registerHandler(client, c, m)
	case "establish":
		return establishHandler(client, c, m)
	case "connect":
		return connectHandler(client, c, m)
	case "key":
		return keyHandler(client, c, m)
	case "message":
		return messageHandler(client, c, m)
	}
	return nil, nil
}

func CreateMessageCallback(client Client) func(Conns, Conn, *Message) {
	return func(cs Conns, c Conn, m *Message) {
		res, err := route(client, cs, c, m)
		if err != nil {
			fmt.Println(err)
			client.GetLog().Fatal(err)
		}

		if res != nil {
			c.Send(res)
		}
	}
}

func GenPort() string {
	return ":" + strconv.Itoa(rand.Intn(65535-10000)+10000)
}
