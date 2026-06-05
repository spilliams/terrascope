package tfgraph

import "testing"

func TestNameParts(t *testing.T) {
	cases := map[string][]string{
		"[root] module.foo[\"asdf.ghjkl\"].module.bar.aws_route53_record.extra[\"alpha.beta.domain.com\"] (expand)": {"module", "foo[\"asdf.ghjkl\"]", "module", "bar", "aws_route53_record", "extra[\"alpha.beta.domain.com\"]"},
		"[root] module.tld.aws_acm_certificate.main (expand)":                                                       {"module", "tld", "aws_acm_certificate", "main"},
		"[root] meta.count-boundary (EachMode fixup)":                                                               {"meta", "count-boundary"},
		"[root] module.tld.provider[\"registry.terraform.io/hashicorp/aws\"].acm_provider":                          {"module", "tld", "provider[\"registry.terraform.io/hashicorp/aws\"]", "acm_provider"},
		"[root] output.name": {"output", "name"},
		"[root] root":        {"root"},
		"[root] provider[\"registry.terraform.io/hashicorp/aws\"] (close)": {"provider[\"registry.terraform.io/hashicorp/aws\"]"},
	}

	for c, expected := range cases {
		t.Run(c, func(t *testing.T) {
			actual := nameParts(c)
			if len(actual) != len(expected) {
				t.Logf("expected: %v", expected)
				t.Logf("actual:   %v", actual)
				t.Fatal("expected and actual don't have the same length")
			}
			for i, name := range expected {
				if actual[i] != name {
					t.Logf("expected: %v", expected)
					t.Logf("actual:   %v", actual)
					t.Fatalf("expected and actual differ at index %d", i)
				}
			}
		})
	}
}

func TestClusterPath(t *testing.T) {
	cases := map[string][]string{
		"[root] aws_route53_record.blog (expand)":                {"aws_route53_record.blog"},
		"[root] data.aws_route53_zone.site (expand)":             {"data.aws_route53_zone.site"},
		"[root] module.blog.aws_route53_zone.site (expand)":      {"module.blog", "aws_route53_zone.site"},
		"[root] module.blog.data.aws_route53_zone.site (expand)": {"module.blog", "data.aws_route53_zone.site"},
	}

	for c, expected := range cases {
		t.Run(c, func(t *testing.T) {
			actual := clusterPath(c)
			if len(actual) != len(expected) {
				t.Logf("expected: %v", expected)
				t.Logf("actual:   %v", actual)
				t.Fatal("expected and actual don't have the same length")
			}
			for i, name := range expected {
				if actual[i] != name {
					t.Logf("expected: %v", expected)
					t.Logf("actual:   %v", actual)
					t.Fatalf("expected and actual differ at index %d", i)
				}
			}
		})
	}
}
