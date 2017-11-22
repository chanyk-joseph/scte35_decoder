package common

import (
	bits "github.com/chanyk-joseph/gobits"
)

//SpliceNull | splice_command_type = 0x00
type SpliceNull struct {
}

//SpliceSchedule | splice_command_type = 0x04
type SpliceSchedule struct {
	SpliceCount    uint8            `json:"splice_count"`
	ScheduleEvents *[]ScheduleEvent `json:"schedule_events,omitempty"`
}

type ScheduleEvent struct {
	SpliceEventID              uint32 `json:"splice_event_id,omitempty"`
	SpliceEventCancelIndicator bool   `json:"splice_event_cancel_indicator,omitempty"`

	OutOfNetworkIndicator *bool `json:"out_of_network_indicator,omitempty"`
	ProgramSpliceFlag     *bool `json:"program_splice_flag,omitempty"`
	DurationFlag          *bool `json:"duration_flag,omitempty"`

	UTCSpliceTime      *uint32              `json:"utc_splice_time,omitempty"`
	ComponentCount     *uint8               `json:"component_count,omitempty"`
	ScheduleComponents *[]ScheduleComponent `json:"schedule_components,omitempty"`

	BreakDuration *BreakDuration `json:"break_duration,omitempty"`

	UniqueProgramID *uint16 `json:"unique_program_id,omitempty"`
	AvailNum        *byte   `json:"avail_num,omitempty"`
	AvailsExpected  *byte   `json:"avails_expected,omitempty"`
}

type ScheduleComponent struct {
	ComponentTag  byte   `json:"component_tag"`
	UTCSpliceTime uint32 `json:"utc_splice_time"`
}

//SpliceInsert | splice_command_type = 0x05
type SpliceInsert struct {
	SpliceEventID              uint32 `json:"splice_event_id"`
	SpliceEventCancelIndicator bool   `json:"splice_event_cancel_indicator"`

	OutOfNetworkIndicator *bool `json:"out_of_network_indicator,omitempty"`
	ProgramSpliceFlag     *bool `json:"program_splice_flag,omitempty"`
	DurationFlag          *bool `json:"duration_flag,omitempty"`
	SpliceImmediateFlag   *bool `json:"splice_immediate_flag,omitempty"`

	SpliceTime *SpliceTime `json:"splice_time,omitempty"`

	ComponentCount   *uint8             `json:"component_count,omitempty"`
	InsertComponents *[]InsertComponent `json:"insert_components,omitempty"`

	BreakDuration *BreakDuration `json:"break_duration,omitempty"`

	UniqueProgramID *uint16 `json:"unique_program_id,omitempty"`
	AvailNum        *byte   `json:"avail_num,omitempty"`
	AvailsExpected  *byte   `json:"avails_expected,omitempty"`
}

//TimeSignal | splice_command_type = 0x06
type TimeSignal struct {
	SpliceTime *SpliceTime `json:"splice_time"`
}

//BandwidthReservation | splice_command_type = 0x07
type BandwidthReservation struct {
}

//PrivateCommand | splice_command_type = 0xff
type PrivateCommand struct {
	Identifier       uint32  `json:"identifier"`
	PrivateByteInHex *string `json:"private_byte_in_hex,omitempty"`
}

//SpliceTime is a command object used in multiple splice command objects
type SpliceTime struct {
	TimeSpecifiedFlag bool    `json:"time_specified_flag"`
	PTSTime           *uint64 `json:"pts_time,omitempty"` //33 bits
}

type BreakDuration struct {
	AutoReturn bool   `json:"auto_return"`
	Duration   uint64 `json:"duration"` //33 bits
}

type InsertComponent struct {
	ComponentTag byte        `json:"component_tag"`
	SpliceTime   *SpliceTime `json:"splice_time,omitempty"`
}

//DecodeFromRawBytes parses input []byte to PrivateCommand object
func (privateCommand *PrivateCommand) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	privateCommand.Identifier, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	bytesLeft := len(input) - 4
	_, privateCommand.PrivateByteInHex, err = bits.HexString(input, numOfParsedBits, bytesLeft*8)
	numOfParsedBits += (bytesLeft * 8)

	return numOfParsedBits, err
}

//DecodeFromRawBytes parses input []byte to SpliceInsert object
func (spliceInsert *SpliceInsert) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte
	var tmpUsedBits int

	spliceInsert.SpliceEventID, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	spliceInsert.SpliceEventCancelIndicator, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	numOfParsedBits += 7 //reserved 7 bits

	if !spliceInsert.SpliceEventCancelIndicator {
		_, spliceInsert.OutOfNetworkIndicator, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, spliceInsert.ProgramSpliceFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, spliceInsert.DurationFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, spliceInsert.SpliceImmediateFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		numOfParsedBits += 4 //reserved 4 bits

		if *spliceInsert.ProgramSpliceFlag && !*spliceInsert.SpliceImmediateFlag {
			spliceInsert.SpliceTime = &SpliceTime{}
			tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)
			tmpUsedBits, err = spliceInsert.SpliceTime.DecodeFromRawBytes(tmpBytes)
			numOfParsedBits += tmpUsedBits
		}
		if !(*spliceInsert.ProgramSpliceFlag) {
			_, spliceInsert.ComponentCount, err = bits.Uint8(input, numOfParsedBits)
			numOfParsedBits += 8

			var insertComponents []InsertComponent
			for i := 0; i < int(*spliceInsert.ComponentCount); i++ {
				comp := &InsertComponent{}
				tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)
				tmpUsedBits, err = comp.DecodeFromRawBytes(tmpBytes, *spliceInsert.SpliceImmediateFlag)
				numOfParsedBits += tmpUsedBits

				insertComponents = append(insertComponents, *comp)
			}
			spliceInsert.InsertComponents = &insertComponents
		}

		if *spliceInsert.DurationFlag {
			b := &BreakDuration{}
			tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)
			tmpUsedBits, err = b.DecodeFromRawBytes(tmpBytes)
			spliceInsert.BreakDuration = b
			numOfParsedBits += tmpUsedBits
		}

		_, spliceInsert.UniqueProgramID, err = bits.Uint16(input, numOfParsedBits)
		numOfParsedBits += 16

		_, spliceInsert.AvailNum, err = bits.Byte(input, numOfParsedBits)
		numOfParsedBits += 8

		_, spliceInsert.AvailsExpected, err = bits.Byte(input, numOfParsedBits)
		numOfParsedBits += 8
	}

	if err != nil {
		return 0, err
	}
	return numOfParsedBits, err
}

//DecodeFromRawBytes parses input []byte to ScheduleEvent object
func (scheduleEvent *ScheduleEvent) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte
	var tmpUsedBits int

	scheduleEvent.SpliceEventID, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	scheduleEvent.SpliceEventCancelIndicator, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	numOfParsedBits += 7 //reserved 7 bits

	if !scheduleEvent.SpliceEventCancelIndicator {
		_, scheduleEvent.OutOfNetworkIndicator, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, scheduleEvent.ProgramSpliceFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		_, scheduleEvent.DurationFlag, err = bits.Bool(input, numOfParsedBits)
		numOfParsedBits++

		numOfParsedBits += 4 //reserved 5 bits

		if *scheduleEvent.ProgramSpliceFlag {
			_, scheduleEvent.UTCSpliceTime, err = bits.Uint32(input, numOfParsedBits)
			numOfParsedBits += 32
		}
		if !(*scheduleEvent.ProgramSpliceFlag) {
			_, scheduleEvent.ComponentCount, err = bits.Uint8(input, numOfParsedBits)
			numOfParsedBits += 8

			var scheduleComponents []ScheduleComponent
			for i := 0; i < int(*scheduleEvent.ComponentCount); i++ {
				comp := &ScheduleComponent{}
				tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 40)
				tmpUsedBits, err = comp.DecodeFromRawBytes(tmpBytes)
				numOfParsedBits += tmpUsedBits

				scheduleComponents = append(scheduleComponents, *comp)
			}
			scheduleEvent.ScheduleComponents = &scheduleComponents
		}

		if *scheduleEvent.DurationFlag {
			b := &BreakDuration{}
			tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)
			tmpUsedBits, err = b.DecodeFromRawBytes(tmpBytes)
			scheduleEvent.BreakDuration = b
			numOfParsedBits += tmpUsedBits
		}

		_, scheduleEvent.UniqueProgramID, err = bits.Uint16(input, numOfParsedBits)
		numOfParsedBits += 16

		_, scheduleEvent.AvailNum, err = bits.Byte(input, numOfParsedBits)
		numOfParsedBits += 8

		_, scheduleEvent.AvailsExpected, err = bits.Byte(input, numOfParsedBits)
		numOfParsedBits += 8
	}

	return numOfParsedBits, nil
}

//DecodeFromRawBytes parses input []byte to SpliceSchedule object
func (spliceSchedule *SpliceSchedule) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte
	var tmpUsedBits int

	spliceSchedule.SpliceCount, _, err = bits.Uint8(input, numOfParsedBits)
	numOfParsedBits += 8

	tmpEvents := []ScheduleEvent{}
	for i := 0; i < int(spliceSchedule.SpliceCount); i++ {
		scheduleEvent := &ScheduleEvent{}
		tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)
		tmpUsedBits, err = scheduleEvent.DecodeFromRawBytes(tmpBytes)

		tmpEvents = append(tmpEvents, *scheduleEvent)
		numOfParsedBits += tmpUsedBits
	}
	spliceSchedule.ScheduleEvents = &tmpEvents

	if err != nil {
		return 0, err
	}
	return numOfParsedBits, err
}

//DecodeFromRawBytes parses input []byte to InsertComponent object
func (insertComponent *InsertComponent) DecodeFromRawBytes(input []byte, spliceImmediateFlag bool) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	insertComponent.ComponentTag, _, err = bits.Byte(input, numOfParsedBits)
	numOfParsedBits += 8

	if !spliceImmediateFlag {
		insertComponent.SpliceTime = &SpliceTime{}

		bitsUsed := 0
		tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 0)
		bitsUsed, err = insertComponent.SpliceTime.DecodeFromRawBytes(tmpBytes)
		numOfParsedBits += bitsUsed
	}

	if err != nil {
		return 0, err
	}
	return numOfParsedBits, err
}

//DecodeFromRawBytes parses input []byte to ScheduleComponent object
func (scheduleComponent *ScheduleComponent) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	scheduleComponent.ComponentTag, _, err = bits.Byte(input, numOfParsedBits)
	numOfParsedBits += 8

	scheduleComponent.UTCSpliceTime, _, err = bits.Uint32(input, numOfParsedBits)
	numOfParsedBits += 32

	if err != nil {
		return 0, err
	}
	return numOfParsedBits, nil
}

//DecodeFromRawBytes parses input []byte to BreakDuration object
func (breakDuration *BreakDuration) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	breakDuration.AutoReturn, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	numOfParsedBits += 6 //reserved 6 bits

	tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 33)
	tmpBytes, _ = bits.ShiftRight(tmpBytes, 7)
	tmpBytes = append([]byte{0x00, 0x00, 0x00}, tmpBytes...)
	breakDuration.Duration, _, err = bits.Uint64(tmpBytes, 0)
	numOfParsedBits += 33

	if err != nil {
		return 0, err
	}
	return numOfParsedBits, err
}

//DecodeFromRawBytes parses input []byte to SpliceTime object
func (spliceTime *SpliceTime) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	var tmpBytes []byte

	spliceTime.TimeSpecifiedFlag, _, err = bits.Bool(input, numOfParsedBits)
	numOfParsedBits++

	if spliceTime.TimeSpecifiedFlag {
		numOfParsedBits += 6 //reserved 6 bits

		tmpBytes, _, err = bits.SubBits(input, numOfParsedBits, 33)
		tmpBytes, _ = bits.ShiftRight(tmpBytes, 7)
		tmpBytes = append([]byte{0x00, 0x00, 0x00}, tmpBytes...)
		_, spliceTime.PTSTime, err = bits.Uint64(tmpBytes, 0)
		numOfParsedBits += 33
	} else {
		numOfParsedBits += 7 //reserved 7 bits
	}

	return numOfParsedBits, err
}

//DecodeFromRawBytes parses input []byte to TimeSignal object
func (timeSignal *TimeSignal) DecodeFromRawBytes(input []byte) (numOfParsedBits int, err error) {
	timeSignal.SpliceTime = &SpliceTime{}
	numOfParsedBits, err = timeSignal.SpliceTime.DecodeFromRawBytes(input)

	return numOfParsedBits, err
}
