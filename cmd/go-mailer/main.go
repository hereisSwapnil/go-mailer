package main

import (
	"log"
	"sync"

	"github.com/hereisSwapnil/go-mailer/internal/consumer"
	"github.com/hereisSwapnil/go-mailer/internal/producer"
	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func main() {
	EmailChannel := make(chan types.Recipient)

	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(EmailChannel)
		if err := producer.LoadRecipientsUsingCsv("demo/emails.csv", EmailChannel); err != nil {
			log.Fatalf("Error loading recipients: %v", err)
		}
	}()

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			if err := consumer.EmailWorker(workerId, EmailChannel); err != nil {
				log.Fatalf("Error processing recipients: %v", err)
			}
		}(i)
	}

	wg.Wait()
}