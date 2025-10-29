package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	rdb "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrCodeNotFound = errors.New("验证码不存在")
	ErrCodeExpired  = errors.New("验证码已过期")
	ErrCodeUsed     = errors.New("验证码已使用")
	ErrCodeInvalid  = errors.New("验证码错误")
)

// Lua脚本：验证并标记验证码为已使用（原子操作）
var verifyAndMarkUsedScript = `
local key = KEYS[1]
local code = ARGV[1]
local now = tonumber(ARGV[2])
local expire_seconds = tonumber(ARGV[3])  -- 过期时间（秒，10分钟=600秒）

-- 获取验证码信息
local data = redis.call('HGETALL', key)
if #data == 0 then
    return {0, '验证码不存在'}
end

-- 解析验证码数据
local stored_code = nil
local used = '0'
local created_at = 0

for i = 1, #data, 2 do
    if data[i] == 'code' then
        stored_code = data[i + 1]
    elseif data[i] == 'used' then
        used = data[i + 1]
    elseif data[i] == 'created_at' then
        created_at = tonumber(data[i + 1])
    end
end

-- 检查是否已使用
if used == '1' then
    return {0, '验证码已使用'}
end

-- 检查是否过期（10分钟）
local elapsed = now - created_at
if elapsed > expire_seconds then
    return {0, '验证码已过期'}
end

-- 验证验证码
if stored_code ~= code then
    return {0, '验证码错误'}
end

-- 标记为已使用
redis.call('HSET', key, 'used', '1', 'used_at', now)
-- 标记后5分钟后删除记录（用于审计）
redis.call('EXPIRE', key, 300)
return {1, 'success'}
`

// Lua脚本：检查发送频率（原子操作）
var checkSendFrequencyScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local min_interval = tonumber(ARGV[2]) -- 最小间隔（秒）

-- 获取上次发送时间
local data = redis.call('HGETALL', key)
if #data == 0 then
    return {1, '0'} -- 可以发送
end

local last_send_time = 0
for i = 1, #data, 2 do
    if data[i] == 'created_at' then
        last_send_time = tonumber(data[i + 1])
        break
    end
end

if last_send_time == 0 or (now - last_send_time) >= min_interval then
    return {1, '0'} -- 可以发送
else
    return {0, tostring(min_interval - (now - last_send_time))} -- 需要等待的秒数
end
`

// SMSCodeStore Redis验证码存储
type SMSCodeStore struct {
	client *rdb.Client
}

// NewSMSCodeStore 创建Redis验证码存储
func NewSMSCodeStore() *SMSCodeStore {
	if Client == nil {
		panic("Redis客户端未初始化，请先调用 redis.Init()")
	}
	return &SMSCodeStore{
		client: Client,
	}
}

// getKey 获取Redis键名
func (s *SMSCodeStore) getKey(phone, codeType string) string {
	return fmt.Sprintf("sms:code:%s:%s", phone, codeType)
}

// StoreCode 存储验证码（10分钟过期）
func (s *SMSCodeStore) StoreCode(ctx context.Context, phone, code, codeType string, clientIP string) error {
	if s.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	key := s.getKey(phone, codeType)
	now := time.Now().Unix()
	expireSeconds := 600 // 10分钟

	// 使用Hash存储验证码信息
	fields := map[string]interface{}{
		"code":       code,
		"phone":      phone,
		"type":       codeType,
		"used":       "0",
		"created_at": now,
		"ip":         clientIP,
	}

	if err := s.client.HSet(ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 设置过期时间（10分钟）
	if err := s.client.Expire(ctx, key, time.Duration(expireSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("设置验证码过期时间失败: %w", err)
	}

	logger.Info("验证码已存储到Redis",
		zap.String("phone", phone),
		zap.String("type", codeType))

	return nil
}

// VerifyCode 验证验证码（使用Lua脚本保证原子性）
func (s *SMSCodeStore) VerifyCode(ctx context.Context, phone, code, codeType string) error {
	if s.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	key := s.getKey(phone, codeType)
	now := time.Now().Unix()
	expireSeconds := 600 // 10分钟（600秒）

	// 加载Lua脚本
	script := rdb.NewScript(verifyAndMarkUsedScript)

	result, err := script.Run(ctx, s.client, []string{key}, code, now, expireSeconds).Result()
	if err != nil {
		if err == rdb.Nil {
			return ErrCodeNotFound
		}
		logger.Error("验证码验证失败", zap.Error(err))
		return fmt.Errorf("验证码验证失败: %w", err)
	}

	// 解析Lua脚本返回结果
	results, ok := result.([]interface{})
	if !ok || len(results) != 2 {
		return fmt.Errorf("验证码验证结果格式错误")
	}

	success := results[0].(int64)
	message := results[1].(string)

	if success == 0 {
		switch message {
		case "验证码不存在":
			return ErrCodeNotFound
		case "验证码已使用":
			return ErrCodeUsed
		case "验证码已过期":
			return ErrCodeExpired
		case "验证码错误":
			return ErrCodeInvalid
		default:
			return fmt.Errorf("%s", message)
		}
	}

	logger.Info("验证码验证成功",
		zap.String("phone", phone),
		zap.String("type", codeType))

	return nil
}

// CheckSendFrequency 检查发送频率（使用Lua脚本保证原子性）
func (s *SMSCodeStore) CheckSendFrequency(ctx context.Context, phone, codeType string, minIntervalSeconds int) (bool, int, error) {
	if s.client == nil {
		return false, 0, fmt.Errorf("Redis客户端未初始化")
	}

	key := s.getKey(phone, codeType)
	now := time.Now().Unix()

	// 加载Lua脚本
	script := rdb.NewScript(checkSendFrequencyScript)

	result, err := script.Run(ctx, s.client, []string{key}, now, minIntervalSeconds).Result()
	if err != nil {
		logger.Error("检查发送频率失败", zap.Error(err))
		return false, 0, fmt.Errorf("检查发送频率失败: %w", err)
	}

	// 解析Lua脚本返回结果
	results, ok := result.([]interface{})
	if !ok || len(results) != 2 {
		return false, 0, fmt.Errorf("检查发送频率结果格式错误")
	}

	canSend := results[0].(int64)
	waitSeconds, _ := results[1].(string)
	waitSec := 0
	fmt.Sscanf(waitSeconds, "%d", &waitSec)

	if canSend == 1 {
		return true, 0, nil
	}

	return false, waitSec, nil
}

// DeleteCode 删除验证码
func (s *SMSCodeStore) DeleteCode(ctx context.Context, phone, codeType string) error {
	if s.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	key := s.getKey(phone, codeType)
	return s.client.Del(ctx, key).Err()
}
