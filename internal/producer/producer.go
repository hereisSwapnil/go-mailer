package producer

import (
	"encoding/csv"
	"log/slog"
	"os"
	"strings"

	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func LoadRecipientsUsingCsv(filePath string, emailChannel chan types.Recipient) error {
	slog.Info("Loading recipients from CSV file", "filePath", filePath)
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		slog.Error("CSV file does not exist", "filePath", filePath)
		os.Exit(1)
	}
	if err != nil {
		slog.Error("Failed to access CSV file", "error", err, "filePath", filePath)
		os.Exit(1)
	}
	if fileInfo.Size() == 0 {
		slog.Error("CSV file is empty", "filePath", filePath)
		os.Exit(1)
	}

	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records[1:] {
		emailChannel <- types.Recipient{
			Name:  strings.TrimSpace(record[0]),
			Email: strings.TrimSpace(record[1]),
		}
	}
	return nil
}