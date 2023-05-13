package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type LuckContent struct {
	Luck       []string
	LuckInfo   []string
	GoodThings []string
	BadThings  []string
}

var Luck = []string{"上吉", "中吉", "末吉", "上平", "中平", "末平", "大凶", "中凶", "末凶"}
var LuckInfo = []string{
	"受到幸运之神眷顾的一天，工作生活都格外地顺利，一切都将按照愿望进行，意想不到的好事接二连三地发生",
	"不幸悄然远离，顺势而为能达到好的结果，今天会是强运的一天",
	"上升的运势抵消不幸，带来意料之外的好运",
	"运势平中有升，可选的道路比想象中的多，不管如何选择都不会出现太差的结果，放手去做即可",
	"普普通通的一天，平静也值得享受",
	"平淡之中偶有坎坷，平常心对待",
	"阴云密布的一天，总是发生失误，受到他人失误的影响，总会有这样的时候，向他人求助或者倾诉是个不错的选择",
	"状态不佳的一天，客观因素成为绊脚石，避开对自己不利的事物，发挥自己的优势",
	"不利之中存在转机，抓住机会或许能够能够扭转局面，做力所能及之事"}
var GoodThings = []string{"直行", "饱餐一顿，做喜欢的事情", "与家人朋友交流", "学习新知识", "看一部好的作品", "尝试作出一些改变", "放松自己，减小压力", "做轻松的事情，调整状态", "思考，反省"}
var BadThings = []string{"贪婪", "焦虑", "冷漠", "不思进取", "轻浮", "失控", "暴躁", "愤怒", "伤心"}

func generateJson() {
	luckm := LuckContent{
		Luck:       Luck,
		LuckInfo:   LuckInfo,
		GoodThings: GoodThings,
		BadThings:  BadThings,
	}
	bm, err := json.MarshalIndent(luckm, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("marshall:")
		fmt.Println(string(bm))
	}
	err = ioutil.WriteFile("LuckContent.json", bm, 0600)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func readJson() {
	bum, err := ioutil.ReadFile("LuckContent.json")
	//err = json.Unmarshal(bum, luckum)
	var buffer bytes.Buffer
	err = json.Indent(&buffer, bum, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("unmarshall:")
	fmt.Println(buffer.String())
}
func main() {
	//fmt.Println(Luck)
	//fmt.Println(LuckInfo)
	//fmt.Println(GoodThings)
	//fmt.Println(BadThings)

	//luckum := new(LuckContent)
	readJson()
}
