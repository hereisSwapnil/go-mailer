package main

import (
	"log/slog"
	"sync"

	"github.com/hereisSwapnil/go-mailer/internal/config"
	"github.com/hereisSwapnil/go-mailer/internal/consumer"
	"github.com/hereisSwapnil/go-mailer/internal/producer"
	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func main() {
	cfg := config.LoadConfig()
	slog.Info("Config loaded successfully! 🚀")

	EmailChannel := make(chan types.Recipient)

	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(EmailChannel)
		if err := producer.LoadRecipientsUsingCsv(cfg.Data.CSVFilePath, EmailChannel); err != nil {
			slog.Error("Error loading recipients!", "error", err)
		}
	}()

	numWorkers := 10
	slog.Info("Starting workers...", "numWorkers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			if err := consumer.EmailWorker(workerId, EmailChannel, cfg); err != nil {
				slog.Error("Error processing recipients!", "error", err)
			}
		}(i)
	}

	wg.Wait()
	slog.Info("All recipients processed successfully! 🎉")
}