package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tahmooress/weConnect-task/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	collection *mongo.Collection
}

func New(ctx context.Context) (*MongoDB, error) {
	url := os.Getenv("MONGODB_IP")
	port := os.Getenv("MONGODB_PORT")

	conStr := fmt.Sprintf("mongodb://%s:%s", url, port)

	clientOptions := options.Client().ApplyURI(conStr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongodb: connection error: %s", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongodb ping erorr: %s", err)
	}

	fmt.Println("successfully connected to mongoDB")

	return &MongoDB{
		collection: client.Database("tasks").Collection("task"),
	}, nil
}

func (m *MongoDB) Insert(ctx context.Context, record *entity.Statistics) (string, error) {
	record.ID = primitive.NewObjectID()
	record.CreatedAt = time.Now()
	_, err := m.collection.InsertOne(ctx, record)
	if err != nil {
		return "", fmt.Errorf("mongodb insert erorr: %s", err)
	}

	return record.ID.Hex(), nil
}

func (m *MongoDB) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex("id")
	if err != nil {
		return fmt.Errorf("mongodb invalid id: %s", err)
	}
	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) GetByID(ctx context.Context, id string) (*entity.Statistics, error) {
	_id, err := primitive.ObjectIDFromHex("id")
	if err != nil {
		return nil, fmt.Errorf("mongodb invalid id: %s", err)
	}

	var result entity.Statistics
	err = m.collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("mongodb reading result: %s", err)
	}

	return &result, nil
}

func (m *MongoDB) GetAll(ctx context.Context, page, limit int64) ([]entity.Statistics, error) {
	pagination := newMongoPaginate(page, limit)
	opt := options.FindOptions{Limit: &pagination.limit, Skip: &pagination.skip}
	opt.SetSort(bson.M{"created_at": 1})

	cur, err := m.collection.Find(ctx, nil, &opt)
	if err != nil {
		return nil, fmt.Errorf("mongodb find: %s", err)
	}

	var results []entity.Statistics
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("mongodb find: %s", err)
	}

	return results, nil
}

func (m *MongoDB) Close() error {
	return m.collection.Database().Client().Disconnect(context.TODO())
}

type mongoPaginate struct {
	limit int64
	skip  int64
}

func newMongoPaginate(limit, page int64) *mongoPaginate {
	return &mongoPaginate{
		limit: limit,
		skip:  limit*page - limit,
	}
}
