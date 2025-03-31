package packages

import "net/http"

// HTTP is the import path for the package
const HTTP = "http"

// Welcome to Swash, the powerful programming language that brings the robustness and versatility of the Go language's HTTP package to your fingertips!
// Designed as a port of Go's HTTP package, Swash offers seamless integration for making HTTP requests and handling responses in your applications.
//
// With Swash's HTTP package, you can leverage the battle-tested features of Go's HTTP library to effortlessly handle HTTP requests and interact with APIs.
// Whether you're fetching data, sending data, or manipulating responses, Swash's HTTP package provides a solid foundation for efficient and secure communication.
//
// Just like its Go counterpart, Swash's HTTP package emphasizes performance, reliability, and simplicity.
// You can effortlessly make HTTP requests, handle cookies and sessions, implement middleware, and perform advanced operations like secure HTTPS connections.
// Swash's HTTP package is designed to empower you with the tools you need to interact with web servers and build data-driven applications.
//
// Get ready to experience the power of Go's HTTP package in the Swash programming language.
// Unlock your potential in making HTTP requests, integrating with APIs, and building dynamic applications with ease.
// Let's embark on a journey of seamless HTTP communication!

// HTTPFunctions is a map that associates string keys with any type of value in Swash.
// It serves as a registry for storing and accessing various HTTP functions that can be used in your applications.
// Each key in the map represents a specific HTTP function, and the associated value can be any type.
// This flexibility allows you to register and retrieve functions for handling different HTTP-related logic.
// By using the HTTPFunctions map, you can easily organize and manage your HTTP-related functionalities in a centralized manner.
var HTTPFunctions = map[string]any{
	"GET":  string(http.MethodGet),
	"PUT":  string(http.MethodPut),
	"POST": string(http.MethodPost),

	// Make is the initializer for the http builder functions
	"make": Make,
}
