package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/elastic/beats/libbeat/beat"

	"we.com/vera.jiang/alertbeat/beater"
	"we.com/vera.jiang/alertbeat/conf"
	"we.com/vera.jiang/alertbeat/input"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()
	cfg := flag.String("k", "alertbeat.yml", "configuration file")
	help := flag.Bool("h", false, "help")
	flag.Parse()
	handleHelp(*help)
	handleConfig(*cfg)
}

func main() {
	go input.Start()
	err := beat.Run("alertbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	err := conf.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
