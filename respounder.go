package main

import (
	"crypto/sha1"
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

	Version         = 1.1
	TimeoutSec      = 3
	BcastAddr       = "224.0.0.252"
	LLMNRPort       = 5355
	DefaultHostname = "aweirdcomputername"
)

const (
	def          = 0x00
	newHostname  = 0x01
	randHostname = 0x02
)

var (
	// stdout is default output
	outFile = os.Stdout

	// default logger is set to abyss
	logger = log.New(ioutil.Discard, "", 0)

	// argument flags
	jsonPtr = flag.Bool("json", false,
		`Prints a JSON to STDOUT if a responder is detected on
        the network. Other text is sent to STDERR`)

	debugPtr = flag.Bool("debug", false,
		`Creates a debug.log file with a trace of the program`)

	hostnamePtr = flag.String("hostname", DefaultHostname,
		`Hostname to search for`)
	randHostnamePtr = flag.Bool("rhostname", false,
		`Searches for a hostname comprised of random string instead
        of the default hostname ("`+DefaultHostname+`")`)

	hostnameType byte
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	initFlags()
	flag.Parse()

	if *hostnamePtr != "aweirdcomputername" {
		hostnameType = newHostname
	} else if *randHostnamePtr {
		hostnameType = randHostname
	} else {
		hostnameType = def
	}

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
	var cName string
	responderIP := ""
	// 2 byte random transaction id eg. 0x8e53
	randomTransactionID := fmt.Sprintf("%04x", rand.Intn(65535))

	switch hostnameType {
	case def, newHostname:
		cName = string(*hostnamePtr)
	case randHostname:
		cName = randomHostname()
	}

	cNameLen := fmt.Sprintf("%02x", len(cName))
	encCName := hex.EncodeToString([]byte(cName))
	// LLMNR request in raw bytes
	llmnrRequest := randomTransactionID +
		"00000001000000000000" + cNameLen + encCName + "0000010001"
	n, _ := hex.DecodeString(llmnrRequest)

	remoteAddr := net.UDPAddr{IP: net.ParseIP(BcastAddr), Port: LLMNRPort}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip})
	if err != nil {
		fmt.Println("Couldn't bind to a UDP interface. Bailing out!")
		logger.Printf("Bind error: %+v\nSource IP: %v\n", err, ip)
		fmt.Println(err)
		logger.Printf("LLMNR request payload was: %x\n", llmnrRequest)
	}

	defer conn.Close()
	_, _ = conn.WriteToUDP(n, &remoteAddr)

	conn.SetReadDeadline(time.Now().Add(TimeoutSec * time.Second))
	buffer := make([]byte, 1024)
	bytes, clientIP, err := conn.ReadFromUDP(buffer)
	if err == nil { // no timeout (or any other) error
		responderIP = strings.Split(clientIP.String(), ":")[0]
		logger.Printf("LLMNR request payload was: %x\n", n)
		logger.Printf("Data received on %s from responder IP %s: %x\n",
			ip, clientIP, buffer[:bytes])
	} else {
		logger.Printf("Error getting response:  %s\n", err)
	}
	return responderIP
}

// Calculate random hostname by taking random lenght
// of the SHA1 of current time.
func randomHostname() string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	h := sha1.New()
	h.Write([]byte(currentTime))
	bs := h.Sum(nil)
	randomSlice := bs[:(rand.Intn(len(bs)-3) + 3)]
	randomName := fmt.Sprintf("%x\n", randomSlice)
	return randomName
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
		fmt.Fprintf(os.Stderr, "Usage: $ respounder [-json] [-debug] [-hostname testhostname | -rhostname]")
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
