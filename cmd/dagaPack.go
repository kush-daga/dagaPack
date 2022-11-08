package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func CreateAsset(fileName string) (asset CodeAsset) {
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

	babel.Init(2)
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

type CodeAssetWithMapping struct {
	asset   CodeAsset
	mapping map[string]int
}

type DepGraph []CodeAssetWithMapping

func CreateGraph(entryFile string) (depGraph DepGraph) {
	mainAsset := CreateAsset(entryFile)
	mainAssetMapping := map[string]int{}

	mainAssetWithMapping := CodeAssetWithMapping{
		mainAsset,
		mainAssetMapping,
	}
	depGraph = []CodeAssetWithMapping{mainAssetWithMapping}

	for i := 0; i < len(depGraph); i++ {
		assetInQueue := depGraph[i]
		dirName := filepath.Dir(assetInQueue.asset.FileName)

		for _, relPath := range assetInQueue.asset.Dependencies {

			absPath, err := filepath.Abs(filepath.Join(dirName, relPath))
			fmt.Println("ABS PATH:", absPath)

			if err != nil {
				panic(err)
			}

			childAsset := CreateAsset(absPath)
			assetInQueue.mapping[relPath] = childAsset.Id

			var tempMapping = map[string]int{}
			childAssetWithMapping := CodeAssetWithMapping{childAsset, tempMapping}
			depGraph = append(depGraph, childAssetWithMapping)
		}
	}

	// j3, _ := json.Marshal(queue[2].asset)
	fmt.Println(depGraph[0])
	fmt.Println(depGraph[1])
	fmt.Println(depGraph[2])

	// // fmt.Println(asset, mapping)
	// depGraph = DepGraph{
	// 	{asset, mapping},
	// }

	// var m map[string]any
	// ja, _ := json.Marshal(asset)
	// json.Unmarshal(ja, &m)
	// m["mapping"] = mapping

	// js, _ := json.Marshal(m)

	// fmt.Println(string(js))

	return depGraph
}

func convertGraphElToMap(a CodeAsset, m map[string]int) (res map[string]any, resString string) {
	res = map[string]any{}
	ja, _ := json.Marshal(a)

	json.Unmarshal(ja, &res)
	res["mapping"] = m

	js, _ := json.Marshal(res)

	return res, string(js)
}
