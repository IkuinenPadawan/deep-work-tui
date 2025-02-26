package storage

import (
	"deep-work-tui/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type jsonTimeblock struct {
	Task             string `json:"Task"`
	Starttime        string `json:"Starttime"`
	Endtime          string `json:"Endtime"`
	StartMinuteOfDay int    `json:"StartMinuteOfDay"`
	EndMinuteOfDay   int    `json:"EndMinuteOfDay"`
	DurationMinutes  int    `json:"DurationMinutes"`
}

type jsonDayRecord struct {
	Date               string          `json:"date"`
	Timeblocks         []jsonTimeblock `json:"timeblocks"`
	TotalTime          int64           `json:"total_time_seconds"`
	TotalTimeFormatted string          `json:"total_time_formatted"`
	StartTime          time.Time       `json:"start_time"`
	EndTime            time.Time       `json:"end_time"`
}

type DayRecord struct {
	Date       string             `json:"date"`
	Timeblocks []models.Timeblock `json:"timeblocks"`
	TotalTime  time.Duration      `json:"total_time"`
	StartTime  time.Time          `json:"start_time"`
	EndTime    time.Time          `json:"end_time"`
}

func formatDuration(d time.Duration) string {
	totalMinutes := int(d.Minutes())
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func SaveDayRecord(timeblocks []models.Timeblock, startTime time.Time) error {
	recordsDir := "records"
	if err := os.MkdirAll(recordsDir, 0755); err != nil {
		return fmt.Errorf("failed to create records directory: %w", err)
	}

	jsonTimeblocks := make([]jsonTimeblock, len(timeblocks))
	for i, tb := range timeblocks {
		startHour, startMin, _ := tb.Starttime.Clock()
		endHour, endMin, _ := tb.Endtime.Clock()

		startMinutes := startHour*60 + startMin
		endMinutes := endHour*60 + endMin

		durationMinutes := endMinutes - startMinutes
		if durationMinutes < 0 {
			durationMinutes += 24 * 60
		}

		jsonTimeblocks[i] = jsonTimeblock{
			Task:             tb.Task,
			Starttime:        tb.Starttime.Format("15:04"),
			Endtime:          tb.Endtime.Format("15:04"),
			StartMinuteOfDay: startMinutes,
			EndMinuteOfDay:   endMinutes,
			DurationMinutes:  durationMinutes,
		}
	}

	jsonRecord := jsonDayRecord{
		Date:               time.Now().Format("2006-01-02"),
		Timeblocks:         jsonTimeblocks,
		TotalTime:          int64(time.Since(startTime).Seconds()),
		TotalTimeFormatted: formatDuration(time.Since(startTime)),
		StartTime:          startTime,
		EndTime:            time.Now(),
	}

	recordJson, err := json.MarshalIndent(jsonRecord, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	filename := filepath.Join(recordsDir, time.Now().Format("2006-01-02")+".json")

	if err := os.WriteFile(filename, recordJson, 0644); err != nil {
		return fmt.Errorf("failed to write record file: %w", err)
	}

	return nil
}
