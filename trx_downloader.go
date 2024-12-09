package main

import (
	"fmt"
	"karinto/trx-downloader/cache"
)

func main() {
	fmt.Println("Hello, World!")
    cache.WriteCache("test", "testvalue")
}
