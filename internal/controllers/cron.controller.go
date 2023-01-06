package controllers

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var current_status string
var job_id cron.EntryID

func ScheduleBackup(c *cron.Cron, sched string, local_patch string, remote_patch string, key_patch string) uint {
	job_id, _ = c.AddFunc(sched, func() { CreateBackup(local_patch, remote_patch, key_patch) })
	fmt.Println("Created a new Job Schedule with ID - ", int(job_id))
	return uint(job_id)
}

func RemoveSchedule(c *cron.Cron, job_id cron.EntryID) {
	c.Remove(job_id)
	fmt.Println("Removed Job from Scheduler - ", job_id)
}
