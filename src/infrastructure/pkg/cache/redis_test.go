package cache_test

import (
	"context"
	"testing"
	"time"

	repo "everlasting/src/domain/sharedkernel/cache"
	"everlasting/src/infrastructure/pkg/cache"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
)

type (
	RedisCacheGetSuite struct {
		suite.Suite
		miniredis *miniredis.Miniredis
		redis     *cache.RedisCache
	}
)

func (suite *RedisCacheGetSuite) SetupTest() {
	var err error
	suite.miniredis, err = miniredis.Run()
	if err != nil {
		suite.Errorf(err, "miniredis")
	}

	suite.redis = cache.NewRedisCache(redis.NewClient(&redis.Options{
		Addr:     suite.miniredis.Addr(),
		Password: "",
	}))
}

func (suite *RedisCacheGetSuite) TestGetWithNilResult() {
	ctx := context.Background()

	var res int
	err := suite.redis.Get(ctx, "invalid_key", res)
	suite.ErrorIs(err, repo.ErrCacheKeyNotFound)
}

func (suite *RedisCacheGetSuite) TestSetAndGet() {
	ctx := context.Background()
	expectedResult := "asu"

	// Set value to redis
	suite.redis.Set(ctx, "key1", expectedResult, time.Second*10)

	var res string
	err := suite.redis.Get(ctx, "key1", &res)

	// Assertion
	suite.NoError(err)
	suite.Equal(expectedResult, res)
}

func (suite *RedisCacheGetSuite) TestSetError() {
	ctx := context.Background()
	err := suite.redis.Set(ctx, "key1", make(chan int32), 0)
	suite.Error(err)
}

func (suite *RedisCacheGetSuite) TestSetDeleteAndExist() {
	ctx := context.Background()
	expectedResult := 200

	// Set value to redis
	suite.redis.Set(ctx, "key1", expectedResult, time.Second*10)
	suite.redis.Set(ctx, "key2", expectedResult, time.Second*10)

	// Delete value to redis
	suite.redis.Delete(ctx, "key1")

	// Check key1, expected not found
	found, _ := suite.redis.Exist(ctx, "key1")
	suite.Equal(int64(0), found)

	// Check key2, expected found
	found, _ = suite.redis.Exist(ctx, "key2")
	suite.Equal(int64(1), found)
}

func (suite *RedisCacheGetSuite) TestGetWithExpiredKey() {
	ctx := context.Background()
	expectedResult := 200

	// Set value to redis with a short expiration time
	suite.redis.Set(ctx, "key1", expectedResult, time.Millisecond*10)

	// Wait for the key to expire
	time.Sleep(time.Millisecond * 20)

	// Attempt to get the expired key, expected result is ErrCacheKeyNotFound
	var res int
	err := suite.redis.Get(ctx, "key1", &res)

	// Assertion
	suite.NoError(err)
	suite.Equal(expectedResult, res)
}

func (suite *RedisCacheGetSuite) TestExistWithInvalidKey() {
	ctx := context.Background()

	// Attempt to check existence with an invalid key, expected result is not found
	found, _ := suite.redis.Exist(ctx, "invalid_key")

	// Assertion
	suite.Equal(int64(0), found)
}

func (suite *RedisCacheGetSuite) TestExistWithError() {
	// Close the Redis connection to simulate an error
	_ = suite.redis.Close()

	ctx := context.Background()

	// Attempt to check existence, expected result is an error
	found, err := suite.redis.Exist(ctx, "key1")

	// Assertion
	suite.Error(err)
	suite.Equal(int64(0), found)
}

func (suite *RedisCacheGetSuite) TestClose() {
	err := suite.redis.Close()
	suite.NoError(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheGetSuite))
}
