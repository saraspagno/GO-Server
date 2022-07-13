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
func echo(c net.Conn, shout string, delay time.Duration, wait *sync.WaitGroup) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
	defer wait.Done() //decrements the waitGroup counter by one: echo action finished so less goroutine
}

func handleConn(c net.Conn) {
	//create waitGroup per connection
	wait := sync.WaitGroup{}
	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Println("inside the scan")
		wait.Add(1)
		//wait added in parameter of echo function
		go echo(c, input.Text(), 2*time.Second, &wait)
	}
	//wait blocks until the WaitGroup counter is zero
	wait.Wait()
	//number of goroutines is zero
	log.Println("reverb3. ")
	// NOTE: ignoring potential errors from input.Err()
	//close the write half of the TCP connection of netcat4
	c.Close()
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
