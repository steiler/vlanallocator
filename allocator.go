package main

import (
	"fmt"

	. "github.com/steiler/vlanallocator/entities"
	"github.com/steiler/vlanallocator/utils"
)

// VlanAllocator is the facade of the VlanAllocator component
type VlanAllocator struct {
	vlanManager         *VlanManager
	bridgeDomainManager *BridgeDomainManager
	targetManager       *TargetManager
}

func NewVlanAllocator() *VlanAllocator {
	return &VlanAllocator{
		vlanManager:         NewVlanManager(),
		bridgeDomainManager: NewBridgeDomainManager(),
		targetManager:       NewTargetManager(),
	}
}

func (va *VlanAllocator) CreateBridgeDomain(labels map[string]string) {
	// convert maps[string]string into labels
	l := NewLabels(labels)
	// create a new BridgeDomain
	va.bridgeDomainManager.New(l)
}

func (va *VlanAllocator) CreateVlan(id uint8, labels map[string]string) {
	// convert maps[string]string into labels
	l := NewLabels(labels)
	// cretae a new Vlan
	va.vlanManager.New(id, l)
}

func (va *VlanAllocator) CreateTarget(labels map[string]string) {
	// convert maps[string]string into labels
	l := NewLabels(labels)
	// cretae a new Vlan
	va.targetManager.New(l)
}

func (va *VlanAllocator) AssignVlanToBridgeDomain(vid uint8, vlanLabels map[string]string, bridgeDomainLabels map[string]string) error {
	// createt vlan labels instance
	vll := NewLabels(vlanLabels)

	// retrieve vlan instance from vlan manager
	vlan := va.vlanManager.GetExisting(vll)
	if vlan == nil {
		return fmt.Errorf("vlan %s does not exsit", vll.StringOneLine())
	}

	// createt bridgedomain labels instance
	bdl := NewLabels(bridgeDomainLabels)

	// retrieve BridgeDomain instance from BridgeDomainManager
	bd := va.bridgeDomainManager.GetExisting(bdl)
	if bd == nil {
		return fmt.Errorf("bridgeDomain %s does not exsit", bdl.StringOneLine())
	}

	return bd.AssignVlan(vlan)
}

func (va *VlanAllocator) AssignVlanToTarget(vid uint8, vlanLabels map[string]string, targetLabels map[string]string) error {
	// createt vlan labels instance
	vll := NewLabels(vlanLabels)

	// retrieve vlan instance from vlan manager
	vlan := va.vlanManager.GetExisting(vll)
	if vlan == nil {
		return fmt.Errorf("vlan %s does not exsit", vll.StringOneLine())
	}

	// createt target labels instance
	tl := NewLabels(targetLabels)

	// retrieve Target instance from TargetManager
	t := va.targetManager.GetExisting(tl)
	if t == nil {
		return fmt.Errorf("target %s does not exsit", tl.StringOneLine())
	}

	return t.AssignVlan(vlan)
}

func (va *VlanAllocator) AssignTargetToBridgeDomain(targetLabels map[string]string, bridgeDomainLabels map[string]string) error {
	// createt Target labels instance
	tl := NewLabels(targetLabels)
	// retrieve Target instance from TargetManager
	target := va.targetManager.GetExisting(tl)
	if target == nil {
		return fmt.Errorf("target %s does not exsit", tl.StringOneLine())
	}
	// createt bridgedomain labels instance
	bdl := NewLabels(bridgeDomainLabels)

	// retrieve BridgeDomain instance from BridgeDomainManager
	bd := va.bridgeDomainManager.GetExisting(bdl)
	if bd == nil {
		return fmt.Errorf("bridgeDomain %s does not exsit", bdl.StringOneLine())
	}

	// setting the bridgereference on the target
	err := target.SetBridgeDomain(bd)
	if err != nil {
		return err
	}

	// adding the target to the bridge
	return bd.AddTarget(target)
}

func (va *VlanAllocator) String(indent int) string {
	// get whitespace for indention
	white := utils.GetWhitespaces(indent)

	// prepare sting that represents the Vlan
	result := fmt.Sprintf("%sVlanAllocator:\n", white)
	result = result + va.vlanManager.String(indent+IndentSpace)
	result = result + va.bridgeDomainManager.String(indent+IndentSpace)
	result = result + va.targetManager.String(indent+IndentSpace)

	return result
}
