package cronjob

import (
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// InitCronJob initializes Cron Job for Backup purpose
func InitCronJob(svc service.Service) {
	zap.S().Info("Cron Job Initiated!")
	c := cron.New()

	BackupAllProfilesJob(svc, c)

	zap.S().Info("Cron Job Started...")
	c.Start()
}

// BackupAllProfilesJob returns a service func that runs cron job
func BackupAllProfilesJob(svc service.Service, cron *cron.Cron) {
	cron.AddFunc("0 0 * * *", func() { svc.BackupAllProfiles() }) // FOR EVERY MIDNIGHT BACKUP
}
