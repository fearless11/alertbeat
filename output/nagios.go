package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/astaxie/beego/httplib"

	"io/ioutil"

	"we.com/vera.jiang/alertbeat/conf"
)

type Message struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Status  int    `json:"status"`
	Output  string `json:"output"`
}

func notifyNagiosAPI(proj, alertmsg, ttype string) error {

	url := conf.Config.Nagios.Addr

	var group string
	if ttype == "web" {
		group = "WebGroup"
	} else if ttype == "basic" {
		group = "BasicGroup"
	} else {
		group = "JavaBizGroup"
	}

	data := Message{
		Host:    group,
		Service: proj,
		Status:  1,
		Output:  alertmsg,
	}
	msg, _ := json.Marshal(data)

	if conf.Config.Debug {
		log.Println("[DEBUG] notifyNagiosAPI:", string(msg))
	}

	req := httplib.Post(url)
	req.Header("Content-Type", "application/json")
	req.Body(msg)
	resq, err := req.Response()
	if err != nil {
		return err
	}
	contents, _ := ioutil.ReadAll(resq.Body)
	defer resq.Body.Close()

	if !strings.Contains(string(contents), "true") {
		return fmt.Errorf("nagiosAPIReturn:%v", string(contents))
	}
	return nil
}

func notify(proj, ttype, msg string) {

	if err := createTemplateFile(proj, ttype); err != nil {
		log.Println("[ERROR] createTemplateFile:", err)
	}

	//reloadNagios()

	if err := notifyNagiosAPI(proj, msg, ttype); err != nil {
		log.Println("[ERROR] notifyNagiosAPI:", err)
		return
	}
}

func reloadNagios() {
	reload := "/etc/init.d/nagios reload"
	cmd := exec.Command("/bin/sh", "-c", reload)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("[ERROR] reloadNagios:", out.String(), err)
	}

	if conf.Config.Debug {
		log.Println("[DEBUG] reloadNagios end ...")
	}
}

func T8TNotify(alert *conf.T8TAlert) {
	msg := fmt.Sprintf("%s%s", alert.Labels, alert.Annotations)
	ttype := alert.Labels["type"]
	proj := alert.Labels["proj"]
	go notify(proj, ttype, msg)
}

func BasicNotify(alert *conf.BasicAlert) {
	msg := fmt.Sprintf("%s", alert.Content)
	go notify(alert.AlarmID, "basic", msg)
}
