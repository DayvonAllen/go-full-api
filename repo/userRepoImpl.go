package repo

import (
	"context"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/events"
	"example.com/app/util"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"strconv"
	"time"
)

type UserRepoImpl struct {
	users []domain.User
	user domain.User
	userDto domain.UserDto
	userDtoList []domain.UserDto
}

func (u UserRepoImpl) FindAll(id primitive.ObjectID, page string, ctx context.Context) (*[]domain.UserDto, error) {
	currentUser, err := u.FindByID(id)

	if err != nil {
		return nil, err
	}

	findOptions := options.FindOptions{}
	perPage := 20
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil,  fmt.Errorf("must input a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	// Get all users
	cur, err := database.GetInstance().UserCollection.Find(ctx, bson.M{
		"profileIsViewable": true,
		"$and": []interface{}{
			bson.M{"_id": bson.M{ "$ne": id }},
			bson.M{"_id": bson.M{"$nin": currentUser.BlockByList}},
			bson.M{"_id": bson.M{"$nin": currentUser.BlockList}},
		},
	}, &findOptions)

	if err != nil {
		return nil, err
	}

	var results []domain.UserDto
	if err = cur.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	u.userDtoList = results

	return &u.userDtoList, nil
}

func (u UserRepoImpl) FindAllBlockedUsers(id primitive.ObjectID) (*[]domain.UserDto, error) {
	currentUser, err := u.FindByID(id)

	if err != nil {
	 	return nil, err
	 }

	query := bson.M{"_id": bson.M{"$in": currentUser.BlockList}}

	// Get all users
	cur, err := database.GetInstance().UserCollection.Find(context.TODO(), query)

	if err != nil {
		return nil, fmt.Errorf("error processing data")
	}

	var results []domain.UserDto
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	u.userDtoList = results

	return &u.userDtoList, nil
}

func (u UserRepoImpl) Create(user *domain.User) error {
	cur, err := database.GetInstance().UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},
	})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		user.Id = primitive.NewObjectID()
		_, err = database.GetInstance().UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			return fmt.Errorf("error processing data")
		}

		go func() {
			err := events.SendKafkaMessage(user, 201)
			if err != nil {
				fmt.Println("Error publishing...")
				return
			}
		}()

		return nil
	}

	return fmt.Errorf("user already exists")
}

func (u UserRepoImpl) FindByID(id primitive.ObjectID) (*domain.UserDto, error) {
	err := database.GetInstance().UserCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error with the database")
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) FindByUsername(username string) (*domain.UserDto, error) {
	err := database.GetInstance().UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return  nil, fmt.Errorf("error processing data")
	}

	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateByID(id primitive.ObjectID, user *domain.User) (*domain.UserDto, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"tokenHash", user.TokenHash}, {"tokenExpiresAt", user.TokenExpiresAt}}}}

	database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileIsViewable", user.ProfileIsViewable}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.ProfileIsViewable = user.ProfileIsViewable

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			fmt.Println("Error publishing...")
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"acceptMessages", user.AcceptMessages}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.AcceptMessages = user.AcceptMessages

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentBadgeUrl", user.CurrentBadgeUrl}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.CurrentBadgeUrl = user.CurrentBadgeUrl

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profilePictureUrl", user.ProfilePictureUrl}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.ProfilePictureUrl = user.ProfilePictureUrl

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileBackgroundPictureUrl", user.ProfileBackgroundPictureUrl}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.ProfileBackgroundPictureUrl = user.ProfileBackgroundPictureUrl

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentTagLine", user.CurrentTagLine}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.CurrentTagLine = user.CurrentTagLine

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdateVerification(id primitive.ObjectID, user *domain.UpdateVerification) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"isVerified", user.IsVerified}}}}

	err := database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		return err
	}

	u.userDto.IsVerified = user.IsVerified

	mappedUser := domain.UserDtoMapper(u.userDto)

	go func() {
		err := events.HandleKafkaMessage(err, mappedUser, 200)
		if err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	return  nil
}

func (u UserRepoImpl) UpdatePassword(id primitive.ObjectID, password string) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"password", password}, {"tokenHash", ""}, {"tokenExpiresAt", 0}, {"updatedAt", time.Now()}}}}

	database.GetInstance().UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return  nil
}

func (u UserRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
	cur, err := database.GetInstance().FlagCollection.Find(context.TODO(), bson.M{
		"$and": []interface{}{
			bson.M{"flaggerID": flag.FlaggerID},
			bson.M{"flaggedUsername": flag.FlaggedUsername},
		},

	})

	if err != nil {
		return  fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		flag.Id = primitive.NewObjectID()
		_, err = database.GetInstance().FlagCollection.InsertOne(context.TODO(), &flag)

		if err != nil {
			return err
		}

		filter := bson.D{{"username", flag.FlaggedUsername}}
		update := bson.M{"$push": bson.M{"flagCount": flag.Id}}

		_, err = database.GetInstance().UserCollection.UpdateOne(context.TODO(),
			filter, update)
		if err != nil {
			return err
		}

		return nil
	}

	return  fmt.Errorf("you've already flagged this user")
}

func (u UserRepoImpl) BlockUser(id primitive.ObjectID, username string) error {

	err := database.GetInstance().UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if id == u.userDto.Id {
		return fmt.Errorf("you can't block yourself")
	}

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	for _, foundId := range u.userDto.BlockByList {
		if foundId == id {
			return fmt.Errorf("already blocked")
		}
	}

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := database.GetInstance().StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", id}}
		update := bson.M{"$push": bson.M{"blockList": u.userDto.Id}}

		_, err = database.GetInstance().UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		filter = bson.D{{"_id", u.userDto.Id}}
		update = bson.M{"$push": bson.M{"blockByList": id}}

		_, err = database.GetInstance().UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to block user")
	}

	return  nil
}

func (u UserRepoImpl) UnBlockUser(id primitive.ObjectID, username string) error {

	err := database.GetInstance().UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if id == u.userDto.Id {
		return fmt.Errorf("you can't block or unblock yourself")
	}

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	newBlockList, userIsBlocked := util.GenerateNewBlockList(id, u.userDto.BlockByList)

	if !userIsBlocked {
		return fmt.Errorf("this user is not blocked")
	}

	currentUser := new(domain.UserDto)

	// todo better query
	err = database.GetInstance().UserCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&currentUser)

	blockList, userIsBlocked := util.GenerateNewBlockList(u.userDto.Id, currentUser.BlockList)

	if !userIsBlocked {
		return fmt.Errorf("this user is not blocked")
	}

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := database.GetInstance().StartSession()

	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", id}}
		update := bson.M{"$set": bson.M{"blockList": blockList}}

		_, err = database.GetInstance().UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		filter = bson.D{{"_id", u.userDto.Id}}
		update = bson.M{"$set": bson.M{"blockByList": newBlockList}}

		_, err = database.GetInstance().UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		return fmt.Errorf("failed to unblock user")
	}

	return  nil
}

func (u UserRepoImpl) DeleteByID(id primitive.ObjectID) error {
	_, err := database.GetInstance().UserCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}