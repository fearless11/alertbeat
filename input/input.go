package input

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"io/ioutil"

	"we.com/vera.jiang/alertbeat/conf"
	"we.com/vera.jiang/alertbeat/parse"
)

func Start() {
	listen := conf.Config.Web
	http.HandleFunc("/", sayHi)
	http.HandleFunc("/v1/t1", handleT8TAlertMsg)
	http.HandleFunc("/v1/basic", handleBasicAlarm)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayHi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", "hi")
}

func handleT8TAlertMsg(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	err = parse.T8TParse(string(body))
	if err != nil {
		io.WriteString(w, err.Error())
	}
}

func handleBasicAlarm(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	err = parse.BasicParse(string(body))
	if err != nil {
		io.WriteString(w, err.Error())
	}
}
