package main

import (
	"fmt"
	"slices"
	"strings"
)

// Location: elements of game state
type Location struct {
	name                         string
	gotoDescription              string
	gotoConditions               map[string]GotoCondition
	lookupDescription            string
	conditionalLookupDescription map[string]string
	connectedLocations           []string
	objects                      []*ObjectInLocation
	items                        []*ItemInLocation
}

// коридор, улица: дверь открыта
func (loc *Location) isGotoConditionExists(target string) bool {
	_, isExists := loc.gotoConditions[target]
	return isExists
}

// condition, negativeReaction := currLocObj.getGotoCondition(targetLocation) // дверь открыта, дверь закрыта
func (loc *Location) getGotoCondition(target string) (condition, negReaction string) {
	gc, isExists := loc.gotoConditions[target]
	if !isExists {
		panic(fmt.Errorf("Location.getGotoCondition(target) failed, target %s not registered", target))
	}
	return gc.condition, gc.negativeReaction
}

func (from *Location) isLocationsConnected(to *Location) bool {
	return slices.Contains(from.connectedLocations, to.name)
}

func (loc *Location) getGoToDescription() string {
	return loc.gotoDescription
}

func (loc *Location) getLookupDescription(template string) string {
	var _ = `
"на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор",
после надевания рюкзака:
на столе: ключи, конспекты. можно пройти - коридор

	`

	// items
	itemsBuf := ""

	items := make([]*ItemInLocation, 0, len(loc.items))
	for _, itm := range loc.items {
		if !itm.collected {
			items = append(items, itm)
		}
	}

	if len(items) != 0 {
		currentPrefix := items[0].prefix
		itemsBuf = currentPrefix
		for _, itm := range items {
			if itm.prefix == currentPrefix {
				if len(itemsBuf) > len(currentPrefix) {
					itemsBuf = itemsBuf + ", "
				}
				itemsBuf = itemsBuf + itm.name
			} else {
				// next collection
				currentPrefix = itm.prefix
				itemsBuf = itemsBuf + ", " + currentPrefix
				itemsBuf = itemsBuf + itm.name
			}
		}
	} else {
		// no items
		itemsBuf = "пустая комната"
	}

	// connectedLocations
	connBuf := ""

	if len(loc.connectedLocations) > 0 {
		connBuf = connBuf + "можно пройти - "
		for _, conn := range loc.connectedLocations {
			if !strings.HasSuffix(connBuf, " - ") {
				connBuf = connBuf + ", "
			}
			connBuf = connBuf + conn
		}
	}

	if template == "" {
		return fmt.Sprintf(loc.lookupDescription, itemsBuf, connBuf)
	}
	return fmt.Sprintf(template, itemsBuf, connBuf)
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

func (loc *Location) removeItem(name string) error {
	for _, itm := range loc.items {
		if itm.name == name {
			itm.collected = true
		}
	}
	return nil
}

// if location.isItemExists(t.item) {
func (loc *Location) isItemExists(name string) bool {
	for _, itm := range loc.items {
		if itm.name == name && !itm.collected {
			return true
		}
	}
	return false
}

type GotoCondition struct {
	condition        string // дверь открыта
	negativeReaction string // дверь закрыта
	action           string // применить ключи дверь
}

type ItemInLocation struct {
	collected bool

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
