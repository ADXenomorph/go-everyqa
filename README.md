# go-everyqa

https://everyqa.io/

EveryQA console tool to test go sdk

## Usage

Just go run
```
TOKEN=... PROJECT_ID=... SPRINT_ID=123 go run main.go
```

or build and do the same

```
go build
TOKEN=... PROJECT_ID=... SPRINT_ID=123 ./go-everyqa
```

Running it without any commands will show help message:
```
user:current - Get current user info
case:get - Get a list of test cases
run:get - Get test runs
run:create - Create test run
run:close - Close Test run
test:create - Create test for selected test case
test:create-action - Create resut for selected test
```