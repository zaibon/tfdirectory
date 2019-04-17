package tfdirectory

import (
	"context"
)

type Farmer struct {
	Email        string   `json:"email" bson:"email"`
	Organization string   `json:"iyo_organization" bson:"_id"`
	Name         string   `json:"name" bson:"name"`
	WalletAddrs  []string `json:"wallet_addresses" bson:"wallet_addresses"`
}

type FarmerQuery struct {
	Organization string
}

type FarmerService interface {
	Insert(ctx context.Context, farmer *Farmer) error
	Update(ctx context.Context, farmer *Farmer) error
	GetByID(ctx context.Context, ID string) (*Farmer, error)
	List(ctx context.Context) ([]*Farmer, error)
}
