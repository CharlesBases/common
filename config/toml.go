package toml

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

var once sync.Once

func Config(path string, config interface{}) {
	once.Do(func() {
		filePath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		if _, err := toml.DecodeFile(filePath, config); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	})
}
