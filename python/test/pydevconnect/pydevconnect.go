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
)

const (
	CMD_VERSION = 501
)

// Test utility to connect to a pydevd server to validate that it is working.
// Protocol: https://github.com/fabioz/PyDev.Debugger/blob/main/_pydevd_bundle/pydevd_comm.py
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Check that pydevd is running.\n")
		fmt.Printf("use: %s host:port\n", os.Args[0])
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatalf("unable to connect to %s: %v", os.Args[1], err)
	}

	msgID := 1 // client requests should be odd
	fmt.Fprintf(conn, "%d\t%d\t%s\n", CMD_VERSION, msgID, "test")

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("error receiving response: %v", err)
	}

	// check response
	tsv := strings.Split(response, "\t")
	if len(tsv) != 3 {
		log.Fatalf("invalid response: expecting three tab-separated components: %q", response)
	}

	if code, err := strconv.Atoi(tsv[0]); err != nil {
		log.Fatalf("could not parse response code: %q", tsv[0])
	} else if code != CMD_VERSION {
		log.Fatalf("expected CMD_VERSION(%d) response code but received %q", CMD_VERSION, tsv[0])
	}

	if reponseID, err := strconv.Atoi(tsv[1]); err != nil {
		log.Fatalf("could not parse response ID: %q", tsv[1])
	} else if reponseID != msgID {
		log.Fatalf("expected response ID %d but received %q", msgID, tsv[1])
	}

	if decoded, err := url.QueryUnescape(tsv[2]); err != nil {
		log.Fatalf("could not decode response text: %q: %v", tsv[2], err)
	} else {
		fmt.Printf("version: %s", decoded) // no \n required as response includes trailing \n
	}
}
