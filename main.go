package main

import (
	"github.com/quincy0/temproal-workflows/starter"
	"github.com/quincy0/temproal-workflows/worker"
)

func main() {
	// run the setup of cron workflows
	starter.InitCron()

	// initiate temporal worker
	worker.Init()
}
