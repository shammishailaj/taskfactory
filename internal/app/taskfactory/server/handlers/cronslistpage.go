package handlers

import (
	"github.com/shammishailaj/taskfactory/pkg/schemas"
	"github.com/shammishailaj/taskfactory/pkg/utils"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

type CronsListPage struct {
	u      *utils.Utils
	semVer *schemas.SemanticVersion
}

func NewCronsListPage(u *utils.Utils) *CronsListPage {
	return &CronsListPage{
		u: u,
	}
}

func (s *CronsListPage) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data := schemas.CronsListPageData{
		AppName:         strings.ToTitle(strings.ToLower(viper.Get("app_name").(string))),
		AppDesc:         viper.Get("app_desc").(string),
		RefreshInterval: int64(viper.Get("crons_list_page_refresh_interval").(int)),
		URL:             "/get/crons",
	}

	s.u.RenderTemplate(w, "", "cronslist", data)
}
