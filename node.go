package tfdirectory

import (
	"context"
	"time"
)

type Resource struct {
	CRU int64 `json:"cru" bson:"cru"`
	MRU int64 `json:"mru" bson:"mru"`
	SRU int64 `json:"sru" bson:"sru"`
	HRU int64 `json:"hru" bson:"hru"`
}

type Location struct {
	Continent string  `json:"continent" bson:"continent"`
	Country   string  `json:"contry" bson:"contry"`
	City      string  `json:"city" bson:"city"`
	Latitude  float32 `json:"latitude" bson:"latitude"`
	Longitude float32 `json:"longitude" bson:"longitude"`
}

type Node struct {
	NodeID           string    `json:"node_id" bson:"node_id"`
	OSVersion        string    `json:"os_version" bson:"os_version"`
	RobotURL         string    `json:"robot_address" bson:"robot_address"`
	TotalResources   *Resource `json:"total_resources" bson:"total_resources"`
	ReservedResoures *Resource `json:"reserved_resources" bson:"reserved_resources"`
	UsedResources    *Resource `json:"used_resources" bson:"used_resources"`
	Location         *Location `json:"location" bson:"location"`
	Uptime           int64     `json:"uptime" bson:"uptime"`

	FarmerID string `json:"farmer_id" bson:"farmer_id"`

	Parameters []string `json:"parameters" bson:"parameters"`

	Updated time.Time `json:"updated" bson:"updated"`
	Created time.Time `json:"created" bson:"created"`
}

type NodeQuery struct {
	location Location
	resource Resource
}

type NodeService interface {
	Register(ctx context.Context, node *Node) error
	GetByID(ctx context.Context, id string) (*Node, error)
	Search(ctx context.Context, query NodeQuery) ([]*Node, error)
}
