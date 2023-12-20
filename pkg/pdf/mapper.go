package pdf

import (
	"fmt"
	"strings"
	"time"
)

var (
	layouts = []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	ErrParseTime = fmt.Errorf("not a time")
)

type Mapper func(string) string

func Map(mappers map[string]Mapper, input, mapperType string) string {
	if strings.TrimSpace(mapperType) == "" {
		return input
	}

	mapper, ok := mappers[mapperType]
	if !ok {
		return input
	}
	return mapper(input)
}

func MapGender(input string) string {
	switch input {
	case "male":
		return "мужской"
	case "female":
		return "женский"
	default:
		return input
	}
}

func parseTime(input string) (time.Time, error) {
	for _, layout := range layouts {
		time, err := time.Parse(layout, input)
		if err == nil {
			return time, nil
		}
	}
	return time.Now(), ErrParseTime
}

func MapToDDMMYYYY(input string) string {
	date, err := parseTime(input)
	if err != nil {
		return input
	}
	return fmt.Sprintf("%02d.%02d.%d", date.Day(), date.Month(), date.Year())
}

func MapToMMYYYY(input string) string {
	date, err := parseTime(input)
	if err != nil {
		return input
	}
	return fmt.Sprintf("%02d.%d", date.Month(), date.Year())
}
