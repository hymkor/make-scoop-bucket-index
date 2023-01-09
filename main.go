package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Manifest struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

var flagBucketDir = flag.String("b","./bucket","Bucket directory path")

func mains() error {
	var dirName = *flagBucketDir

	files, err := os.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		if filepath.Ext(name) != ".json" {
			continue
		}
		jsonPath := filepath.Join(dirName, name)

		jsonBin, err := os.ReadFile(jsonPath)
		if err != nil {
			return err
		}

		var manifest Manifest
		err = json.Unmarshal(jsonBin, &manifest)
		if err != nil {
			return err
		}

		fmt.Printf("* [%s](%s) %s - %s\r\n",
			name[0:len(name)-5],
			manifest.Homepage,
			manifest.Version,
			manifest.Description)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
