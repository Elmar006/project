package api

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Elmar006/project/internal/db"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("Repeat rule is empty")
	}

	layout := db.TimeDateFormat
	start, err := time.ParseInLocation(layout, dstart, now.Location())
	if err != nil {
		return "", errors.New("Invalid start date format")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("Ivalid repeat format")
	}

	rule := parts[0]

	switch rule {
	case "d":
		return nextDateDays(now, start, parts)
	case "y":
		return nextDateYears(now, start)
	case "w":
		return nextDateWeekdays(now, start, parts)
	case "m":
		return nextDateMonths(now, start, parts)
	default:
		return "", errors.New("Unsupported repeat format")
	}
}

func nextDateDays(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("Invalid days format")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("Invalid days number")
	}

	if days < 1 || days > 400 {
		return "", errors.New("Days must be between 1 and 400")
	}

	neutralTime := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	current := neutralTime(start)
	nowPars := neutralTime(now)

	current = current.AddDate(0, 0, days)
	for !current.After(nowPars) {
		current = current.AddDate(0, 0, days)
	}

	return current.Format(db.TimeDateFormat), nil
}

func nextDateYears(now, start time.Time) (string, error) {
	neutralTime := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	current := neutralTime(start)
	nowNorm := neutralTime(now)

	addOneYear := func(t time.Time) time.Time {
		originMonth := t.Month()
		originDay := t.Day()
		nextYear := t.Year() + 1
		candidate := time.Date(nextYear, originMonth, originDay, 0, 0, 0, 0, t.Location())

		if candidate.Month() == originMonth && candidate.Day() == originDay {
			return candidate
		} else {
			nextMonth := originMonth + 1
			if nextMonth > 12 {
				nextMonth = 1
				nextYear++
			}
			return time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, t.Location())
		}
	}

	current = addOneYear(current)

	for !current.After(nowNorm) {
		current = addOneYear(current)
	}

	return current.Format(db.TimeDateFormat), nil
}

func nextDateWeekdays(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("Invalid weekdays format")
	}

	weekdayStrs := strings.Split(parts[1], ",")
	weekdays := make(map[time.Weekday]bool)
	for _, str := range weekdayStrs {
		n, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return "", err
		}
		if n < 1 || n > 7 {
			return "", errors.New("Weekdays must be between 1 and 7")
		}

		var wd time.Weekday
		if n == 7 {
			wd = time.Sunday
		} else {
			wd = time.Weekday(n)
		}
		weekdays[wd] = true
	}

	if len(weekdays) == 0 {
		return "", errors.New("No weekdays specified")
	}

	neutralTime := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	current := neutralTime(start)
	nowNorm := neutralTime(now)

	if current.After(now) && weekdays[current.Weekday()] {
		return current.Format(db.TimeDateFormat), nil
	}

	current = current.AddDate(0, 0, 1)
	for {
		if weekdays[current.Weekday()] && current.After(nowNorm) {
			return current.Format(db.TimeDateFormat), nil
		}
		current = current.AddDate(0, 0, 1)
	}
}

func nextDateMonths(now, start time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("Invalid weekdays format")
	}

	daysStrs := strings.Split(parts[1], ",")
	days := make([]int, 0, len(daysStrs))
	seenDays := make(map[int]bool)

	for _, str := range daysStrs {
		n, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return "", err
		}
		if n > 0 {
			if n > 31 {
				return "", errors.New("day must be between 1 and 31")
			}
		} else if n < 0 {
			if n != -1 && n != -2 {
				return "", errors.New("special day must be -1 or -2")
			}
		}

		if !seenDays[n] {
			seenDays[n] = true
			days = append(days, n)
		}
	}

	months := make(map[int]bool)
	if len(parts) == 3 {
		monthStrs := strings.Split(parts[2], ",")

		for _, str := range monthStrs {
			n, err := strconv.Atoi(strings.TrimSpace(str))
			if err != nil {
				return "", err
			}
			if n < 1 || n > 12 {
				return "", errors.New("No month specified")
			}
			months[n] = true
		}
	}

	allMonth := len(months) == 0

	lastDayOfMonth := func(t time.Time) int {
		nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		return nextMonth.AddDate(0, 0, -1).Day()
	}

	matchesRule := func(t time.Time) bool {
		if !allMonth && !months[int(t.Month())] {
			return false
		}

		currentDay := t.Day()
		lastDay := lastDayOfMonth(t)

		for _, d := range days {
			if d > 0 {
				if currentDay == d {
					return true
				}
			} else if d == -1 {
				if currentDay == lastDay {
					return true
				}
			} else if d == -2 {
				if currentDay == lastDay-1 {
					return true
				}
			}
		}
		return false
	}

	neutralTime := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	current := neutralTime(start)
	nowNorm := neutralTime(now)

	if current.After(nowNorm) && matchesRule(current) {
		return current.Format(db.TimeDateFormat), nil
	}

	current = current.AddDate(0, 0, 1)

	for {
		if matchesRule(current) && current.After(nowNorm) {
			return current.Format(db.TimeDateFormat), nil
		}
		current = current.AddDate(0, 0, 1)
	}
}
