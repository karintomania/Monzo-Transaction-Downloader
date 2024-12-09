package main

import (
	"karinto/trx-downloader/config"
)

func main() {
    config.InitConfig()

    RefreshToken()

}
