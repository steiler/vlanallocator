package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func main() {

	va := NewVlanAllocator()

	entities := map[string]map[string]string{
		"bd1": {
			"name": "bd1",
			"type": "bd",
		},
		"vlan1": {
			"name": "vlan1",
			"type": "vlan",
		},
		"target-router1-eth0": {
			"name":      "target1-router1-eth0",
			"type":      "interface",
			"device":    "router1",
			"interface": "eth0",
		},
		"target-router2": {
			"name":      "target1-router2",
			"type":      "router",
			"device":    "router2",
			"interface": "eth0",
		},
		"target-router3-eth0": {
			"name":      "target1-router3-eth0",
			"type":      "interface",
			"device":    "router3",
			"interface": "eth0",
		},
		"target-router4-eth0": {
			"name":      "target1-router4-eth0",
			"type":      "interface",
			"device":    "router4",
			"interface": "eth0",
		},
	}

	va.CreateBridgeDomain(entities["bd1"])

	va.CreateVlan(6, entities["vlan1"])

	va.CreateTarget(entities["target-router1-eth0"])
	va.CreateTarget(entities["target-router2"])
	va.CreateTarget(entities["target-router3-eth0"])
	va.CreateTarget(entities["target-router4-eth0"])

	err := va.AssignVlanToBridgeDomain(6, entities["vlan1"], entities["bd1"])
	if err != nil {
		log.Error(err)
	}

	err = va.AssignVlanToBridgeDomain(6, entities["vlan1"], entities["bd1"])
	if err != nil {
		log.Error(err)
	}

	err = va.AssignVlanToTarget(7, entities["vlan1"], entities["target-router2"])
	if err != nil {
		log.Error(err)
	}

	err = va.AssignTargetToBridgeDomain(entities["target-router2"], entities["bd1"])
	if err != nil {
		log.Error(err)
	}

	fmt.Println(va.String(0))
}
