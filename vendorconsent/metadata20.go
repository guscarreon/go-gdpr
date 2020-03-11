package vendorconsent

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/prebid/go-gdpr/consentconstants"
)

// Parse the metadata from the consent string.
// This returns an error if the input is too short to answer questions about that data.
func parseMetadata20(data []byte) (consentMetadata20, error) {
	if len(data) < 30 {
		return nil, fmt.Errorf("vendor consent strings are at least 30 bytes long. This one was %d", len(data))
	}
	metadata := consentMetadata20(data)
	if metadata.MaxVendorID() < 1 {
		return nil, fmt.Errorf("the consent string encoded a MaxVendorID of %d, but this value must be greater than or equal to 1", metadata.MaxVendorID())
	}
	if metadata.Version() < 1 {
		return nil, fmt.Errorf("the consent string encoded a Version of %d, but this value must be greater than or equal to 1", metadata.Version())
	}
	if metadata.VendorListVersion() == 0 {
		return nil, errors.New("the consent string encoded a VendorListVersion of 0, but this value must be greater than or equal to 1")

	}
	return consentMetadata20(data), nil
}

// consemtMetadata20 implements the parts of the VendorConsents interface which are common
// to BitFields and RangeSections. This relies on Parse to have done some validation already,
// to make sure that functions on it don't overflow the bounds of the byte array.
type consentMetadata20 []byte

func (c consentMetadata20) Version() uint8 {
	// Stored in bits 0-5
	return uint8(c[0] >> 2)
}

func (c consentMetadata20) Created() time.Time {
	_ = c[5]
	// Stored in bits 6-41.. which is [000000xx xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx xx000000] starting at the 1st byte
	deciseconds := int64(binary.BigEndian.Uint64([]byte{
		0x0,
		0x0,
		0x0,
		(c[0]&0x3)<<2 | c[1]>>6,
		c[1]<<2 | c[2]>>6,
		c[2]<<2 | c[3]>>6,
		c[3]<<2 | c[4]>>6,
		c[4]<<2 | c[5]>>6,
	}))
	return time.Unix(deciseconds/decisPerOne, (deciseconds%decisPerOne)*nanosPerDeci)
}

func (c consentMetadata20) LastUpdated() time.Time {
	// Stored in bits 42-77... which is [00xxxxxx xxxxxxxx xxxxxxxx xxxxxxxx xxxxxx00 ] starting at the 6th byte
	deciseconds := int64(binary.BigEndian.Uint64([]byte{
		0x0,
		0x0,
		0x0,
		(c[5] >> 2) & 0x0f,
		c[5]<<6 | c[6]>>2,
		c[6]<<6 | c[7]>>2,
		c[7]<<6 | c[8]>>2,
		c[8]<<6 | c[9]>>2,
	}))
	return time.Unix(deciseconds/decisPerOne, (deciseconds%decisPerOne)*nanosPerDeci)
}

func (c consentMetadata20) CmpID() uint16 {
	// Stored in bits 78-89... which is [000000xx xxxxxxxx xx000000] starting at the 10th byte
	leftByte := ((c[9] & 0x03) << 2) | c[10]>>6
	rightByte := (c[10] << 2) | c[11]>>6
	return binary.BigEndian.Uint16([]byte{leftByte, rightByte})
}

func (c consentMetadata20) CmpVersion() uint16 {
	// Stored in bits 90-101.. which is [00xxxxxx xxxxxx00] starting at the 12th byte
	leftByte := (c[11] >> 2) & 0x0f
	rightByte := (c[11] << 6) | c[12]>>2
	return binary.BigEndian.Uint16([]byte{leftByte, rightByte})
}

func (c consentMetadata20) ConsentScreen() uint8 {
	// Stored in bits 102-107.. which is [000000xx xxxx0000] starting at the 13th byte
	return uint8(((c[12] & 0x03) << 4) | c[13]>>4)
}

func (c consentMetadata20) ConsentLanguage() string {
	// Stored in bits 108-119... which is [0000xxxx xxxxxxxx] starting at the 14th byte.
	// Each letter is stored as 6 bits, with A=0 and Z=25
	leftChar := ((c[13] & 0x0f) << 2) | c[14]>>6
	rightChar := c[14] & 0x3f
	return string([]byte{leftChar + 65, rightChar + 65}) // Unicode A-Z is 65-90
}

func (c consentMetadata20) VendorListVersion() uint16 {
	// The vendor list version is stored in bits 120 - 131
	rightByte := ((c[16] & 0xf0) >> 4) | ((c[15] & 0x0f) << 4)
	leftByte := c[15] >> 4
	return binary.BigEndian.Uint16([]byte{leftByte, rightByte})
}

func (c consentMetadata20) MaxVendorID() uint16 {
	// The max vendor ID is stored in bits 213 - 228 [00000xxx xxxxxxxx xxxxx000]
	leftByte := byte((c[26]&0x07)<<5 + (c[27]&0xf8)>>3)
	rightByte := byte((c[27]&0x07)<<5 + (c[28]&0xf8)>>3)
	return binary.BigEndian.Uint16([]byte{leftByte, rightByte})
}

func (c consentMetadata20) PurposeAllowed(id consentconstants.Purpose) bool {
	// Purposes are stored in bits 152 - 175. The interface contract only defines behavior for ints in the range [1, 24]...
	// so in the valid range, this won't even overflow a uint8.
	return isSet(c, uint(id)+151)
}
