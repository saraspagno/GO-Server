// Eva Hallermeier 337914121
// Shana Sebban 337912182
// Sara Spagnoletto 345990808

// Question 8.8

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

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	var wg sync.WaitGroup
	input := bufio.NewScanner(c)
	timeout := time.NewTimer(10 * time.Second)
	//making a channel for received text
	ch := make(chan string)

	go func() {
		for input.Scan() {
			ch <- input.Text()
		}
		close(ch)
	}()

	for {
		select {
		case text, ok := <-ch:
			if ok {
				wg.Add(1)
				//resetting the timeout
				timeout.Reset(10 * time.Second)
				go func() {
					defer wg.Done()
					echo(c, text, 2*time.Second)
				}()
			} else {
				wg.Wait()
				c.Close()
				return
			}
		// if timeout exceeded
		case <-timeout.C:
			timeout.Stop()
			c.Close()
			log.Println("reverb4.")
			return
		}
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
