package utilities

import (
	"fmt"
	"os"
	"strings"
)

func LoadEnv() {
	dat, err := os.ReadFile(".env")

	if err != nil {
		fmt.Println("Unable to load env variables", err)
		return
	}

	for _, line := range strings.Split(string(dat), "\n") {

		if len(line) == 0 || len(strings.Split(line, "=")) != 2 {
			continue
		}

		key, value := strings.Split(line, "=")[0], strings.Split(line, "=")[1]
		os.Setenv(key, value)
	}
}
