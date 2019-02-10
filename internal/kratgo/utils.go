package kratgo

import (
	"fmt"
	"os"
	"path"
)

func getLogOutput(output string) (*os.File, error) {
	if output == "" {
		return nil, fmt.Errorf("Invalid log output: '%s'", output)

	} else if output == "console" {
		return os.Stderr, nil
	}

	dirPath, _ := path.Split(output)
	if err := os.MkdirAll(dirPath, os.ModeDir); err != nil {
		return nil, err
	}
	return os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0755)

}
