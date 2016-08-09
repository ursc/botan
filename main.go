package main

import (
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	app := NewApplication("botan.toml", "botan.json", "botan.map")
	app.Start()
}
