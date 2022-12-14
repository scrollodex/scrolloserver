package safeget

import (
	"fmt"
	"strconv"
)

// String returns the value as a string.
func String(f map[string]interface{}, k string) string {
	switch v := f[k].(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		fmt.Printf("DEBUG: STRING %v = %v\n", k, v)
		panic(fmt.Sprintf("unimplemented type: %T", v))
	}
}

// Int returns the value as a int, truncating or converting if needed.
func Int(f map[string]interface{}, k string) int {
	switch v := f[k].(type) {
	case nil:
		return 0
	case int:
		return v
	case float64:
		return int(v)
	case float32:
		return int(v)
	case string:
		var err error
		var r int
		r, err = strconv.Atoi(v)
		if err != nil {
			return -1
		}
		return r
	default:
		fmt.Printf("DEBUG: INT %v = %v\n", k, v)
		panic(fmt.Sprintf("unimplemented type: %T", v))
	}
}
