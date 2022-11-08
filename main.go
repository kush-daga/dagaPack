package main

import (
	"github.com/kush-daga/dagaPack/cmd"
)

func main() {
	// asset := cmd.CreateAsset("./example/src/entry.js")
	cmd.CreateGraph("./example/src/entry.js")
	// fmt.Println(asset.Id, asset.Code)
}
