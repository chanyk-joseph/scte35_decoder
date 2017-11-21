# scte35_decoder
scte35_decoder is a raw bytes parser for SCTE35 signal based on SCTE35 2013/2017 schema<br/>
2013: http://www.scte.org/documents/pdf/standards/ANSI_SCTE%2035%202013.pdf <br/>
2017: http://www.scte.org/SCTEDocs/Standards/PublicReview/SCTE%2035%202017.pdf 

## Usage
```go
package main

import (
	"encoding/hex"
	"fmt"

	SCTE35_2017 "github.com/chanyk-joseph/scte35_decoder/2017"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	data, err := hex.DecodeString("fc304700000000000000fff00506fe1909d1f9002f0223435545490000000a7f9f01144e6174696f6e616c5f4261636b4f75745f456e64310000f0085053394b546524dd8c7fef2b10a4")
	check(err)

	obj := &SCTE35_2017.SCTE35{}
	_, err = obj.DecodeFromRawBytes(data)
	check(err)

	fmt.Println("Schema Version: ", obj.SchemaVersion())
	fmt.Println("Table ID: ", obj.TableID)
	fmt.Println("splice_command_type: ", obj.SpliceCommandType)
	fmt.Println("CRC32 In Hex: ", obj.CRC32InHex)
	fmt.Println("Entire SCTE35 Structure: \n", obj.JSON("	"))

	fmt.Println("==============================================================================")

	jsonStr := obj.JSON()
	obj2 := &SCTE35_2017.SCTE35{}
	err = obj2.DecodeFromJSON(jsonStr)
	check(err)
	fmt.Println("Schema Version: ", obj2.SchemaVersion())
	fmt.Println("Table ID: ", obj2.TableID)
	fmt.Println("splice_command_type: ", obj2.SpliceCommandType)
	fmt.Println("CRC32 In Hex: ", obj2.CRC32InHex)
	fmt.Println("Entire SCTE35 Structure: \n", obj2.JSON("	"))
}
```

Sample Output
```
Schema Version:  v2017
Table ID:  252
splice_command_type:  6
CRC32 In Hex:  ef2b10a4
Entire SCTE35 Structure: 
 {
	"table_id": 252,
	"section_syntax_indicator": false,
	"private_indicator": false,
	"section_length": 71,
	"protocol_version": 0,
	"encrypted_packet": false,
	"encryption_algorithm": 0,
	"pts_adjustment": 0,
	"cw_index": 0,
	"tier": 4095,
	"splice_command_length": 5,
	"splice_command_type": 6,
	"descriptor_loop_length": 47,
	"alignment_stuffing_in_hex": "8c7f",
	"crc_32_in_hex": "ef2b10a4",
	"time_signal": {
		"splice_time": {
			"time_specified_flag": true,
			"pts_time": 420073977
		}
	},
	"splice_descriptors": [
		{
			"splice_descriptor_tag": 2,
			"descriptor_length": 35,
			"identifier": 1129661769,
			"segmentation_descriptor": {
				"segmentation_event_id": 10,
				"segmentation_event_cancel_indicator": false,
				"program_segmentation_flag": true,
				"segmentation_duration_flag": false,
				"delivery_not_restricted_flag": false,
				"web_delivery_allowed_flag": true,
				"no_regional_blackout_flag": true,
				"archive_allowed_flag": true,
				"device_restrictions": 3,
				"segmentation_upid_type": 1,
				"segmentation_upid_length": 20,
				"segmentation_upid_in_hex": "4e6174696f6e616c5f4261636b4f75745f456e64",
				"segmentation_type_id": 49,
				"segment_num": 0,
				"segments_expected": 0
			}
		},
		{
			"splice_descriptor_tag": 240,
			"descriptor_length": 8,
			"identifier": 1347631435,
			"private_byte_in_hex": "546524dd"
		}
	]
}
==============================================================================
Schema Version:  v2017
Table ID:  252
splice_command_type:  6
CRC32 In Hex:  ef2b10a4
Entire SCTE35 Structure: 
 {
	"table_id": 252,
	"section_syntax_indicator": false,
	"private_indicator": false,
	"section_length": 71,
	"protocol_version": 0,
	"encrypted_packet": false,
	"encryption_algorithm": 0,
	"pts_adjustment": 0,
	"cw_index": 0,
	"tier": 4095,
	"splice_command_length": 5,
	"splice_command_type": 6,
	"descriptor_loop_length": 47,
	"alignment_stuffing_in_hex": "8c7f",
	"crc_32_in_hex": "ef2b10a4",
	"time_signal": {
		"splice_time": {
			"time_specified_flag": true,
			"pts_time": 420073977
		}
	},
	"splice_descriptors": [
		{
			"splice_descriptor_tag": 2,
			"descriptor_length": 35,
			"identifier": 1129661769,
			"segmentation_descriptor": {
				"segmentation_event_id": 10,
				"segmentation_event_cancel_indicator": false,
				"program_segmentation_flag": true,
				"segmentation_duration_flag": false,
				"delivery_not_restricted_flag": false,
				"web_delivery_allowed_flag": true,
				"no_regional_blackout_flag": true,
				"archive_allowed_flag": true,
				"device_restrictions": 3,
				"segmentation_upid_type": 1,
				"segmentation_upid_length": 20,
				"segmentation_upid_in_hex": "4e6174696f6e616c5f4261636b4f75745f456e64",
				"segmentation_type_id": 49,
				"segment_num": 0,
				"segments_expected": 0
			}
		},
		{
			"splice_descriptor_tag": 240,
			"descriptor_length": 8,
			"identifier": 1347631435,
			"private_byte_in_hex": "546524dd"
		}
	]
}
```