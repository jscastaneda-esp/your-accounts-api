package schedule

import (
	"context"
	"time"
	"your-accounts-api/shared/infrastructure/injection"

	"codnect.io/chrono"
	"github.com/gofiber/fiber/v2/log"
)

var taskScheduler chrono.TaskScheduler
var tasks []chrono.ScheduledTask = []chrono.ScheduledTask{}

func Start() {
	taskScheduler = chrono.NewDefaultTaskScheduler()

	taskCleanLogsOrphan, err := taskScheduler.ScheduleAtFixedRate(func(ctx context.Context) {
		if err := injection.LogApp.DeleteOrphan(context.Background()); err != nil {
			log.Error(err)
		}
	}, 12*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	taskCleanLogsOld, err := taskScheduler.ScheduleAtFixedRate(func(ctx context.Context) {
		if err := injection.LogApp.DeleteOld(context.Background()); err != nil {
			log.Error(err)
		}
	}, 12*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	tasks = append(tasks, taskCleanLogsOrphan, taskCleanLogsOld)
}

func Stop() {
	for _, task := range tasks {
		task.Cancel()
	}

	shutdownChannel := taskScheduler.Shutdown()
	<-shutdownChannel
}
