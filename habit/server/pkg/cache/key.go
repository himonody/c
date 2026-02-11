package cache

import (
	"fmt"
	"time"
)

// ============================================
// Redis Key 定义规范
// ============================================
// 命名格式: {模块}:{业务}:{标识}
// 例如: user:info:123, admin:token:456

// ============================================
// 过期时间配置
// ============================================
const (
	// 用户相关缓存过期时间
	UserCacheExpiration     = 30 * time.Minute // 用户信息缓存
	UserTokenExpiration     = 24 * time.Hour   // 用户Token
	UserSessionExpiration   = 24 * time.Hour   // 用户会话

	// 管理员相关缓存过期时间
	AdminCacheExpiration    = 30 * time.Minute // 管理员信息缓存
	AdminTokenExpiration    = 24 * time.Hour   // 管理员Token
	AdminSessionExpiration  = 24 * time.Hour   // 管理员会话

	// Token黑名单过期时间
	TokenBlacklistExpiration = 24 * time.Hour

	// 系统配置缓存过期时间
	SystemConfigExpiration = 1 * time.Hour // 系统配置缓存

	// 挑战统计缓存过期时间
	ChallengeStatExpiration = 5 * time.Minute // 平台累计统计缓存

	// 其他缓存过期时间
	DefaultCacheExpiration = 10 * time.Minute  // 默认缓存时间
)

// ============================================
// 用户相关 Redis Key
// ============================================

// UserCacheKey 用户信息缓存 - 通过ID查询
// Key: user:info:{userID}
func UserCacheKey(userID int64) string {
	return fmt.Sprintf("user:info:%d", userID)
}

// UserByUsernameCacheKey 用户信息缓存 - 通过用户名查询
// Key: user:username:{username}
func UserByUsernameCacheKey(username string) string {
	return fmt.Sprintf("user:username:%s", username)
}

// UserTokenKey 用户当前有效Token (单点登录)
// Key: user:token:{userID}
func UserTokenKey(userID int64) string {
	return fmt.Sprintf("user:token:%d", userID)
}

// UserSessionKey 用户会话信息
// Key: user:session:{userID}
func UserSessionKey(userID int64) string {
	return fmt.Sprintf("user:session:%d", userID)
}

// UserOnlineKey 用户在线状态
// Key: user:online:{userID}
func UserOnlineKey(userID int64) string {
	return fmt.Sprintf("user:online:%d", userID)
}

// ============================================
// 管理员相关 Redis Key
// ============================================

// AdminCacheKey 管理员信息缓存 - 通过ID查询
// Key: admin:info:{adminID}
func AdminCacheKey(adminID int64) string {
	return fmt.Sprintf("admin:info:%d", adminID)
}

// AdminByUsernameCacheKey 管理员信息缓存 - 通过用户名查询
// Key: admin:username:{username}
func AdminByUsernameCacheKey(username string) string {
	return fmt.Sprintf("admin:username:%s", username)
}

// AdminTokenKey 管理员当前有效Token (单点登录)
// Key: admin:token:{adminID}
func AdminTokenKey(adminID int64) string {
	return fmt.Sprintf("admin:token:%d", adminID)
}

// AdminSessionKey 管理员会话信息
// Key: admin:session:{adminID}
func AdminSessionKey(adminID int64) string {
	return fmt.Sprintf("admin:session:%d", adminID)
}

// ============================================
// Token黑名单相关 Redis Key
// ============================================

// TokenBlacklistKey 用户Token黑名单
// Key: blacklist:token:{token}
func TokenBlacklistKey(token string) string {
	return fmt.Sprintf("blacklist:token:%s", token)
}

// AdminTokenBlacklistKey 管理员Token黑名单
// Key: blacklist:admin:token:{token}
func AdminTokenBlacklistKey(token string) string {
	return fmt.Sprintf("blacklist:admin:token:%s", token)
}

// ============================================
// 业务相关 Redis Key (后续扩展)
// ============================================

// UserWalletKey 用户钱包缓存
// Key: wallet:user:{userID}
func UserWalletKey(userID int64) string {
	return fmt.Sprintf("wallet:user:%d", userID)
}

// UserInviteCodeKey 用户邀请码缓存
// Key: invite:code:{inviteCode}
func UserInviteCodeKey(inviteCode string) string {
	return fmt.Sprintf("invite:code:%s", inviteCode)
}

// UserInviteListKey 用户邀请列表缓存
// Key: invite:list:{userID}
func UserInviteListKey(userID int64) string {
	return fmt.Sprintf("invite:list:%d", userID)
}

// LeaderboardKey 排行榜缓存
// Key: leaderboard:{type}
// type: daily, weekly, monthly, all
func LeaderboardKey(leaderboardType string) string {
	return fmt.Sprintf("leaderboard:%s", leaderboardType)
}

// UserRankKey 用户排名缓存
// Key: rank:user:{userID}:{type}
func UserRankKey(userID int64, rankType string) string {
	return fmt.Sprintf("rank:user:%d:%s", userID, rankType)
}

// ChallengeKey 挑战信息缓存
// Key: challenge:{challengeID}
func ChallengeKey(challengeID int64) string {
	return fmt.Sprintf("challenge:%d", challengeID)
}

// ChallengeTotalStatKey 平台累计统计缓存
// Key: challenge:stat:total
func ChallengeTotalStatKey() string {
	return "challenge:stat:total"
}

// UserChallengeKey 用户挑战进度缓存
// Key: challenge:user:{userID}:{challengeID}
func UserChallengeKey(userID, challengeID int64) string {
	return fmt.Sprintf("challenge:user:%d:%d", userID, challengeID)
}

// CalendarKey 用户打卡日历缓存
// Key: calendar:user:{userID}:{year}:{month}
func CalendarKey(userID int64, year, month int) string {
	return fmt.Sprintf("calendar:user:%d:%d:%d", userID, year, month)
}

// WithdrawLockKey 提现操作锁
// Key: withdraw:lock:{userID}
func WithdrawLockKey(userID int64) string {
	return fmt.Sprintf("withdraw:lock:%d", userID)
}

// WithdrawRequestKey 提现申请缓存
// Key: withdraw:request:{requestID}
func WithdrawRequestKey(requestID int64) string {
	return fmt.Sprintf("withdraw:request:%d", requestID)
}

// ============================================
// 限流相关 Redis Key
// ============================================

// RateLimitKey 接口限流
// Key: rate:limit:{ip}:{endpoint}
func RateLimitKey(ip, endpoint string) string {
	return fmt.Sprintf("rate:limit:%s:%s", ip, endpoint)
}

// LoginAttemptKey 登录尝试次数
// Key: login:attempt:{username}
func LoginAttemptKey(username string) string {
	return fmt.Sprintf("login:attempt:%s", username)
}

// ============================================
// 验证码相关 Redis Key
// ============================================

// CaptchaKey 图形验证码
// Key: captcha:{captchaID}
func CaptchaKey(captchaID string) string {
	return fmt.Sprintf("captcha:%s", captchaID)
}

// SMSCodeKey 短信验证码
// Key: sms:code:{phone}
func SMSCodeKey(phone string) string {
	return fmt.Sprintf("sms:code:%s", phone)
}

// EmailCodeKey 邮箱验证码
// Key: email:code:{email}
func EmailCodeKey(email string) string {
	return fmt.Sprintf("email:code:%s", email)
}

// ============================================
// 系统配置相关 Redis Key
// ============================================

// SystemConfigKey 系统配置缓存 - 通过配置键查询
// Key: system:config:key:{configKey}
func SystemConfigKey(configKey string) string {
	return fmt.Sprintf("system:config:key:%s", configKey)
}

// SystemConfigIDKey 系统配置缓存 - 通过ID查询
// Key: system:config:id:{configID}
func SystemConfigIDKey(configID int64) string {
	return fmt.Sprintf("system:config:id:%d", configID)
}

// GlobalConfigKey 全局配置缓存
// Key: global:config
func GlobalConfigKey() string {
	return "global:config"
}

// ============================================
// 统计相关 Redis Key
// ============================================

// DailyStatsKey 每日统计
// Key: stats:daily:{date}
func DailyStatsKey(date string) string {
	return fmt.Sprintf("stats:daily:%s", date)
}

// UserStatsKey 用户统计数据
// Key: stats:user:{userID}
func UserStatsKey(userID int64) string {
	return fmt.Sprintf("stats:user:%d", userID)
}

// OnlineUsersCountKey 在线用户数
// Key: stats:online:count
func OnlineUsersCountKey() string {
	return "stats:online:count"
}

// ============================================
// 辅助函数
// ============================================

// GetExpirationByKey 根据Key类型返回对应的过期时间
func GetExpirationByKey(key string) time.Duration {
	// 这里可以根据key的前缀判断应该使用什么过期时间
	// 示例实现，可根据实际需求调整
	return DefaultCacheExpiration
}
