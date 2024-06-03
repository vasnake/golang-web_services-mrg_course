package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SellerStruct struct {
	ID    int
	Name  string
	Deals int
}

type GoodiesItemStruct struct {
	ID      int
	Name    string
	InStock int
	Seller  int // SellerStruct ref
	Catalog int // CatalogStruct ref
}

// Catalog is a tree, Catalog contains items: links(parentCatalog, childrenCatalogsList, itemsList)
// json(id, name, childsOption, itemsOption)
type CatalogStruct struct {
	ID       int
	Name     string
	Parent   int   // CatalogStruct ref or 0
	Children []int // CatalogStruct ref
	Items    []int // GoodiesItem ref
}

func loadCatalogTree(data map[string]any) ([]CatalogStruct, []GoodiesItemStruct, error) {
	catalogAny, exists := data["catalog"]
	if !exists {
		return nil, nil, fmt.Errorf("can't find catalog in given data")
	}
	catalogMap, isMap := catalogAny.(map[string]any)
	if !isMap {
		return nil, nil, fmt.Errorf("catalog type is not of map[string]any")
	}

	var catalogRows = make([]CatalogStruct, 0, 16)
	var itemRows = make([]GoodiesItemStruct, 0, 16)
	var catalogId int = 0
	err := loadCatalogTreeRecursive(catalogMap, catalogId, &catalogRows, &itemRows)

	// link parent-children ids
	for i, p := range catalogRows { // parent
		// show("parent.id: ", p.ID)
		for _, c := range catalogRows { // child
			if c.Parent == p.ID {
				// show("\tchild.parent == parent.id, p.id, c.id: ", p.ID, c.ID)
				p.Children = append(p.Children, c.ID)
			}
		}
		catalogRows[i] = p
	}

	return catalogRows, itemRows, err
}

func loadCatalogTreeRecursive(catalogMap map[string]any, parent int, catalogRows *[]CatalogStruct, itemRows *[]GoodiesItemStruct) error {
	id, err := loadIntFromMap(catalogMap, "id")
	if err != nil {
		return fmt.Errorf("load catalog id failed, %w", err)
	}
	name, err := loadStringFromMap(catalogMap, "name")
	if err != nil {
		return fmt.Errorf("load catalog name failed, %w", err)
	}

	// current catalog row
	var catalog = CatalogStruct{ID: id, Name: name, Parent: parent, Children: make([]int, 0, 16), Items: make([]int, 0, 16)}

	// items list
	itemsAny, itemsExists := catalogMap["items"]
	if itemsExists {
		itemsSlice, isSlice := itemsAny.([]any)
		if !isSlice {
			return fmt.Errorf("items is not a slice, %#v", itemsAny)
		}
		for _, itemAny := range itemsSlice {
			itemMap, isMap := itemAny.(map[string]any)
			if !isMap {
				return fmt.Errorf("item type is not of map[string]any")
			}
			item, err := loadItem(itemMap)
			if err != nil {
				return fmt.Errorf("loadItem failed, %w", err)
			}
			item.Catalog = catalog.ID
			catalog.Items = append(catalog.Items, item.ID)
			*itemRows = append(*itemRows, item)
		}
	}

	*catalogRows = append(*catalogRows, catalog)

	// sub-catalogs
	childrenAny, childrenExists := catalogMap["childs"]
	if childrenExists {
		childrenSlice, isSlice := childrenAny.([]any)
		if !isSlice {
			return fmt.Errorf("children is not a slice, %#v", childrenAny)
		}
		for _, childAny := range childrenSlice {
			childMap, isMap := childAny.(map[string]any)
			if !isMap {
				return fmt.Errorf("child type is not of map[string]any")
			}
			err = loadCatalogTreeRecursive(childMap, catalog.ID, catalogRows, itemRows)
		}
	}

	return nil
}

func loadItem(data map[string]any) (GoodiesItemStruct, error) {
	item := GoodiesItemStruct{}

	id, err := loadIntFromMap(data, "id")
	if err != nil {
		return item, fmt.Errorf("load item id failed, %w", err)
	}
	item.ID = id

	name, err := loadStringFromMap(data, "name")
	if err != nil {
		return item, fmt.Errorf("load item name failed, %w", err)
	}
	item.Name = name

	inStock, err := loadIntFromMap(data, "in_stock")
	if err != nil {
		return item, fmt.Errorf("load item in_stock failed, %w", err)
	}
	item.InStock = inStock

	seller, err := loadIntFromMap(data, "seller_id")
	if err != nil {
		return item, fmt.Errorf("load item seller_id failed, %w", err)
	}
	item.Seller = seller

	return item, nil
}

func loadSellers(data map[string]any) ([]SellerStruct, error) {
	sellersAny, exists := data["sellers"]
	if !exists {
		return nil, fmt.Errorf("can't find sellers in given data")
	}
	sellersSlice, isSliceOfAny := sellersAny.([]any)
	if !isSliceOfAny {
		return nil, fmt.Errorf("sellers type is not of []any")
	}

	var result = make([]SellerStruct, 0, 16)
	for _, sellerAny := range sellersSlice {
		sellerMap, isMap := sellerAny.(map[string]any)
		if !isMap {
			return nil, fmt.Errorf("seller type is not of map[string]any")
		}

		id, err := loadIntFromMap(sellerMap, "id")
		if err != nil {
			return nil, fmt.Errorf("load seller id failed, %w", err)
		}
		deals, err := loadIntFromMap(sellerMap, "deals")
		if err != nil {
			return nil, fmt.Errorf("load seller deals failed, %w", err)
		}
		name, err := loadStringFromMap(sellerMap, "name")
		if err != nil {
			return nil, fmt.Errorf("load seller name failed, %w", err)
		}

		result = append(result, SellerStruct{ID: id, Name: name, Deals: deals})
	}

	return result, nil
}

func loadTestData(jsonFileName string) (data map[string]any, err error) {
	data = make(map[string]any, 16)
	bytes, err := os.ReadFile(jsonFileName)
	if err == nil {
		err = json.Unmarshal(bytes, &data)
	}
	return data, err
}
