package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"

	"golang.org/x/exp/utf8string"
)

type Bank struct{
	Code string `json:"code"`
	Name string `json:"name"`
	Kana string `json:"kana"`
	Hira string `json:"hira"`
	Roma string `json:"roma"`
}

type ZenginCode struct{
	Bank
	// Branches map[string]Bank `json:"branches"`
}

type HiraganaList struct{
	Parent string
	ChildList []string
}

type BankSelectDataChildren struct{
	Label string
	BankList []ZenginCode
}

type BankSelectData struct{
	Label string
	Child []BankSelectDataChildren
}

func main() {
	url := "https://raw.githubusercontent.com/zengin-code/zengin-js/master/lib/zengin-data.js"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	// bodyから、module.export=と、末尾カンマを外す
	r := regexp.MustCompile(`^module.exports\s?=\s?|;\n?$`)
	replaceBodyByte := []byte(r.ReplaceAllString(string(body), ""))

	var jsonData map[string]ZenginCode
	err = json.Unmarshal(replaceBodyByte, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// object => Array
	var bankData []ZenginCode
	for _, bank := range jsonData {
		bankData = append(bankData, bank)
	}

	sort.SliceStable(bankData,func(i, j int) bool {
		return bankData[i].Hira < bankData[j].Hira
	})

	// 選択UI用バンクデータ
	hiraganaList := []HiraganaList{
		{Parent: "あ行",ChildList: []string{"あ","い","う","え","お"}},
		{Parent: "か行",ChildList: []string{"かが","きぎ","くぐ","けげ","こご"}},
		{Parent: "さ行",ChildList: []string{"さざ","しじ","すず","せぜ","そぞ"}},
		{Parent: "た行",ChildList: []string{"ただ","ちぢ","つづ","てで","とど"}},
		{Parent: "な行",ChildList: []string{"な","に","ぬ","ね","の"}},
		{Parent: "は行",ChildList: []string{"はばぱ","ひびぴ","ふぶぷ","へべぺ","ほぼぽ",}},
		{Parent: "ま行",ChildList: []string{"ま","み","む","め","も"}},
		{Parent: "や行",ChildList: []string{"や","ゆ","よ"}},
		{Parent: "ら行",ChildList: []string{"ら","り","る","れ","ろ"}},
		{Parent: "わ行",ChildList: []string{"わ","を","ん"}},
	}

	var selectBankData []BankSelectData
	for _, hiragana := range hiraganaList {

		var childData []BankSelectDataChildren
		for _, hiraganaChild := range hiragana.ChildList{

			labelText := utf8string.NewString(hiraganaChild)

			child := BankSelectDataChildren{
				Label: labelText.Slice(0,1),
				BankList: []ZenginCode{},
			}

			// 銀行一覧から該当の銀行を取得
			for _, bank := range bankData{
				r := regexp.MustCompile("^[" + hiraganaChild + "]")
				if r.MatchString(bank.Hira){
					child.BankList = append(child.BankList, bank)
				}
			}

			childData = append(childData, child)
		}

		data := BankSelectData{
			Label: hiragana.Parent,
			Child: childData,
		}

		selectBankData = append(selectBankData, data)
	}

	// ファイル出力
	outputJSON, err := json.Marshal(selectBankData)
	if err != nil {
		log.Fatal(err)
	}

	content := []byte(outputJSON)
	os.WriteFile("output.json", content, os.ModePerm)
}