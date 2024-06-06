package main

import (
	"context"
	"fmt"
	// "strconv"
)

// storage to gql converter
type StorageGQLAdapter struct {
	shopStorage      *ShopStorage
	shoppingCartsSvc *ShoppingCartService
}

func (sa *StorageGQLAdapter) New() *StorageGQLAdapter {
	return &StorageGQLAdapter{
		shopStorage:      (&ShopStorage{}).New(),
		shoppingCartsSvc: (&ShoppingCartService{}).New(),
	}
}

func (sa *StorageGQLAdapter) GetShoppingCartItems(ctx context.Context) ([]*CartItem, error) {
	sess, err := SessionFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("StorageGQLAdapter.GetShoppingCartItems failed, can't find user session: %w", err)
	}

	cart := sa.GetShoppingCartByUserID(sess.GetUserID())
	result := make([]*CartItem, 0, len(cart.items))

	for idx, item := range cart.items {
		result = append(
			result,
			itemRow2CartItem(item, cart.quantities[idx]),
		)
	}

	return result, nil
}

// GetShoppingCartByUserID: get existing or create new
func (sa *StorageGQLAdapter) GetShoppingCartByUserID(userID string) *ShoppingCart {
	return sa.shoppingCartsSvc.GetShoppingCartByUserID(userID)
}

func (sa *StorageGQLAdapter) AddToShoppingCart(ctx context.Context, itemID int, quantity int) error {
	sess, err := SessionFromContext(ctx)
	if err != nil {
		return fmt.Errorf("StorageGQLAdapter.AddToShoppingCart failed, can't find user session: %w", err)
	}
	show("got session: ", sess)

	itemRow, err := sa.shopStorage.GetItemByID(itemID)
	if err != nil {
		return fmt.Errorf("StorageGQLAdapter.AddToShoppingCart failed, can't find item: %w", err)
	}
	show("got shop item: ", itemRow)

	cart := sa.GetShoppingCartByUserID(sess.GetUserID())
	show("got shopping cart: ", cart)

	cart.AddItem(itemRow, quantity)
	return nil
}

func (sa *StorageGQLAdapter) GetSellerByID(sid int) (*Seller, error) {
	sellerRow, err := sa.shopStorage.GetSellerByID(sid)
	if err != nil {
		return nil, fmt.Errorf("StorageGQLAdapter.GetSellerByID failed, can't find seller row: %w", err)
	}

	items := sa.shopStorage.FindItemsBySellerID(sid)
	return sellerRow2Seller(sellerRow, items), nil
}

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

func catalogRow2Catalog(cr *CatalogStruct) *Catalog {
	return &Catalog{
		ID:             cr.ID, // strconv.Itoa(cr.ID),
		Name:           cr.Name,
		ParentID:       cr.Parent,
		ChildrenIDList: cr.Children,
		ItemsIDList:    cr.Items,
	}
}

func itemRow2CartItem(itm *GoodiesItemStruct, quantity int) *CartItem {
	return &CartItem{
		Quantity: quantity,
		Item:     itemRow2Item(itm),
	}
}

func itemRow2Item(itm *GoodiesItemStruct) *Item {
	return &Item{
		ID:          itm.ID,                      // int `json:"id"`
		Name:        itm.Name,                    // string `json:"name,omitempty"`
		CatalogID:   itm.Catalog,                 // int `json:"-"`
		SellerID:    itm.Seller,                  // int `json:"-"`
		InStockText: inStockMapping(itm.InStock), // string `json:"inStockText"`
		// InCart: int `json:"inCart"`
	}
}

func sellerRow2Seller(ss *SellerStruct, items []int) *Seller {
	return &Seller{
		ID:          ss.ID,    // int `json:"id"`
		Name:        ss.Name,  // string `json:"name,omitempty"`
		Deals:       ss.Deals, // int `json:"deals"`
		ItemsIDList: items,    // []int `json:"-"`
	}
}

func inStockMapping(itemsCount int) string {
	switch {
	case itemsCount <= 1:
		return "мало"
	case itemsCount > 3:
		return "много"
	default:
		return "хватает"
	}
}
