package main

import (
	"fmt"
	"restic-host-api/connectors"
	"restic-host-api/controllers"
	"restic-host-api/models"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {

	controllers.LoadConfig()

	route := gin.Default()

	c := cron.New()
	c.Start()

	route.GET("/jobs", func(ctx *gin.Context) { controllers.GetAllJobs(ctx) })
	route.GET("/jobs/:id", func(ctx *gin.Context) { controllers.GetJob(ctx) })

	route.POST("/jobs", func(ctx *gin.Context) { controllers.CreateJob(ctx, c) })

	route.PATCH("/jobs/:id", func(ctx *gin.Context) { controllers.UpdateJob(ctx, c) })

	route.DELETE("/jobs/:id", func(ctx *gin.Context) { controllers.DeleteJob(ctx, c) })

	var jobs []models.Job
	connectors.ConnectDB(viper.GetViper().GetString("database.path")).Find(&jobs)

	for _, job := range jobs {

		backup_repo := viper.GetViper().GetString("endpoint.type") + ":" + "https://" + viper.GetViper().GetString("endpoint.username") + ":" + viper.GetViper().GetString("endpoint.password") + "@" + viper.GetViper().GetString("endpoint.url") + "/" + viper.GetViper().GetString("repo.name") + "/"
		fmt.Println(string(backup_repo))

		// Recreate Job Schedule after restarting app
		job_id := controllers.ScheduleBackup(c, job.Schedule, job.Path, backup_repo, viper.GetViper().GetString("keyfile.path"))

		// Update Job in Database
		UpdateJob := models.Job{JobID: job_id}
		connectors.ConnectDB(viper.GetViper().GetString("database.path")).Model(&job).Updates(UpdateJob)
	}

	route.Run()

}
