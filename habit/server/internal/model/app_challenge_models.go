package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// AppChallenge 打卡挑战活动配置
type AppChallenge struct {
	ID                   int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// 结算设置
	IsAutoSettle         bool      `gorm:"column:is_auto_settle;not null;default:1" json:"is_auto_settle"`                    // 1:自动结算, 0:手动结算
	SettleTime           string    `gorm:"column:settle_time;type:time;not null;default:'06:10:00'" json:"settle_time"`       // 每日结算时间点
	CycleDays            int       `gorm:"column:cycle_days;not null;default:1" json:"cycle_days"`                         // 挑战总天数
	StartTime            string    `gorm:"column:start_time;type:time;not null;default:'06:00:00'" json:"start_time"`         // 打卡开始时间
	EndTime              string    `gorm:"column:end_time;type:time;not null;default:'06:10:00'" json:"end_time"`             // 打卡结束时间
	// 资金设置 (单位: 分，避免精度丢失)
	MaxDepositAmount     int       `gorm:"column:max_deposit_amount;not null;default:0" json:"max_deposit_amount"`         // 挑战金上限
	MinWithdrawAmount    int       `gorm:"column:min_withdraw_amount;not null;default:0" json:"min_withdraw_amount"`        // 最低提现
	MaxDailyProfit       int       `gorm:"column:max_daily_profit;not null" json:"max_daily_profit"`                       // 个人每日最高收益上限
	ExcessTaxRate        int       `gorm:"column:excess_tax_rate;not null;default:98" json:"excess_tax_rate"`              // 超过部分扣除比例(%)
	MinDailyProfit       int       `gorm:"column:min_daily_profit;not null" json:"min_daily_profit"`                       // 个人每日最低收益
	DailyPlatformSubsidy int       `gorm:"column:daily_platform_subsidy;not null" json:"daily_platform_subsidy"`         // 每日平台补贴
	UncheckDeductRate    int       `gorm:"column:uncheck_deduct_rate;not null;default:100" json:"uncheck_deduct_rate"`    // 未打卡扣除金比例(%)
	MinUncheckUsers      int       `gorm:"column:min_uncheck_users;not null;default:2" json:"min_uncheck_users"`          // 人数不足不扣除的阈值
	// 佣金设置 (单位: 分)
	CommissionFollow     int       `gorm:"column:commission_follow;not null" json:"commission_follow"`                     // 推荐关注佣金
	CommissionJoin       int       `gorm:"column:commission_join;not null" json:"commission_join"`                         // 推荐参加挑战佣金
	CommissionL1         int       `gorm:"column:commission_l1;not null" json:"commission_l1"`                             // 一级推荐佣金
	CommissionL2         int       `gorm:"column:commission_l2;not null" json:"commission_l2"`                             // 二级推荐佣金
	CommissionL3         int       `gorm:"column:commission_l3;not null" json:"commission_l3"`                             // 三级推荐佣金
	UpdatedAt            *time.Time `gorm:"column:updated_at" json:"updated_at"`                                          // 更新时间
}

func (AppChallenge) TableName() string {
	return "app_challenge"
}

// AppUserChallenge 用户参与挑战记录
type AppUserChallenge struct {
	ID              int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 用户挑战ID
	UserID          int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`                          // 用户ID
	ChallengeID     int64           `gorm:"column:challenge_id;not null;default:0" json:"challenge_id"`                  // 活动配置ID
	PoolID          int64           `gorm:"column:pool_id;not null;default:0" json:"pool_id"`                          // 奖池ID
	// ChallengeAmount 用户挑战金额
	ChallengeAmount decimal.Decimal `gorm:"column:challenge_amount;type:decimal(30,2);not null;default:0.00" json:"challenge_amount"` // 用户挑战金额
	// PreRecharge 用户预充值挑战金额
	PreRecharge     decimal.Decimal `gorm:"column:pre_recharge;type:decimal(30,2);not null;default:0.00" json:"pre_recharge"`     // 用户预充值挑战金额
	// StartDate 活动开始日期
	StartDate       time.Time       `gorm:"column:start_date;type:date;not null" json:"start_date"`                    // 活动开始日期 YYYYMMDD
	// EndDate 活动结束日期
	EndDate         time.Time       `gorm:"column:end_date;type:date;not null" json:"end_date"`                        // 活动结束日期 YYYYMMDD
	// Status 状态：1进行中 2成功 3失败
	Status          int8            `gorm:"column:status;not null;default:1" json:"status"`                            // 状态 1进行中 2成功 3失败
	// FailReason 失败原因：0无 1已打卡 2未打卡
	FailReason      int8            `gorm:"column:fail_reason;not null;default:2" json:"fail_reason"`                  // 失败原因 0无 1已打卡 2未打卡
	CreatedAt       time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`    // 报名时间
	FinishedAt      *time.Time      `gorm:"column:finished_at" json:"finished_at"`                                    // 完成时间
}

func (AppUserChallenge) TableName() string {
	return "app_user_challenge"
}

// AppUserChallengeCheckin 每日打卡记录（心情 + 内容）
type AppUserChallengeCheckin struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 打卡ID
	UserChallengeID int64      `gorm:"column:user_challenge_id;not null;default:0" json:"user_challenge_id"`        // 用户挑战ID
	UserID          int64      `gorm:"column:user_id;not null;default:0" json:"user_id"`                          // 用户ID
	// CheckinDate 打卡日期
	CheckinDate     time.Time  `gorm:"column:checkin_date;type:date;not null" json:"checkin_date"`                 // 打卡日期 YYYYMMDD
	// CheckinTime 打卡时间
	CheckinTime     *time.Time `gorm:"column:checkin_time" json:"checkin_time"`                                  // 打卡时间
	// MoodCode 心情枚举：1开心 2平静 3一般 4疲惫 5低落 6爆棚
	MoodCode        int8       `gorm:"column:mood_code;not null;default:0" json:"mood_code"`                     // 心情枚举 1开心 2平静 3一般 4疲惫 5低落 6爆棚
	// MoodText 心情文字描述
	MoodText        string     `gorm:"column:mood_text;type:text;not null" json:"mood_text"`                     // 用户心情文字描述（最多200字）
	// ContentType 内容类型：1图片 2视频广告
	ContentType     int8       `gorm:"column:content_type;not null;default:1" json:"content_type"`                 // 打卡内容类型 1图片 2视频广告
	// Status 状态：1打卡成功 2未打卡 3打卡失败
	Status          int8       `gorm:"column:status;not null;default:2" json:"status"`                           // 状态 1打卡成功 2未打卡 3打卡失败
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`    // 记录创建时间
}

func (AppUserChallengeCheckin) TableName() string {
	return "app_user_challenge_checkin"
}

// AppChallengeCheckinImage 打卡图片表
type AppChallengeCheckinImage struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 图片ID
	CheckinID int64     `gorm:"column:checkin_id;not null;default:0" json:"checkin_id"`                    // 打卡ID
	UserID    int64     `gorm:"column:user_id;not null;default:0" json:"user_id"`                          // 用户ID
	// ImageURL 图片URL
	ImageURL  string    `gorm:"column:image_url;size:500;not null;default:''" json:"image_url"`            // 图片URL
	// ImageHash 图片Hash（防重复）
	ImageHash string    `gorm:"column:image_hash;size:64;not null;default:''" json:"image_hash"`            // 图片Hash（防重复）
	// SortNo 图片顺序
	SortNo    int8      `gorm:"column:sort_no;not null;default:1" json:"sort_no"`                          // 图片顺序
	// Status 状态：1正常 2屏蔽 3审核中
	Status    int8      `gorm:"column:status;not null;default:1" json:"status"`                            // 状态 1正常 2屏蔽 3审核中
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`    // 上传时间
}

func (AppChallengeCheckinImage) TableName() string {
	return "app_challenge_checkin_image"
}

// AppChallengeCheckinVideoAd 视频广告打卡记录
type AppChallengeCheckinVideoAd struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 视频广告打卡ID
	CheckinID    int64           `gorm:"column:checkin_id;not null;default:0" json:"checkin_id"`                    // 打卡ID
	UserID       int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`                          // 用户ID
	// AdPlatform 广告平台，如：csj、gdt、unity
	AdPlatform   string          `gorm:"column:ad_platform;size:50;not null;default:''" json:"ad_platform"`        // 广告平台，如：csj、gdt、unity
	// AdUnitID 广告位ID
	AdUnitID     string          `gorm:"column:ad_unit_id;size:100;not null;default:''" json:"ad_unit_id"`          // 广告位ID
	// AdOrderNo 广告联盟订单号（唯一）
	AdOrderNo    string          `gorm:"column:ad_order_no;size:100;not null;default:''" json:"ad_order_no"`        // 广告联盟返回的订单号（唯一）
	// VideoDuration 视频时长（秒）
	VideoDuration int            `gorm:"column:video_duration;not null;default:0" json:"video_duration"`             // 视频时长（秒）
	// WatchDuration 实际观看时长（秒）
	WatchDuration int            `gorm:"column:watch_duration;not null;default:0" json:"watch_duration"`             // 实际观看时长（秒）
	// RewardAmount 该广告产生的收益
	RewardAmount decimal.Decimal `gorm:"column:reward_amount;type:decimal(30,2);not null;default:0.00" json:"reward_amount"` // 该广告产生的收益
	// VerifyStatus 校验状态：0待校验 1成功 2失败
	VerifyStatus int8            `gorm:"column:verify_status;not null;default:0" json:"verify_status"`               // 校验状态：0待校验 1成功 2失败
	CreatedAt    time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`      // 创建时间
	VerifiedAt   *time.Time      `gorm:"column:verified_at" json:"verified_at"`                                    // 校验时间
}

func (AppChallengeCheckinVideoAd) TableName() string {
	return "app_challenge_checkin_video_ad"
}

// AppUserChallengeSettlement 用户每日挑战收益结算表
type AppUserChallengeSettlement struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 主键ID
	UserChallengeID int64          `gorm:"column:user_challenge_id;not null;default:0" json:"user_challenge_id"`        // 关联用户挑战ID
	UserID         int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`                          // 用户ID
	CheckinID      int64           `gorm:"column:checkin_id;not null;default:0" json:"checkin_id"`                    // 关联打卡ID
	SettleDate     time.Time       `gorm:"column:settle_date;type:date;not null" json:"settle_date"`                 // 结算日期
	// 资金明细 (统一使用 DECIMAL(30,2))
	BaseProfit     decimal.Decimal `gorm:"column:base_profit;type:decimal(30,2);not null;default:0.00" json:"base_profit"`    // 分到的原始金额(未扣除前)
	PlatformSubsidy decimal.Decimal `gorm:"column:platform_subsidy;type:decimal(30,2);not null;default:0.00" json:"platform_subsidy"` // 分到的平台补贴
	TotalRawProfit decimal.Decimal `gorm:"column:total_raw_profit;type:decimal(30,2);not null;default:0.00" json:"total_raw_profit"` // 扣除前总收益
	TaxDeduction   decimal.Decimal `gorm:"column:tax_deduction;type:decimal(30,2);not null;default:0.00" json:"tax_deduction"` // 超过上限扣除的金额
	FinalProfit    decimal.Decimal `gorm:"column:final_profit;type:decimal(30,2);not null;default:0.00" json:"final_profit"`   // 实际到账收益
	// 状态记录
	IsSettled      int8            `gorm:"column:is_settled;not null;default:1" json:"is_settled"`                     // 结算状态 1已结算 2待审核 3异常
	SettleAt       time.Time       `gorm:"column:settle_at;not null;default:CURRENT_TIMESTAMP" json:"settle_at"`        // 结算执行时间
}

func (AppUserChallengeSettlement) TableName() string {
	return "app_user_challenge_settlement"
}

// AppChallengePoolDaily 全平台奖池日结算表
type AppChallengePoolDaily struct {
	ID                    int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 主键ID
	PoolDate              time.Time       `gorm:"column:pool_date;type:date;not null;uniqueIndex:uk_date" json:"pool_date"`    // 资金归属日期
	TotalUsers            int             `gorm:"column:total_users;not null;default:0" json:"total_users"`                    // 当日参与总人数
	SuccessUsers          int             `gorm:"column:success_users;not null;default:0" json:"success_users"`                // 当日成功打卡人数
	FailUsers             int             `gorm:"column:fail_users;not null;default:0" json:"fail_users"`                      // 当日未打卡人数
	// 资金汇总 (统一使用 DECIMAL(30,2))
	FailDeductPool        decimal.Decimal `gorm:"column:fail_deduct_pool;type:decimal(30,2);not null;default:0.00" json:"fail_deduct_pool"` // 未打卡用户扣除的总金额
	PlatformSubsidyPool   decimal.Decimal `gorm:"column:platform_subsidy_pool;type:decimal(30,2);not null;default:0.00" json:"platform_subsidy_pool"` // 当日平台总投入补贴
	TotalDistributable    decimal.Decimal `gorm:"column:total_distributable;type:decimal(30,2);not null;default:0.00" json:"total_distributable"` // 当日可分配总奖金
	AvgProfit             decimal.Decimal `gorm:"column:avg_profit;type:decimal(30,2);not null;default:0.00" json:"avg_profit"` // 单人理论应分金额
	SystemTaxIncome       decimal.Decimal `gorm:"column:system_tax_income;type:decimal(30,2);not null;default:0.00" json:"system_tax_income"` // 系统回收金额(阶梯扣除部分)
	CreatedAt             time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`      // 创建时间
}

func (AppChallengePoolDaily) TableName() string {
	return "app_challenge_pool_daily"
}

// AppChallengeDailyStat 每日运营统计
type AppChallengeDailyStat struct {
	// StatDate 统计日期
	StatDate      time.Time       `gorm:"column:stat_date;type:date;primaryKey" json:"stat_date"`                    // 统计日期
	// JoinUserCnt 参与人数
	JoinUserCnt   int             `gorm:"column:join_user_cnt;not null;default:0" json:"join_user_cnt"`                // 参与人数
	// SuccessUserCnt 成功人数
	SuccessUserCnt int            `gorm:"column:success_user_cnt;not null;default:0" json:"success_user_cnt"`            // 成功人数
	// FailUserCnt 失败人数
	FailUserCnt   int             `gorm:"column:fail_user_cnt;not null;default:0" json:"fail_user_cnt"`                  // 失败人数
	// JoinAmount 参与总金额
	JoinAmount    decimal.Decimal `gorm:"column:join_amount;type:decimal(30,2);not null;default:0.00" json:"join_amount"`    // 参与总金额
	// SuccessAmount 成功金额
	SuccessAmount decimal.Decimal `gorm:"column:success_amount;type:decimal(30,2);not null;default:0.00" json:"success_amount"` // 成功金额
	// FailAmount 失败金额
	FailAmount    decimal.Decimal `gorm:"column:fail_amount;type:decimal(30,2);not null;default:0.00" json:"fail_amount"`      // 失败金额
	// PlatformBonus 平台补贴
	PlatformBonus decimal.Decimal `gorm:"column:platform_bonus;type:decimal(30,2);not null;default:0.00" json:"platform_bonus"` // 平台补贴
	// PoolAmount 奖池金额
	PoolAmount    decimal.Decimal `gorm:"column:pool_amount;type:decimal(30,2);not null;default:0.00" json:"pool_amount"`      // 奖池金额
	CreatedAt     time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`          // 创建时间
	UpdatedAt     time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`          // 更新时间
}

func (AppChallengeDailyStat) TableName() string {
	return "app_challenge_daily_stat"
}

// AppChallengeTotalStat 平台累计统计
type AppChallengeTotalStat struct {
	// ID 固定主键（通常为 1）
	ID               int8            `gorm:"column:id;primaryKey;not null;default:1" json:"id"`                      // 固定主键（通常为 1）
	// TotalUserCnt 累计用户数
	TotalUserCnt     int             `gorm:"column:total_user_cnt;not null;default:0" json:"total_user_cnt"`            // 累计用户数
	// TotalJoinCnt 累计参与人次
	TotalJoinCnt     int             `gorm:"column:total_join_cnt;not null;default:0" json:"total_join_cnt"`              // 累计参与人次
	// TotalSuccessCnt 累计成功人次
	TotalSuccessCnt  int             `gorm:"column:total_success_cnt;not null;default:0" json:"total_success_cnt"`        // 累计成功人次
	// TotalFailCnt 累计失败人次
	TotalFailCnt     int             `gorm:"column:total_fail_cnt;not null;default:0" json:"total_fail_cnt"`              // 累计失败人次
	// TotalJoinAmount 累计参与金额
	TotalJoinAmount  decimal.Decimal `gorm:"column:total_join_amount;type:decimal(30,2);not null;default:0.00" json:"total_join_amount"` // 累计参与金额
	// TotalSuccessAmount 累计成功金额
	TotalSuccessAmount decimal.Decimal `gorm:"column:total_success_amount;type:decimal(30,2);not null;default:0.00" json:"total_success_amount"` // 累计成功金额
	// TotalFailAmount 累计失败金额
	TotalFailAmount  decimal.Decimal `gorm:"column:total_fail_amount;type:decimal(30,2);not null;default:0.00" json:"total_fail_amount"` // 累计失败金额
	// TotalPlatformBonus 累计平台补贴
	TotalPlatformBonus decimal.Decimal `gorm:"column:total_platform_bonus;type:decimal(30,2);not null;default:0.00" json:"total_platform_bonus"` // 累计平台补贴
	// TotalPoolAmount 累计奖池金额
	TotalPoolAmount  decimal.Decimal `gorm:"column:total_pool_amount;type:decimal(30,2);not null;default:0.00" json:"total_pool_amount"` // 累计奖池金额
	// UpdatedAt 更新时间
	UpdatedAt        time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`          // 更新时间
}

func (AppChallengeTotalStat) TableName() string {
	return "app_challenge_total_stat"
}

// AppChallengeRankDaily 排行榜日快照
type AppChallengeRankDaily struct {
	ID       int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 排行榜ID
	// RankDate 排行日期
	RankDate time.Time       `gorm:"column:rank_date;type:date;not null" json:"rank_date"`                      // 排行日期
	// RankType 排行类型：1邀请 2收益 3毅力
	RankType int8            `gorm:"column:rank_type;not null;default:0" json:"rank_type"`                     // 排行类型：1邀请 2收益 3毅力
	UserID   int64           `gorm:"column:user_id;not null;default:0" json:"user_id"`                          // 用户ID
	// Value 排行值
	Value    decimal.Decimal `gorm:"column:value;type:decimal(30,2);not null;default:0.00" json:"value"`           // 排行值
	// RankNo 排名
	RankNo   int             `gorm:"column:rank_no;not null;default:0" json:"rank_no"`                          // 排名
}

func (AppChallengeRankDaily) TableName() string {
	return "app_challenge_rank_daily"
}
