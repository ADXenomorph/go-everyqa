package service

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"

	everyqa "github.com/everyqa/public-api/go"
	"github.com/olekukonko/tablewriter"
)

type EveryQA struct {
	api *everyqa.APIClient
}

func NewEveryQA(api *everyqa.APIClient) *EveryQA {
	return &EveryQA{api: api}
}

func handleApiError(err error) {
	switch t := err.(type) {
	case nil:
		return
	case everyqa.GenericSwaggerError:
		showErrorAndExit(fmt.Sprintf("api error:\n\t Error: %v,\n\t Model: %v", t.Error(), t.Model()))
	default:
		handleErr(err)
	}
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

func (svc *EveryQA) GetCurrentUser(ctx context.Context) {
	user, _, err := svc.api.UsersApi.User(ctx)
	handleApiError(err)

	fmt.Printf(
		"Current user: \nID: %s\nLast name: %s\nFirst name: %s\n",
		user.UserId,
		user.LastName,
		user.FirstName,
	)
}

func (svc *EveryQA) GetCases(ctx context.Context, projectId string) {
	cases, _, err := svc.api.CasesApi.GetAllCasesByProjectId(ctx, projectId)
	handleApiError(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name"})

	for _, c := range cases {
		table.Append([]string{
			strconv.Itoa(int(c.Id)), c.Name,
		})
	}
	table.Render()
}

func (svc *EveryQA) GetRuns(ctx context.Context, projectId string) {
	runs, _, err := svc.api.RunsApi.GetAllRunsByProjectId(ctx, projectId)
	handleApiError(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name"})

	for _, run := range runs {
		table.Append([]string{
			strconv.Itoa(int(run.Id)), run.Name,
		})
	}
	table.Render()
}

func (svc *EveryQA) CreateRun(ctx context.Context, projectId string, sprintId string, assignTo string, name string) {
	dto := everyqa.ModelCreateTestRunDto{
		AssignTo: assignTo,
		Name:     name,
		SprintId: sprintId,
	}
	run, _, err := svc.api.RunsApi.CreateRun(ctx, projectId, dto)
	handleApiError(err)

	printRun(run)
}

func (svc *EveryQA) CloseRun(ctx context.Context, projectId string, runId int32) {
	run, _, err := svc.api.RunsApi.CloseRunById(ctx, projectId, runId)
	handleApiError(err)

	printRun(run)
}

func printRun(run everyqa.TestrunTestRun) {
	fmt.Printf(
		"Closed run: \n\tID: %s\n\tName: %s\n\tAssignTo: %s\n\tSprintId: %s\n\tStatus: %s\n",
		strconv.Itoa(int(run.Id)),
		run.Name,
		run.AssignedTo,
		run.SprintId,
		run.Status,
	)
}

func (svc *EveryQA) CreateTest(ctx context.Context, projectId int32, runId int32) {
	exec, _, err := svc.api.TestsApi.CreateTestByCaseId(ctx, projectId, runId)
	handleApiError(err)

	fmt.Printf("Execution: \n%+v", exec)
}

func (svc *EveryQA) CreateTestAction(ctx context.Context, projectId string, runId int32, testId int32, notes string, statusId int32) {
	dto := everyqa.ModelAddActionToTestDto{
		Notes:    notes,
		StatusId: statusId,
	}
	action, _, err := svc.api.ActionApi.CreateActionByTestId(ctx, projectId, runId, testId, dto)
	handleApiError(err)

	fmt.Printf("Action: \n%+v", action)
}
