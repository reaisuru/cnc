package packages

// JSON is the require path
const JSON = "encoding/json"

// Swash JSON Package (Ported from Golang)
//
// This package provides essential JSON functionality for the Swash programming language.
// It offers functions for encoding and decoding JSON data.
// The package is a port of the original JSON package written in Golang, ensuring reliability and performance.
// Developers can seamlessly integrate this package into their Swash codebase to handle JSON operations with ease.

// JSONFunctions is a map[string]any that serves as a registry for all JSON-related functions.
// The keys in this map represent the names of the available JSON functions, while the values are the corresponding function objects.
// This map allows easy access to various JSON operations, such as encoding and decoding
// By maintaining all JSON functions in a single map, it provides a centralized and convenient way to utilize JSON functionality in Swash.
var JSONFunctions = map[string]any{

	// encode converts the given data into a JSON string representation.
	"encode": JsonEncode,

	// decode converts the provided JSON string into an map representation for Go2Swash.
	"decode": JsonDecode,
}
