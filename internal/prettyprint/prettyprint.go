package prettyprint

import (
	"encoding/json"
	"fmt"
)

func Print(value any) {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
