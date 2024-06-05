package main

import "fmt"

type ShopStorage struct {
	sellersRows []SellerStruct
	itemsRows   []GoodiesItemStruct
	catalogRows []CatalogStruct
}

func (ss *ShopStorage) FindItemsBySellerID(sid int) (itemIDList []int) {
	result := make([]int, 0, 16)
	for _, ir := range ss.itemsRows {
		if ir.Seller == sid {
			result = append(result, ir.ID)
		}
	}
	return result
}

func (ss *ShopStorage) GetSellerByID(sid int) (*SellerStruct, error) {
	for _, sr := range ss.sellersRows {
		if sr.ID == sid {
			return &sr, nil
		}
	}
	return nil, fmt.Errorf("ShopStorage.GetSellerByID failed, can't find sellerid %d in seller rows", sid)
}

func (ss *ShopStorage) GetCatalogByID(cid int) (*CatalogStruct, error) {
	for _, cr := range ss.catalogRows {
		if cr.ID == cid {
			return &cr, nil
		}
	}
	return nil, fmt.Errorf("ShopStorage.GetCatalogByID failed, can't find catalogid %d in catalog rows", cid)
}

func (ss *ShopStorage) GetItemByID(iid int) (*GoodiesItemStruct, error) {
	for _, ir := range ss.itemsRows {
		if ir.ID == iid {
			return &ir, nil
		}
	}
	return nil, fmt.Errorf("ShopStorage.GetItemByID failed, can't find itemid %d in item rows", iid)
}
