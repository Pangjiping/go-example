package limiter

import (
	"sync"
	"time"
)

// 并发访问同一个user_id/ip需要加锁
var recordMu map[string]*sync.RWMutex

func init() {
	recordMu = make(map[string]*sync.RWMutex)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// record 上次访问时的时间戳和令牌数
type record struct {
	last  time.Time
	token int
}

// TokenBucket 令牌桶的具体实现
type TokenBucket struct {
	BucketSize int                // 令牌桶的容量，最多可以存放多少个令牌
	TokenRate  time.Duration      // 多长时间生成一个令牌
	records    map[string]*record // 报错user_id/ip的访问记录
}

func NewTokenBucket(bucketSize int, tokenRate time.Duration) *TokenBucket {
	return &TokenBucket{
		BucketSize: bucketSize,
		TokenRate:  tokenRate,
		records:    make(map[string]*record),
	}
}

// getUidOrIp 获取请求用户的user_id/ip
func (t *TokenBucket) getUidOrIp() string {
	return "127.0.0.1"
}

// getRecord 获取这个user_id/ip上次访问的时间戳和令牌数
func (t *TokenBucket) getRecord(uidOrIp string) *record {
	if r, ok := t.records[uidOrIp]; ok {
		return r
	}
	return &record{}
}

func (t *TokenBucket) storeRecord(uidOrIp string, r *record) {
	t.records[uidOrIp] = r
}

func (t *TokenBucket) validate(uidOrIp string) bool {
	rl, ok := recordMu[uidOrIp]
	if !ok {
		var mu sync.RWMutex
		rl = &mu
		recordMu[uidOrIp] = rl
	}

	rl.Lock()
	defer rl.Unlock()

	r := t.getRecord(uidOrIp)
	now := time.Now()
	if r.last.IsZero() {
		// 第一次访问初始化为最大令牌数
		r.last, r.token = now, t.BucketSize
	} else {
		if r.last.Add(t.TokenRate).Before(now) {
			// 如果与上一次请求隔了token rate
			// 增加令牌，更新last
			r.token += max(int(now.Sub(r.last)/t.TokenRate), t.BucketSize)
			r.last = now
		}
	}

	var result bool
	// 如果令牌数大于1，取走一个令牌，validate结果为true
	if r.token > 0 {
		r.token--
		result = true
	}

	// 保存最新的record
	t.storeRecord(uidOrIp, r)
	return result
}

// IsLimited 是否被限流
func (t *TokenBucket) IsLimited() bool {
	return !t.validate(t.getUidOrIp())
}
