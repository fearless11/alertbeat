package beater

import (
	"fmt"
	"log"
	"time"

	"we.com/vera.jiang/alertbeat/conf"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"
)

type Alertbeat struct {
	client publisher.Client
}

var (
	T8TToggle   chan int
	BasicToggle chan int
	BasicNagios chan string
	T8TMsg      chan *conf.T8TAlert
	BasicMsg    chan *conf.BasicAlert
)

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	T8TToggle = make(chan int)
	BasicToggle = make(chan int)
	BasicNagios = make(chan string)
	T8TMsg = make(chan *conf.T8TAlert)
	BasicMsg = make(chan *conf.BasicAlert)
	bt := &Alertbeat{}
	return bt, nil
}

func (bt *Alertbeat) Run(b *beat.Beat) error {
	// logp.Info("parsebeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()
	for {
		var basicalertmsg, t8talertmsg string
		select {
		case <-BasicToggle:
			basicmsg := <-BasicMsg
			basicnagios := <-BasicNagios
			basicalertmsg = fmt.Sprintf("%s%s", basicmsg.AlarmID, basicmsg.Content)
			bt.publisherBasic(basicmsg, basicnagios)
		case <-T8TToggle:
			t8tmsg := <-T8TMsg
			t8talertmsg = fmt.Sprintf("%s%s", t8tmsg.Labels, t8tmsg.Annotations)
			bt.publisherT8T(t8tmsg)
		}
		if len(basicalertmsg) == 0 && len(t8talertmsg) == 0 {
			break
		}
	}
	return nil
}

func (bt *Alertbeat) publisherT8T(msg *conf.T8TAlert) {
	alertmsg := fmt.Sprintf("%s%s", msg.Labels, msg.Annotations)
	if conf.Config.Debug {
		log.Println("[DEBUG] publisherT8T:", alertmsg)
	}

	event := common.MapStr{
		"@timestamp": common.Time(time.Now().UTC()),
		"type":       msg.Labels["type"],
		"proj":       msg.Labels["proj"],
		"alertname":  msg.Labels["alertname"],
		"lv":         msg.Labels["level"],
		"env":        msg.Labels["env"],
		"from":       msg.Labels["from"],
		"count":      msg.Annotations["count"],
		"host":       msg.Annotations["host"],
		"interface":  msg.Annotations["interface"],
		"nagios":     msg.Labels["nagios"],
		"msg":        alertmsg,
	}
	go bt.client.PublishEvent(event)
}

func (bt *Alertbeat) publisherBasic(msg *conf.BasicAlert, nagios string) {

	if conf.Config.Debug {
		alertmsg := fmt.Sprintf("%s%s", msg.AlarmID, msg.Content)
		log.Println("[DEBUG] publisherBasic:", alertmsg)
	}

	event := common.MapStr{
		"@timestamp": common.Time(time.Now().UTC()),
		"type":       "basic",
		"proj":       msg.AlarmID,
		"msg":        msg.Content,
		"nagios":     nagios,
	}
	go bt.client.PublishEvent(event)
}

func (bt *Alertbeat) Stop() {
	bt.client.Close()
	close(T8TToggle)
	close(T8TMsg)
	close(BasicMsg)
	close(BasicToggle)
	close(BasicNagios)
}
