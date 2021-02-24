package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Branches map[string]Bank `json:"branches"`
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
    jsonFromFile, err := ioutil.ReadFile("./zengin-data.json")
    if err != nil {
        log.Fatal(err)
    }

	var jsonData map[string]ZenginCode
    err = json.Unmarshal(jsonFromFile, &jsonData)
    if err != nil {
        log.Fatal(err)
    }

    // fmt.Println(jsonData["0001"].Name)
    // fmt.Println(jsonData["0001"].Branches["001"].Name)

    // sort.SliceStable(jsonData,func(i, j int) bool {
    //     return jsonData[i].Hira < jsonData[j].Hira
    // })

    // object => Array
    var bankData []ZenginCode
    for _, bank := range jsonData {
        bankData = append(bankData, bank)
    }
    fmt.Println(len(bankData))
    fmt.Println(bankData[0])

    sort.SliceStable(bankData,func(i, j int) bool {
        return bankData[i].Hira < bankData[j].Hira
    })
    // fmt.Println(bankData[0])
    // fmt.Println(bankData[0])



    // 選択UI用バンクデータ
    hiraganaList := []HiraganaList{
        {Parent: "あ行",ChildList: []string{"あ","い","う","え","お"}},
        {Parent: "か行",ChildList: []string{"か","き","く","け","こ"}},
        {Parent: "さ行",ChildList: []string{"さ","し","す","せ","そ"}},
        {Parent: "た行",ChildList: []string{"た","ち","つ","て","と"}},
        {Parent: "な行",ChildList: []string{"な","に","ぬ","ね","の"}},
        {Parent: "は行",ChildList: []string{"は","ひ","ふ","へ","ほ"}},
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
                Label: hiraganaChild,BankList: []ZenginCode{},
            }

            // 銀行一覧から該当の銀行を取得
            for _, bank := range bankData{
                r := regexp.MustCompile("^" + hiraganaChild)
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
    
    os.Stdout.Write(outputJSON)

    content := []byte(outputJSON)
    ioutil.WriteFile("output.json", content, os.ModePerm)

}