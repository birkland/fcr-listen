package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-stomp/stomp"
	"log"
	"strings"
	"time"
)

var (
	host    string
	port    int
	subname string
)

func main() {

	log.SetFlags(log.Flags() | log.Lshortfile)

	flag.IntVar(&port, "port", 61613, "STOMP connection port")
	flag.StringVar(&host, "host", "localhost", "STOMP connection host or IP")
	flag.StringVar(&subname, "subscribe", "/topic/fedora", "Queue or topic to subscribe to")

	flag.Parse()

	for err := perform(listen); err != nil; err = perform(listen) {
		log.Println("Task error, restarting", err)
	}
}

func perform(task func(*stomp.Conn) error) (err error) {
	var conn *stomp.Conn

	// Try infinitely to connect
	for conn, err = connect(); err != nil; conn, err = connect() {
		log.Println("Could not connect!", err)
		time.Sleep(1 * time.Second)
	}

	defer func() {
		if e := conn.Disconnect(); e != nil {
			log.Println("Error disconnecting: ", e)
		}
	}()

	return task(conn)
}

func connect() (*stomp.Conn, error) {
	return stomp.Dial("tcp", fmt.Sprint(host, ":", port),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.HeartBeat(10*time.Second, 30*time.Second),
		stomp.ConnOpt.HeartBeatError(10*time.Second))
}

func listen(conn *stomp.Conn) (err error) {

	sub, err := conn.Subscribe(subname, stomp.AckClient)
	if err != nil {
		return fmt.Errorf("Could not subscribe to %s: %v", subname, err)
	}
	log.Println("Subscribed to", subname)

	for msg := range sub.C {

		fmt.Println()

		for i := 0; i < msg.Header.Len(); i++ {
			printHeader(fmt.Sprint(msg.Header.GetAt(i)))
		}

		printBody(string(msg.Body[:]))

		if msg.ShouldAck() {
			err = conn.Ack(msg)
			if err != nil {
				log.Println("Could not Ack message", msg.Destination, err)
			}
		}
	}

	return
}

func printHeader(txt string) {
	line := color.New(color.FgCyan)

	if strings.HasPrefix(txt, "org.fcrepo") {
		line.Add(color.Bold)
	}

	colorPrint(line, txt)
}

func printBody(txt string) {
	colorPrint(color.New(color.FgRed, color.FgYellow, color.Bold), txt)
}

func colorPrint(c *color.Color, txt string) {
	if _, err := c.Println(txt); err != nil {
		fmt.Println(txt)
	}
}
