package db

import (
	"context"
	"out/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var timeOnline int64 = 5

var projectStage = bson.D{
	{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}}}}

func (db *MongoDB) AddStream(ctx context.Context, coll string, data any) error {
	c := db.DB.Collection(coll)
	_, err := c.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) GetAllUniqueUsers(ctx context.Context, coll string) (*models.ViewerCount, error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: time.Now().Unix() - timeOnline*60},
			}},
		}},
	}

	totalGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.A{"clientid", "$clientid"}},
		}}}

	uniqueGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "viewerCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}}

	c := db.DB.Collection(coll)
	cursor, err := c.Aggregate(context.TODO(), mongo.Pipeline{matchStage, totalGroupStage, uniqueGroupStage, projectStage})
	if err != nil {
		return nil, err
	}

	var results []models.ViewerCount
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.ViewerCount{}, nil
	}
	return &results[0], nil
}

func (db *MongoDB) GetStreamsForUser(ctx context.Context, coll string, cID string) (*models.StreamCount, error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "clientid", Value: cID},
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: time.Now().Unix() - timeOnline*60},
			}},
		}},
	}

	totalGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "content", Value: "$content"},
				{Key: "clientid", Value: "$clientid"},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}},
	}

	uniqueGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "clientid", Value: "$_id.clientid"},
			}},
			{Key: "streamCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}},
	}

	c := db.DB.Collection(coll)
	cursor, err := c.Aggregate(context.TODO(), mongo.Pipeline{matchStage, totalGroupStage, uniqueGroupStage, projectStage})
	if err != nil {
		return nil, err
	}

	var results []models.StreamCount
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.StreamCount{}, nil
	}
	return &results[0], nil
}

func (db *MongoDB) GetUniqueClientPerContent(ctx context.Context, coll string, content string) (*models.ViewerCount, error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "content", Value: content},
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: time.Now().Unix() - timeOnline*60},
			}},
		}},
	}

	totalGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "content", Value: "$content"},
				{Key: "clientid", Value: "$clientid"},
			}},
		}},
	}

	uniqueGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "content", Value: "$_id.content"},
			}},
			{Key: "viewerCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}},
	}

	c := db.DB.Collection(coll)
	cursor, err := c.Aggregate(context.TODO(), mongo.Pipeline{matchStage, totalGroupStage, uniqueGroupStage, projectStage})
	if err != nil {
		return nil, err
	}

	var results []models.ViewerCount
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.ViewerCount{}, nil
	}
	return &results[0], nil
}

func (db *MongoDB) GetUniqueClientsPerState(ctx context.Context, coll string) ([]bson.M, error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: time.Now().Unix() - timeOnline*60},
			}},
		}},
	}

	totalGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "geo", Value: "$geo"},
				{Key: "clientid", Value: "$clientid"},
			}},
		}},
	}

	uniqueGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "geo", Value: "$_id.geo"},
			}},
			{Key: "viewerCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}},
	}

	c := db.DB.Collection(coll)
	cursor, err := c.Aggregate(context.TODO(), mongo.Pipeline{matchStage, totalGroupStage, uniqueGroupStage})
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *MongoDB) GetUniqueClientsPerDevice(ctx context.Context, coll string) ([]bson.M, error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: time.Now().Unix() - timeOnline*60},
			}},
		}},
	}

	totalGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "device", Value: "$device"},
				{Key: "clientid", Value: "$clientid"},
			}},
		}},
	}

	uniqueGroupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "device", Value: "$_id.device"},
			}},
			{Key: "viewerCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}},
	}

	c := db.DB.Collection(coll)
	cursor, err := c.Aggregate(context.TODO(), mongo.Pipeline{matchStage, totalGroupStage, uniqueGroupStage})
	if err != nil {
		return nil, err
	}

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
