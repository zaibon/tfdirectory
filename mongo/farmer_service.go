package mongo

import (
	"context"
	"fmt"

	"github.com/zaibon/tfdirectory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FarmerService struct {
	collection *mongo.Collection
}

var _ tfdirectory.FarmerService = (*FarmerService)(nil)

func farmerModelIndex() mongo.IndexModel {
	indexOpts := options.Index()
	indexOpts.SetUnique(true)
	indexOpts.SetBackground(true)
	indexOpts.SetSparse(true)
	return mongo.IndexModel{
		Keys: bson.M{
			"iyo_organization": int32(1),
		},
		Options: indexOpts,
	}
}

// NewFarmereService creates a FarmerService that connects to the mongodb instance
// represented by Session, using the data dbName and the collection collectionName
func NewFarmerService(session *Session, dbName string, collectionName string) *FarmerService {
	collection := session.GetCollection(dbName, collectionName)
	_, err := collection.Indexes().CreateOne(context.Background(), farmerModelIndex())
	if err != nil {
		panic(fmt.Sprintf("fail to create index for farmer model: %+v", err))
	}
	return &FarmerService{collection}
}

func (fs *FarmerService) Insert(ctx context.Context, farmer *tfdirectory.Farmer) error {
	_, err := fs.collection.InsertOne(ctx, farmer)
	return err
}

func (fs *FarmerService) Update(ctx context.Context, farmer *tfdirectory.Farmer) error {
	filter := bson.D{{"_id", farmer.Organization}}
	update := bson.D{
		{"$set", bson.D{
			{"email", farmer.Email},
			{"name", farmer.Name},
			{"wallet_addresses", farmer.WalletAddrs},
		}},
	}

	updateResult, err := fs.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return fmt.Errorf("farmer not found")
	}
	if updateResult.ModifiedCount > 1 {
		panic(fmt.Sprintf("multiple node match on farmer organization (%s), should never happen", farmer.Organization))
	}
	return nil
}

func (fs *FarmerService) GetByID(ctx context.Context, ID string) (*tfdirectory.Farmer, error) {
	var farmer tfdirectory.Farmer
	result := fs.collection.FindOne(ctx, bson.M{
		"_id": ID,
	})
	if err := result.Err(); err != nil {
		return nil, err
	}
	err := result.Decode(&farmer)
	return &farmer, err
}

func (fs *FarmerService) List(ctx context.Context) ([]*tfdirectory.Farmer, error) {
	farmers := make([]*tfdirectory.Farmer, 0, 10)

	cur, err := fs.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var farmer tfdirectory.Farmer
		if err := cur.Decode(&farmer); err != nil {
			return nil, err
		}
		if err := cur.Err(); err != nil {
			return nil, err
		}
		farmers = append(farmers, &farmer)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return farmers, nil
}
