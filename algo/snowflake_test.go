package algo

import (
	"testing"

	"github.com/CharlesBases/common/log"
)

func TestNextID(t *testing.T) {
	defer log.Flush()

	tarceID := NextID()

	log.Debug("十六进制:", DecHex(tarceID))
	log.Debug("十进制:", tarceID)
	log.Debug("二进制:", DecBin(tarceID))
}
