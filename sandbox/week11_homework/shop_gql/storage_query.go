package main

import "fmt"

type ShopStorage struct {
	sellersRows []SellerStruct
	itemsRows   []GoodiesItemStruct
	catalogRows []CatalogStruct
}

func (ss *ShopStorage) New() *ShopStorage {
	return &ShopStorage{
		sellersRows: make([]SellerStruct, 0, 16),
		itemsRows:   make([]GoodiesItemStruct, 0, 16),
		catalogRows: make([]CatalogStruct, 0, 16),
	}
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

type ShoppingCart struct {
	userID     string
	items      []*GoodiesItemStruct
	quantities []int
}

func (sc *ShoppingCart) New() *ShoppingCart {
	return &ShoppingCart{
		userID:     "",
		items:      make([]*GoodiesItemStruct, 0, 16),
		quantities: make([]int, 0, 16),
	}
}

func (sc *ShoppingCart) findItemIndex(item *GoodiesItemStruct) int {
	for idx, row := range sc.items {
		if row.ID == item.ID {
			return idx
		}
	}
	return -1
}

func (sc *ShoppingCart) AddItem(item *GoodiesItemStruct, quantity int) {
	idx := sc.findItemIndex(item)
	if idx < 0 {
		// not in cart
		sc.items = append(sc.items, item)
		sc.quantities = append(sc.quantities, quantity)
	} else {
		// in cart
		sc.quantities[idx] = sc.quantities[idx] + quantity
	}
}

type ShoppingCartService struct {
	shoppingCartList []*ShoppingCart
}

func (scs *ShoppingCartService) New() *ShoppingCartService {
	return &ShoppingCartService{
		shoppingCartList: make([]*ShoppingCart, 0, 16),
	}
}

// GetShoppingCartByUserID: find existing or create new
func (scs *ShoppingCartService) GetShoppingCartByUserID(userID string) *ShoppingCart {
	for _, sc := range scs.shoppingCartList {
		if sc.userID == userID {
			return sc
		}
	}

	sc := (&ShoppingCart{}).New()
	sc.userID = userID

	scs.shoppingCartList = append(scs.shoppingCartList, sc)

	return sc
}
