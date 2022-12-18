package schemas

import (
	"fmt"
	"strings"
)

type ScheduleResponse struct {
	Args     []string `json:"args"`
	Command  string   `json:"command"`
	Error    string   `json:"error"`
	Schedule string   `json:"schedule"`
}

func (s *ScheduleResponse) String() string {
	return fmt.Sprintf("Command: %s %s Schedule: %s Error: %s", s.Command, strings.Join(s.Args, " "), s.Schedule, s.Error)
}

type ScheduleRequest struct {
	CronSchedule    string   `json:"cron_schedule"`
	CronCommand     string   `json:"cron_command"`
	CronCommandArgs []string `json:"cron_command_args"`
}

func (s *ScheduleRequest) String() string {
	return fmt.Sprintf("Command: %s %s Schedule: %s", s.CronCommand, strings.Join(s.CronCommandArgs, " "), s.CronSchedule)
}
