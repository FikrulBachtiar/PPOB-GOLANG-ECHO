package utils

import (
	"fmt"
	"strings"
	"time"
)

func ConvertStringToDurationType(durationType string) (time.Duration, error) {
	switch strings.ToLower(durationType) {
	case "second":
		return time.Second, nil;
	case "minutes":
		return time.Minute, nil;
	case "hour":
		return time.Hour, nil;
	default:
		return 0, fmt.Errorf("invalid duration type");
	}
}