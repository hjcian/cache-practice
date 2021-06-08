# Go parameters
#= reminders
#= GOFLAGS="-count=1" ---> turn off test caching
#= -covermode=count   ---> how many times did each statement run?
#=			=atomic   ---> like count, but counts precisely in parallel programs
GOENV=CGO_ENABLED=0 GOFLAGS="-count=1"
GOCMD=$(GOENV) go
GOGET=go get
GOTEST=$(GOCMD) test -covermode=atomic -coverprofile=./coverage.out -v -timeout=20m

.EXPORT_ALL_VARIABLES:
APP_PORT?=8080

.PHONY: run
run:
	@${GOCMD} run main.go


.PHONY: tidy
tidy:
	go mod tidy