package usta

import (
	"os"
)

// var instead of const so we can change it in tests
var baseURL = "https://leagues.ustanorcal.com"

func useMockData() bool {
	return os.Getenv("USE_MOCK_DATA") != ""
}

func ptrTo[T any](val T) *T {
	return &val
}
