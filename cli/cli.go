package cli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	everyqa "github.com/everyqa/public-api/go"

	"github.com/ADXenomorph/go-everyqa/service"
)

type CommandLine struct {
	api *everyqa.APIClient
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	token := requireEnv("TOKEN")
	project := requireEnv("PROJECT_ID")
	sprint := requireEnv("SPRINT_ID")

	cfg := everyqa.NewConfiguration()
	api := everyqa.NewAPIClient(cfg)
	cli.api = api

	ctx := context.Background()
	ctx = context.WithValue(ctx, everyqa.ContextAccessToken, token)

	everyqaService := service.NewEveryQA(api)

	action := os.Args[1]
	switch action {
	case "user:current":
		everyqaService.GetCurrentUser(ctx)
	case "case:get":
		everyqaService.GetCases(ctx, project)
	case "run:get":
		everyqaService.GetRuns(ctx, project)
	case "run:create":
		cmd := flag.NewFlagSet(action, flag.ExitOnError)
		assignTo := cmd.String("a", "", "User to assign test run to")
		name := cmd.String("n", "", "Name of test run")
		err := cmd.Parse(os.Args[2:])
		
		handleErr(err)

		if *assignTo == "" || *name == "" {
			cmd.Usage()
			runtime.Goexit()
		}
		everyqaService.CreateRun(ctx, project, sprint, *assignTo, *name)
	case "run:close":
		cmd := flag.NewFlagSet(action, flag.ExitOnError)
		runId := cmd.Int("r", 0, "Run id to close")
		err := cmd.Parse(os.Args[2:])
		
		handleErr(err)
		
		if *runId == 0 {
			cmd.Usage()
			runtime.Goexit()
		}
		everyqaService.CloseRun(ctx, project, int32(*runId))
	case "test:create":
		log.Fatal("test:create doesnt work")
		// cmd := flag.NewFlagSet("test:create", flag.ExitOnError)
		// runId := cmd.Int("r", 0, "Run id to create test for")
		// err := cmd.Parse(os.Args[2:])
		// if err != nil {
		// 	log.Panic(err)
		// }
		// if *runId == 0 {
		// 	cmd.Usage()
		// 	runtime.Goexit()
		// }
		// cli.createTest(ctx, project, int32(*runId))
	case "test:create-action":
		cmd := flag.NewFlagSet(action, flag.ExitOnError)
		runId := cmd.Int("r", 0, "Run id to create action for")
		testId := cmd.Int("t", 0, "Test id to create action for")
		notes := cmd.String("n", "", "Notes for action")
		statusId := cmd.Int("s", 0, "Status ID for action")
		err := cmd.Parse(os.Args[2:])
		
		handleErr(err)

		if *runId == 0 || *testId == 0 || *notes == "" || *statusId == 0 {
			cmd.Usage()
			runtime.Goexit()
		}
		everyqaService.CreateTestAction(ctx, project, int32(*runId), int32(*testId), *notes, int32(*statusId))
	default:
		cli.printUsage()
		runtime.Goexit()
	}
}

func requireEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		showErrorAndExit(fmt.Sprintf("%s env is not set", name))
	}

	return val
}

func handleErr(err error) {
	if err != nil {
		showErrorAndExit(err.Error())
	}
}

func showErrorAndExit(msg string) {
	fmt.Println(msg)
	runtime.Goexit()
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("    user:current - Get current user info")
	fmt.Println("    case:get - Get a list of test cases")
	fmt.Println("    run:get - Get test runs")
	fmt.Println("    run:create - Create test run")
	fmt.Println("    run:close - Close Test run")
	fmt.Println("    test:create - Create test for selected test case")
	fmt.Println("    test:create-action - Create resut for selected test")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}
