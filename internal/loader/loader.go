package loader

import (
	"fmt"
)

func Show(line string) {
	fmt.Print("\r\033[K")
	fmt.Print(line)
}

func Clear() {
	fmt.Print("\r\033[K")
}
