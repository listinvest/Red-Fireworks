package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
)


const (
	// Stolen from https://godoc.org/golang.org/x/net/internal/iana,
	// can't import "internal" packages
	ProtocolICMP = 1
	//ProtocolIPv6ICMP = 58
)

// EncryptDecrypt runs a XOR encryption on the input string, encrypting it if it hasn't already been,
// and decrypting it if it has, using the key provided.
func EncryptDecrypt(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i % len(key)])
	}

	return output
}

//func createHash(key string) string {
//	hasher := md5.New()
//	hasher.Write([]byte(key))
//	return hex.EncodeToString(hasher.Sum(nil))
//}

func encrypt(data []byte, passphrase string) []byte {
	return data
}

func decrypt(data []byte, passphrase string) []byte {
	return data
}

func sendICMP(addr string,data string) (*net.IPAddr, error) {
	// Start listening for icmp replies
	c, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// Resolve any DNS (if used) and get the real IP of the target
	dst, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		panic(err)
		return nil, err
	}

	encrptedData := encrypt([]byte(data), key)
	fmt.Println("Sent:")
	fmt.Println(encrptedData)
	// Make a new ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: 0xffff, Seq: 1, //<< uint(seq), // TODO
			Data: []byte(encrptedData),
		},
	}

	b, err := m.Marshal(nil)
	if err != nil {
		return dst, err
	}

	// Send it
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return dst, err
	} else if n != len(b) {
		return dst, fmt.Errorf("got %v; want %v", n, len(b))
	}

	return dst, err
}

func recICMP() (error) {
	// Start listening for icmp replies
	c, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		return err
	}
	defer c.Close()
	// Wait for a reply
	reply := make([]byte, 1024*2)
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return  err
	}
	// Parse
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return  nil
	default:
		b, _ := rm.Body.Marshal(0)
		rData := decrypt(b, key)
		fmt.Println("Rec:")
		fmt.Println(rData)
		return fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}

}


func sendResponse(m string, conn *net.UDPConn, addr *net.UDPAddr) {
	_,err := conn.WriteToUDP([]byte(m), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}


func main() {
	
	/// Play with main according to your needs: 
	/// Who do u want to start first ? Rec first ? balbalbal 
	
	
	p2 := func(addr string,data string){
		dst, err := sendICMP(addr, data)
		if err != nil {
			log.Printf("Ping %s (%s): %s\n", addr, err)
			return
		}
		log.Printf("Ping %s (%s):\n", addr, dst)
	}
	
	
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP: net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		//p2("127.0.0.1","ICMP - What's up?")
		
		_,remoteaddr,err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err !=  nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse("UPD - Good job, received.",ser, remoteaddr)
		_ = recICMP()
		p2("127.0.0.1","ICMP - Good job, received.")
	}


}







