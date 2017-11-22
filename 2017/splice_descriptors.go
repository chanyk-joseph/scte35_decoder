package schema_2017

import (
	bits "github.com/chanyk-joseph/gobits"
	common "github.com/chanyk-joseph/scte35_decoder/common"
)

type SegmentationDescriptor struct {
	common.SegmentationDescriptor

	SubSegmentNum       *uint8 `json:"sub_segment_num,omitempty"`
	SubSegmentsExpected *uint8 `json:"sub_segments_expected,omitempty"`
}

func (segDesc *SegmentationDescriptor) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	segDesc.SegmentationEventID, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	segDesc.SegmentationEventCancelIndicator, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	numOfParsedBits += 7 //reserved 7 bits

	if !segDesc.SegmentationEventCancelIndicator {
		_, segDesc.ProgramSegmentationFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, segDesc.SegmentationDurationFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, segDesc.DeliveryNotRestrictedFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		if !*segDesc.DeliveryNotRestrictedFlag {
			_, segDesc.WebDeliveryAllowedFlag, err = bits.Bool(input, numOfParsedBits)
			numOfParsedBits++

			_, segDesc.NoRegionalBlackoutFlag, err = bits.Bool(input, numOfParsedBits)
			numOfParsedBits++

			_, segDesc.ArchiveAllowedFlag, err = bits.Bool(input, numOfParsedBits)
			numOfParsedBits++

			tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 2)
			tmpBytes, _ = bits.ShiftRight(tmpBytes, 6)
			_, segDesc.DeviceRestrictions, err = bits.Uint8(tmpBytes, 0)
			numOfParsedBits += 2
		} else {
			numOfParsedBits += 5
		}

		if !*segDesc.ProgramSegmentationFlag {
			_, segDesc.ComponentCount, err = bits.Uint8(input, numOfParsedBits)
			numOfParsedBits += 8

			var components []common.SegmentationComponent
			for i := 0; i < int(*segDesc.ComponentCount); i++ {
				segComp := common.SegmentationComponent{}

				segComp.ComponentTag, _, err = bits.Byte(input, numOfParsedBits)
				numOfParsedBits += 8

				numOfParsedBits += 7 //reserved 7 bits

				tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 33)
				tmpBytes, _ = bits.ShiftRight(tmpBytes, 7)
				tmpBytes = append([]byte{0x00, 0x00, 0x00}, tmpBytes...)
				segComp.PTSOffset, _, err = bits.Uint64(tmpBytes, 0)
				numOfParsedBits += 33

				components = append(components, segComp)
			}
			segDesc.SegmentationComponents = &components
		}

		if *segDesc.SegmentationDurationFlag {
			tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 40)
			tmpBytes = append([]byte{0x00, 0x00, 0x00}, tmpBytes...)
			_, segDesc.SegmentationDuration, err = bits.Uint64(tmpBytes, 0)
			numOfParsedBits += 40
		}

		_, segDesc.SegmentationUpidType, err = bits.Byte(input, numOfParsedBits)
		numOfParsedBits += 8

		_, segDesc.SegmentationUpidLength, err = bits.Uint8(input, numOfParsedBits)
		numOfParsedBits += 8

		if int(*segDesc.SegmentationUpidLength) > 0 {
			_, segDesc.SegmentationUpidInHex, err = bits.HexString(input, numOfParsedBits, int(*segDesc.SegmentationUpidLength)*8)
			numOfParsedBits += (int(*segDesc.SegmentationUpidLength) * 8)
		}

		_, segDesc.SegmentationTypeID, err = bits.Uint8(input, numOfParsedBits)
		numOfParsedBits += 8

		_, segDesc.SegmentNum, err = bits.Uint8(input, numOfParsedBits)
		numOfParsedBits += 8

		_, segDesc.SegmentsExpected, err = bits.Uint8(input, numOfParsedBits)
		numOfParsedBits += 8

		if *segDesc.SegmentationTypeID == 0x34 || *segDesc.SegmentationTypeID == 0x36 {
			_, segDesc.SubSegmentNum, err = bits.Uint8(input, numOfParsedBits)
			numOfParsedBits += 8

			_, segDesc.SubSegmentsExpected, err = bits.Uint8(input, numOfParsedBits)
			numOfParsedBits += 8
		}
	}

	return numOfParsedBits, err
}
