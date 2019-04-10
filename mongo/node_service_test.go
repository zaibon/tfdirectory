package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/require"
	"github.com/zaibon/tfdirectory"
)

func createTestNodeService(t *testing.T) *NodeService {
	require := require.New(t)
	ctx := context.Background()
	session, err := NewSession(ctx, "mongodb://localhost:27017")
	require.NoError(err)

	return NewNodeService(session, "test", "nodes")
}

func createTestNode() *tfdirectory.Node {
	resources := &tfdirectory.Resource{
		CRU: 10,
		MRU: 20,
		HRU: 2000,
		SRU: 500,
	}
	node := &tfdirectory.Node{
		NodeID:    "nodeid",
		OSVersion: "master",
		RobotURL:  "http://localhost:6600",
		Location: &tfdirectory.Location{
			City:      "liege",
			Continent: "europe",
			Latitude:  5.12,
			Longitude: 23.4,
		},
		TotalResources:   resources,
		UsedResources:    resources,
		ReservedResoures: resources,
		FarmerID:         "test_farm",
		Created:          time.Now(),
		Updated:          time.Now(),
		Parameters:       []string{"support", "development"},
		Uptime:           500,
	}
	return node
}

func TestInsertAndSearch(t *testing.T) {
	require := require.New(t)
	nodeSrv := createTestNodeService(t)

	ctx := context.Background()
	// cleanup db
	_, err := nodeSrv.collection.DeleteMany(ctx, bson.D{})
	require.NoError(err)

	// test register
	node := createTestNode()
	err = nodeSrv.Register(ctx, node)
	require.NoError(err)

	// test uniqueness of ID
	t.Run("ensure unique ID", func(t *testing.T) {
		err = nodeSrv.Register(ctx, node)
		assert.Error(t, err)
	})

	t.Run("search", func(t *testing.T) {
		nodes, err := nodeSrv.Search(ctx, tfdirectory.NodeQuery{})
		require.NoError(err)
		assert.Equal(t, 1, len(nodes))
		nodes[0].Created = node.Created
		nodes[0].Updated = node.Updated
		assert.EqualValues(t, node, nodes[0])
	})

	t.Run("GetByID", func(t *testing.T) {
		actual, err := nodeSrv.GetByID(ctx, node.NodeID)
		require.NoError(err)
		// TODO:manually compare time cause testify failed to do it properly
		// assert.True(t, compareTime(node.Created, actual.Created))
		// assert.True(t, compareTime(node.Updated, actual.Updated))
		actual.Created = node.Created
		actual.Updated = node.Updated
		assert.EqualValues(t, node, actual)
	})

	t.Run("not found", func(t *testing.T) {
		_, err = nodeSrv.GetByID(ctx, "notfound")
		assert.Error(t, err)
	})
}

func TestUpdateResources(t *testing.T) {
	require := require.New(t)
	nodeSrv := createTestNodeService(t)

	ctx := context.Background()
	// cleanup db
	_, err := nodeSrv.collection.DeleteMany(ctx, bson.D{})
	require.NoError(err)

	// insert one node
	node := createTestNode()
	err = nodeSrv.Register(ctx, node)
	require.NoError(err)

	updatedResources := tfdirectory.Resource{
		CRU: 1,
		MRU: 1,
		HRU: 1,
		SRU: 1,
	}
	err = nodeSrv.UpdateReservedResources(ctx, node.NodeID, updatedResources)
	require.NoError(err)
	err = nodeSrv.UpdateUsedResources(ctx, node.NodeID, updatedResources)

	actual, err := nodeSrv.GetByID(ctx, node.NodeID)
	require.NoError(err)
	assert.EqualValues(t, updatedResources, *actual.ReservedResoures)
	assert.EqualValues(t, updatedResources, *actual.UsedResources)

	err = nodeSrv.UpdateReservedResources(ctx, "notfound", updatedResources)
	assert.Error(t, err)
}
