package main

import (
	"karinto/trx-downloader/monzo"
	"log"
	"time"
)

func main() {
    since := time.Now().AddDate(0, -1, 0)

    log.Printf("Refresh token")
    // monzo.RefreshToken()

    log.Printf("Download Transactions")
    transactions := monzo.DownloadTransactions(since)

    log.Printf("Write CSV file")
    path, err := monzo.EncodeTransactionsCsv(transactions)
    if err != nil {
        log.Fatalf("Failed to encode transactions to CSV: %v", err)
    }

    log.Printf("Transactions are saved to %s", path)
}
