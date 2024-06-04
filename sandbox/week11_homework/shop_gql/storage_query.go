package main

import "fmt"

type ShopStorage struct {
	sellersRows []SellerStruct
	itemsRows   []GoodiesItemStruct
	catalogRows []CatalogStruct
}

func (ss *ShopStorage) GetCatalogByID(cid int) (*CatalogStruct, error) {
	for _, c := range ss.catalogRows {
		if c.ID == cid {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("ShopStorage.GetCatalogByID failed, can't find catalogid %d in catalog rows", cid)
}

func (ss *ShopStorage) GetItemByID(iid int) (*GoodiesItemStruct, error) {
	for _, item := range ss.itemsRows {
		if item.ID == iid {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("ShopStorage.GetItemByID failed, can't find itemid %d in item rows", iid)
}
