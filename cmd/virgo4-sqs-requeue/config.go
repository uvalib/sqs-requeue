package main

import (
	"flag"
	"log"
)

// ServiceConfig defines all of the service configuration parameters
type ServiceConfig struct {
	InQueueName       string
	OutQueueName      string
	MessageBucketName string
	PollTimeOut       int64
}

// LoadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func LoadConfiguration() *ServiceConfig {

	var cfg ServiceConfig

	flag.StringVar(&cfg.InQueueName, "inqueue", "", "Inbound queue name")
	flag.StringVar(&cfg.OutQueueName, "outqueue", "", "Output directory name")
	flag.StringVar(&cfg.MessageBucketName, "bucket", "", "Oversize message bucket name")
	flag.Int64Var(&cfg.PollTimeOut, "pollwait", 15, "Poll wait time (in seconds)")

	flag.Parse()

	if len(cfg.InQueueName) == 0 {
		log.Fatalf("InQueueName cannot be blank")
	}
	if len(cfg.OutQueueName) == 0 {
		log.Printf("OutQueueName cannot be blank")
	}
	if len(cfg.MessageBucketName) == 0 {
		log.Fatalf("MessageBucketName cannot be blank")
	}

	log.Printf("[CONFIG] InQueueName          = [%s]", cfg.InQueueName)
	log.Printf("[CONFIG] OutQueueName         = [%s]", cfg.OutQueueName)
	log.Printf("[CONFIG] MessageBucketName    = [%s]", cfg.MessageBucketName)
	log.Printf("[CONFIG] PollTimeOut          = [%d]", cfg.PollTimeOut)

	return &cfg
}

//
// end of file
//
