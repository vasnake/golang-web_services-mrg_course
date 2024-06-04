// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package main

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

type Mutation struct {
}

type Query struct {
}
