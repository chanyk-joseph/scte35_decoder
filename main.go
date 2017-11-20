package main

import (
	"encoding/hex"
	"fmt"

	SCTE35_2013 "github.com/chanyk-joseph/scte35_decoder/2013"
	SCTE35_2017 "github.com/chanyk-joseph/scte35_decoder/2017"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("test")

	// data, err := hex.DecodeString("FC304B00000000000000FFF00506FF8B90196B00350229435545490000000A7FDE00005265C001154E6174696F6E616C5F4261636B4F75745F47727032300000F0085053394B546524A78EB611D2") //time_signal
	data, err := hex.DecodeString("FC00490000000000000000000F0500000002004081E56EC735000000000029021543554549000000027FCF00002932E0080132340000021043554549000000017F9F0801313500007F0D0304") //splice_insert
	check(err)

	scte35a := &SCTE35_2017.SCTE35{}
	parsedBits, err := scte35a.DecodeFromRawBytes(data)
	check(err)
	fmt.Println("parsedBits: ", parsedBits)
	fmt.Println(scte35a.JSON("	"))

	fmt.Println("===============================")

	scte35b := &SCTE35_2013.SCTE35{}
	parsedBits, err = scte35b.DecodeFromRawBytes(data)
	check(err)
	fmt.Println("parsedBits: ", parsedBits)
	fmt.Println(scte35b.JSON("	"))
}
