package main

import "fmt"

type ShopStorage struct {
	sellersRows []SellerStruct
	itemsRows   []GoodiesItemStruct
	catalogRows []CatalogStruct
}

// catalogRow, err := sa.shopStorage.GetCatalogByID(cid)
func (ss *ShopStorage) GetCatalogByID(cid int) (*CatalogStruct, error) {
	for _, c := range ss.catalogRows {
		if c.ID == cid {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("ShopStorage.GetCatalogByID failed, can't find catalogid %d in catalog rows", cid)
}
