package pollers

import (
	"bufio"
	"github.com/freeformz/shh/mm"
	"github.com/freeformz/shh/utils"
	"os"
	"strings"
	"time"
)

const (
	LOAD_DATA = "/proc/loadavg"
)

type Load struct {
	measurements chan<- *mm.Measurement
}

func NewLoadPoller(measurements chan<- *mm.Measurement) Load {
	return Load{measurements: measurements}
}

func (poller Load) Poll(tick time.Time) {
	ctx := utils.Slog{"poller": poller.Name(), "fn": "Poll", "tick": tick}

	file, err := os.Open(LOAD_DATA)
	if err != nil {
		ctx.FatalError(err, "opening "+LOAD_DATA)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		ctx.FatalError(err, "reading line from "+LOAD_DATA)
	}
	fields := strings.Fields(line)
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"1m"}, utils.Atofloat64(fields[0])}
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"5m"}, utils.Atofloat64(fields[1])}
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"15m"}, utils.Atofloat64(fields[2])}
	entities := strings.Split(fields[3], "/")
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"scheduling", "entities", "executing"}, utils.Atofloat64(entities[0])}
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"scheduling", "entities", "total"}, utils.Atofloat64(entities[1])}
	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"pid", "last"}, utils.Atofloat64(fields[4])}
}

func (poller Load) Name() string {
	return "load"
}

func (poller Load) Exit() {}
