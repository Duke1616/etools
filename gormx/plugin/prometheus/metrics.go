package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

type prometheusPlugin struct {
	vector *prometheus.SummaryVec
}

func (p prometheusPlugin) Name() string {
	//TODO implement me
	panic("implement me")
}

func (p prometheusPlugin) Initialize(db *gorm.DB) error {
	//TODO implement me
	panic("implement me")
}
