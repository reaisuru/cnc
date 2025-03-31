package packages

import "encoding/json"

// JsonDecode is ported into the JSONFunctions map
func JsonDecode(content string) map[string]any {
	destination := make(map[string]any)
	if err := json.Unmarshal([]byte(content), &destination); err != nil {
		return make(map[string]any)
	}

	return destination
}
