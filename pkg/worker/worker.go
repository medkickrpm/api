package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func TriggerCPTWorker() {
	funcList := []func(...uint) error{
		processCPTCode99453,
		processCPTCode99454,
		processCPTCode99457,
		processCPTCode99458,
		processCPTCode99490,
		processCPTCode99439,
		processCPTCode99426,
		processCPTCode99427,
		processCPTCode99484,
	}

	for _, f := range funcList {
		if err := f(); err != nil {
			fmt.Println(err)
		}
	}
}

func RunCPTWorkerForPatient(service string, patientID uint) {
	var funcList []func(...uint) error

	if service == "RPM" {
		funcList = []func(...uint) error{
			processCPTCode99457,
			processCPTCode99458,
		}
	} else if service == "CCM" {
		funcList = []func(...uint) error{
			processCPTCode99490,
			processCPTCode99439,
		}
	} else if service == "PCM" {
		funcList = []func(...uint) error{
			processCPTCode99426,
			processCPTCode99427,
		}
	} else if service == "BHI" {
		funcList = []func(...uint) error{
			processCPTCode99484,
		}
	} else {
		return
	}

	for _, f := range funcList {
		if err := f(patientID); err != nil {
			fmt.Println(err)
		}
	}
}

func RunCPTWorker() {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	s := gocron.NewScheduler(loc)

	job99453, _ := s.Tag("99453").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99453(); err != nil {
			fmt.Println(err)
		}
	})

	job99454, _ := s.Tag("99454").Cron("30 23 16-31 * *").Do(func() {
		if err := processCPTCode99454(); err != nil {
			fmt.Println(err)
		}
	})

	job99457, _ := s.Tag("99457").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99457(); err != nil {
			fmt.Println(err)
		}
	})

	job99458, _ := s.Tag("99458").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99458(); err != nil {
			fmt.Println(err)
		}
	})

	job99490, _ := s.Tag("99490").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99490(); err != nil {
			fmt.Println(err)
		}
	})

	job99439, _ := s.Tag("99439").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99439(); err != nil {
			fmt.Println(err)
		}
	})

	job99426, _ := s.Tag("99426").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99426(); err != nil {
			fmt.Println(err)
		}
	})

	job99427, _ := s.Tag("99427").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99427(); err != nil {
			fmt.Println(err)
		}
	})

	job99484, _ := s.Tag("99484").Every(1).Day().At("23:30").Do(func() {
		if err := processCPTCode99484(); err != nil {
			fmt.Println(err)
		}
	})

	s.RunAllWithDelay(time.Second * 2)

	_, _ = s.Tag("Worker").Every(6).Hour().Do(func() {
		fmt.Println("Run CPT Worker At:")
		fmt.Println("	RPM 99453: ", job99453.NextRun())
		fmt.Println("	RPM 99454: ", job99454.NextRun())
		fmt.Println("	RPM 99457: ", job99457.NextRun())
		fmt.Println("	RPM 99458: ", job99458.NextRun())
		fmt.Println("	CCM 99490: ", job99490.NextRun())
		fmt.Println("	CCM 99439: ", job99439.NextRun())
		fmt.Println("	PCM 99426: ", job99426.NextRun())
		fmt.Println("	PCM 99427: ", job99427.NextRun())
		fmt.Println("	BHI 99484: ", job99484.NextRun())
	})

	s.StartBlocking()
}
