package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoImpl struct {
	users []domain.User
	user domain.User
}

var dbConnection = database.GetInstance()

func (u UserRepoImpl) Create(user *domain.User) error {
	user.Id = primitive.NewObjectID()
	_, err := dbConnection.Collection.InsertOne(context.TODO(), &user)

	if err != nil {
		return err
	}

	return nil
}

func (u UserRepoImpl) FindAll() (*[]domain.User, error) {
	// Get all users
	cur, err := dbConnection.Collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.User
		err := cur.Decode(&elem)

		if err != nil {
			return nil, err
		}

		u.users = append(u.users, elem)
	}

	err = cur.Err()
	if err != nil {
		return nil, err
	}

	// Close the cursor once finished
	err = cur.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	return &u.users, nil
}

func (u UserRepoImpl) FindByID(id primitive.ObjectID) (*domain.User, error) {
	err := dbConnection.Collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&u.user)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &u.user, nil
}

func (u UserRepoImpl) UpdateByID(id primitive.ObjectID, user *domain.User) (*domain.User, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"Email", user.Email}}}}

	err := database.GetInstance().Collection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.users)

	if err != nil {
		return nil, err
	}

	return &u.user, nil
}

func (u UserRepoImpl) DeleteByID(id primitive.ObjectID) error {
	_, err := database.GetInstance().Collection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}