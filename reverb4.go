// Eva Hallermeier 337914121
// Shana Sebban 337912182
// Sara Spagnoletto 345990808

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

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
	input := bufio.NewScanner(c)
	timeout := time.NewTimer(10 * time.Second)
	text := make(chan string)

	var wg sync.WaitGroup
	go func() {
		for input.Scan() {
			text <- input.Text()
		}
		close(text)
	}()

	for {
		select {
		case t, ok := <-text:
			if ok {
				wg.Add(1)
				timeout.Reset(10 * time.Second)
				go func() {
					defer wg.Done()
					echo(c, t, 1*time.Second)
				}()
			} else {
				wg.Wait()
				c.Close()
				return
			}
		case <-timeout.C:
			timeout.Stop()
			c.Close()
			fmt.Println("disconnect silent client")
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
