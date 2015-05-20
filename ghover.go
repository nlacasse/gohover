package ghover

import (
	"time"

	"bitbucket.org/gmcbay/i2c"
	gpio "github.com/stianeikeland/go-rpio"
)

var (
	dict = map[byte]string{
		0x22: "Right Swipe",
		0x24: "Left Swipe",
		0x28: "Up Swipe",
		0x30: "Down Swipe",
		0x41: "South Tap",
		0x42: "West Tap",
		0x44: "North Tap",
		0x48: "East Tap",
		0x50: "Center Tap",
	}
)

type ghover struct {
	address byte
	bus     *i2c.I2CBus
	ts      gpio.Pin
	reset   gpio.Pin
}

func newGhover(address, ts, reset byte) (*ghover, error) {
	bus, err := i2c.Bus(1) // rpi1 = 0, rpi2 = 1
	if err != nil {
		return nil, err
	}

	g := ghover{
		address: address,
		bus:     bus,
		ts:      gpio.Pin(ts),
		reset:   gpio.Pin(reset),
	}

	if err := gpio.Open(); err != nil {
		return nil, err
	}

	g.ts.Input()
	g.reset.Output()
	g.reset.Low()

	time.Sleep(5 * time.Millisecond)

	g.reset.High()
	g.reset.Input()

	time.Sleep(5 * time.Millisecond)

	return &g, nil
}

// TODO(nlacasse): Make this more go like.
func (g *ghover) IsReady() bool {
	if g.ts.Read() == gpio.High {
		return true
	}

	g.ts.Output()
	g.ts.Low()
	return false
}

func (g *ghover) GetEvent() (byte, error) {
	bytes, err := g.bus.ReadByteBlock(g.address, 0, 18)
	if err != nil {
		return 0, err
	}

	gestureEvent := bytes[10]
	touchEvent := (((bytes[14] & 0xe0) >> 5) | ((bytes[15] & 0x3) << 3))

	var event byte
	if gestureEvent > 1 {
		event = (1 << (bytes[10] - 1)) | 0x20
	} else if touchEvent > 0 {
		event = touchEvent | 0x80
	}
	return event, nil
}

func (g *ghover) SetRelease() {
	g.ts.High()
	g.ts.Input()
}

func (g *ghover) Close() error {
	return g.Close()
}
