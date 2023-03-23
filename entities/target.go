package entities

import (
	"fmt"

	"github.com/steiler/vlanallocator/utils"
)

type TargetManager struct {
	targets []*Target
}

func NewTargetManager() *TargetManager {
	return &TargetManager{
		targets: []*Target{},
	}
}

// GetExisting retrieves an existing Target from the Manager.
// returns the existing Target with matching labels, nil otherwise
func (tm *TargetManager) GetExisting(labels *Labels) *Target {
	// iterate through known Targets
	for _, t := range tm.targets {
		// return instance if labels match exactly
		if t.MatchOnLabels(labels) {
			return t
		}
	}
	return nil
}

// GetExistingOrNew is a wrapper around GetExisting() and New().
// It will try to find the Target via the label set if there is no match, it will return
// a new Target with the given instance.
func (tm *TargetManager) GetExistingOrNew(labels *Labels) *Target {
	// retrieve existing Target
	t := tm.GetExisting(labels)
	// if nil returned, create new instance
	if t == nil {
		t = tm.New(labels)
	}
	// return the result
	return t
}

// New will create and return a new Target instance that is registered in the TargetManager
func (tm *TargetManager) New(labels *Labels) *Target {
	// create new Target
	t := newTarget(labels)
	// append to known instances
	tm.targets = append(tm.targets, t)
	// return new instance
	return t
}

func (tm *TargetManager) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sTargetManager: \n", white)
	for _, t := range tm.targets {
		result = result + t.String(indent+IndentSpace)
	}
	return result
}

type Target struct {
	labels       *Labels
	bridgeDomain *BridgeDomain
}

func newTarget(l *Labels) *Target {
	return &Target{
		labels: l,
	}
}

func (t *Target) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sTarget: \n", white)
	if labelString := t.labels.String(indent + IndentSpace); len(labelString) > 0 {
		result = result + labelString
	}
	return result
}

// MatchOnLabels performs an equality match between the Target assigned
// and the provided labels.
func (t *Target) MatchOnLabels(l *Labels) bool {
	return t.labels.Equals(l)
}

func (t *Target) AssignVlan(v *Vlan) error {
	// make sure target is assigned to a BridgeDomain
	if t.bridgeDomain == nil {
		return fmt.Errorf("target %s not assigned to any BridgeDomain yet", t.labels.StringOneLine())
	}
	// deligate the assignment to the Bridgedomain
	return t.bridgeDomain.AssignVlan(v)
}

func (t *Target) SetBridgeDomain(bd *BridgeDomain) error {
	if t.bridgeDomain != nil {
		// already assigned to the given BD
		if t.bridgeDomain == bd {
			return nil
		}
		// already assigned to a different BD
		return fmt.Errorf("target [ %s ] is already assigned to a different bridgedomain [ %s ]", t.labels.StringOneLine(), t.bridgeDomain.labels.StringOneLine())
	}
	t.bridgeDomain = bd
	return nil
}

func (t *Target) LabelString() string {
	return t.labels.StringOneLine()
}
