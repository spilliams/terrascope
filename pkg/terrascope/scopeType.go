package terrascope

// ScopeType represents a single scope available to a project
type ScopeType struct {
	Name         string `hcl:"name"`
	Description  string `hcl:"description,optional"`
	DefaultValue string `hcl:"default,optional"`
	// Validations  []*ProjectScopeValidation `hcl:"validation,block"`
}

func (sc *ScopeType) String() string {
	return sc.Name
}
