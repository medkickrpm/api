package utils

import (
	"fmt"
	"time"
)

func ConvertDateFormat(inputDateString, outputFormat string) (string, error) {
	// List of potential layouts to try for parsing the date
	layoutList := []string{
		"2006-01-02",
		"2/1/2006",
		// Add more layouts as needed
	}

	var parsedTime time.Time
	var err error

	// Try parsing the input date string with different layouts
	for i, layout := range layoutList {
		parsedTime, err = time.Parse(layout, inputDateString)
		fmt.Println(i, layout, parsedTime, err)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("error parsing date: %v", err)
	}

	// Format the parsed time into the desired output format
	outputDateString := parsedTime.Format(outputFormat)

	return outputDateString, nil
}
