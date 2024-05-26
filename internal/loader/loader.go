package loader

import (
	"fmt"
	"sync"
	"time"
)

var (
	debounceDuration = 100 * time.Millisecond // Set your desired debounce duration
	timer            *time.Timer
	mu               sync.Mutex
)

// Show prints the line with debouncing to limit the rate of updates.
func Show(line string) {
	mu.Lock()
	defer mu.Unlock()

	if timer != nil {
		timer.Stop()
	}

	timer = time.AfterFunc(debounceDuration, func() {
		mu.Lock()
		defer mu.Unlock()
		clear()
		fmt.Print(line)
	})
}

// Clear prints an empty line to clear the current line immediately.
func Clear() {
	mu.Lock()
	defer mu.Unlock()
	clear()
}

func clear() {
	fmt.Print("\r\033[K")
}
