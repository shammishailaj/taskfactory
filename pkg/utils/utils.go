package utils

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

type Utils struct {
	Log       *log.Logger
	AwsRegion string
	Crons     *cron.Cron
}

func NewUtils(l *log.Logger) *Utils {
	return &Utils{
		Log:       l,
		AwsRegion: "",
	}
}

func (u *Utils) SetAwsRegion(awsRegion string) {
	u.AwsRegion = awsRegion
}

func (u *Utils) SetCron(crons *cron.Cron) {
	u.Crons = crons
}
