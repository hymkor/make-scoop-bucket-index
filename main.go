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

	"github.com/hymkor/go-sortedkeys"
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

	others := map[string][][4]string{}
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
				others[m[1]] = append(others[m[1]], [...]string{
					title, manifest.Homepage, manifest.Version, manifest.Description,
				})
				continue
			}
		}

		fmt.Printf("* [%s](%s) %s - %s\r\n",
			title,
			manifest.Homepage,
			manifest.Version,
			manifest.Description)
	}
	for p := sortedkeys.New(others); p.Range(); {
		fmt.Printf("\r\n%s\r\n", p.Key)
		for _, repo := range p.Value {
			fmt.Printf("* [%s](%s) %s - %s\r\n", repo[0], repo[1], repo[2], repo[3])
		}
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
