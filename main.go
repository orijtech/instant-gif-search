// Copyright 2017 orijtech. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/orijtech/otils"
)

func main() {
	var http1 bool
	var http1Port int
	flag.BoolVar(&http1, "http1", false, "run the server in HTTP1 mode")
	flag.IntVar(&http1Port, "port", 9887, "the port to run the HTTP1 server on")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	if http1 {
		addr := fmt.Sprintf(":%d", http1Port)
		log.Printf("http1Mode: running on address: %q\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Redirecting non-https traffic
	go func() {
		nonHTTPSAddr := ":80"
		nonHTTPSHandler := otils.RedirectAllTrafficTo("https://rspgifs.orijtech.com")
		err := http.ListenAndServe(nonHTTPSAddr, nonHTTPSHandler)
		if err != nil {
			log.Printf("http-redirector: failed to bind the redirector: %v", err)
		}
	}()

	domains := []string{
		"www.instant-gif.orijtech.com",
		"instant-gif.orijtech.com",
		"www.rspgifs.orijtech.com",
		"rspgifs.orijtech.com",
	}

	log.Printf("Running in http2 mode")
	log.Fatal(http.Serve(autocert.NewListener(domains...), http.DefaultServeMux))
}
