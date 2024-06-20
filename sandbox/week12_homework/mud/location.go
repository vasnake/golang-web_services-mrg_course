package main

import (
	"fmt"
	"slices"
)

// Location: elements of game state
type Location struct {
	name               string
	gotoDescription    string
	lookupDescription  string
	connectedLocations []string
	objects            []*ObjectInLocation
	items              []*ItemInLocation

	// lookupDescription: "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор",
	// после надевания рюкзака:
	// на столе: ключи, конспекты. можно пройти - коридор
}

func (from *Location) isLocationsConnected(to *Location) bool {
	return slices.Contains(from.connectedLocations, to.name)
}

func (loc *Location) getGoToDescription() string {
	return loc.gotoDescription
}

func (loc *Location) getLookupDescription() string {
	var _ = `
		lookupDescription:  
			"на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор",
			"на столе: ключи, конспекты. можно пройти - коридор"

		connectedLocations: []string{"коридор"},
		objects:            []*ObjectInLocation{},
		items: []*ItemInLocation{
			{name: "ключи", prefix: "на столе: "},
			{name: "конспекты", prefix: "на столе: "},
			{name: "рюкзак", prefix: "на стуле - "},
		},
	`
	// return loc.lookupDescription
	panic("TODO: func (loc *Location) getLookupDescription() string")
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

type ItemInLocation struct {
	// lookupDescription: "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор",
	// после надевания рюкзака:
	// на столе: ключи, конспекты. можно пройти - коридор

	// ключи, рюкзак, ...
	name string

	// `на столе: `, `на стуле - `, ...
	prefix string
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
