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

func livestreamCacheKey(id int64) string {
	return "livestream:" + strconv.FormatInt(id, 10)
}

func livestreamIDByCacheKey(key string) int64 {
	if !strings.HasPrefix(key, "livestream:") {
		log.Fatalf("invalid cache key(prefix), %s", key)
	}
	parsed, err := strconv.ParseInt(key[11:], 10, 64)
	if err != nil {
		log.Fatalf("invalid cache key(int), '%s'", key)
	}
	return parsed
}

func getLivestreamsByCache(ids []int64) (map[int64]*Livestream, error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = livestreamCacheKey(id)
	}
	items := make(map[string]*memcache.Item, len(keys))
	for _, key := range keys {
		item, err := memcacheClient.Get(key)
		if err != nil && err != memcache.ErrCacheMiss {
			return nil, err
		}
		if err != memcache.ErrCacheMiss {
			items[key] = item
		}
	}
	livestreams := make(map[int64]*Livestream, len(items))
	for key, item := range items {
		livestream := &grpc.Livestream{}
		if err := proto.Unmarshal(item.Value, livestream); err != nil {
			return nil, err
		}
		tags := make([]Tag, len(livestream.Tags))
		for i, tag := range livestream.Tags {
			tags[i] = Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
		livestreams[livestreamIDByCacheKey(key)] = &Livestream{
			ID:           livestream.ID,
			Owner:        User{ID: livestream.OwnerID},
			Title:        livestream.Title,
			Description:  livestream.Description,
			PlaylistUrl:  livestream.PlaylistUrl,
			ThumbnailUrl: livestream.ThumbnailUrl,
			Tags:         tags,
			StartAt:      livestream.StartAt,
			EndAt:        livestream.EndAt,
		}
	}
	return livestreams, nil
}

func cacheLivestreams(livestreams []*Livestream) error {
	for _, livestream := range livestreams {
		tags := make([]*grpc.Tag, len(livestream.Tags))
		for i, tag := range livestream.Tags {
			tags[i] = &grpc.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}

		pb, err := proto.Marshal(&grpc.Livestream{
			ID:           livestream.ID,
			OwnerID:      livestream.Owner.ID,
			Title:        livestream.Title,
			Description:  livestream.Description,
			PlaylistUrl:  livestream.PlaylistUrl,
			ThumbnailUrl: livestream.ThumbnailUrl,
			Tags:         tags,
			StartAt:      livestream.StartAt,
			EndAt:        livestream.EndAt,
		})
		if err != nil {
			return err
		}
		if err := memcacheClient.Set(&memcache.Item{
			Key:   livestreamCacheKey(livestream.ID),
			Value: pb,
		}); err != nil {
			return err
		}
	}
	return nil
}

func invalidateLivestreamCache(id int64) error {
	err := memcacheClient.Delete(livestreamCacheKey(id))
	if err == memcache.ErrCacheMiss {
		return nil
	}
	return err
}
