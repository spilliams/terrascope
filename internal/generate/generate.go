// Package generate contains structs and functions for generating terrascope
// files
package generate

type Runner interface {
	Run() error
}
