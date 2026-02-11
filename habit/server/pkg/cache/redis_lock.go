package cache

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var unlockScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
`)

var renewScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("pexpire", KEYS[1], ARGV[2])
else
	return 0
end
`)

func TryLock(ctx context.Context, client *redis.Client, key string, ttl time.Duration) (string, bool, error) {
	if client == nil {
		return "", false, redis.ErrClosed
	}
	token := uuid.NewString()
	ok, err := client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return "", false, err
	}
	return token, ok, nil
}

func Unlock(ctx context.Context, client *redis.Client, key, token string) error {
	if client == nil {
		return redis.ErrClosed
	}
	_, err := unlockScript.Run(ctx, client, []string{key}, token).Result()
	return err
}

type RedisLock struct {
	client *redis.Client
	key    string
	token  string
	ttl    time.Duration
	stopCh chan struct{}
}

func (l *RedisLock) StopRenew() {
	select {
	case <-l.stopCh:
		return
	default:
		close(l.stopCh)
	}
}

func (l *RedisLock) Unlock(ctx context.Context) error {
	l.StopRenew()
	return Unlock(ctx, l.client, l.key, l.token)
}

func AcquireLockWithRenew(ctx context.Context, client *redis.Client, key string, ttl time.Duration) (*RedisLock, bool, error) {
	if ttl <= 0 {
		return nil, false, errors.New("ttl must be positive")
	}

	token, ok, err := TryLock(ctx, client, key, ttl)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	lock := &RedisLock{
		client: client,
		key:    key,
		token:  token,
		ttl:    ttl,
		stopCh: make(chan struct{}),
	}

	interval := ttl / 3
	if interval < 200*time.Millisecond {
		interval = 200 * time.Millisecond
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_, _ = renewScript.Run(ctx, client, []string{key}, token, int64(ttl/time.Millisecond)).Result()
			case <-lock.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return lock, true, nil
}
