package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	cnt := make(map[string]int)
	
	for _,str := range strings.Fields(s) {
		cnt[str]++
	}
	
	return cnt
}

func main() {
	wc.Test(WordCount)
}
