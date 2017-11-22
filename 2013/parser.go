package schema_2013

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"unsafe"

	bits "github.com/chanyk-joseph/gobits"
	common "github.com/chanyk-joseph/scte35_decoder/common"
)

//SCTE35(splice_info_section) is implemented based on SCTE35 2013
//http://www.scte.org/documents/pdf/standards/ANSI_SCTE%2035%202013.pdf
type SCTE35 struct {
	common.SCTE35

	//Available Splice Commands
	SpliceNull           *common.SpliceNull           `json:"splice_null,omitempty"`
	SpliceSchedule       *common.SpliceSchedule       `json:"splice_schedule,omitempty"`
	SpliceInsert         *common.SpliceInsert         `json:"splice_insert,omitempty"`
	TimeSignal           *common.TimeSignal           `json:"time_signal,omitempty"`
	BandwidthReservation *common.BandwidthReservation `json:"bandwidth_reservation,omitempty"`
	PrivateCommand       *common.PrivateCommand       `json:"private_command,omitempty"`

	SpliceDescriptors []SpliceDescriptor `json:"splice_descriptors"`
}

type SpliceDescriptor struct {
	common.SpliceDescriptor

	AvailDescriptor        *common.AvailDescriptor        `json:"avail_descriptor,omitempty"`
	DTMFDescriptor         *common.DTMFDescriptor         `json:"dtmf_descriptor,omitempty"`
	SegmentationDescriptor *common.SegmentationDescriptor `json:"segmentation_descriptor,omitempty"`
}

func (spliceDesc *SpliceDescriptor) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	spliceDesc.SpliceDescriptorTag, _, err = bits.Byte(input, numOfParsedBits)
	numOfParsedBits += 8

	spliceDesc.DescriptorLength, _, err = bits.Uint8(input, numOfParsedBits)
	numOfParsedBits += 8

	spliceDesc.Identifier, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	numOfBitsLeft := int(spliceDesc.DescriptorLength-4) * 8 // -4 for identifier
	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, numOfBitsLeft)
	spliceDescriptorUsedBits := 0
	switch spliceDesc.SpliceDescriptorTag {
	case 0x00:
		availDesc := &common.AvailDescriptor{}
		spliceDescriptorUsedBits, err = availDesc.DecodeFromRawBytes(tmpBytes)

		spliceDesc.AvailDescriptor = availDesc
	case 0x01:
		dtmfDesc := &common.DTMFDescriptor{}
		spliceDescriptorUsedBits, err = dtmfDesc.DecodeFromRawBytes(tmpBytes)

		spliceDesc.DTMFDescriptor = dtmfDesc
	case 0x02:
		segDesc := &common.SegmentationDescriptor{}
		spliceDescriptorUsedBits, err = segDesc.DecodeFromRawBytes(tmpBytes)

		spliceDesc.SegmentationDescriptor = segDesc
	}
	numOfParsedBits += spliceDescriptorUsedBits

	numOfBitsLeftForPrivateBytes := int(spliceDesc.DescriptorLength-4)*8 - spliceDescriptorUsedBits
	if numOfBitsLeftForPrivateBytes < 0 {
		spliceDesc = &SpliceDescriptor{}
		return 0, errors.New("The number of bytes used by splice descriptor is more than descriptor_length")
	}
	if numOfBitsLeftForPrivateBytes > 0 {
		_, spliceDesc.PrivateByteInHex, err = bits.HexString(input, numOfParsedBits, numOfBitsLeftForPrivateBytes)
		numOfParsedBits += numOfBitsLeftForPrivateBytes
	}

	return numOfParsedBits, err
}

func (scte35 *SCTE35) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	scte35.TableID, _, err = bits.Uint8(input, numOfParsedBits)
	numOfParsedBits += 8

	scte35.SectionSyntaxIndicator, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	scte35.PrivateIndicator, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	numOfParsedBits += 2 //reserved 2 bits

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 12)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 4)
	scte35.SectionLength, _, err = bits.Uint16(tmpBytes, 0)
	numOfParsedBits += 12

	scte35.ProtocolVersion, _, err = bits.Uint8(input, numOfParsedBits)
	numOfParsedBits += 8

	scte35.EncryptedPacket, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 6)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 2)
	scte35.EncryptionAlgorithm, _, err = bits.Byte(tmpBytes, 0)
	numOfParsedBits += 6

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 33)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 7)
	tmpBytes = append([]byte{0x00, 0x00, 0x00}, tmpBytes...)
	scte35.PTSAdjustment, _, err = bits.Uint64(tmpBytes, 0)
	numOfParsedBits += 33

	scte35.CWIndex, _, err = bits.Uint8(input, numOfParsedBits)
	numOfParsedBits += 8

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 12)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 4)
	scte35.Tier, _, err = bits.Uint16(tmpBytes, 0)
	numOfParsedBits += 12

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 12)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 4)
	scte35.SpliceCommandLength, _, err = bits.Uint16(tmpBytes, 0)
	numOfParsedBits += 12

	scte35.SpliceCommandType, _, err = bits.Byte(input, numOfParsedBits)
	numOfParsedBits += 8

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, int(scte35.SpliceCommandLength*8))
	numOfCommandBits := 0
	switch scte35.SpliceCommandType {
	case 0x00:
		spliceNull := &common.SpliceNull{}
		scte35.SpliceNull = spliceNull
	case 0x04:
		spliceSchedule := &common.SpliceSchedule{}
		numOfCommandBits, err = spliceSchedule.DecodeFromRawBytes(tmpBytes)

		scte35.SpliceSchedule = spliceSchedule
	case 0x05:
		spliceInsert := &common.SpliceInsert{}
		numOfCommandBits, err = spliceInsert.DecodeFromRawBytes(tmpBytes)

		scte35.SpliceInsert = spliceInsert
	case 0x06:
		timeSignal := &common.TimeSignal{}
		numOfCommandBits, err = timeSignal.DecodeFromRawBytes(tmpBytes)

		scte35.TimeSignal = timeSignal
	case 0x07:
		bandwidthReservation := &common.BandwidthReservation{}
		scte35.BandwidthReservation = bandwidthReservation
	case 0xff:
		privateCommand := &common.PrivateCommand{}
		numOfCommandBits, err = privateCommand.DecodeFromRawBytes(tmpBytes)

		scte35.PrivateCommand = privateCommand
	default:
		return 0, errors.New("Unsupported Splice Command Type: " + strconv.Itoa(int(scte35.SpliceCommandType)))
	}
	if err != nil {
		return 0, errors.New("Unable To Parse Splice Command: " + hex.EncodeToString(tmpBytes) + "\n" + err.Error())
	}
	if int(scte35.SpliceCommandLength*8) != numOfCommandBits {
		return 0, errors.New("The number of bits(" + string(numOfCommandBits) + ") used by the splice command is not equal to the expected value: " + string(scte35.SpliceCommandLength*8))
	}
	numOfParsedBits += numOfCommandBits

	scte35.DescriptorLoopLength, _, err = bits.Uint16(input, numOfParsedBits)
	numOfParsedBits += 16

	numOfBitsForDescriptors := int(scte35.DescriptorLoopLength) * 8
	endBitPos := numOfParsedBits + numOfBitsForDescriptors
	for numOfParsedBits < endBitPos {
		spliceDescriptor := &SpliceDescriptor{}
		tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)

		descUsedBits := 0
		descUsedBits, err = spliceDescriptor.DecodeFromRawBytes(tmpBytes)
		if err != nil {
			return 0, err
		}
		numOfParsedBits += descUsedBits

		scte35.SpliceDescriptors = append(scte35.SpliceDescriptors, *spliceDescriptor)
	}

	inputBitLen := bits.Len(input)
	bitRequiredForCRC32 := 32
	if scte35.EncryptedPacket {
		bitRequiredForCRC32 += 32
	}
	if inputBitLen < numOfParsedBits+bitRequiredForCRC32 {
		scte35 = &SCTE35{}
		return 0, errors.New("Parse Error: Not Enough Bits For CRC32 Field, Input Bytes(Hex): " + hex.EncodeToString(input))
	}
	if (inputBitLen-numOfParsedBits-bitRequiredForCRC32)%8 != 0 {
		scte35 = &SCTE35{}
		return 0, errors.New("Parse Error: The number of bits left for alignment_stuffing is not divisible by 8: " + hex.EncodeToString(input))
	}

	if inputBitLen-numOfParsedBits-bitRequiredForCRC32 > 0 {
		_, scte35.AlignmentStuffingInHex, err = bits.HexString(input, numOfParsedBits, inputBitLen-numOfParsedBits-bitRequiredForCRC32)
		numOfParsedBits += (inputBitLen - numOfParsedBits - bitRequiredForCRC32)
	}

	if scte35.EncryptedPacket {
		_, scte35.ECRC32InHex, err = bits.HexString(input, numOfParsedBits, 32)
		numOfParsedBits += 32
	}

	scte35.CRC32InHex, _, err = bits.HexString(input, numOfParsedBits, 32)
	numOfParsedBits += 32

	if numOfParsedBits != bits.Len(input) {
		scte35 = &SCTE35{}
		return 0, errors.New("Parse Error: The number of used bits for constructing the SCTE35 is less than the input")
	}

	return numOfParsedBits, nil
}

func (scte35 *SCTE35) UnmarshalJSON(bytes []byte) (err error) {
	type Alias SCTE35
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(scte35),
	}
	err = json.Unmarshal(bytes, &aux)

	scte35 = (*SCTE35)(unsafe.Pointer(aux))
	return err
}

func (scte35 *SCTE35) DecodeFromJSON(jsonStr string) (err error) {
	err = scte35.UnmarshalJSON([]byte(jsonStr))
	return err
}

func (scte35 *SCTE35) JSON(indent ...string) (result string) {
	var buf []byte
	var err error

	if len(indent) == 0 {
		buf, err = json.Marshal(scte35)
	} else {
		buf, err = json.MarshalIndent(scte35, "", indent[0])
	}

	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (scte35 *SCTE35) SchemaVersion() string {
	return "v2013"
}
