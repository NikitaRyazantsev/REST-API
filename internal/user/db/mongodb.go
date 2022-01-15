package db

import (
	"context"
	"errors"
	"fmt"
	"project/internal/user"
	"project/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create database structure
type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

// NewStorage - Initialize new storage
func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

// Create - create new user in database
func (d *db) Create(ctx context.Context, user user.User) (string, error) {

	// Push user to collection
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	// Get ID of new user from database
	d.logger.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)

	// Check result
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

// GetUserFriends - get all friends from one user
func (d *db) GetUserFriends(ctx context.Context, id string) ([]string, error) {
	var u user.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u.Friends, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u.Friends, err
		}
		return u.Friends, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u.Friends, fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", id, err)
	}

	return u.Friends, nil
}

func (d *db) UpdateAge(ctx context.Context, id string, age string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", id)
	}
	filter := bson.M{"_id": objectID}
	updateAge := bson.D{
		{"$set", bson.D{{"age", age}}},
	}

	result, err := d.collection.UpdateOne(ctx, filter, updateAge)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return err
	}

	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", id)
	}

	filter := bson.M{"_id": objectID}

	var u user.User
	res := d.collection.FindOne(ctx, filter)
	if err = res.Decode(&u); err != nil {
		return fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", id, err)
	}

	updateFilter := bson.M{"friends": u.Username}
	updateResult, err := d.collection.UpdateMany(ctx, updateFilter, bson.D{
		{"$pull", bson.D{{"friends", u.Username}}},
	})
	if err != nil {
		return fmt.Errorf("failed delete from other users friends. error: %v", err)
	}
	d.logger.Tracef("Modified %d documents", updateResult.ModifiedCount)

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		return err
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func (d *db) MakeFriends(ctx context.Context, firstUserID string, secondUserID string) (firstUser user.User, secondUser user.User, err error) {
	firstObjectID, err := primitive.ObjectIDFromHex(firstUserID)
	if err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", firstUserID)
	}
	secondObjectID, err := primitive.ObjectIDFromHex(secondUserID)
	if err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", secondUserID)
	}

	filter := bson.M{"_id": firstObjectID}
	res := d.collection.FindOne(ctx, filter)
	if err = res.Decode(&firstUser); err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", firstUserID, err)
	}

	filter = bson.M{"_id": secondObjectID}
	res = d.collection.FindOne(ctx, filter)
	if err = res.Decode(&secondUser); err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", secondUserID, err)
	}

	updateFilter := bson.M{"_id": firstObjectID}
	updateResult, err := d.collection.UpdateOne(ctx, updateFilter, bson.D{
		{"$push", bson.D{{"friends", secondUser.Username}}},
	})
	if err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to add friend to friends. error: %v", err)
	}
	d.logger.Tracef("Modified %d documents", updateResult.ModifiedCount)

	updateFilter = bson.M{"_id": secondObjectID}
	updateResult, err = d.collection.UpdateOne(ctx, updateFilter, bson.D{
		{"$push", bson.D{{"friends", firstUser.Username}}},
	})
	if err != nil {
		return firstUser, secondUser, fmt.Errorf("failed to add friend to friends. error: %v", err)
	}
	d.logger.Tracef("Modified %d documents", updateResult.ModifiedCount)

	return firstUser, secondUser, nil
}
