module github.com/grafana/alloy/otelcol/extension/alloyengine

go 1.25.1

replace github.com/grafana/alloy => ../..

replace github.com/grafana/alloy/syntax => ../../syntax

require (
	github.com/grafana/alloy v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.9.1
	go.opentelemetry.io/collector/component v1.44.0
	go.opentelemetry.io/collector/extension v1.44.0
	go.uber.org/zap v1.27.0
)