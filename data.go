package main

import (
	"time"
	"os"
	"math/rand"
	"bufio"
	"fmt"
	"github.com/valyala/fasthttp"
)

type spoof struct {
	count int
}

type header struct {
	advanced bool
	a_headers map[string]string
	proxies   []string // New field to store the loaded proxies
}

func NewSpoof(count int) *spoof {
	return &spoof{
		count: count,
	}
}

func NewHeader(advanced bool, a_headers map[string]string, proxies []string) *header { // Added proxies parameter
	return &header{
		advanced: advanced,
		a_headers: a_headers,
		proxies:   proxies, // Initialize proxies field
	}
}

var (
	Source = rand.NewSource(time.Now().Unix())
	Random = rand.New(Source)
)

func randomIp() string {
	return fmt.Sprintf("%d.%d.%d.%d", Random.Intn(255), Random.Intn(255), Random.Intn(255), Random.Intn(255))
}

func (s spoof) spoofS(req *fasthttp.Request) {
	rangeip := []string{}
	file, err := os.Open("src/rate_headers.txt")
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		rangeip = append(rangeip, scan.Text())
	}
	for i := 0; i < s.count; i++ {
		ip := randomIp()
		req.Header.Add(string(rangeip[Random.Intn(len(rangeip))]), ip)
	}
}

func (h header) headers(req *fasthttp.Request) {
	uas := []string{}
	file, err := os.Open("src/user-agents.txt")
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		uas = append(uas, scan.Text())
	}
	req.Header.Add("User-Agent", uas[Random.Intn(len(uas))])
	if h.advanced == true {
		for head, value := range h.a_headers {
			req.Header.Add(head, value)
		}
	}

	// Add a random proxy from the loaded proxies
	if len(h.proxies) > 0 {
		proxy := h.proxies[Random.Intn(len(h.proxies))]
		req.SetRequestURI(target.url)
		req.SetHost(proxy)
	}
}
