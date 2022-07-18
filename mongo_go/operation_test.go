package mongo_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetCollection(t *testing.T) {
	collection, err := getCollection()
	assert.Nil(t, err)
	assert.NotNil(t, collection)
}

func Test_OperationForMongoDB(t *testing.T) {
	err := InsertOne2Mongo()
	assert.Nil(t, err)
	err = InsertMany2Mongo()
	assert.Nil(t, err)

	err = FindFromMongo()
	assert.Nil(t, err)

	err = DeleteFromMongo()
	assert.Nil(t, err)
}
