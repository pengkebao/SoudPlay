package main

import (
	"fmt"
	"os"
	"soundPlay/Sound"
	"soundPlay/Tts"
)

func main() {
	file, err := Tts.Conver("测试开始")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = Sound.Play(file, true)
	if err != nil {
		fmt.Println(err)
	}
	os.Remove(file)
}
