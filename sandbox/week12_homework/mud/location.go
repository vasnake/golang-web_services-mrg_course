package main

import (
	"fmt"
	"slices"
)

// Location: elements of game state
type Location struct {
	name              string
	gotoDescription   string
	lookupDescription string
	objects           []*ObjectInLocation
}

func (loc *Location) getGoToDescription() string {
	return loc.gotoDescription
}

func (loc *Location) getLookupDescription() string {
	return loc.lookupDescription
}

func (loc *Location) getName() string {
	return loc.name
}

func (loc *Location) getObject(name string) (*ObjectInLocation, error) {
	idx := slices.IndexFunc(loc.objects, func(x *ObjectInLocation) bool {
		return x.name == name
	})
	if idx < 0 {
		return nil, fmt.Errorf("Location.getObject, object '%s' not exists in location '%s'", name, loc.name)
	}

	return loc.objects[idx], nil
}

type ObjectInLocation struct {
	name           string
	compatibleWith []string
	currentState   string
	activatedState string // {9, "применить ключи дверь", "дверь открыта"}
}

func (o *ObjectInLocation) isCompatibleWith(itemName string) bool {
	return slices.Contains(o.compatibleWith, itemName)
}

// obj.activateWith(a.applyItem)
func (o *ObjectInLocation) activateWith(itemName string) {
	if !o.isCompatibleWith(itemName) {
		panic(fmt.Sprintf("ObjectInLocation.activateWith, %s not compatible with %s", o.name, itemName))
	}
	o.currentState = o.activatedState
}

// player.commandReaction(obj.getState())
func (o *ObjectInLocation) getState() string {
	// {9, "применить ключи дверь", "дверь открыта"}
	return o.currentState
}
