package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	SCTE35 "github.com/chanyk-joseph/scte35_decoder/common"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("test")

	data, err := hex.DecodeString("00000001004081E563C7E500000000") //splice insert
	check(err)

	fmt.Println("data len: ", len(data))
	obj := &SCTE35.SpliceInsert{}
	usedBits, err := obj.ParseFromBytes(data)
	check(err)
	fmt.Println("usedBits: ", usedBits)

	buf, err := json.MarshalIndent(obj, "", "	")
	check(err)
	fmt.Println(string(buf))

}
