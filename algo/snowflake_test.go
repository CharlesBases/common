package algo

import (
	"fmt"
	"testing"
)

func TestGetTraceID(t *testing.T) {
	tarceID := GetTraceID()
	fmt.Println(tarceID)
}
