package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"strings"

	"we.com/vera.jiang/alertbeat/beater"
	"we.com/vera.jiang/alertbeat/conf"
	"we.com/vera.jiang/alertbeat/output"
)

func filterCharacters(body string) string {
	re, _ := regexp.Compile("[|]")
	return fmt.Sprintf("%s", re.ReplaceAllString(body, ""))
}

func BasicParse(body string) error {
	msg := filterCharacters(body)
	var alert *conf.BasicAlert
	err := json.Unmarshal([]byte(msg), &alert)
	if err != nil {
		return err
	}

	alertmsg := fmt.Sprintf("%s%s", alert.AlarmID)
	if len(alertmsg) == 0 {
		return fmt.Errorf("%v", "parse alert message fail")
	}

	nagios := "0"
	if checkBasicAlert(alert) {
		output.BasicNotify(alert)
		nagios = "1"
		if conf.Config.Debug {
			log.Println("[DEBUG] BasicNotify:", alert.AlarmID, alert.Content)
		}
	}

	beater.BasicToggle <- 1
	beater.BasicMsg <- alert
	beater.BasicNagios <- nagios
	return nil
}

func T8TParse(body string) error {
	msg := filterCharacters(body)
	var alert *conf.T8TAlert
	err := json.Unmarshal([]byte(msg), &alert)
	if err != nil {
		return err
	}

	alertmsg := fmt.Sprintf("%s%s", alert.Labels, alert.Annotations)
	if len(alertmsg) == 0 {
		return fmt.Errorf("%v", "parse alert message fail")
	}

	if checkT8TAlert(alert) {
		output.T8TNotify(alert)
		if conf.Config.Debug {
			log.Println("[DEBUG] T8TNotify:", alert.Labels, alert.Annotations)
		}
	}

	beater.T8TToggle <- 1
	beater.T8TMsg <- alert
	return nil
}

func checkKey(str, key string) bool {

	reg := regexp.MustCompile(key)
	l := reg.FindAllString(str, -1)
	if len(l) == 0 {
		return false
	}
	return true
}

func checkT8TAlert(alert *conf.T8TAlert) bool {

	alert.Labels["nagios"] = "1"
	u := "unknown"

	t, exist := alert.Labels["type"]
	if !exist {
		alert.Labels["type"] = u
		t = u
	}

	lv, exist := alert.Labels["level"]
	if !exist {
		alert.Labels["level"] = "unknown"
		lv = "unknown"
	}

	if t == "java" {
		_, exist = alert.Annotations["level"]
		if exist {
			alert.Labels["level"] = alert.Annotations["level"]
			lv = alert.Annotations["level"]
		}
	}

	if t == "web" {
		_, exist = alert.Annotations["proj"]
		if exist {
			alert.Labels["proj"] = alert.Annotations["proj"]
		}
	}

	if !checkKey(lv, `(?i:warn|CRITICAL|unknown)`) {
		alert.Labels["nagios"] = "0"
		return false
	}

	return true
}

func checkBasicAlert(alert *conf.BasicAlert) bool {

	for _, id := range conf.Config.Ignore {
		if alert.AlarmID == id {
			return false
		}
	}

	content := alert.Content
	alert.Content = strings.Replace(content, "\n", ";", -1)
	if checkKey(content, `Type: RECOVERY .* State: OK `) {
		return false
	}
	return true
}
