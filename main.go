package main

import (
	"dbgod_workOrder_job/jobs"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	jobs.StartJob()
}
