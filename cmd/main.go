package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	babel "github.com/jvatic/goja-babel"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

var ID = 0

type CodeAsset struct {
	Id           int
	Dependencies []string
	FileName     string
	Code         string
}

func createAsset(fileName string) (asset CodeAsset) {
	content, err := os.ReadFile(fileName)

	if err != nil {
		panic(err.Error())
	}

	ast, err := js.Parse(parse.NewInputString(string(content)), js.Options{WhileToFor: true})

	if err != nil {
		panic(err)
	}

	dependencies := []string{}
	for _, ele := range ast.List {
		switch ele.(type) {
		case *js.ImportStmt:
			relativePath := (strings.Split(strings.Split(ele.JS(), "from ")[1], "\"")[1])
			dependencies = append(dependencies, relativePath)
		}
	}

	id := ID
	ID += 1

	code, _ := babel.Transform(strings.NewReader(ast.JS()), map[string]interface{}{
		"presets": []string{
			"env",
		},
	})

	buf := new(strings.Builder)
	_, err = io.Copy(buf, code)

	if err != nil {
		panic(err)
	}

	return CodeAsset{
		Id:           id,
		Dependencies: dependencies,
		FileName:     fileName,
		Code:         buf.String(),
	}
}

func main() {
	asset := createAsset("../example/src/entry.js")
	fmt.Println(asset.Id, asset.Code)
}
