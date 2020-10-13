package izapple2

/*
ThunderClock`, real time clock card.

See:
	https://ia800706.us.archive.org/22/items/ThunderClock_Plus/ThunderClock_Plus.pdf
	https://prodos8.com/docs/technote/01/
	https://www.semiee.com/file/backup/NEC-D1990.pdf


uPD1990AC hookup:
	bit 0 = data in
	bit 1 = CLK
	bit 2 = STB
	bit 3 = C0
	bit 4 = C1
	bit 5 = C2
	bit 7 = data out
*/

// CardThunderClockPlus represents a ThunderClock+ card
type CardThunderClockPlus struct {
	cardBase
	microPD1990ac
}

// NewCardThunderClockPlus creates a new CardThunderClockPlus
func NewCardThunderClockPlus() *CardThunderClockPlus {
	var c CardThunderClockPlus
	c.name = "ThunderClock+ Card"

	data, err := loadResource("<internal>/ThunderclockPlusROM.bin")
	if err != nil {
		panic(err)
	}
	c.loadRom(data)

	return &c
}

func (c *CardThunderClockPlus) assign(a *Apple2, slot int) {
	c.addCardSoftSwitchR(0, func(*ioC0Page) uint8 {
		bit := c.microPD1990ac.out()
		// Get the next data bit from uPD1990AC on the MSB
		if bit {
			return 0x80
		}
		return 0
	}, "THUNDERCLOCKR")

	c.addCardSoftSwitchW(0, func(_ *ioC0Page, value uint8) {
		dataIn := (value & 0x01) == 1
		clock := ((value >> 1) & 0x01) == 1
		strobe := ((value >> 2) & 0x01) == 1
		command := (value >> 3) & 0x07
		/* fmt.Printf("[cardThunderClock] dataIn %v, clock %v, strobe %v, command %v.\n",
		dataIn, clock, strobe, command) */

		c.microPD1990ac.in(clock, strobe, command, dataIn)
	}, "THUNDERCLOCKW")

	c.cardBase.assign(a, slot)
}
