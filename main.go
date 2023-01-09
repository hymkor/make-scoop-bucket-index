package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

type Manifest struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

var (
	flagBucketDir   = flag.String("b", "./bucket", "Bucket directory path")
	flagConcatinate = flag.Bool("c", false, "Concatinate the output with the contents from STDIN")

	flagShowNotMatchingUser = flag.String("nmu", "", "Show not matchin user")

	rxGitHubUrl = regexp.MustCompile(`^https://github.com/([\w-]+)/([\w-]+)`)
)

func mains() error {
	var dirName = *flagBucketDir

	if *flagConcatinate {
		io.Copy(os.Stdout, os.Stdin)
	}

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

		title := name[0 : len(name)-5]
		if m := rxGitHubUrl.FindStringSubmatch(manifest.Homepage); m != nil {
			if *flagShowNotMatchingUser != "" && *flagShowNotMatchingUser != m[1] {
				title = m[1] + " / " + m[2]
			}
		}

		fmt.Printf("* [%s](%s) %s - %s\r\n",
			title,
			manifest.Homepage,
			manifest.Version,
			manifest.Description)
	}
	return nil
}

var version string

func main() {
	flag.Parse()

	fmt.Fprintf(os.Stderr, "%s %s for %s/%s by %s\n",
		os.Args[0], version, runtime.GOOS, runtime.GOARCH, runtime.Version())

	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
