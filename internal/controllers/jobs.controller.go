package controllers

import (
	"fmt"
	"net/http"
	"restic-host-api/connectors"
	"restic-host-api/models"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type CreateJobInput struct {
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path" binding:"required"`
	JobID    uint   `json:"jobid"`
	Schedule string `json:"schedule" binding:"required"`
}

type UpdateJobInput struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	JobID    uint   `json:"jobid"`
	Schedule string `json:"schedule"`
}

// GET /jobs
// List of all Jobs
func GetAllJobs(context *gin.Context) {
	var jobs []models.Job
	connectors.ConnectDB(viper.GetViper().GetString("database.path")).Find(&jobs)

	context.JSON(http.StatusOK, gin.H{"data": jobs})
}

// GET /jobs/:id
// Get Job by id
func GetJob(context *gin.Context) {
	var job models.Job
	if err := connectors.ConnectDB(viper.GetViper().GetString("database.path")).Where("id = ?", context.Param("id")).First(&job).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Job not found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": job})
}

// POST /jobs
// Create Job
func CreateJob(context *gin.Context, c *cron.Cron) {
	var input CreateJobInput
	// Validate input
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Add Job to Cron Scheduler
	backup_repo := viper.GetViper().GetString("endpoint.type") + ":" + "https://" + viper.GetViper().GetString("endpoint.username") + ":" + viper.GetViper().GetString("endpoint.password") + "@" + viper.GetViper().GetString("endpoint.url") + "/" + viper.GetViper().GetString("repo.name") + "/"
	fmt.Print(backup_repo)

	job_id := ScheduleBackup(c, input.Schedule, input.Path, backup_repo, viper.GetViper().GetString("keyfile.path"))

	Job := models.Job{Name: input.Name, Path: input.Path, Schedule: input.Schedule, JobID: job_id}
	// Create Job in database
	connectors.ConnectDB(viper.GetViper().GetString("database.path")).Create(&Job)
	context.JSON(http.StatusOK, gin.H{"data": Job})
}

// PATCH /jobs/:id
// Edit Job by ID
func UpdateJob(context *gin.Context, c *cron.Cron) {
	// Check if Job exist
	var job models.Job
	if err := connectors.ConnectDB(viper.GetViper().GetString("database.path")).Where("id = ?", context.Param("id")).First(&job).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Job not found!"})
		return
	}
	// Validate input
	var input UpdateJobInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update Schedule of Job
	RemoveSchedule(c, cron.EntryID(job.JobID))
	backup_repo := viper.GetViper().GetString("endpoint.type") + ":" + "https://" + viper.GetViper().GetString("endpoint.username") + ":" + viper.GetViper().GetString("endpoint.password") + "@" + viper.GetViper().GetString("endpoint.url") + "/" + viper.GetViper().GetString("repo.name") + "/"

	job_id := ScheduleBackup(c, input.Schedule, input.Path, backup_repo, viper.GetViper().GetString("keyfile.path"))

	// Get actual Job values to variable
	UpdateJob := models.Job{Name: input.Name, Path: input.Path, Schedule: input.Schedule, JobID: job_id}

	// Update Job in Database
	connectors.ConnectDB(viper.GetViper().GetString("database.path")).Model(&job).Updates(UpdateJob)

	context.JSON(http.StatusOK, gin.H{"data": UpdateJob})
}

// DELETE /jobs/:id
// Remove Job by ID
func DeleteJob(context *gin.Context, c *cron.Cron) {
	// Check if Job exist
	var job models.Job
	if err := connectors.ConnectDB(viper.GetViper().GetString("database.path")).Where("id = ?", context.Param("id")).First(&job).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Job not found!"})
		return
	}
	// Remove Job from Scheduler
	RemoveSchedule(c, cron.EntryID(job.JobID))

	// Remove Job drom database
	connectors.ConnectDB(viper.GetViper().GetString("database.path")).Delete(&job)

	context.JSON(http.StatusOK, gin.H{"data": true})
}
