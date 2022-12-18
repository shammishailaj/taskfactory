package handlers

import (
	"fmt"
	jsonresp "github.com/shammishailaj/taskfactory/pkg/http/response/json"
	"github.com/shammishailaj/taskfactory/pkg/schemas"
	"github.com/shammishailaj/taskfactory/pkg/utils"
	"net/http"
)

type CronsList struct {
	u      *utils.Utils
	semVer *schemas.SemanticVersion
}

func NewCronsList(u *utils.Utils) *CronsList {
	return &CronsList{
		u: u,
	}
}

func (s *CronsList) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	c := &utils.TaskFactory{
		Tasks: s.u.Crons,
		U:     s.u,
	}

	// Schedule the command
	cronsList := c.GetCronJobs()

	cronsListResponse := make([]schemas.CronListResponse, len(cronsList))

	for key, cron := range cronsList {
		s.u.Log.Infof("Job Entry #%d:=======================", key)
		cronsListResponse[key].Job = fmt.Sprintf("%d", cron.Job)
		cronsListResponse[key].Schedule = fmt.Sprintf("%d", cron.Schedule)
		cronsListResponse[key].Previous = cron.Prev.String()
		cronsListResponse[key].Next = cron.Next.String()
		s.u.Log.Infof("%s", cronsListResponse[key].String())
		s.u.Log.Infof("\n")
	}

	jsonresp.OK(w, r, cronsListResponse)
}
