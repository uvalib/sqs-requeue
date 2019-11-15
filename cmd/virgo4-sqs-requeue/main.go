package main

import (
	"log"
	"os"
	"time"

	"github.com/uvalib/virgo4-sqs-sdk/awssqs"
)

//
// main entry point
//
func main() {

	log.Printf("===> %s service staring up <===", os.Args[0])

	// Get config params and use them to init service context. Any issues are fatal
	cfg := LoadConfiguration()

	// load our AWS_SQS helper object
	aws, err := awssqs.NewAwsSqs(awssqs.AwsSqsConfig{MessageBucketName: cfg.MessageBucketName})
	if err != nil {
		log.Fatal(err)
	}

	// get the queue handle from the queue name
	inQueueHandle, err := aws.QueueHandle(cfg.InQueueName)
	if err != nil {
		log.Fatal(err)
	}

	outQueueHandle, err := aws.QueueHandle(cfg.OutQueueName)
	if err != nil {
		log.Fatal(err)
	}

	count := uint(0)
	for {

		// wait for a batch of messages
		inMessages, err := aws.BatchMessageGet(inQueueHandle, awssqs.MAX_SQS_BLOCK_COUNT, time.Duration(cfg.PollTimeOut)*time.Second)
		if err != nil {
			log.Fatal(err)
		}

		// did we receive any?
		sz := len(inMessages)
		if sz != 0 {

			//log.Printf( "Received %d messages", sz )
			count += uint(sz)

			// make our outbound buffer
			outMessages := make([]awssqs.Message, 0, count)
			for _, m := range inMessages {
				outMessages = append(outMessages, *m.ContentClone())
			}

			opStatus, err := aws.BatchMessagePut(outQueueHandle, outMessages)
			if err != nil && err != awssqs.ErrOneOrMoreOperationsUnsuccessful {
				log.Fatal(err)
			}

			// check the operation results
			for ix, op := range opStatus {
				if op == false {
					log.Printf("WARNING: message %d failed to send to outbound queue", ix)
				}
			}

			// delete them all
			opStatus, err = aws.BatchMessageDelete(inQueueHandle, inMessages)
			if err != nil && err != awssqs.ErrOneOrMoreOperationsUnsuccessful {
				log.Fatal(err)
			}

			// check the operation results
			for ix, op := range opStatus {
				if op == false {
					log.Printf("WARNING: message %d failed to delete", ix)
				}
			}

		} else {
			log.Printf("Processed %d messages, terminating", count)
			break
		}
	}
}

//
// end of file
//
