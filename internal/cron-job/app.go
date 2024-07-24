package cronjob

import (
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// InitCronJob initializes Cron Job for Backup purpose
func InitCronJob(svc service.Service) {
	zap.S().Info("Create new cron")
	c := cron.New()

	// Add a cron job to run every midnight => "0 0 0 * * *"
	// c.AddFunc("*/1 * * * *", func() { svc.BackupAllProfiles() })  FOR EVERY MINUTE BACKUP
	c.AddFunc("0 0 0 * * *", func() { svc.BackupAllProfiles() }) //FOR EVERY MIDNIGHT BACKUP

	zap.S().Info("Start cron")
	c.Start()
}
