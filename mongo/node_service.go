package mongo

import (
	"context"
	"fmt"

	"github.com/zaibon/tfdirectory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NodeService struct {
	collection *mongo.Collection
}

var _ tfdirectory.NodeService = (*NodeService)(nil)

func nodeModelIndex() mongo.IndexModel {
	indexOpts := options.Index()
	indexOpts.SetUnique(true)
	indexOpts.SetBackground(true)
	indexOpts.SetSparse(true)
	return mongo.IndexModel{
		Keys: bson.M{
			"node_id": int32(1),
		},
		Options: indexOpts,
	}
}

// NewNodeService creates a NodeService that connects to the mongodb instance
// represented by Session, using the data dbName and the collection collectionName
func NewNodeService(session *Session, dbName string, collectionName string) *NodeService {
	collection := session.GetCollection(dbName, collectionName)
	_, err := collection.Indexes().CreateOne(context.Background(), nodeModelIndex())
	if err != nil {
		panic(fmt.Sprintf("fail to create index for node model: %+v", err))
	}
	return &NodeService{collection}
}

// Register inserts a new node into mongodb
func (n *NodeService) Register(ctx context.Context, node *tfdirectory.Node) error {
	_, err := n.collection.InsertOne(ctx, node)
	return err
}

// GetByID retrieve a node by its ID
func (n *NodeService) GetByID(ctx context.Context, ID string) (*tfdirectory.Node, error) {
	var node tfdirectory.Node
	result := n.collection.FindOne(ctx, bson.M{
		"_id": ID,
	})
	if err := result.Err(); err != nil {
		return nil, err
	}
	err := result.Decode(&node)
	return &node, err
}

//Search lists some nodes based on the query arguments
func (n *NodeService) Search(ctx context.Context, query tfdirectory.NodeQuery) ([]*tfdirectory.Node, error) {
	nodes := make([]*tfdirectory.Node, 0, 10)

	// TODO: convert query to valid mongodb filter
	q := bson.D{}
	if query.Farmer != "" {
		q = append(q, bson.E{"farmer", query.Farmer})
	}
	if query.Location.Country != "" {
		q = append(q, bson.E{"location.country", query.Location.Country})
	}
	if query.Resource.CRU != 0 {
		q = append(q, bson.E{"total_resources.cru", bson.D{{"$gte", query.Resource.CRU}}})
	}
	if query.Resource.MRU != 0 {
		q = append(q, bson.E{"total_resources.mru", bson.D{{"$gte", query.Resource.MRU}}})
	}
	if query.Resource.HRU != 0 {
		q = append(q, bson.E{"total_resources.hru", bson.D{{"$gte", query.Resource.HRU}}})
	}
	if query.Resource.SRU != 0 {
		q = append(q, bson.E{"total_resources.sru", bson.D{{"$gte", query.Resource.SRU}}})
	}
	cur, err := n.collection.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var node tfdirectory.Node
		if err := cur.Decode(&node); err != nil {
			return nil, err
		}
		if err := cur.Err(); err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (n *NodeService) UpdateUsedResources(ctx context.Context, ID string, resource tfdirectory.Resource) error {
	return n.updateResource(ctx, ID, usedResources, resource)
}
func (n *NodeService) UpdateReservedResources(ctx context.Context, ID string, resource tfdirectory.Resource) error {
	return n.updateResource(ctx, ID, reservedResources, resource)
}

func (n *NodeService) updateResource(ctx context.Context, ID string, resourceType ResourceType, resource tfdirectory.Resource) error {
	if resourceType != reservedResources && resourceType != usedResources {
		return fmt.Errorf("can only update reserved_resources or used_resources, not %s", resourceType)
	}

	filter := bson.D{{"_id", ID}}
	update := bson.D{
		{"$set", bson.D{
			{fmt.Sprintf("%s.cru", string(resourceType)), resource.CRU},
			{fmt.Sprintf("%s.mru", string(resourceType)), resource.MRU},
			{fmt.Sprintf("%s.sru", string(resourceType)), resource.SRU},
			{fmt.Sprintf("%s.hru", string(resourceType)), resource.HRU},
		}},
	}

	updateResult, err := n.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount < 1 {
		return fmt.Errorf("node not found")
	}
	if updateResult.ModifiedCount > 1 {
		panic(fmt.Sprintf("multiple node match on one ID (%s), should never happen", ID))
	}
	return nil
}

type ResourceType string

const (
	reservedResources ResourceType = "reserved_resources"
	usedResources     ResourceType = "used_resources"
)
