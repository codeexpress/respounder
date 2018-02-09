package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	Banner = `

     .´/
    / (           .----------------.
    [ ]░░░░░░░░░░░|// RESPOUNDER //|
    ) (           '----------------'
    '-'
`

	Version    = 1.0
	TimeoutSec = 3
	BcastAddr  = "224.0.0.252"
	LLMNRPort  = 5355
)

var (
	// stdout is default output
	outFile = os.Stdout

	// default logger is set to abyss
	logger = log.New(ioutil.Discard, "", 0)

	// argument flags
	jsonPtr = flag.Bool("json", false,
		`Prints a JSON to STDOUT if a responder is detected on
	network. Other text is sent to STDERR`)

	debugPtr = flag.Bool("debug", false,
		`Creates a debug.log file with a trace of the program`)
)

func main() {
	initFlags()

	fmt.Fprintln(os.Stderr, Banner)

	interfaces, _ := net.Interfaces()
	logger.Println("======== Starting RESPOUNDER ========")
	logger.Printf("List of all interfaces: \n %+v\n", interfaces)

	var resultMap []map[string]string

	for _, inf := range interfaces {
		detailsMap := checkResponderOnInterface(inf)
		if len(detailsMap) > 0 {
			resultMap = append(resultMap, detailsMap)
		}
	}

	if *debugPtr {
		fmt.Fprintln(os.Stderr, "Debug file 'debug.log' created.")
	}

	if *jsonPtr {
		resultJSON, _ := json.Marshal(resultMap)
		fmt.Println(string(resultJSON))
	}
	logger.Println("======== Ending RESPOUNDER Session ========")
}

// Test presence of responder on a given interface
func checkResponderOnInterface(inf net.Interface) map[string]string {
	var json map[string]string
	addrs, _ := inf.Addrs()
	logger.Printf("List of all addresses on interface [%s]: %+v\n",
		inf.Name, addrs)
	ip := getValidIPv4Addr(addrs)
	logger.Printf("Bind IP address for interface %+v is %+v\n",
		inf.Name, ip)

	if ip != nil {
		fmt.Fprintf(outFile, "%-10s Sending probe from %s...\t",
			"["+inf.Name+"]", ip)
		responderIP := sendLLMNRProbe(ip)
		if responderIP != "" {
			fmt.Fprintf(outFile, "responder detected at %s\n", responderIP)
			json = map[string]string{
				"interface":   inf.Name,
				"sourceIP":    ip.String(),
				"responderIP": responderIP,
			}
		} else {
			fmt.Fprintln(outFile, "responder not detected")
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

	remoteAddr := net.UDPAddr{IP: net.ParseIP(BcastAddr), Port: LLMNRPort}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip})
	if err != nil {
		fmt.Println("Couldn't bind to a UDP interface. Bailing out!")
		logger.Printf("Bind error: %+v\nSource IP: %v\n", err, ip)
		fmt.Println(err)
	}

	defer conn.Close()
	_, _ = conn.WriteToUDP(n, &remoteAddr)

	conn.SetReadDeadline(time.Now().Add(TimeoutSec * time.Second))
	buffer := make([]byte, 1024)
	bytes, clientIP, err := conn.ReadFromUDP(buffer)
	if err == nil { // no timeout (or any other) error
		responderIP = strings.Split(clientIP.String(), ":")[0]
		logger.Printf("Data received on %s from responder IP %s: %x\n",
			ip, clientIP, buffer[:bytes])
	} else {
		logger.Printf("Error getting response:  %s\n", err)
	}
	return responderIP
}

// From all the IP addresses of this interface,
// extract the IPv4 address where we'll bind to
func getValidIPv4Addr(addrs []net.Addr) net.IP {
	var ip net.IP
	for _, addr := range addrs { // amongst all addrs,
		ip = addr.(*net.IPNet).IP.To4() // pick the IPv4 addr
		if ip != nil && ip.String() != "127.0.0.1" {
			break
		}
	}
	return ip
}

// parses cmd line flag and set appropriate variables
func initFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Respounder version %1.1f\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: $ respounder [-json] [-debug]")
		fmt.Fprintf(os.Stderr, "\n\nFlags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *jsonPtr {
		outFile = os.Stderr
	}
	if *debugPtr {
		f, err := os.OpenFile("debug.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		logger = log.New(f, "", 0)
		logger.SetPrefix("[" + time.Now().Format("02-Jan-2006 15:04:05 MST") + "]: ")
	}
}
