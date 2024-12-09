package main

import (
	"flag"
	"karinto/trx-downloader/monzo"
	"log"
	"time"
)

func main() {
	var flagDays int
	flag.IntVar(&flagDays, "d", 30, "Download transactions from the past specified number of days. For example, -d 30 downloads the last 30 days' transactions.")

	flag.Parse()

	callMonzo(flagDays)
}

func callMonzo(days int) {
	since := time.Now().AddDate(0, 0, -1*days)

	log.Printf("Refresh token")
	monzo.RefreshToken()

	log.Printf("Download Transactions")
	transactions := monzo.DownloadTransactions(since)

	log.Printf("Write CSV file")
	path, err := monzo.EncodeTransactionsCsv(transactions)
	if err != nil {
		log.Fatalf("Failed to encode transactions to CSV: %v", err)
	}

	log.Printf("Transactions are saved to %s", path)
}
