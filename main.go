package main

import (
	"fmt"

	"github.com/kush-daga/dagaPack/cmd"
)

func main() {
	asset := cmd.CreateAsset("./example/src/entry.js")
	fmt.Println(asset.Id, asset.Code)
}
