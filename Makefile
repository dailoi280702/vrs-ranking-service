.PHONY: gen
gen:
	@mockery

.PHONY: test
test:
	@mkdir tmp 2>/dev/null || true
	@go test -race -v -coverprofile=tmp/coverage.out -count=1 ./usecase/... 
	@go tool cover -func=tmp/coverage.out
	@go tool cover -html=tmp/coverage.out -o tmp/coverage.html

.PHONY: coverage
coverage:
	@open tmp/coverage.html

.PHONY: docs
docs:
	@swag init
	@swag fmt

.PHONY: dev
dev:
	@air main.go
