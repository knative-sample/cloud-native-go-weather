all: weather

weather:
	@echo "run weather"
	go run cmd/weather/main.go --zipkin-endpoint="http://tracing-analysis-dc-qd.aliyuncs.com/adapt_a92srsbtkl@6580e31a949f4eb_a92srsbtkl@53df7ad2afe8301/api/v2/spans"

city:
	@echo "run city"
	go run cmd/city/main.go --port=9090 --zipkin-endpoint="http://tracing-analysis-dc-qd.aliyuncs.com/adapt_a92srsbtkl@6580e31a949f4eb_a92srsbtkl@53df7ad2afe8301/api/v2/spans"


detail:
	@echo "run detail"
	go run cmd/detail/main.go --port=9091 --zipkin-endpoint="http://tracing-analysis-dc-qd.aliyuncs.com/adapt_a92srsbtkl@6580e31a949f4eb_a92srsbtkl@53df7ad2afe8301/api/v2/spans"



