package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	banner = `

     .´/                              
    / (           .----------------.    
    [ ]░░░░░░░░░░░|// RESPOUNDER //|    
    ) (           '----------------'` + "\n" +
		"    '-' \n"

	timeoutSec = 3
)

func main() {
	fmt.Fprintln(os.Stderr, banner)

	interfaces, _ := net.Interfaces()
	for _, inf := range interfaces {
		checkResponderOnInterface(inf)
	}
}

func checkResponderOnInterface(inf net.Interface) string {
	json := ""
	addrs, _ := inf.Addrs()
	if len(addrs) > 0 {
		ip := addrs[0].(*net.IPNet).IP
		if ip.String() != "127.0.0.1" {
			fmt.Printf("%-10s Sending probe from %s...\t", "["+inf.Name+"]", ip)
			responderAddr := sendLLMNRProbe(ip)
			if responderAddr != "" {
				fmt.Printf("responder detected at %s\n", responderAddr)
			} else {
				fmt.Println("responder not detected")
			}
		}
	}
	return json
}

// Creates and sends a LLMNR request to the UDP multicast address.
func sendLLMNRProbe(ip net.IP) string {
	responderIP := ""
	// 2 byte random transaction id eg. 0x8e53
	rand.Seed(time.Now().UnixNano())
	randomTransactionId := fmt.Sprintf("%04x", rand.Intn(65535))

	// LLMNR request in raw bytes
	// TODO: generate a new computer name evertime instead of the
	// hardcoded value 'awierdcomputername'
	llmnrRequest := randomTransactionId +
		"0000000100000000000012617769657264636f6d70757465726e616d650000010001"
	n, _ := hex.DecodeString(llmnrRequest)

	remoteAddr := net.UDPAddr{IP: net.ParseIP("224.0.0.252"), Port: 5355}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip})
	if err != nil {
		fmt.Println("Couldn't bind to a UDP interface. Bailing out!")
	}

	defer conn.Close()
	_, _ = conn.WriteToUDP(n, &remoteAddr)

	conn.SetReadDeadline(time.Now().Add(timeoutSec * time.Second))
	buffer := make([]byte, 1024)
	_, clientIP, err := conn.ReadFromUDP(buffer)
	if err == nil { // no timeout (or any other) error
		responderIP = strings.Split(clientIP.String(), ":")[0]
	}
	return responderIP
}
