package handlers

import (
	"encoding/json"
	jsonresp "github.com/shammishailaj/taskfactory/pkg/http/response/json"
	"github.com/shammishailaj/taskfactory/pkg/schemas"
	"github.com/shammishailaj/taskfactory/pkg/utils"
	"net/http"
)

type Schedule struct {
	u      *utils.Utils
	semVer *schemas.SemanticVersion
}

func NewSchedule(u *utils.Utils) *Schedule {
	return &Schedule{
		u: u,
	}
}

func (s *Schedule) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Read the schedule and command from the request body
	var requestBody schemas.ScheduleRequest

	requestUnmarshallErr := json.NewDecoder(r.Body).Decode(&requestBody)
	if requestUnmarshallErr != nil {
		s.u.Log.Errorf("Error Unmarshalling request JSON: %s", requestUnmarshallErr.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.u.Log.Infof("Scheduled command details in request:\n %s\n", requestBody.String())

	c := &utils.TaskFactory{
		Tasks: s.u.Crons,
		U:     s.u,
	}

	// Schedule the command
	scheduleCmdErr := c.ScheduleCommand(requestBody.CronSchedule, requestBody.CronCommand, requestBody.CronCommandArgs)

	response := &schemas.ScheduleResponse{
		Args:     requestBody.CronCommandArgs,
		Command:  requestBody.CronCommand,
		Error:    "nil",
		Schedule: requestBody.CronSchedule,
	}

	s.u.Log.Infof("Scheduled command details:\n %s\n", response.String())

	if scheduleCmdErr != nil {
		response.Error = scheduleCmdErr.Error()
		jsonresp.ServerErrorString(w, r, response)
		return
	}

	jsonresp.OK(w, r, s)
}
