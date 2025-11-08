package producer

import (
	"encoding/csv"
	"os"

	"github.com/hereisSwapnil/go-mailer/internal/types"
)

func LoadRecipientsUsingCsv(filePath string, emailChannel chan types.Recipient) (error) {
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
			Name: record[0],
			Email: record[1],
		}
	}
	return nil
}