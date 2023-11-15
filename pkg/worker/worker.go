package worker

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func RunCPTWorkerForPatient(patientID uint) {
	if err := processCPTCode99457(patientID); err != nil {
		fmt.Println(err)
	}
}

func RunCPTWorker() {
	s := gocron.NewScheduler(time.UTC)

	job99453, _ := s.Tag("99453").Every(1).Day().At("05:30").Do(func() { // 1:00 AM CT USA
		if err := processCPTCode99453(); err != nil {
			fmt.Println(err)
		}
	})

	job99454, _ := s.Tag("99454").Cron("30 5 16-31 * *").Do(func() { // 11:00 PM CT USA
		if err := processCPTCode99454(); err != nil {
			fmt.Println(err)
		}
	})

	job99457, _ := s.Tag("99457").Every(1).Day().At("05:30").Do(func() { // 1:00 AM CT USA
		if err := processCPTCode99457(); err != nil {
			fmt.Println(err)
		}
	})

	s.RunAllWithDelay(time.Second * 2)

	_, _ = s.Tag("Worker").Every(6).Hour().Do(func() {
		fmt.Println("Run CPT Worker At:")
		fmt.Println("	99453: ", job99453.NextRun())
		fmt.Println("	99454: ", job99454.NextRun())
		fmt.Println("	99457: ", job99457.NextRun())
	})

	s.StartBlocking()
}
