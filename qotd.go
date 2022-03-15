package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/mdns"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	RFC865MaxLength = 512
)

var log = logrus.New()

func init() {
	rand.Seed(time.Now().UnixNano())
	log.Formatter = new(logrus.TextFormatter)
}

func main() {
	app := cli.NewApp()
	app.Name = "QOTD"
	app.Usage = "Run a QOTD Server"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "port,p", Value: "3333", Usage: "port to bind the server to"},
		cli.BoolFlag{Name: "strict", Usage: "quotes served in RFC 865 strict mode"},
		cli.BoolFlag{Name: "no-tcp", Usage: "server does not listen on tcp"},
		cli.BoolFlag{Name: "no-udp", Usage: "server does not listen on udp"},
		cli.BoolFlag{Name: "no-mdns", Usage: "server does not advertise over mdns"},
	}

	app.Action = func(c *cli.Context) {
		if len(c.Args()) != 1 {
			log.Fatal("Server must be started with a path to a file with quotes")
		}

		port := c.String("port")
		fileName := c.Args()[0]
		quotes := loadQuotes(fileName)
		strictMode := c.Bool("strict")
		startUdp := !c.Bool("no-udp")
		startTcp := !c.Bool("no-tcp")
		advertiseService := !c.Bool("no-mdns")

		if strictMode {
			port = "17"
			startTcp = true
			startUdp = true
		}

		if advertiseService {
			advertisedService := advertiseQOTDService(startTcp, startUdp, port)
			defer advertisedService.Shutdown()
		}

		if startUdp {
			go listenForUdp(port, quotes, strictMode)
		}

		if startTcp {
			go listenForTcp(port, quotes, strictMode)
		}

		if startTcp || startUdp {
			// Keep this busy
			for {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			log.Fatal("Server not started on TCP or UDP, don't pass both --no-tcp and --no-udp")
		}
	}

	app.Run(os.Args)
}

func advertiseQOTDService(advertiseTcp bool, advertiseUdp bool, port string) *mdns.Server {
	host, _ := os.Hostname()
	println(host)
	service := &mdns.MDNSService{
		Instance: host,
		Service:  "_qotd._tcp",
		IPs:      []net.IP{net.IP([]byte{0, 0, 0, 0})},
		Port:     3333,
		TXT:      []string{"Local web server"},
	}
	// service.Init()

	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	return server
}

func listenForTcp(port string, quotes []string, strictMode bool) {
	tcp, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer tcp.Close()
	log.Info("TCP: QOTD Server Started on Port " + port)
	for {
		conn, err := tcp.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go serveRandomQuote(conn, quotes, strictMode)
	}
}

func listenForUdp(port string, quotes []string, strictMode bool) {
	udpService := ":" + port
	updAddr, udpErr := net.ResolveUDPAddr("udp", udpService)
	if udpErr != nil {
		log.Fatal("Error Resolving UDP Address: ", udpErr.Error())
		os.Exit(1)
	}
	updSock, udpErr := net.ListenUDP("udp", updAddr)
	if udpErr != nil {
		log.Fatal("Error listening: ", udpErr.Error())
		os.Exit(1)
	}
	log.Info("UDP: QOTD Server Started on Port " + port)
	defer updSock.Close()
	for {
		serveUDPRandomQuote(updSock, quotes, strictMode)
	}
}

func serveUDPRandomQuote(conn *net.UDPConn, quotes []string, strictMode bool) {
	requestUUID, err := uuid.NewV4()
	buf := make([]byte, 512)

	_, addr, err := conn.ReadFromUDP(buf[0:])

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  addr.String(),
	}).Info("UDP Request Received")

	quote, quoteId := randomQuoteFormattedForDelivery(quotes, strictMode)
	conn.WriteToUDP([]byte(quote), addr)

	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  addr.String(),
	}).Info("UDP Quote #" + strconv.Itoa(quoteId) + " Served")
}

func serveRandomQuote(conn net.Conn, quotes []string, strictMode bool) {
	requestUUID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("TCP Request Received")

	quote, quoteId := randomQuoteFormattedForDelivery(quotes, strictMode)
	conn.Write([]byte(quote))
	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("TCP Quote #" + strconv.Itoa(quoteId) + " Served")

	conn.Close()
	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("Connection Closed")
}

func randomQuoteFormattedForDelivery(quotes []string, strictMode bool) (string, int) {
	quoteId := rand.Intn(len(quotes))
	var quote = quotes[quoteId]
	if strictMode && len(quote) > RFC865MaxLength {
		// 3 bytes for ..., 2 bytes for closing \r\n
		quote = string([]byte(quote)[0 : RFC865MaxLength-6])
		quote = quote + "..."
	}

	quote = quote + "\r\n"
	return quote, quoteId
}

func loadQuotes(path string) []string {
	u, err := url.Parse(path)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return loadQuotesFromFile(path)
	} else {
		return loadQuotesFromHTTP(path)
	}
}

func loadQuotesFromHTTP(address string) []string {
	resp, err := http.Get(address)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	rawData, err := ioutil.ReadAll(resp.Body)
	quotes := strings.Split(string(rawData), "\n%\n")
	return quotes
}

func loadQuotesFromFile(fileName string) []string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	quotes := strings.Split(string(file), "\n%\n")
	return quotes
}
