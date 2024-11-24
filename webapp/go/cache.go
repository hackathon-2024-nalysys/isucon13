package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/isucon/isucon13/webapp/go/grpc"
	"google.golang.org/protobuf/proto"
)

func userCacheKey(id int64) string {
	return "user:" + strconv.FormatInt(id, 10)
}

func userIDByCacheKey(key string) int64 {
	if !strings.HasPrefix(key, "user:") {
		log.Fatalf("invalid cache key(prefix), %s", key)
	}
	parsed, err := strconv.ParseInt(key[5:], 10, 64)
	if err != nil {
		log.Fatalf("invalid cache key(int), %s", key)
	}
	return parsed
}

func getUsersByCache(ids []int64) (map[int64]*User, error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = userCacheKey(id)
	}
	items, err := memcacheClient.GetMulti(keys)
	if err != nil {
		return nil, err
	}
	users := make(map[int64]*User, len(items))
	for key, item := range items {
		user := &grpc.User{}
		if err := proto.Unmarshal(item.Value, user); err != nil {
			return nil, err
		}
		users[userIDByCacheKey(key)] = &User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Description: user.Description,
			Theme: Theme{
				ID:       user.Theme.ID,
				DarkMode: user.Theme.DarkMode,
			},
			IconHash: user.IconHash,
		}
	}
	return users, nil
}

func cacheUsers(users []*User) error {
	for _, user := range users {
		pb, err := proto.Marshal(&grpc.User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Description: user.Description,
			Theme: &grpc.Theme{
				ID:       user.Theme.ID,
				DarkMode: user.Theme.DarkMode,
			},
			IconHash: user.IconHash,
		})
		if err != nil {
			return err
		}
		if err := memcacheClient.Set(&memcache.Item{
			Key:   userCacheKey(user.ID),
			Value: pb,
		}); err != nil {
			return err
		}
	}
	return nil
}

func invalidateUserCache(id int64) error {
	err := memcacheClient.Delete(userCacheKey(id))
	if err == memcache.ErrCacheMiss {
		return nil
	}
	return err
}
