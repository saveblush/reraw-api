package utils

import (
	"fmt"

	"github.com/saveblush/reraw-api/internal/core/config"
)

var (
	patternKeyMode = "%s-%s-%s"
)

// Pointer pointer
func Pointer[Value any](v Value) *Value {
	return &v
}

// SetKey set key
// set key สำหรับเก็บ cache
func SetKey(sign, user string) string {
	return fmt.Sprintf(patternKeyMode, config.CF.App.ProjectId, sign, user)
}
