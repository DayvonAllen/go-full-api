package repo

import (
	"context"
	cache2 "example.com/app/cache"
	"example.com/app/database"
	"example.com/app/domain"
	"example.com/app/events"
	"example.com/app/util"
	"fmt"
	"github.com/go-redis/cache/v8"
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
	users        []domain.User
	user         domain.User
	userDto      domain.UserDto
	userDtoList  []domain.UserDto
	userResponse domain.UserResponse
}

func (u UserRepoImpl) FindAll(id primitive.ObjectID, page string, ctx context.Context, rdb *cache.Cache, username string) (*domain.UserResponse, error) {

	var data domain.UserDto

	err := rdb.Get(ctx, util.GenerateKey(username, "finduserbyusername"), &data)

	var currentUser *domain.UserDto

	if err == nil {
		cache2.RedisCachePool.Put(rdb)
		currentUser = &data
		fmt.Println("Found in Cache in find all users...")
	} else {
		currentUser, err = u.FindByID(id, rdb, ctx)

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			fmt.Println("Did not find in Cache in find all users...")
			return nil, err
		}
	}

	conn := database.MongoConnectionPool.Get().(*database.Connection)

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return nil, fmt.Errorf("page must be a number")
	}
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	// Get all users
	cur, err := conn.UserCollection.Find(ctx, bson.M{
		"profileIsViewable": true,
		"$and": []interface{}{
			bson.M{"_id": bson.M{"$ne": id}},
			bson.M{"_id": bson.M{"$nin": currentUser.BlockByList}},
			bson.M{"_id": bson.M{"$nin": currentUser.BlockList}},
		},
	}, &findOptions)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return nil, err
	}

	var results []domain.UserDto
	if err = cur.All(ctx, &results); err != nil {
		database.MongoConnectionPool.Put(conn)
		log.Fatal(err)
	}

	u.userDtoList = results

	u.userResponse = domain.UserResponse{Users: u.userDtoList, CurrentPage: page}

	database.MongoConnectionPool.Put(conn)

	return &u.userResponse, nil
}

func (u UserRepoImpl) FindAllBlockedUsers(id primitive.ObjectID, rdb *cache.Cache, ctx context.Context, username string) (*[]domain.UserDto, error) {
	var data domain.UserDto

	err := rdb.Get(ctx, util.GenerateKey(username, "finduserbyusername"), &data)

	var currentUser *domain.UserDto

	if err == nil {
		cache2.RedisCachePool.Put(rdb)
		currentUser = &data
		fmt.Println("Found in Cache in find all blocked users...")
	} else {
		currentUser, err = u.FindByID(id, rdb, ctx)

		if err != nil {
			fmt.Println("Did not find in Cache in find all blocked users...")
			cache2.RedisCachePool.Put(rdb)
			return nil, err
		}
	}

	conn := database.MongoConnectionPool.Get().(*database.Connection)

	query := bson.M{"_id": bson.M{"$in": currentUser.BlockList}}

	// Get all users
	cur, err := conn.UserCollection.Find(context.TODO(), query)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return nil, fmt.Errorf("error processing data")
	}

	var results []domain.UserDto
	if err = cur.All(context.TODO(), &results); err != nil {
		database.MongoConnectionPool.Put(conn)
		panic(err)
	}

	u.userDtoList = results

	database.MongoConnectionPool.Put(conn)
	return &u.userDtoList, nil
}

func (u UserRepoImpl) Create(user *domain.User) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	cur, err := conn.UserCollection.Find(context.TODO(), bson.M{
		"$or": []interface{}{
			bson.M{"email": user.Email},
			bson.M{"username": user.Username},
		},
	})

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("error processing data")
	}
	found := cur.Next(context.TODO())
	if !found {
		user.Id = primitive.NewObjectID()
		_, err = conn.UserCollection.InsertOne(context.TODO(), &user)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return fmt.Errorf("error processing data")
		}

		go func() {
			err := events.SendKafkaMessage(user, 201)
			if err != nil {
				fmt.Println("Error publishing...")
				return
			}
		}()

		database.MongoConnectionPool.Put(conn)
		return nil
	}
	err = cur.Decode(&u.userDto)
	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return err
	}

	err = cur.Close(context.TODO())

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return err
	}

	if u.userDto.Username == user.Username {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("username is taken")
	}

	database.MongoConnectionPool.Put(conn)
	return fmt.Errorf("email is taken")
}

func (u UserRepoImpl) FindByID(id primitive.ObjectID, rdb *cache.Cache, ctx context.Context) (*domain.UserDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error with the database")
	}

	go func() {
		err = rdb.Set(&cache.Item{
			Ctx:   ctx,
			Key:   util.GenerateKey(u.userDto.Username, "finduserbyusername"),
			Value: u.userDto,
			TTL:   time.Hour,
		})

		if err != nil {
			fmt.Println("Found in cache in find by ID...")
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}
		cache2.RedisCachePool.Put(rdb)
		fmt.Println("Cached in find by ID...")
		return
	}()

	database.MongoConnectionPool.Put(conn)
	return &u.userDto, nil
}

func (u UserRepoImpl) FindByUsername(username string, rdb *cache.Cache, ctx context.Context) (*domain.UserDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	err := conn.UserCollection.FindOne(context.TODO(), bson.M{"username": username, "$and":
	[]interface{}{
		bson.M{"profileIsViewable": true,
		},
	}}).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("cannot find user")
		}
		return nil, fmt.Errorf("error processing data")
	}

	go func() {
		err = rdb.Set(&cache.Item{
			Ctx:   ctx,
			Key:   util.GenerateKey(username, "finduserbyusername"),
			Value: u.userDto,
			TTL:   time.Hour,
		})

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}
		cache2.RedisCachePool.Put(rdb)
		fmt.Println("Cached in find by username...")
		return
	}()

	database.MongoConnectionPool.Put(conn)
	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateByID(id primitive.ObjectID, user *domain.User) (*domain.UserDto, error) {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"tokenHash", user.TokenHash}, {"tokenExpiresAt", user.TokenExpiresAt}}}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	database.MongoConnectionPool.Put(conn)
	return &u.userDto, nil
}

func (u UserRepoImpl) UpdateProfileVisibility(id primitive.ObjectID, user *domain.UpdateProfileVisibility, rdb *cache.Cache, ctx context.Context) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileIsViewable", user.ProfileIsViewable}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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

	go func() {

		fmt.Println(util.GenerateKey(u.userDto.Username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(u.userDto.Username, "finduserbyusername"))

		if err != nil {
			fmt.Println("Not in cache, update profile visibility")
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}
		fmt.Println("Removed from cache, update profile visibility")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)

	return nil
}

func (u UserRepoImpl) UpdateMessageAcceptance(id primitive.ObjectID, user *domain.UpdateMessageAcceptance, rdb *cache.Cache, ctx context.Context) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"acceptMessages", user.AcceptMessages}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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

	go func() {

		fmt.Println(util.GenerateKey(u.userDto.Username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(u.userDto.Username, "finduserbyusername"))

		if err != nil {
			fmt.Println("Not in cache, update message acceptance")
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, update message acceptance")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdateCurrentBadge(id primitive.ObjectID, user *domain.UpdateCurrentBadge, rdb *cache.Cache, ctx context.Context) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentBadgeUrl", user.CurrentBadgeUrl}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
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

	go func() {

		fmt.Println(util.GenerateKey(u.userDto.Username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(u.userDto.Username, "finduserbyusername"))

		if err != nil {
			fmt.Println("Not in cache, update message acceptance")
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, update current badge")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdateProfilePicture(id primitive.ObjectID, user *domain.UpdateProfilePicture, rdb *cache.Cache, ctx context.Context) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profilePictureUrl", user.ProfilePictureUrl}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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

	go func() {
		fmt.Println(util.GenerateKey(u.userDto.Username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(u.userDto.Username, "finduserbyusername"))

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, update profile picture")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdateProfileBackgroundPicture(id primitive.ObjectID, user *domain.UpdateProfileBackgroundPicture, rdb *cache.Cache, ctx context.Context) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"profileBackgroundPictureUrl", user.ProfileBackgroundPictureUrl}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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

	go func() {
		fmt.Println(util.GenerateKey(u.userDto.Username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(u.userDto.Username, "finduserbyusername"))

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, update profile background picture")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdateCurrentTagline(id primitive.ObjectID, user *domain.UpdateCurrentTagline, rdb *cache.Cache, ctx context.Context) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"currentTagLine", user.CurrentTagLine}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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

	go func() {

		fmt.Println(util.GenerateKey(u.userDto.Username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(u.userDto.Username, "finduserbyusername"))

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, update current tag")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdateVerification(id primitive.ObjectID, user *domain.UpdateVerification) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"isVerified", user.IsVerified}}}}

	err := conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&u.userDto)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
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
		database.MongoConnectionPool.Put(conn)
		return err
	}

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdatePassword(id primitive.ObjectID, password string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"password", password}, {"tokenHash", ""}, {"tokenExpiresAt", 0}, {"updatedAt", time.Now()}}}}

	conn.UserCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UpdateFlagCount(flag *domain.Flag) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	cur, err := conn.FlagCollection.Find(context.TODO(), bson.M{
		"$and": []interface{}{
			bson.M{"flaggerID": flag.FlaggerID},
			bson.M{"flaggedUsername": flag.FlaggedUsername},
		},
	})

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("error processing data")
	}

	if !cur.Next(context.TODO()) {
		flag.Id = primitive.NewObjectID()
		_, err = conn.FlagCollection.InsertOne(context.TODO(), &flag)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return err
		}

		filter := bson.D{{"username", flag.FlaggedUsername}}
		update := bson.M{"$push": bson.M{"flagCount": flag.Id}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)
		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return err
		}

		database.MongoConnectionPool.Put(conn)
		return nil
	}

	database.MongoConnectionPool.Put(conn)
	return fmt.Errorf("you've already flagged this user")
}

func (u UserRepoImpl) BlockUser(id primitive.ObjectID, username string, rdb *cache.Cache, ctx context.Context, currentUsername string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if id == u.userDto.Id {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("you can't block yourself")
	}

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	for _, foundId := range u.userDto.BlockByList {
		if foundId == id {
			database.MongoConnectionPool.Put(conn)
			return fmt.Errorf("already blocked")
		}
	}

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		panic(err)
	}

	defer session.EndSession(context.Background())

	// execute this code in a logical transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", id}}
		update := bson.M{"$push": bson.M{"blockList": u.userDto.Id}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return nil, err
		}

		filter = bson.D{{"_id", u.userDto.Id}}
		update = bson.M{"$push": bson.M{"blockByList": id}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return nil, err
		}

		database.MongoConnectionPool.Put(conn)
		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("failed to block user")
	}

	go func() {
		fmt.Println(util.GenerateKey(currentUsername, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(currentUsername, "finduserbyusername"))

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, block user")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) UnblockUser(id primitive.ObjectID, username string, rdb *cache.Cache, ctx context.Context, currentUsername string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	err := conn.UserCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&u.userDto)

	if id == u.userDto.Id {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("you can't block or unblock yourself")
	}

	if err != nil {
		database.MongoConnectionPool.Put(conn)

		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return err
	}

	newBlockList, userIsBlocked := util.GenerateNewBlockList(id, u.userDto.BlockByList)

	if !userIsBlocked {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("this user is not blocked")
	}

	currentUser := new(domain.UserDto)

	// todo better query
	err = conn.UserCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&currentUser)

	blockList, userIsBlocked := util.GenerateNewBlockList(u.userDto.Id, currentUser.BlockList)

	if !userIsBlocked {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("this user is not blocked")
	}

	// sets mongo's read and write concerns
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// set up for a transaction
	session, err := conn.StartSession()

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		panic(err)
	}

	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {

		filter := bson.D{{"_id", id}}
		update := bson.M{"$set": bson.M{"blockList": blockList}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return nil, err
		}

		filter = bson.D{{"_id", u.userDto.Id}}
		update = bson.M{"$set": bson.M{"blockByList": newBlockList}}

		_, err = conn.UserCollection.UpdateOne(context.TODO(),
			filter, update)

		if err != nil {
			database.MongoConnectionPool.Put(conn)
			return nil, err
		}

		database.MongoConnectionPool.Put(conn)
		return nil, err
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return fmt.Errorf("failed to unblock user")
	}

	go func() {
		fmt.Println(util.GenerateKey(currentUsername, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(currentUsername, "finduserbyusername"))

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, unblock user")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func (u UserRepoImpl) DeleteByID(id primitive.ObjectID, rdb *cache.Cache, ctx context.Context, username string) error {
	conn := database.MongoConnectionPool.Get().(*database.Connection)

	_, err := conn.UserCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

	if err != nil {
		database.MongoConnectionPool.Put(conn)
		return err
	}

	u.user.Id = id

	go func() {
		err := events.HandleKafkaMessage(err, &u.user, 204)
		if err != nil {
			return
		}
	}()

	go func() {
		fmt.Println(util.GenerateKey(username, "finduserbyusername"))
		err := rdb.Delete(ctx, util.GenerateKey(username, "finduserbyusername"))

		if err != nil {
			cache2.RedisCachePool.Put(rdb)
			panic(err)
		}

		fmt.Println("Removed from cache, delete by ID")
		cache2.RedisCachePool.Put(rdb)

		return
	}()

	database.MongoConnectionPool.Put(conn)
	return nil
}

func NewUserRepoImpl() UserRepoImpl {
	var userRepoImpl UserRepoImpl

	return userRepoImpl
}
