package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// AppChallenge 打卡挑战活动配置
type AppChallenge struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// DayCount 挑战天数（如 1/7/21）
	DayCount     int             `gorm:"column:day_count;not null;default:0" json:"day_count"`
	// Amount 单人挑战金额
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(30,2);not null;default:0.00" json:"amount"`
	// CheckinStart 每日打卡开始时间（分钟/秒等具体单位按业务约定）
	CheckinStart uint16          `gorm:"column:checkin_start;not null;default:0" json:"checkin_start"`
	// CheckinEnd 每日打卡结束时间（分钟/秒等具体单位按业务约定）
	CheckinEnd   uint16          `gorm:"column:checkin_end;not null;default:0" json:"checkin_end"`
	// PlatformBonus 平台补贴金额
	PlatformBonus decimal.Decimal `gorm:"column:platform_bonus;type:decimal(30,2);not null;default:0.00" json:"platform_bonus"`
	// Status 状态：1启用 2停用
	Status       int8            `gorm:"column:status;not null;default:1" json:"status"`
	// Sort 排序
	Sort         int8            `gorm:"column:sort;not null;default:1" json:"sort"`
	CreatedAt    time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (AppChallenge) TableName() string {
	return "app_challenge"
}

// AppUserChallenge 用户参与挑战记录
type AppUserChallenge struct {
	ID              int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID          int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`
	ChallengeID     int64           `gorm:"column:challenge_id;not null;default:0" json:"challenge_id"`
	PoolID          int64           `gorm:"column:pool_id;not null;default:0" json:"pool_id"`
	// ChallengeAmount 用户挑战金额
	ChallengeAmount decimal.Decimal `gorm:"column:challenge_amount;type:decimal(30,2);not null;default:0.00" json:"challenge_amount"`
	// StartDate 活动开始日期
	StartDate       time.Time       `gorm:"column:start_date;type:date;not null" json:"start_date"`
	// EndDate 活动结束日期
	EndDate         time.Time       `gorm:"column:end_date;type:date;not null" json:"end_date"`
	// Status 状态：1进行中 2成功 3失败
	Status          int8            `gorm:"column:status;not null;default:1" json:"status"`
	// FailReason 失败原因：0无 1已打卡 2未打卡
	FailReason      int8            `gorm:"column:fail_reason;not null;default:2" json:"fail_reason"`
	CreatedAt       time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	FinishedAt      *time.Time      `gorm:"column:finished_at" json:"finished_at"`
}

func (AppUserChallenge) TableName() string {
	return "app_user_challenge"
}

// AppUserChallengeCheckin 每日打卡记录（心情 + 内容）
type AppUserChallengeCheckin struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserChallengeID int64      `gorm:"column:user_challenge_id;not null;default:0" json:"user_challenge_id"`
	UserID          int64      `gorm:"column:user_id;not null;default:0" json:"user_id"`
	// CheckinDate 打卡日期
	CheckinDate     time.Time  `gorm:"column:checkin_date;type:date;not null" json:"checkin_date"`
	// CheckinTime 打卡时间
	CheckinTime     *time.Time `gorm:"column:checkin_time" json:"checkin_time"`
	// MoodCode 心情枚举：1开心 2平静 3一般 4疲惫 5低落 6爆棚
	MoodCode        int8       `gorm:"column:mood_code;not null;default:0" json:"mood_code"`
	// MoodText 心情文字描述
	MoodText        string     `gorm:"column:mood_text;type:text;not null" json:"mood_text"`
	// ContentType 内容类型：1图片 2视频广告
	ContentType     int8       `gorm:"column:content_type;not null;default:1" json:"content_type"`
	// Status 状态：1打卡成功 2未打卡 3打卡失败
	Status          int8       `gorm:"column:status;not null;default:2" json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (AppUserChallengeCheckin) TableName() string {
	return "app_user_challenge_checkin"
}

// AppChallengeCheckinImage 打卡图片表
type AppChallengeCheckinImage struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CheckinID int64     `gorm:"column:checkin_id;not null;default:0" json:"checkin_id"`
	UserID    int64     `gorm:"column:user_id;not null;default:0" json:"user_id"`
	// ImageURL 图片URL
	ImageURL  string    `gorm:"column:image_url;size:500;not null;default:''" json:"image_url"`
	// ImageHash 图片Hash（防重复）
	ImageHash string    `gorm:"column:image_hash;size:64;not null;default:''" json:"image_hash"`
	// SortNo 图片顺序
	SortNo    int8      `gorm:"column:sort_no;not null;default:1" json:"sort_no"`
	// Status 状态：1正常 2屏蔽 3审核中
	Status    int8      `gorm:"column:status;not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (AppChallengeCheckinImage) TableName() string {
	return "app_challenge_checkin_image"
}

// AppChallengeCheckinVideoAd 视频广告打卡记录
type AppChallengeCheckinVideoAd struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CheckinID    int64           `gorm:"column:checkin_id;not null;default:0" json:"checkin_id"`
	UserID       int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`
	// AdPlatform 广告平台，如：csj、gdt、unity
	AdPlatform   string          `gorm:"column:ad_platform;size:50;not null;default:''" json:"ad_platform"`
	// AdUnitID 广告位ID
	AdUnitID     string          `gorm:"column:ad_unit_id;size:100;not null;default:''" json:"ad_unit_id"`
	// AdOrderNo 广告联盟订单号（唯一）
	AdOrderNo    string          `gorm:"column:ad_order_no;size:100;not null;default:''" json:"ad_order_no"`
	// VideoDuration 视频时长（秒）
	VideoDuration int            `gorm:"column:video_duration;not null;default:0" json:"video_duration"`
	// WatchDuration 实际观看时长（秒）
	WatchDuration int            `gorm:"column:watch_duration;not null;default:0" json:"watch_duration"`
	// RewardAmount 该广告产生的收益
	RewardAmount decimal.Decimal `gorm:"column:reward_amount;type:decimal(30,2);not null;default:0.00" json:"reward_amount"`
	// VerifyStatus 校验状态：0待校验 1成功 2失败
	VerifyStatus int8            `gorm:"column:verify_status;not null;default:0" json:"verify_status"`
	CreatedAt    time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	VerifiedAt   *time.Time      `gorm:"column:verified_at" json:"verified_at"`
}

func (AppChallengeCheckinVideoAd) TableName() string {
	return "app_challenge_checkin_video_ad"
}

// AppChallengePool 活动奖池表
type AppChallengePool struct {
	ID          int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ChallengeID int64           `gorm:"column:challenge_id;not null;default:0" json:"challenge_id"`
	// StartDate 活动开始日期
	StartDate   *time.Time      `gorm:"column:start_date" json:"start_date"`
	// EndDate 活动结束日期
	EndDate     *time.Time      `gorm:"column:end_date" json:"end_date"`
	// TotalAmount 奖池当前总金额
	TotalAmount decimal.Decimal `gorm:"column:total_amount;type:decimal(30,2);not null;default:0.00" json:"total_amount"`
	// Settled 是否已结算：0否 1是
	Settled     int8            `gorm:"column:settled;not null;default:0" json:"settled"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (AppChallengePool) TableName() string {
	return "app_challenge_pool"
}

// AppChallengePoolFlow 奖池资金流水
type AppChallengePoolFlow struct {
	ID        int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PoolID    int64           `gorm:"column:pool_id;not null;default:0" json:"pool_id"`
	UserID    int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`
	Amount    decimal.Decimal `gorm:"column:amount;type:decimal(30,2);not null;default:0.00" json:"amount"`
	// Type 类型：1报名 2失败 3平台补贴 4结算
	Type      int8            `gorm:"column:type;not null;default:0" json:"type"`
	CreatedAt time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (AppChallengePoolFlow) TableName() string {
	return "app_challenge_pool_flow"
}

// AppChallengeSettlement 挑战结算结果
type AppChallengeSettlement struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserChallengeID int64          `gorm:"column:user_challenge_id;not null;default:0" json:"user_challenge_id"`
	UserID         int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`
	// Reward 最终获得金额
	Reward         decimal.Decimal `gorm:"column:reward;type:decimal(30,2);not null;default:0.00" json:"reward"`
	CreatedAt      time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (AppChallengeSettlement) TableName() string {
	return "app_challenge_settlement"
}

// AppChallengeDailyStat 每日运营统计
type AppChallengeDailyStat struct {
	// StatDate 统计日期
	StatDate      time.Time       `gorm:"column:stat_date;type:date;primaryKey" json:"stat_date"`
	// JoinUserCnt 参与人数
	JoinUserCnt   int             `gorm:"column:join_user_cnt;not null;default:0" json:"join_user_cnt"`
	// SuccessUserCnt 成功人数
	SuccessUserCnt int            `gorm:"column:success_user_cnt;not null;default:0" json:"success_user_cnt"`
	// FailUserCnt 失败人数
	FailUserCnt   int             `gorm:"column:fail_user_cnt;not null;default:0" json:"fail_user_cnt"`
	// JoinAmount 参与总金额
	JoinAmount    decimal.Decimal `gorm:"column:join_amount;type:decimal(30,2);not null;default:0.00" json:"join_amount"`
	// SuccessAmount 成功金额
	SuccessAmount decimal.Decimal `gorm:"column:success_amount;type:decimal(30,2);not null;default:0.00" json:"success_amount"`
	// FailAmount 失败金额
	FailAmount    decimal.Decimal `gorm:"column:fail_amount;type:decimal(30,2);not null;default:0.00" json:"fail_amount"`
	// PlatformBonus 平台补贴
	PlatformBonus decimal.Decimal `gorm:"column:platform_bonus;type:decimal(30,2);not null;default:0.00" json:"platform_bonus"`
	// PoolAmount 奖池金额
	PoolAmount    decimal.Decimal `gorm:"column:pool_amount;type:decimal(30,2);not null;default:0.00" json:"pool_amount"`
	CreatedAt     time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (AppChallengeDailyStat) TableName() string {
	return "app_challenge_daily_stat"
}

// AppChallengeTotalStat 平台累计统计
type AppChallengeTotalStat struct {
	// ID 固定主键（通常为 1）
	ID               int8            `gorm:"column:id;primaryKey;not null;default:1" json:"id"`
	// TotalUserCnt 累计用户数
	TotalUserCnt     int             `gorm:"column:total_user_cnt;not null;default:0" json:"total_user_cnt"`
	// TotalJoinCnt 累计参与人次
	TotalJoinCnt     int             `gorm:"column:total_join_cnt;not null;default:0" json:"total_join_cnt"`
	// TotalSuccessCnt 累计成功人次
	TotalSuccessCnt  int             `gorm:"column:total_success_cnt;not null;default:0" json:"total_success_cnt"`
	// TotalFailCnt 累计失败人次
	TotalFailCnt     int             `gorm:"column:total_fail_cnt;not null;default:0" json:"total_fail_cnt"`
	// TotalJoinAmount 累计参与金额
	TotalJoinAmount  decimal.Decimal `gorm:"column:total_join_amount;type:decimal(30,2);not null;default:0.00" json:"total_join_amount"`
	// TotalSuccessAmount 累计成功金额
	TotalSuccessAmount decimal.Decimal `gorm:"column:total_success_amount;type:decimal(30,2);not null;default:0.00" json:"total_success_amount"`
	// TotalFailAmount 累计失败金额
	TotalFailAmount  decimal.Decimal `gorm:"column:total_fail_amount;type:decimal(30,2);not null;default:0.00" json:"total_fail_amount"`
	// TotalPlatformBonus 累计平台补贴
	TotalPlatformBonus decimal.Decimal `gorm:"column:total_platform_bonus;type:decimal(30,2);not null;default:0.00" json:"total_platform_bonus"`
	// TotalPoolAmount 累计奖池金额
	TotalPoolAmount  decimal.Decimal `gorm:"column:total_pool_amount;type:decimal(30,2);not null;default:0.00" json:"total_pool_amount"`
	// UpdatedAt 更新时间
	UpdatedAt        time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (AppChallengeTotalStat) TableName() string {
	return "app_challenge_total_stat"
}

// AppChallengeRankDaily 排行榜日快照
type AppChallengeRankDaily struct {
	ID       int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// RankDate 排行日期
	RankDate time.Time       `gorm:"column:rank_date;type:date;not null" json:"rank_date"`
	// RankType 排行类型：1邀请 2收益 3毅力
	RankType int8            `gorm:"column:rank_type;not null;default:0" json:"rank_type"`
	UserID   int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`
	// Value 排行值
	Value    decimal.Decimal `gorm:"column:value;type:decimal(30,2);not null;default:0.00" json:"value"`
	// RankNo 排名
	RankNo   int             `gorm:"column:rank_no;not null;default:0" json:"rank_no"`
}

func (AppChallengeRankDaily) TableName() string {
	return "app_challenge_rank_daily"
}
