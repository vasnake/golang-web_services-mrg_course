package main

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dataAdapter *StorageGQLAdapter
}

func (r *Resolver) New() *Resolver {
	return &Resolver{
		dataAdapter: (&StorageGQLAdapter{}).New(),
	}
}
