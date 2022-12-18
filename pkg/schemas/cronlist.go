package schemas

import "fmt"

type CronListResponse struct {
	JobID    int    `json:"jobID"`
	Job      string `json:"job"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Schedule string `json:"schedule"`
}

func (c *CronListResponse) String() string {
	return fmt.Sprintf("\nJobID: %d\nSchedule: %s\nJob: %s\nPrevious: %s\nNext: %s\n", c.JobID, c.Schedule, c.Job, c.Previous, c.Next)
}
