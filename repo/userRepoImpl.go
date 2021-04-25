package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepoImpl struct {
	users []domain.User
	user domain.User
	userDto domain.UserDto
	userDtoList []domain.UserDto
}

var dbConnection = database.GetInstance()

func (u UserRepoImpl) FindAll() (*[]domain.UserDto, error) {
	// Get all users
	cur, err := dbConnection.Collection("users").Find(context.TODO(), bson.M{"profileIsViewable": true})

	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.UserDto
		err := cur.Decode(&elem)

		if err != nil {
			return nil, err
		}

		u.userDtoList = append(u.userDtoList, elem)
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

	return &u.userDtoList, nil
}

func (u UserRepoImpl) Create(user *domain.User) error {
	cur, err := dbConnection.Collection("users").Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},

	})

	if err != nil {
		return err
	}

	if !cur.Next(context.TODO()) {
		user.Id = primitive.NewObjectID()
		_, err = dbConnection.Collection("users").InsertOne(context.TODO(), &user)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) FindByID(id primitive.ObjectID) (*domain.UserDto, error) {
	err := dbConnection.Collection("users").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) FindByUsername(username string) (*domain.UserDto, error) {
	err := dbConnection.Collection("users").FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateByID(id primitive.ObjectID, user *domain.User) (*domain.UserDto, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"tokenHash", user.TokenHash}, {"tokenExpiresAt", user.TokenExpiresAt}}}}

	err := database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return nil, err
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileIsViewable", user.ProfileIsViewable}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"acceptMessages", user.AcceptMessages}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentBadgeUrl", user.CurrentBadgeUrl}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profilePictureUrl", user.ProfilePictureUrl}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileBackgroundPictureUrl", user.ProfileBackgroundPictureUrl}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentTagLine", user.CurrentTagLine}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateVerification(id primitive.ObjectID, user *domain.UpdateVerification) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"isVerified", user.IsVerified}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdatePassword(id primitive.ObjectID, password string) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"password", password}, {"tokenHash", ""}, {"tokenExpiresAt", 0}}}}

	database.GetInstance().Collection("users").FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateFlagCount(flag *domain.Flag) error {

	cur, err := dbConnection.Collection("flags").Find(context.TODO(), bson.M{
		"$and": []interface{}{
			bson.M{"flaggerID": flag.FlaggerID},
			bson.M{"flaggedUsername": flag.FlaggedUsername},
		},

	})

	if err != nil {
		return err
	}

	if !cur.Next(context.TODO()) {
		flag.Id = primitive.NewObjectID()
		_, err = dbConnection.Collection("flags").InsertOne(context.TODO(), &flag)

		if err != nil {
			return err
		}

		filter := bson.D{{"username", flag.FlaggedUsername}}
		update := bson.M{"$push": bson.M{"flagCount":flag}}

		_, err = database.GetInstance().Collection("users").UpdateOne(context.TODO(),
			filter, update)
		if err != nil {
			return err
		}

		return nil
	}

	return  fmt.Errorf("you've already flagged this user")
}

func (u UserRepoImpl) BlockUser(id primitive.ObjectID, username string) error {

	err := dbConnection.Collection("users").FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.user)

	if id == u.user.Id {
		return fmt.Errorf("you can't block yourself")
	}

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	for _, foundId := range u.user.BlockByList {
		if foundId == id {
			return fmt.Errorf("already blocked")
		}
	}

	filter := bson.D{{"_id", id}}
	update := bson.M{"$push": bson.M{"blockList": u.user.Id}}

	_, err = database.GetInstance().Collection("users").UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	filter = bson.D{{"_id", u.user.Id}}
	update = bson.M{"$push": bson.M{"blockByList": id}}

	_, err = database.GetInstance().Collection("users").UpdateOne(context.TODO(),
		filter, update)

	if err != nil {
		return err
	}

	return  nil
}


func (u UserRepoImpl) DeleteByID(id primitive.ObjectID) error {
	_, err := database.GetInstance().Collection("users").DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}