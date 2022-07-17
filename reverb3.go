// Eva Hallermeier 337914121
// Shana Sebban 337912182
// Sara Spagnoletto 345990808

//Question 8.4

// TCP echo server.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

//wait group added in parameters of echo function
func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	//create waitGroup per connection
	var wg sync.WaitGroup
	input := bufio.NewScanner(c)
	for input.Scan() {
		wg.Add(1)
		//added defer function to close after echo
		go func() {
			defer wg.Done()
			echo(c, input.Text(), 2*time.Second)
		}()
	}
	//wait blocks until the WaitGroup counter is zero
	wg.Wait()
	//number of goroutines is zero
	log.Println("reverb3.")
	//closing only the half write
	if con, ok := c.(*net.TCPConn); ok {
		con.CloseWrite()
	} else {
		c.Close()
	}
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
