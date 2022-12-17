package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
	"text/template"
)

var mainTmpl = "{{define \"main\" }} {{ template \"base\" . }} {{ end }}"

func loadTemplate(templateType string, templateName string) *template.Template {
	path, pathErr := filepath.Abs("./web/template/" + viper.GetString("theme") + "/" + templateType)
	if pathErr != nil {
		log.Errorf("Error getting absolute path for theme templates. %s", pathErr.Error())
	} else {
		log.Infof("Absolute path of template [%s] of type [%s] = [%s]", templateName, templateType, path)
	}

	base, baseErr := filepath.Abs("./web/template/base.gohtml")
	if baseErr != nil {
		log.Errorf("Error getting absolute path for base.gohtml. %s", baseErr.Error())
	} else {
		log.Infof("Absolute path of base template file = %s", base)
	}

	partials, partialsErr := filepath.Glob(path + "/partials/*.gohtml")
	if partialsErr != nil {
		log.Errorf("Error reading Partials. %s", partialsErr.Error())
	} else {
		log.Infof("partials = %#v", partials)
	}

	page := path + "/pages/" + templateName + ".gohtml"

	mainTemplate := template.New("main")
	mainTemplate, mainTemplateErr := mainTemplate.Parse(mainTmpl)
	if mainTemplateErr != nil {
		log.Errorf("Error parsing \"main\" template. %s", mainTemplateErr.Error())
	} else {
		log.Infof("mainTemplate = %#v", mainTemplate)
	}
	files := append(partials, base, page)
	currentTmpl, currentTmplErr := mainTemplate.Clone()
	if currentTmplErr != nil {
		log.Errorf("Error cloning mainTemplate. %s", currentTmplErr.Error())
	} else {
		log.Infof("Current template, currentTmpl = %#v", currentTmpl)
	}

	currentTemplate, currentTemplateErr := currentTmpl.ParseFiles(files...)
	if currentTmplErr != nil {
		log.Errorf("Error parsing currentTmpl. %s", currentTemplateErr.Error())
	} else {
		log.Infof("current template, currentTemplate = %#v", currentTemplate)
	}
	return template.Must(currentTemplate, currentTemplateErr)
}

// RenderTemplate serves HTML template
func (u *Utils) RenderTemplate(res http.ResponseWriter, templateType string, templateName string, data interface{}) {
	u.Log.Infof("RenderTemplate() called with templateType = %s, templateName = %s", templateType, templateName)
	tmpl := loadTemplate(templateType, templateName)
	err := tmpl.Execute(res, data)
	if err != nil {
		u.Log.Errorf("Template Error: %s", err.Error())
	} else {
		u.Log.Infof("Template %s  of type %s executed successfully", templateName, templateType)
	}
}
