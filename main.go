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
	// data, err := hex.DecodeString("FC304B00000000000000FFF00506FF8B90196B00350229435545490000000A7FDE00005265C001154E6174696F6E616C5F4261636B4F75745F47727032300000F0085053394B546524A78EB611D2") //time_signal
	data, err := hex.DecodeString("FC00490000000000000000000F0500000002004081E56EC735000000000029021543554549000000027FCF00002932E0080132340000021043554549000000017F9F0801313500007F0D0304") //splice_insert
	// data, err := hex.DecodeString("fc302500000000000000fff01405000000017feffe2d142b00fe0123d3080001010100007f157a49")
	// data, err := hex.DecodeString("fc304700000000000000fff00506fe1909d1f9002f0223435545490000000a7f9f01144e6174696f6e616c5f4261636b4f75745f456e64310000f0085053394b546524dd8c7fef2b10a4")
	check(err)

	obj := &SCTE35_2013.SCTE35{}
	_, err = obj.DecodeFromRawBytes(data)
	check(err)

	fmt.Println("Schema Version: ", obj.SchemaVersion())
	fmt.Println("Table ID: ", obj.TableID)
	fmt.Println("splice_command_type: ", obj.SpliceCommandType)
	fmt.Println("CRC32 In Hex: ", obj.CRC32InHex)
	fmt.Println("Entire SCTE35 Structure: \n", obj.JSON("	"))

	fmt.Println("===============================================================================")

	jsonStr := obj.JSON()
	obj2 := &SCTE35_2013.SCTE35{}
	err = obj2.DecodeFromJSON(jsonStr)
	check(err)
	fmt.Println("Schema Version: ", obj2.SchemaVersion())
	fmt.Println("Table ID: ", obj2.TableID)
	fmt.Println("splice_command_type: ", obj2.SpliceCommandType)
	fmt.Println("CRC32 In Hex: ", obj2.CRC32InHex)
	fmt.Println("Entire SCTE35 Structure: \n", obj2.JSON("	"))

	_ = &SCTE35_2017.SCTE35{}

	// parsers := [...]common.Parser{&SCTE35_2017.SCTE35{}, &SCTE35_2013.SCTE35{}}
	// for i := 0; i < len(parsers); i++ {
	// 	parser := parsers[i]

	// 	parsedBits, err := parser.DecodeFromRawBytes(data)
	// 	fmt.Println("SchemaVersion: ", parser.SchemaVersion())
	// 	check(err)
	// 	fmt.Println("parsedBits: ", parsedBits)
	// 	if err == nil {
	// 		fmt.Println(parser.JSON("	"))
	// 	}
	// 	fmt.Println("===============================")
	// }
}
