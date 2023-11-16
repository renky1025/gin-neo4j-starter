package customizedmetric

import "github.com/penglongli/gin-metrics/ginmetrics"

func InitMetric() {
	// add customized metrics
	gaugeMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        "gauge_metric",
		Description: "an gauge type metric",
		Labels:      []string{"label1"},
	}

	// Add metric to global monitor object
	_ = ginmetrics.GetMonitor().AddMetric(gaugeMetric)
	//_ = ginmetrics.GetMonitor().GetMetric("digitvalue_gauge_metric").SetGaugeValue([]string{"label_value1"}, 0.1)
	_ = ginmetrics.GetMonitor().GetMetric("digitvalue_gauge_metric").Inc([]string{"label_value1"}) // value + 1
	//_ = ginmetrics.GetMonitor().GetMetric("digitvalue_gauge_metric").Add([]string{"label_value1"}, 0.2)
}
