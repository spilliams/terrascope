package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const data = `scope "foo" "top" {
	account = "12345"
	scope "bar" "next" {}
}`

func main() {
	cfg, err := parse(data)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%+v", cfg.Scopes)
}

func parse(data string) (*config, error) {
	cfg := &config{}
	err := hclsimple.Decode("main.hcl", []byte(data), nil, cfg)
	if err != nil {
		log.Println("error decoding scope data")
		return nil, err
	}
	return cfg, nil
}

type config struct {
	Scopes []*scope `hcl:"scope,block"`
}

type scope struct {
	Type     string         `hcl:"type,label"`
	Name     string         `hcl:"name,label"`
	Children []*scope       `hcl:"scope,block"`
	Attrs    hcl.Attributes `hcl:"attrs,remain"`
}

func (s *scope) String() string {
	attrs := []string{}
	for k := range s.Attrs {
		attrs = append(attrs, k)
	}
	return fmt.Sprintf("<scope.%s.%s: %d children, attrs: %s>", s.Type, s.Name, len(s.Children), attrs)
}
