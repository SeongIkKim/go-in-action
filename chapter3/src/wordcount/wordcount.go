package main

import (
	"fmt"
	"os"

	"github.com/webgenie/go-in-action/chapter3/words"
)

func main() {
	filename := os.Args[1]

	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("파일을 열 때 오류가 발생했습니다.", err)
		return
	}

	text := string(contents)

	count := words.CountWords(text)
	fmt.Printf("총 %d개의 단어를 발견했습니다.\n", count)
}
