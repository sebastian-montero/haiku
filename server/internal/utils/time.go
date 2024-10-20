package utils

import "time"

func GetStringTime() string { return time.Now().Format(time.RFC3339) }
