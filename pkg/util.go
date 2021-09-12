package pkg

import (
	"fmt"
	"os"
)

func getCgoupValueByPath(path string) int64 {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	var value int64
	n, err := fmt.Sscanf(string(data), "%d", &value)
	if err != nil || n != 1 {
		return 0
	}
	return value
}