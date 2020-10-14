package izapple2

import "fmt"

/*
Apple II Memory Expansion Card


See:
	http://www.apple-iigs.info/doc/fichiers/a2me.pdf
	http://ae.applearchives.com/files/RamFactor_Manual_1.5.pdf
	http://www.1000bit.it/support/manuali/apple/technotes/memx/tn.memx.1.html

There is a self test in ROM, address Cs0A.

From the RamFactor docs:
	The RamFactor card has five addressable registers, which are addressed
according to the slot number the card is in:
		$C080+slot * 16low byte of RAM address
		$C081+slot * 16middle byte of RAM address
		$C082+slot * 16high byte of RAM address
		$C083+slot * 16data at addressed location
		$C08F+slot * 16Firmware Bank Select
	After power up or Control-Reset, the registers on the card are all in a
disabled state. They will be enabled by addressing any address in the firmware
page $Cs00-CsFF.
	The three address bytes can be both written into and read from. If the card
has one Megabyte or less, reading the high address byte will always return a
value in the range $F0-FF. The top nybble can be any value  when you write it,
but it will always be “F” when you read it. If the card has more than one
Megabyte of RAM, the top nybble will be a meaningful part of the address.

*/
const (
	memoryExpansionSize256  = 256 * 1024
	memoryExpansionSize512  = 512 * 1024
	memoryExpansionSize768  = 768 * 1024
	memoryExpansionSize1024 = 1024 * 1024
	memoryExpansionMask     = 0x000fffff // 10 bits, 1MB
)

// CardMemoryExpansion is a Memory Expansion card
type CardMemoryExpansion struct {
	cardBase
	ram   [memoryExpansionSize1024]uint8
	index int
}

// NewCardMemoryExpansion creates a new VidHD card
func NewCardMemoryExpansion() *CardMemoryExpansion {
	var c CardMemoryExpansion
	c.name = "Memory Expansion Card"
	c.loadRomFromResource("<internal>/MemoryExpansionCard-341-0344a.bin")

	return &c
}

// GetInfo returns smartport info
func (c *CardMemoryExpansion) GetInfo() map[string]string {
	info := make(map[string]string)
	info["size"] = fmt.Sprintf("%vKB", len(c.ram)/1024)
	return info
}

func (c *CardMemoryExpansion) assign(a *Apple2, slot int) {

	// Read pointer position
	c.addCardSoftSwitchR(0, func(*ioC0Page) uint8 {
		return uint8(c.index)
	}, "MEMORYEXLOR")
	c.addCardSoftSwitchR(1, func(*ioC0Page) uint8 {
		return uint8(c.index >> 8)
	}, "MEMORYEXMIR")
	c.addCardSoftSwitchR(2, func(*ioC0Page) uint8 {
		// Top nibble returned is 0xf
		return uint8(c.index>>16) | 0xf0
	}, "MEMORYEXHIR")

	// Set pointer position
	c.addCardSoftSwitchW(0, func(_ *ioC0Page, value uint8) {
		c.index = (c.index &^ 0xff) + int(value)
	}, "MEMORYEXLOW")
	c.addCardSoftSwitchW(1, func(_ *ioC0Page, value uint8) {
		c.index = (c.index &^ 0xff00) + int(value)<<8
	}, "MEMORYEXMIW")
	c.addCardSoftSwitchW(2, func(_ *ioC0Page, value uint8) {
		// Only lo nibble is used
		c.index = (c.index &^ 0xff0000) + int(value&0x0f)<<16
	}, "MEMORYEXHIW")

	// Read data
	c.addCardSoftSwitchR(3, func(*ioC0Page) uint8 {
		var value uint8
		if c.index < len(c.ram) {
			value = c.ram[c.index]
		} else {
			value = 0xde // Ram socket not populated
		}
		c.index = (c.index + 1) & memoryExpansionMask
		return value
	}, "MEMORYEXR")

	// Write data
	c.addCardSoftSwitchW(3, func(_ *ioC0Page, value uint8) {
		if c.index < len(c.ram) {
			c.ram[c.index] = value
		}
		c.index = (c.index + 1) & memoryExpansionMask
	}, "MEMORYEXW")

	// The rest of the softswitches return 255, at least on //e and //c
	for i := uint8(4); i < 16; i++ {
		c.addCardSoftSwitchR(i, func(*ioC0Page) uint8 {
			return 255
		}, "MEMORYEXUNUSEDR")
	}

	c.cardBase.assign(a, slot)
}
