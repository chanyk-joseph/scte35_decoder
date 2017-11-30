package common

type Parser interface {
	DecodeFromRawBytes([]byte) (int, error)
	DecodeFromJSON(string) error
	JSON(...string) string
	SchemaVersion() string
}

type SCTE35 struct {
	TableID                uint8  `json:"table_id"`
	SectionSyntaxIndicator bool   `json:"section_syntax_indicator"`
	PrivateIndicator       bool   `json:"private_indicator"`
	SectionLength          uint16 `json:"section_length"` // 12 bits
	ProtocolVersion        uint8  `json:"protocol_version"`
	EncryptedPacket        bool   `json:"encrypted_packet"`
	EncryptionAlgorithm    byte   `json:"encryption_algorithm"` // 6 bits
	PTSAdjustment          uint64 `json:"pts_adjustment"`       // 33 bits
	CWIndex                uint8  `json:"cw_index"`
	Tier                   uint16 `json:"tier"`                  // 12 bits
	SpliceCommandLength    uint16 `json:"splice_command_length"` // 12 bits
	SpliceCommandType      byte   `json:"splice_command_type"`

	DescriptorLoopLength uint16             `json:"descriptor_loop_length"`
	SpliceDescriptors    []SpliceDescriptor `json:"splice_descriptors"`

	AlignmentStuffingInHex *string `json:"alignment_stuffing_in_hex,omitempty"`
	ECRC32InHex            *string `json:"e_crc_32_in_hex,omitempty"`
	CRC32InHex             string  `json:"crc_32_in_hex"`
}

type SpliceDescriptor struct {
	SpliceDescriptorTag byte   `json:"splice_descriptor_tag"`
	DescriptorLength    uint8  `json:"descriptor_length"`
	Identifier          uint32 `json:"identifier"`

	PrivateByteInHex *string `json:"private_byte_in_hex,omitempty"`
}
