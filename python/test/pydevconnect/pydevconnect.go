/*
Copyright 2021 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Test utility to connect to a pydevd server to validate that it is working.
// Protocol: https://github.com/fabioz/PyDev.Debugger/blob/main/_pydevd_bundle/pydevd_comm.py
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CMD_RUN              = 101
	CMD_LIST_THREADS     = 102
	CMD_THREAD_CREATE    = 103
	CMD_THREAD_KILL      = 104
	CMD_THREAD_RUN       = 106
	CMD_SET_BREAK        = 111
	CMD_WRITE_TO_CONSOLE = 116
	CMD_VERSION          = 501
	CMD_RETURN           = 502
	CMD_ERROR            = 901
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Check that pydevd is running.\n")
		fmt.Printf("use: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	var conn net.Conn
	for i := 0; i < 60; i++ {
		var err error
		conn, err = net.Dial("tcp", os.Args[1])
		if err == nil {
			break
		}
		fmt.Printf("(sleeping) unable to connect to %s: %v\n", os.Args[1], err)
		time.Sleep(2 * time.Second)
	}

	pydb := newPydevdDebugConnection(conn)

	code, response := pydb.makeRequestWithResponse(CMD_VERSION, "pydevconnect")
	if code != CMD_VERSION {
		log.Fatalf("expected CMD_VERSION (%d) response (%q)", code, response)
	}
	if decoded, err := url.QueryUnescape(response); err != nil {
		log.Fatalf("CMD_VERSION response (%q): decoding error: %v", response, err)
	} else {
		fmt.Printf("version: %s", decoded)
	}

	pydb.makeRequest(CMD_RUN, "test")
}

type pydevdDebugConnection struct {
	conn   net.Conn
	reader *bufio.Reader
	msgID  int
}

func newPydevdDebugConnection(c net.Conn) *pydevdDebugConnection {
	return &pydevdDebugConnection{
		conn:   c,
		reader: bufio.NewReader(c),
		msgID:  1,
	}
}

func (c *pydevdDebugConnection) makeRequest(code int, arg string) {
	currMsgID := c.msgID
	c.msgID += 2 // outgoing requests should have odd msgID

	fmt.Printf("Making request: code=%d msgId=%d arg=%q\n", code, currMsgID, arg)
	fmt.Fprintf(c.conn, "%d\t%d\t%s\n", code, currMsgID, arg)
}

func (c *pydevdDebugConnection) makeRequestWithResponse(code int, arg string) (int, string) {
	currMsgID := c.msgID
	c.msgID += 2 // outgoing requests should have odd msgID

	fmt.Printf("Making request: code=%d msgId=%d arg=%q\n", code, currMsgID, arg)
	fmt.Fprintf(c.conn, "%d\t%d\t%s\n", code, currMsgID, arg)

	for {
		response, err := c.reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error receiving response: %v", err)
		}
		fmt.Printf("Received response: %q\n", response)

		// check response
		tsv := strings.Split(response, "\t")
		if len(tsv) != 3 {
			log.Fatalf("invalid response: expecting three tab-separated components: %q", response)
		}

		code, err = strconv.Atoi(tsv[0])
		if err != nil {
			log.Fatalf("could not parse response code: %q", tsv[0])
		}

		responseID, err := strconv.Atoi(tsv[1])
		if err != nil {
			log.Fatalf("could not parse response ID: %q", tsv[1])
		} else if responseID == currMsgID {
			return code, tsv[2]
		}

		// handle commands sent to us
		switch code {
		case CMD_THREAD_CREATE:
			fmt.Printf("CMD_THREAD_CREATE: %s\n", tsv[2:])

		default:
			log.Fatalf("Unknown/unhandled code %d: %q", code, tsv[2:])
		}
	}
}
