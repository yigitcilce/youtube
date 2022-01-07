package youtube

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

type playerConfig []byte
type playerCache struct {
	key       string
	expiredAt time.Time
	config    playerConfig
}

var basejsPattern = regexp.MustCompile(`(/s/player/\w+/player_ias.vflset/\w+/base.js)`)
var signatureRegexp = regexp.MustCompile(`(?m)(?:^|,)(?:signatureTimestamp:)(\d+)`)

const defaultCacheExpiration = time.Minute * time.Duration(5)

// Get : get cache when it has same video id and not expired
func (s playerCache) Get(key string) playerConfig {
	if key == s.key && s.expiredAt.After(time.Now()) {
		return s.config
	}
	return nil
}

// Set : set cache with default expiration
func (s *playerCache) Set(key string, operations playerConfig) {
	s.key = key
	s.config = operations
	s.expiredAt = time.Now().Add(defaultCacheExpiration)
}

func (c *Client) getPlayerConfig(ctx context.Context, videoID string) (playerConfig, error) {
	embedURL := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", videoID)
	embedBody, err := c.httpGetBodyBytes(ctx, embedURL)
	if err != nil {
		return nil, err
	}

	// example: /s/player/f676c671/player_ias.vflset/en_US/base.js
	escapedBasejsURL := string(basejsPattern.Find(embedBody))
	if escapedBasejsURL == "" {
		return nil, errors.New("unable to find basejs URL in playerConfig")
	}

	config := c.playerCache.Get(escapedBasejsURL)
	if config != nil {
		return config, nil
	}

	config, err = c.httpGetBodyBytes(ctx, "https://youtube.com"+escapedBasejsURL)
	if err != nil {
		return nil, err
	}

	c.playerCache.Set(escapedBasejsURL, config)
	return config, nil
}
