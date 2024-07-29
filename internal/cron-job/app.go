package cronjob

import (
	"fmt"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// InitCronJob initializes Cron Job for Backup purpose
func InitCronJob(svc service.Service) {
	zap.S().Info("Cron Job Initiated!")
	c := cron.New()

	// Add a cron job to run every midnight => "0 0 0 * * *"
	cronExpr := fmt.Sprintf("%s %s %s %s %s %s", Seconds, Minutes, Hours, DayOfMonth, Month, DayofWeek)
	c.AddFunc(cronExpr, BackupAllProfilesJob(svc)) // FOR EVERY MIDNIGHT BACKUP

	zap.S().Info("Cron Job Started...")
	c.Start()
}

// BackupAllProfilesJob returns a service func that runs cron job
func BackupAllProfilesJob(svc service.Service) func() {
	return func() {
		svc.BackupAllProfiles()
	}
}
