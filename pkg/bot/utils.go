package bot

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func SessionDir() (string, error) {
	dir, ok := os.LookupEnv("SESSION_DIR")
	if ok {
		return filepath.Abs(dir)
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}

	return filepath.Abs(filepath.Join(dir, ".td"))
}

func ReadWarningsToStdErr(err chan error) {
	go func() {
		for {
			fmt.Fprintln(os.Stderr, <-err)
		}
	}()
}

func match(input, regex string) (val string, ok bool) {
	re := regexp.MustCompile(regex)
	match := re.FindString(input)

	return match, match != ""
}
