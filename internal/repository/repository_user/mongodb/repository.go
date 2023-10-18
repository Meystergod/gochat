package repository_user

import (
	"context"
	"time"

	"github.com/Meystergod/gochat/internal/domain"
	"github.com/Meystergod/gochat/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(storage *mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		collection: storage.Collection(collection),
	}
}

func (userRepository *UserRepository) CreateUser(ctx context.Context, domainUser *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	repositoryUser, err := userToRepository(domainUser)
	if err != nil {
		return utils.EmptyString, errors.Wrap(err, utils.ErrorConvertDomainToRepository.Error())
	}

	result, err := userRepository.collection.InsertOne(ctx, repositoryUser)
	if err != nil {
		return utils.EmptyString, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return utils.EmptyString, errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	return oid.Hex(), nil
}

func (userRepository *UserRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var repositoryUser *User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result := userRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err := result.Decode(&repositoryUser); err != nil {
		return nil, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	domainUser := userToDomain(repositoryUser)

	return &domainUser, nil
}

func (userRepository *UserRepository) GetAllUsers(ctx context.Context) (*[]domain.User, error) {
	var repositoryUsers []User

	filter := bson.M{}

	cursor, err := userRepository.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if err = cursor.All(ctx, &repositoryUsers); err != nil {
		return nil, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	var domainUsers []domain.User

	for _, u := range repositoryUsers {
		domainUsers = append(domainUsers, userToDomain(&u))
	}

	return &domainUsers, nil
}

func (userRepository *UserRepository) UpdateUser(ctx context.Context, domainUser *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	repositoryUser, err := userToRepository(domainUser)
	if err != nil {
		return errors.Wrap(err, utils.ErrorConvertDomainToRepository.Error())
	}

	userByte, err := bson.Marshal(&repositoryUser)
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

	filter := bson.M{"_id": repositoryUser.ID}

	result, err := userRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}

func (userRepository *UserRepository) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, utils.ErrorConvert.Error())
	}

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
