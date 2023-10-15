package mongodb

import (
	"context"
	"time"

	"github.com/Meystergod/gochat/internal/entity/model"
	"github.com/Meystergod/gochat/internal/repository"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(storage *mongo.Database, collection string) repository.UserRepository {
	return &userRepository{
		collection: storage.Collection(collection),
	}
}

func (userRepository *userRepository) CreateUser(ctx context.Context, user *model.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	var oid primitive.ObjectID

	result, err := userRepository.collection.InsertOne(ctx, user)
	if err != nil {
		return oid, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return oid, errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	return oid, nil
}

func (userRepository *userRepository) GetUser(ctx context.Context, oid primitive.ObjectID) (*model.User, error) {
	var user *model.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"_id": oid}

	result := userRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return user, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err := result.Decode(&user); err != nil {
		return user, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return user, nil
}

func (userRepository *userRepository) GetAllUsers(ctx context.Context) (*[]model.User, error) {
	var users []model.User

	filter := bson.M{}

	cursor, err := userRepository.collection.Find(ctx, filter)
	if err != nil {
		return &users, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if err = cursor.All(ctx, &users); err != nil {
		return &users, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return &users, nil
}

func (userRepository *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"_id": user.ID}

	userByte, err := bson.Marshal(user)
	if err != nil {
		return errors.Wrap(err, utils.ErrorMarshal.Error())
	}

	var object bson.M

	err = bson.Unmarshal(userByte, &object)
	if err != nil {
		return errors.Wrap(err, utils.ErrorUnmarshal.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := userRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}

func (userRepository *userRepository) DeleteUser(ctx context.Context, oid primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"_id": oid}

	result, err := userRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}
