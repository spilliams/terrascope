package coverage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	tfjson "github.com/hashicorp/terraform-json"
)

// DeterminationMode is an enum type
type DeterminationMode string

const (
	// PriorStateDeterminationMode tells a Report to use the "prior state" of its plans
	// to determine whether a resource is covered.
	PriorStateDeterminationMode DeterminationMode = "priorState"
	// PlannedValuesDeterminationMode tells a Report to use the "planned values" of its
	// plans to determine whether a resource is covered.
	PlannedValuesDeterminationMode DeterminationMode = "plannedValue"
)

// Report represents a single terraform root or module, and all the associated
// coverage plans or "actuals" that have been recorded against that
// configuration.
type Report struct {
	dir      string
	expected *tfconfig.Module
	actuals  []tfjson.Plan

	Mode DeterminationMode
}

// NewReport makes a new report out of the given directory.
// Returns errors from tfconfig parsing the directory as a module location.
func NewReport(dir string) (*Report, error) {
	module, diags := tfconfig.LoadModule(dir)
	if diags.HasErrors() {
		return nil, fmt.Errorf(diags.Error())
	}

	return &Report{
		dir:      dir,
		expected: module,
		Mode:     PlannedValuesDeterminationMode,
	}, nil
}

// AddCoverage tells the receiver to record new coverage for its resources.
func (r *Report) AddCoverage(p tfjson.Plan) {
	if r.actuals == nil {
		r.actuals = make([]tfjson.Plan, 0)
	}
	r.actuals = append(r.actuals, p)
}

// Coverage returns the percentage of resources covered, as well as a list of
// resources not covered. The way this computes "covered" depends on the
// receiver's Mode.
func (r *Report) Coverage() (float64, []string, error) {
	covered := r.generateExpectedMap()
	fmt.Printf("expected map: %#v\n", covered)

	actualResources, err := r.listActualResources()
	if err != nil {
		return 0, nil, err
	}
	for _, name := range actualResources {
		covered[name] = true
	}
	fmt.Printf("covered map: %#v\n", covered)

	total := float64(len(covered))
	var count float64
	uncovered := make([]string, 0)
	for k, v := range covered {
		if v {
			count++
			continue
		}
		uncovered = append(uncovered, k)
	}
	return (count / total), uncovered, nil
}

func (r *Report) generateExpectedMap() map[string]bool {
	expected := map[string]bool{}
	// fmt.Printf("expected managed resources: %#v\n", r.expected.ManagedResources)
	for name := range r.expected.ManagedResources {
		expected[name] = false
	}

	// fmt.Printf("expected module calls: %#v\n", r.expected.ModuleCalls)
	for name := range r.expected.ModuleCalls {
		expected[fmt.Sprintf("module.%s", name)] = false
	}
	return expected
}

func (r *Report) listActualResources() ([]string, error) {
	set := make(map[string]bool, 0)

	for _, plan := range r.actuals {
		var root *tfjson.StateModule
		switch r.Mode {
		case PlannedValuesDeterminationMode:
			root = plan.PlannedValues.RootModule
			break
		case PriorStateDeterminationMode:
			if plan.PriorState == nil {
				return []string{}, nil
			}
			root = plan.PriorState.Values.RootModule
			break
		default:
			return []string{}, fmt.Errorf("unrecognized determination-mode '%s' See documentation for coverage.DeterminationMode", string(r.Mode))
		}

		resources := listRootResources(root)

		for _, resource := range resources {
			set[resource] = true
		}
	}

	return set2list(set), nil
}

func listRootResources(root *tfjson.StateModule) []string {
	list := make([]string, 0)

	// add all the resources from the root
	for _, change := range root.Resources {
		address := change.Address
		if change.Index != nil {
			// remove the final [foo] group
			parts := strings.Split(address, "[")
			address = strings.Join(parts[0:len(parts)-1], "[")
		}
		list = append(list, address)
	}

	return list
}

// Combine incorporates the coverages of the given Report into the receiver.
func (r *Report) Combine(other *Report) error {
	if r.dir != other.dir {
		return fmt.Errorf("Cannot combine a report about %s into a report about %s (different location)", other.dir, r.dir)
	}

	r.actuals = append(r.actuals, other.actuals...)
	return nil
}

// Reset empties the receiver of its plans. This effectively sets the coverage back to 0.
// NB: This does not set the receiver's Mode to its default value.
func (r *Report) Reset() {
	r.actuals = make([]tfjson.Plan, 0)
}

func set2list(set map[string]bool) []string {
	i := 0
	list := make([]string, len(set))
	for k := range set {
		list[i] = k
		i++
	}
	return list
}
