package controller

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestCronSchedule(t *testing.T) {
	tests := []struct {
		// your fields here
		name          string
		schedule      string
		lastRestart   time.Time
		expectRestart bool
		expectError   bool
	}{
		// your cases here
		{"Overdue Restart", "*/2 * * * *", time.Now().Add(-1 * time.Hour), true, false},
		{"Nothing to do Restart", "*/2 * * * *", time.Now().Add(1 * time.Minute), false, false},
		{"Invalid Schedule", "not-valid-cron", time.Now(), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// your logic here
			schedule, err := cron.ParseStandard(tt.schedule)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			nextRun := schedule.Next(tt.lastRestart)
			isTimeToRestart := !time.Now().Before(nextRun)
			if isTimeToRestart != tt.expectRestart {
				t.Errorf("got %v, want %v", isTimeToRestart, tt.expectRestart)
			}
		})
	}
}
