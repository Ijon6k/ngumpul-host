package availability

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func bootID() string {
	data, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	if err == nil {
		id := strings.TrimSpace(string(data))
		if id != "" {
			return id
		}
	}
	return "ngumpul-host-boot-id"
}

func uptimeSeconds() int64 {
	data, err := os.ReadFile("/proc/uptime")
	if err == nil {
		parts := strings.Fields(string(data))
		if len(parts) > 0 {
			if secs, err := strconv.ParseFloat(parts[0], 64); err == nil {
				return int64(secs)
			}
		}
	}
	return 0
}

func formatUptime(totalSeconds int64) string {
	if totalSeconds <= 0 {
		return "Just started"
	}
	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	mins := (totalSeconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%d hours, %d mins", hours, mins)
	}
	return fmt.Sprintf("%d mins", mins)
}

func formatDuration(sec int64) string {
	if sec < 60 {
		return fmt.Sprintf("%ds", sec)
	}
	mins := sec / 60
	remSec := sec % 60
	if mins < 60 {
		if remSec > 0 {
			return fmt.Sprintf("%dm %ds", mins, remSec)
		}
		return fmt.Sprintf("%dm", mins)
	}
	hours := mins / 60
	remMins := mins % 60
	if remMins > 0 {
		return fmt.Sprintf("%dh %dm", hours, remMins)
	}
	return fmt.Sprintf("%dh", hours)
}

func formatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	diff := time.Since(t)
	if diff < 1*time.Minute {
		return "Just now"
	}
	if diff < 60*time.Minute {
		return fmt.Sprintf("%d mins ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}
	return t.Format("Jan 02, 15:04")
}
