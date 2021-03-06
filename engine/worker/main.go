package main

import (
	"container/list"

	"google.golang.org/grpc"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/cdsclient"
)

type currentWorker struct {
	autoUpdate    bool
	singleUse     bool
	apiEndpoint   string
	token         string
	id            string
	model         sdk.Model
	groupID       int64
	bookedWJobID  int64
	nbActionsDone int
	basedir       string
	manualExit    bool
	logger        struct {
		logChan chan sdk.Log
		llist   *list.List
	}
	exportPort int32
	hatchery   struct {
		name string
	}
	grpc struct {
		address  string
		conn     *grpc.ClientConn
		insecure bool
	}
	currentJob struct {
		wJob             *sdk.WorkflowNodeJobRun
		currentStep      int
		buildVariables   []sdk.Variable
		pkey             string
		gitsshPath       string
		params           []sdk.Parameter
		secrets          []sdk.Variable
		workingDirectory string
	}
	status struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	client cdsclient.Interface
}

func main() {
	sdk.SetAgent(sdk.WorkerAgent)

	w := &currentWorker{}
	cmd := cmdMain(w)
	cmd.AddCommand(cmdExport)
	cmd.AddCommand(cmdUpload(w))
	cmd.AddCommand(cmdArtifacts(w))
	cmd.AddCommand(cmdDownload(w))
	cmd.AddCommand(cmdTmpl(w))
	cmd.AddCommand(cmdCheckSecret(w))
	cmd.AddCommand(cmdTag(w))
	cmd.AddCommand(cmdRun(w))
	cmd.AddCommand(cmdUpdate(w))
	cmd.AddCommand(cmdExit(w))
	cmd.AddCommand(cmdVersion)
	cmd.AddCommand(cmdRegister(w))
	cmd.AddCommand(cmdCache(w))
	cmd.AddCommand(cmdKey(w))

	// last command: doc, this command is hidden
	cmd.AddCommand(cmdDoc(cmd))

	cmd.Execute()
}
