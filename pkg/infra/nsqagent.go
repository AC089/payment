/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-19 13:43:43
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-19 13:43:46
 */
package infra

import (
	"time"

	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

type NsqAgent struct {
	NsqUrl      string `yaml:"nsq_url"`
	Timeout		uint32 `yaml:"timeout"`
	TopicPre	    string `yaml:"topic_pre"`
	producer    *nsq.Producer
}

var _agent *NsqAgent

func CreateNsqInstance() *NsqAgent {
	_agent = new(NsqAgent)
	return _agent
}

func GetNsqInstance() *NsqAgent {
	return _agent
}

func (l *NsqAgent) InitProducer() error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Duration(l.Timeout) * time.Second
	p, err := nsq.NewProducer(l.NsqUrl, cfg)
	if err != nil {
		return err
	}
	l.producer = p
	return nil
}

func (l *NsqAgent) Publish(topic string, msg []byte) {
	RE:
		err := l.producer.Publish(topic, msg)
		if err != nil { //重连
			log.Error(err)
			time.Sleep(time.Second * 3)
			goto RE
		}
}