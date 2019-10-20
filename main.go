package main

import (
	"flag"
	"fmt"
	"github.com/Achillesxu/bingImgDown/binghomeimage"
	"os"
	"strings"
)

var downImage string

const defaultUrl = "https://bing.com"

var intro = fmt.Sprintf("%s now only support download %s homepage image\n", os.Args[0], defaultUrl)

func init() {
	const (
		defaultUrl = ""
		usage      = "download image url, for example: " + defaultUrl
	)
	flag.StringVar(&downImage, "down", defaultUrl, usage)
}

func cmdUsage() {
	retSli := strings.Split(os.Args[0], "/")
	fmt.Printf("Usage of %s:\nfor instance:\t./%s -down %s\n",
		os.Args[0], retSli[len(retSli)-1], defaultUrl)
	flag.PrintDefaults()
}

func main() {
	fmt.Println(intro)
	flag.Parse()
	if strings.HasPrefix(downImage, "https") {
		if downImage != defaultUrl {
			fmt.Printf("%s now dont support this site!", os.Args[0])
		} else {
			binghomeimage.DownLoadBingHomeImage(downImage)
		}
	} else {
		cmdUsage()
	}
}
