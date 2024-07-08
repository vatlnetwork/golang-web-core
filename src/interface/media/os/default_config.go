package media_os_int

import (
	"fmt"
	"golang-web-core/src/application/util"
	"runtime"
)

func DefaultConfig() Config {
	appDir, _ := util.AppDir()
	if runtime.GOOS == "windows" {
		return Config{
			Directory: fmt.Sprintf("%v\\resources\\media", appDir),
		}
	}
	return Config{
		Directory: fmt.Sprintf("%v/resources/media", appDir),
	}
}
