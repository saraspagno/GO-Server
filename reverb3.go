// Eva Hallermeier 337914121
// Shana Sebban 337912182
// Sara Spagnoletto 345990808


/*
Part 1 - Exercice 8.4: modify reverb2.go and use a sync.WaitGroup per connection
to count the number of active echo goroutines.
When it falls to zero, close the write half of the TCP
connection of netcat4.
Verify that netcat4 waits for the final echoes of multiple concurrent shouts,
even after the standard input has been closed (by ctrl-Z as described in class).
*/


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
	"time"
	"sync" //for waitgroup
)

//wait group added in parameters of echo function
func echo(c net.Conn, shout string, delay time.Duration, waitg *sync.WaitGroup) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
	defer waitg.Done() //decrements the waitgroup counter by one: echo action finished so less goroutine
}

func handleConn(c net.Conn) {
	waitg := sync.WaitGroup{} //create waitgroup per connection
	input := bufio.NewScanner(c)
	for input.Scan() {
		waitg.Add(1)
		go echo(c, input.Text(), 2*time.Second, &waitg) //waitg added in parameter of echo function
	}
	//time.Sleep(5 * time.Second)
	waitg.Wait() //Wait blocks until the WaitGroup counter is zero.
	log.Println("reverb3. ") //number of goroutines is zero
	// NOTE: ignoring potential errors from input.Err()
	//close the write half of the TCP connection of netcat4.
	c.Close() // close quitely
}

//!-

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
