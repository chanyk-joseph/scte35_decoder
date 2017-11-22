package common

import (
	bits "github.com/chanyk-joseph/gobits"
)

type AvailDescriptor struct {
	ProviderAvailID uint32 `json:"provider_avail_id"`
}

type DTMFDescriptor struct {
	Preroll   byte   `json:"preroll"`
	DTMFCount uint8  `json:"dtmf_count"`
	DTMFChars string `json:"dtmf_chars"`
}

type SegmentationDescriptor struct {
	SegmentationEventID              uint32 `json:"segmentation_event_id"`
	SegmentationEventCancelIndicator bool   `json:"segmentation_event_cancel_indicator"`

	ProgramSegmentationFlag   *bool `json:"program_segmentation_flag,omitempty"`
	SegmentationDurationFlag  *bool `json:"segmentation_duration_flag,omitempty"`
	DeliveryNotRestrictedFlag *bool `json:"delivery_not_restricted_flag,omitempty"`

	WebDeliveryAllowedFlag *bool  `json:"web_delivery_allowed_flag,omitempty"`
	NoRegionalBlackoutFlag *bool  `json:"no_regional_blackout_flag,omitempty"`
	ArchiveAllowedFlag     *bool  `json:"archive_allowed_flag,omitempty"`
	DeviceRestrictions     *uint8 `json:"device_restrictions,omitempty"` //2 bits

	ComponentCount         *uint8                   `json:"component_count,omitempty"`
	SegmentationComponents *[]SegmentationComponent `json:"segmentation_components,omitempty"` //2 bits

	SegmentationDuration   *uint64 `json:"segmentation_duration,omitempty"` //40 bits
	SegmentationUpidType   *byte   `json:"segmentation_upid_type,omitempty"`
	SegmentationUpidLength *uint8  `json:"segmentation_upid_length,omitempty"`
	SegmentationUpidInHex  *string `json:"segmentation_upid_in_hex,omitempty"`
	SegmentationTypeID     *uint8  `json:"segmentation_type_id,omitempty"`
	SegmentNum             *uint8  `json:"segment_num,omitempty"`
	SegmentsExpected       *uint8  `json:"segments_expected,omitempty"`
}

type SegmentationComponent struct {
	ComponentTag byte   `json:"component_tag,omitempty"`
	PTSOffset    uint64 `json:"pts_offset,omitempty"` //33 bits
}

type TimeDescriptor struct {
	TAI_seconds uint64 `json:"tai_seconds"`
	TAI_ns      uint32 `json:"tai_ns"`
	UTC_offset  uint16 `json:"utc_offset"`
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

			var components []SegmentationComponent
			for i := 0; i < int(*segDesc.ComponentCount); i++ {
				segComp := SegmentationComponent{}

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
	}

	return numOfParsedBits, err
}

func (availDesc *AvailDescriptor) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	availDesc.ProviderAvailID, _, err = bits.Uint32(input, 1)
	return 32, err
}

func (dtmfDesc *DTMFDescriptor) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	dtmfDesc.Preroll, _, err = bits.Byte(input, numOfParsedBits)
	numOfParsedBits += 8

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 3)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 5)
	dtmfDesc.DTMFCount, _, err = bits.Uint8(tmpBytes, 0)
	numOfParsedBits += 3

	numOfParsedBits += 5 // reserved 5 bits

	dtmfDesc.DTMFChars, _, err = bits.String(input, numOfParsedBits, int(dtmfDesc.DTMFCount)*8)
	numOfParsedBits += int(dtmfDesc.DTMFCount) * 8

	return numOfParsedBits, err
}

func (timeDesc *TimeDescriptor) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 48)
	tmpBytes = append([]byte{0x00, 0x00}, tmpBytes...)
	timeDesc.TAI_seconds, _, err = bits.Uint64(tmpBytes, 0)
	numOfParsedBits += 48

	timeDesc.TAI_ns, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	timeDesc.UTC_offset, _, err = bits.Uint16(input, numOfParsedBits)
	numOfParsedBits += 16

	return numOfParsedBits, err
}
