package producer

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

var conf *Config
var sess *session.Session
var seq int64 = 1

func Init(path string) (err error) {
	conf = &Config{}
	err = LoadConfig(conf, path)
	if err != nil {
		return
	}
	sess, err = session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(conf.AWSAccessKeyId, conf.AWSSecretKey, ""),
			Region:      &conf.Region,
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return
	}
	return
	if err != nil {
		return
	}
	return
}

type MqBody struct {
	Seq          int64  `json:"seq"`
	FromApp      string `json:"from_app"`
	FromNode     string `json:"from_node"`
	ToNode       string `json:"to_node"`
	ToApp        string `json:"to_app"`
	Uid          uint64 `json:"uid"`
	FromHostPort string `json:"from_host_port"`
	ToHostPort   string `json:"to_host_port"`
	FromIp       string `json:"from_ip"`
	ToIp         string `json:"to_ip"`
	Count        uint64 `json:"count"`
	IsEnd        bool   `json:"is_end"`
}

// from add , to reduce
func Send(body MqBody) (err error) {
	svc := sqs.New(sess)
	body.Seq = seq
	b, err := json.Marshal(&body)
	if err != nil {
		return
	}
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(b)),
		QueueUrl:    aws.String(conf.QueueURL),
	})
	if err != nil {
		return
	}
	seq++
	return
}
