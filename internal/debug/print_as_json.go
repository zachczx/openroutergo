package debug

import (
	"encoding/json"
	"fmt"
)

// PrintAsJSON prints a value as a formatted JSON string.
//
// WARNING: This function is intended for debugging purposes only. Do not use it in production.
func PrintAsJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("debug print as JSON failed to marshal: %v\n", err)
		return
	}

	fmt.Println(string(b))
}
