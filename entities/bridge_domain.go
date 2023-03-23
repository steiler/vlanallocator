package entities

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/steiler/vlanallocator/utils"
)

type BridgeDomainManager struct {
	bridgeDomains []*BridgeDomain
}

func NewBridgeDomainManager() *BridgeDomainManager {
	return &BridgeDomainManager{
		bridgeDomains: []*BridgeDomain{},
	}
}

// GetExisting retrieves an existing Bridgedomain from the Manager.
// returns the existing BridgeDomain with matching labels, nil otherwise
func (bdm *BridgeDomainManager) GetExisting(labels *Labels) *BridgeDomain {
	// iterate through known bridgedomains
	for _, bd := range bdm.bridgeDomains {
		// return instance if labels match exactly
		if bd.MatchOnLabels(labels) {
			return bd
		}
	}
	return nil
}

// GetExistingOrNew is a wrapper around GetExistingBridgeDomain() and NewBridgeDomain().
// It will try to find the BridgeDomain via the label set if there is no match, it will return
// a new BridgeDomain with the given instance.
func (bdm *BridgeDomainManager) GetExistingOrNew(labels *Labels) (*BridgeDomain, error) {
	// retrieve existing BridgeDomain
	bd := bdm.GetExisting(labels)
	// if nil returned, create new instance
	if bd == nil {
		var err error
		bd, err = bdm.New(labels)
		if err != nil {
			return nil, err
		}
	}
	// return the result
	return bd, nil
}

// NewBridgeDomain will create and return a new BridgeDomain instance that is registered in the BridgeDomainManager
func (bdm *BridgeDomainManager) New(labels *Labels) (*BridgeDomain, error) {
	if bdm.GetExisting(labels) != nil {
		return nil, fmt.Errorf("bridgedomain with given labels (%s) already exists", labels.StringOneLine())
	}

	// create new BridgeDomain
	bd := newBridgeDomain(labels)
	// append to known instances
	bdm.bridgeDomains = append(bdm.bridgeDomains, bd)
	// return new instance
	return bd, nil
}

func (bdm *BridgeDomainManager) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sBridgeDomainManager: \n", white)
	for _, b := range bdm.bridgeDomains {
		result = result + b.String(indent+IndentSpace)
	}
	return result
}

type BridgeDomain struct {
	labels  *Labels
	targets []*Target
	vlans   map[uint8]*Vlan
}

// newBridgeDomain creates a new BridgeDomain
func newBridgeDomain(labels *Labels) *BridgeDomain {
	return &BridgeDomain{
		labels:  labels,
		targets: []*Target{},
		vlans:   map[uint8]*Vlan{},
	}
}

func (bd *BridgeDomain) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sBridgeDomain: \n", white)
	if labelString := bd.labels.String(indent + IndentSpace); len(labelString) > 0 {
		result = result + labelString
	}
	return result
}

// MatchOnLabels performs an equality match between the BridgeDomain assigned
// and the provided labels.
func (bd *BridgeDomain) MatchOnLabels(l *Labels) bool {
	return bd.labels.Equals(l)
}

// AssignVlan assigns a specific Vlan to a BridgeDomain
func (bd *BridgeDomain) AssignVlan(vlan *Vlan) error {
	if _, exists := bd.vlans[vlan.id]; exists {
		return fmt.Errorf("vlan with id %d [ %s ] already exists on BridgeDomain [ %s ]", vlan.id, vlan.labels.StringOneLine(), bd.labels.StringOneLine())
	}
	bd.vlans[vlan.id] = vlan
	return nil
}

func (bd *BridgeDomain) AddTarget(target *Target) error {
	for _, t := range bd.targets {
		if t.MatchOnLabels(target.labels) {
			log.Debug("target [ %s ] already assigned to bridgedomain [ %s ], continuing", target.labels.StringOneLine(), bd.labels.StringOneLine())
			return nil
		}
	}

	bd.targets = append(bd.targets, target)
	return nil
}
