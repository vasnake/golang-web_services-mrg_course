package main

// Сущность каталога
type Catalog struct {
	// ID каталога
	// ID string `json:"id,omitempty"`
	ID int `json:"id"`
	// Название раздела каталога
	Name string `json:"name,omitempty"`
	// Родительский раздел
	ParentID int `json:"-"`
	// Дочерние разделы
	ChildrenIDList []int `json:"-"`
	// Товары в разделе
	ItemsIDList []int `json:"-"`
}

// Сущность товара
type Item struct {
	// ID товара
	// ID string `json:"id,omitempty"`
	ID int `json:"id"`
	// Имя товара
	Name string `json:"name,omitempty"`
	// Родительский каталог
	CatalogID int `json:"-"`
	// Сущность продавец
	SellerID int `json:"-"`
	// Количесто товара в корзине у пользователя
	InCart int `json:"inCart"`
	// Текст, сколько осталось на складе (<=1 - мало; >=2 && <=3 - хватает; >3 - много)
	InStockText string `json:"inStockText"`
}

// Сущность продавца
type Seller struct {
	// ID продавца
	// ID *string `json:"id,omitempty"`
	ID int `json:"id"`
	// Имя продавца
	Name string `json:"name,omitempty"`
	// Количество сделок у продавца (берется из testdata.json)
	Deals int `json:"deals"`
	// Товары продавца
	ItemsIDList []int `json:"-"`
}

/*

// CartInput - параметры метода добавления товара в корзину
type CartInput struct {
	// ID товара, который надо добавить в корзину
	ItemID int `json:"itemID"`
	// Количество товара, которое нужно добавить в корзину
	Quantity int `json:"quantity"`
}

// CartItem - сущность элемента корзины
type CartItem struct {
	// Количество товаров данного типа в корзине
	Quantity int `json:"quantity"`
	// ID товара, который надо добавить в корзину
	Item *Item `json:"item"`
}
*/
