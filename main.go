package main

import "github.com/aleksannder/url-shortener/common"

func main() {
	server := Server{
		cfg: common.GetConfig(),
	}
	server.Run()
}
