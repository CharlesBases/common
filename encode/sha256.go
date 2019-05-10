package encode

import (
	"crypto/sha256"
	"fmt"
)

const (
	sign = "nk76Nk66^$ksdHKn123jby89kKHKb412"
)

func Sha256(password string) string {
	return fmt.Sprintf(`%x`, sha256.Sum256([]byte(fmt.Sprintf(`%s_%s`, password, sign))))
}
