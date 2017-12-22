package conn

const (
	STATUS_CONNECTING = iota
	STATUS_CONNECTED
	STATUS_ERROR
)

const (
	TCP_PINGTICK_PERIOD  = 60
	UDP_PING_TICK_PERIOD = 10
	UDP_GC_PERIOD        = 90
)

const (
	BW_SCALE = 24
	BW_UNIT  = 1 << BW_SCALE
)

const (
	BBR_SCALE = 8
	BBR_UNIT  = 1 << BBR_SCALE
)

const (
	highGain  = BBR_UNIT*2885/1000 + 1
	drainGain = BBR_UNIT * 1000 / 2885
	cwndGain  = BBR_UNIT * 2

	fullBwThresh = BBR_UNIT * 5 / 4
	fullBwCnt    = 3
)

var (
	pacingGain = [...]float64{
		1.25, 0.75, 1, 1, 1, 1, 1, 1,
	}
)

type mode int

const (
	startup mode = iota
	drain
	probeBW
)
