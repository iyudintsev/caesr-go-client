package main

import (
	"fmt"

	"github.com/iyudintsev/caesr-go-client/internal/client"
)

func main() {
	fmt.Println(client.DoSmth().Transcript)
}
