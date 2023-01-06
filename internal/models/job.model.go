package models

type Job struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	JobID    uint   `json:"jobid"`
	Schedule string `json:"schedule"`
}
