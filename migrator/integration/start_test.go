package integration

import (
	"github.com/Duke1616/etools/migrator/integration/startup"
	"testing"
)

func TestStart(t *testing.T) {
	app := startup.InitApp()
	for _, c := range app.Consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	err := app.AdminServer.Start()
	if err != nil {
		panic(err)
	}
}
