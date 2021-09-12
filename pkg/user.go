package pkg

import (
	"fmt"
	"os"
	"strconv"
)

func GetUsers() int {
	f, err := os.Open("/dev/pts")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	c := 0
	for _, v := range files {
		_, err := strconv.Atoi(v.Name())
		if err == nil {
			c++
		}
	}
	return c
}
