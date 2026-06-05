package hcl

import "github.com/hashicorp/hcl/v2"

var s3BackendBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "bucket"},
		{Name: "dynamodb_table"},
		{Name: "encrypt"},
		{Name: "key"},
		{Name: "region"},
		{Name: "profile"},
	},
}
