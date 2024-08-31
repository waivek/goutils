package goutils

import (
	"testing"
	"time"
    "fmt"
)


// TestTimeAgo tests the TimeAgo function
func TestTimeAgo(t *testing.T) {
	// Current time for testing
	now := time.Now().Unix()

	tests := []struct {
		name     string
		epoch    int64
		expected string
	}{
		{"Just now", now, "just now"},
		{"1 minute ago", now - 60, "1 minute ago"},
		{"2 minutes ago", now - 120, "2 minutes ago"},
		{"1 hour ago", now - 3600, "1 hour ago"},
		{"2 hours ago", now - 7200, "2 hours ago"},
		{"1 day ago", now - 86400, "1 day ago"},
		{"2 days ago", now - 86400*2, "2 days ago"},
		{"1 week ago", now - 86400*7, "1 week ago"},
		{"2 weeks ago", now - 86400*14, "2 weeks ago"},
		{"1 month ago", now - 86400*30, "1 month ago"},
		{"2 months ago", now - 86400*60, "2 months ago"},
		{"1 year ago", now - 86400*365, "1 year ago"},
		{"2 years ago", now - 86400*365*2, "2 years ago"},
	}

    fmt.Printf("\ntimeago_test.go\n\n")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeAgo(tt.epoch)
            message := "Timeago(%d) = `%v`; want `%v`"
			if result != tt.expected {
				// t.Errorf("TimeAgo(%d) = %v; want %v", tt.epoch, result, tt.expected)
                t.Errorf(message, tt.epoch, result, tt.expected)
			}
            //fmt.Printf("    [PASS] TimeAgo(%d) = %v; want %v\n", tt.epoch, result, tt.expected)
            message2 := fmt.Sprintf("    [black on green bold] PASS [/] " + message, tt.epoch, result, tt.expected)
            RichPrint(message2, false)
		})
	}
    fmt.Printf("\n")
}

