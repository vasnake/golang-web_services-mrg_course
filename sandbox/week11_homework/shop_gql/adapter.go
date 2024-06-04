package main

import (
	"fmt"
	// "strconv"
)

// storage to gql converter
type StorageGQLAdapter struct {
	shopStorage ShopStorage
}

func catalogRow2Catalog(cr *CatalogStruct) *Catalog {
	return &Catalog{
		ID:             cr.ID, // strconv.Itoa(cr.ID),
		Name:           cr.Name,
		ParentID:       cr.Parent,
		ChildrenIDList: cr.Children,
		ItemsIDList:    cr.Items,
	}
}

func itemRow2Item(ir *GoodiesItemStruct) *Item {
	return &Item{
		ID:          ir.ID,                      // int `json:"id"`
		Name:        ir.Name,                    // string `json:"name,omitempty"`
		CatalogID:   ir.Catalog,                 // int `json:"-"`
		SellerID:    ir.Seller,                  // int `json:"-"`
		InStockText: inStockMapping(ir.InStock), // string `json:"inStockText"`
		// InCart: int `json:"inCart"`
	}
}

func inStockMapping(cnt int) string {
	switch {
	case cnt <= 1:
		return "мало"
	case cnt > 3:
		return "много"
	default:
		return "хватает"
	}
}

// item, err := r.dataAdapter.GetItemByID(iid)
func (sa *StorageGQLAdapter) GetItemByID(iid int) (*Item, error) {
	itemRow, err := sa.shopStorage.GetItemByID(iid)
	if err != nil {
		return nil, fmt.Errorf("StorageGQLAdapter.GetItemByID failed, can't find item row: %w", err)
	}
	return itemRow2Item(itemRow), nil
}

func (sa *StorageGQLAdapter) GetCatalogByID(cid int) (*Catalog, error) {
	catalogRow, err := sa.shopStorage.GetCatalogByID(cid)
	if err != nil {
		return nil, fmt.Errorf("StorageGQLAdapter.GetCatalogByID failed, can't find catalog row: %w", err)
	}

	return catalogRow2Catalog(catalogRow), nil
}

func (sa *StorageGQLAdapter) GetCatalogChildrenByParentID(cid int) ([]*Catalog, error) {
	catalogRow, err := sa.shopStorage.GetCatalogByID(cid)
	if err != nil {
		return nil, fmt.Errorf("StorageGQLAdapter.GetCatalogChildrenByParentID failed, can't find catalog row: %w", err)
	}

	result := make([]*Catalog, 0, len(catalogRow.Children))

	for _, childID := range catalogRow.Children {
		childRow, err := sa.shopStorage.GetCatalogByID(childID)
		panicOnError("StorageGQLAdapter.GetCatalogChildrenByParentID failed to get child by id", err)
		result = append(result, catalogRow2Catalog(childRow))
	}

	return result, nil
}
