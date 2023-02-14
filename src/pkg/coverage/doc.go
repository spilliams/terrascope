// Package coverage is a helper library that seeks to provide coverage
// about a terraform module. You can call it on a terraform root directory
// (where a complete terraform configuration is held), and it will report to you
// which resources are "covered" and which are not.
//
// It does this by comparing its understanding of what's in the module or root
// (given by github.com/hashicorp/terraform-config-inspect/tfconfig) with any
// number of "actuals" representing plan output (in the format of the Plan
// struct from github.com/hashicorp/terraform-json/).
//
// The basic process looks like this:
// 1. make a new Report, and tell it which module to look at.
// 2. Apply that module
// 3. Run a `terraform show` and add that to the Report using `AddCoverage`
// 4. Ask the report for its Coverage, and get both a percentage and a list of
// un-covered resource addresses in response.
//
// See coverage_test.go for implementation examples.
package coverage
