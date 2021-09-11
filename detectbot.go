/*
	Detects meris bot IPs by testing if they have both port 2000 & port 5678 open.
	2021-11-09, Stefan Behte, https://github.com/craig/merisbot-detect
	Enjoy.
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

// connect to all ports
func raw_connect(host string, ports []string) bool {

	// iterate over ports
	for _, port := range ports {
		timeout := time.Second
		sock := host + ":" + port

		conn, err := net.DialTimeout("tcp", sock, timeout)

		if err != nil {
			return false
		}
		if conn != nil {
			// close connection
			defer conn.Close()
			//fmt.Println("Opened", sock)
		}
	}
	return true
}

func testhost(text string) {
	s := []string{"2000", "5678"}

	if raw_connect(text, s) {
		fmt.Println(text)
		//} else {
		//	fmt.Println(text + " ")
	}
}

func main() {

	for {
		var text string

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text = scanner.Text()

			if len(text) > 0 {
				go testhost(text)
			}
			time.Sleep(1 * time.Millisecond)
		}

	}
}
