package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
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
	body, err := ioutil.ReadAll(resp.Body)

	// bodyから、module.export=と、末尾カンマを外す
	r := regexp.MustCompile(`^module.exports\s?=\s?|;\n?$`)
	replaceBodyByte := []byte(r.ReplaceAllString(string(body), ""))

	var jsonData map[string]ZenginCode
	err = json.Unmarshal(replaceBodyByte, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

		// メインバンクのコード一覧
	mainBankCodeList := []string{
		"0001", // みずほ
		"0005", // 三菱UFJ
		"0009", // 三井住友
		"0010", // りそな
		"0017", // 埼玉りそな
		"0033", // ジャパンネット
		"0036", // 楽天
		"9900", // ゆうちょ
	};

	selectBankData := createBankSelectUiData(jsonData,mainBankCodeList)
	mainBankData := createMainBankData(jsonData,mainBankCodeList)

	// 出力先の取得
	fSelectBank := flag.String("output1","./output.json","選択UI用jsonの出力先")
	fMainBank := flag.String("output2","./output_main.json","メインバンクjsonの出力先")
	flag.Parse()

	// ファイル出力（選択UI用バンクデータ）
	outputSelectBank, err := json.Marshal(selectBankData)
	if err != nil {
		log.Fatal(err)
	}
	contentSelectBank := []byte(outputSelectBank)
	ioutil.WriteFile(*fSelectBank, contentSelectBank, os.ModePerm)

	// ファイル出力（選択UI用バンクデータ）
	outputMainBank, err := json.Marshal(mainBankData)
	if err != nil {
		log.Fatal(err)
	}
	contentMainBank := []byte(outputMainBank)
	ioutil.WriteFile(*fMainBank, contentMainBank, os.ModePerm)
}

func getRuneAt(s string, i int) rune {
    rs := []rune(s)
    return rs[i]
}

// 選択UI用バンクデータの作成
func createBankSelectUiData(jsonData map[string]ZenginCode, mainBankCodeList []string)[]BankSelectData{
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

			child := BankSelectDataChildren{
				Label: string(getRuneAt(hiraganaChild, 0)),
				BankList: []ZenginCode{},
			}

			// 銀行一覧から該当の銀行を取得
			for _, bank := range bankData{
				r := regexp.MustCompile("^[" + hiraganaChild + "]")
				if r.MatchString(bank.Hira){
					// メインバンクの場合はnameに「銀行」を付加
					if contains(bank.Code, mainBankCodeList){
						bank.Name = bank.Name + "銀行"
					}
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

	return selectBankData
}

// メインバンクデータの作成
func createMainBankData(jsonData map[string]ZenginCode, mainBankCodeList []string)[]ZenginCode{
	var mainBankData []ZenginCode
	for _, mainBankCode := range mainBankCodeList {
		appendBank := jsonData[mainBankCode]
		// メインバンクの場合はnameに「銀行」を付加
		if contains(appendBank.Code, mainBankCodeList){
			appendBank.Name = appendBank.Name + "銀行"
		}
		mainBankData = append(mainBankData, appendBank)
	}

	return mainBankData
}

func contains(target string, list []string) bool {
	for _, v := range list {
		if target == v {
			return true
		}
	}
	return false
}