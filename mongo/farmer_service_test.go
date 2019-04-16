package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/require"
	"github.com/zaibon/tfdirectory"
)

func createTestFarmerService(t *testing.T) *FarmerService {
	require := require.New(t)
	ctx := context.Background()
	session, err := NewSession(ctx, "mongodb://localhost:27017")
	require.NoError(err)

	return NewFarmerService(session, "test", "nodes")
}

func createTestFarmer() *tfdirectory.Farmer {
	return &tfdirectory.Farmer{
		Organization: "organisation",
		Name:         "Name",
		Email:        "mail@user.com",
		WalletAddrs:  []string{"addr1", "addr2"},
	}
}

func TestInserAndList(t *testing.T) {
	require := require.New(t)
	farmerSrv := createTestFarmerService(t)

	ctx := context.Background()
	// cleanup db
	_, err := farmerSrv.collection.DeleteMany(ctx, bson.D{})
	require.NoError(err)

	// test Insert
	farmer := createTestFarmer()
	err = farmerSrv.Insert(ctx, farmer)
	require.NoError(err)

	// test uniqueness of ID
	t.Run("ensure unique ID", func(t *testing.T) {
		err = farmerSrv.Insert(ctx, farmer)
		assert.Error(t, err)
	})

	// add another famrer
	farmer2 := createTestFarmer()
	farmer2.Organization = "org2"
	err = farmerSrv.Insert(ctx, farmer2)
	require.NoError(err)

	t.Run("list", func(t *testing.T) {
		famers, err := farmerSrv.List(ctx)
		require.NoError(err)
		assert.Equal(t, 2, len(famers))
		assert.EqualValues(t, farmer, famers[0])
	})

	t.Run("GetByID", func(t *testing.T) {
		actual, err := farmerSrv.GetByID(ctx, farmer.Organization)
		require.NoError(err)
		assert.EqualValues(t, farmer, actual)
	})

	t.Run("not found", func(t *testing.T) {
		_, err = farmerSrv.GetByID(ctx, "notfound")
		assert.Error(t, err)
	})
}
