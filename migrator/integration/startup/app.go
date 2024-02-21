package startup

import (
	"github.com/Duke1616/etools/httpx/ginx"
	"github.com/Duke1616/etools/kafka/saramax/events"
)

type App struct {
	Consumers   []events.Consumer
	AdminServer *ginx.Server
}
