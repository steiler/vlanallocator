package entities

import (
	"fmt"

	"github.com/steiler/vlanallocator/utils"
)

type VlanManager struct {
	vlans []*Vlan
}

func NewVlanManager() *VlanManager {
	return &VlanManager{
		vlans: []*Vlan{},
	}
}

// GetExisting retrieves an existing Vlan from the Manager.
// returns the existing Vlan with matching labels, nil otherwise
func (vm *VlanManager) GetExisting(labels *Labels) *Vlan {
	// iterate through known Vlan
	for _, v := range vm.vlans {
		// return instance if labels match exactly
		if v.MatchOnLabels(labels) {
			return v
		}
	}
	return nil
}

// GetExistingOrNew is a wrapper around GetExisting() and New().
// It will try to find the Vlan via the label set if there is no match, it will return
// a new Vlan with the given instance.
func (vm *VlanManager) GetExistingOrNew(id uint8, labels *Labels) *Vlan {
	// retrieve existing Vlan
	t := vm.GetExisting(labels)
	// if nil returned, create new instance
	if t == nil {
		t = vm.New(id, labels)
	}
	// return the result
	return t
}

// New will create and return a new Vlan instance that is registered in theVlanManager
func (vm *VlanManager) New(id uint8, labels *Labels) *Vlan {
	// create new Vlan
	v := newVlan(id, labels)
	// append to known instances
	vm.vlans = append(vm.vlans, v)
	// return new instance
	return v
}

func (vm *VlanManager) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sVlanManager: \n", white)
	for _, v := range vm.vlans {
		result = result + v.String(indent+IndentSpace)
	}
	return result
}

type Vlan struct {
	id           uint8
	labels       *Labels
	bridgeDomain *BridgeDomain
}

func newVlan(id uint8, labels *Labels) *Vlan {
	return &Vlan{
		id:     id,
		labels: labels,
	}
}

func (v *Vlan) GetLabels() *Labels {
	if v.labels == nil {
		v.labels = NewLabels(nil)
	}
	return v.labels
}

func (v *Vlan) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sVlanID: %d\n", white, v.id)
	if labelString := v.labels.String(indent + IndentSpace); len(labelString) > 0 {
		result = result + labelString
	}
	return result
}

// MatchOnLabels performs an equality match between the Vlan assigned
// and the provided labels.
func (v *Vlan) MatchOnLabels(l *Labels) bool {
	return v.labels.Equals(l)
}
