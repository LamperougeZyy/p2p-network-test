package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	serverTCPPort = ":7001"
	serverUDPPort = ":9001"
)

var (
	serverTCPIP = "0.0.0.0"
	serverUDPIP = "127.0.0.1"
	serverIP    = flag.String("serverIP", "", "IP address of rendezvous server")
)

func reDo(c int, f func() error) error {
	var err error
	for i := 0; i < c; i++ {
		err = f()
		if err == nil {
			return err
		}
	}

	return err
}

func reDo(n int, f func() error) error {
	var err error
	for i := 0; i < n; i++ {
		if err = f(); err == nil {
			return err
		}
	}
	return err
}

func main() {
	flag.Parse()

	if *serverIP != "" {
		serverTCPIP = *serverIP
		serverUDPIP = *serverIP
	}
	fmt.Print("\n  UDP Hole Punching v0.0.1 👊\n\n")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error

	var username string
	for username == "" || len(username) > 32 {
		fmt.Println("	Username(<= 32 chars)")
		fmt.Print("	> ")
		_, err = fmt.Scanln(&username)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("\n")
	}

	//创建记录历史消息的文档
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fName := fmt.Sprintf("%s/history-%s.txt", wd, username)
	hf, err := os.Create(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer hf.Close()

	h := shared.NewHistory(hf)

	var c shared.Client
	var sAddr *net.UDPAddr
	var addr *net.UDPAddr
	sAddr, err = net.ResolveUDPAddr("udp", serverUDPIP+serverUDPPort)
	if err != nil {
		log.Fatal(err)
	}

	//建立client实例，建立五次，只要一次成功就继续
	err = reDo(5, func() error {
		addr, err = net.ResolveUDPAddr("udp", shared.GenPort())
		if err != nil {
			log.Fatal(err)
		}

		c, err = udp_client.New(username, addr, sAddr)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	c.OnRegistered(registeredCallback)
	c.OnConnecting(connectingCallback)
	c.OnConnected(createConnectedCallback(h))
	c.OnMessage(createMessageCallback(h))

	fmt.Println("  ID")
	fmt.Printf("  > %s\n\n", c.GetSelf().ID)

	err = c.Start()
	if err != nil {
		log.Fatal(err)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-exit)

	c.Stop()
}
