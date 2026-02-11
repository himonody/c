package response

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// Language codes
const (
	LangZhCN = "zh-CN" // 简体中文
	LangEN   = "en"    // English
	LangMS   = "ms"    // Bahasa Melayu (马来语)
)

// GetLang 从请求头获取语言
func GetLang(c *fiber.Ctx) string {
	// 优先从 X-Lang 头获取
	lang := c.Get("X-Lang")
	if lang == "" {
		// 从 Accept-Language 获取
		lang = c.Get("Accept-Language")
	}

	// 解析语言代码
	if lang != "" {
		// 处理 Accept-Language 格式: zh-CN,zh;q=0.9,en;q=0.8
		parts := strings.Split(lang, ",")
		if len(parts) > 0 {
			lang = strings.TrimSpace(strings.Split(parts[0], ";")[0])
		}
	}

	// 标准化语言代码
	lang = strings.ToLower(lang)
	switch {
	case strings.HasPrefix(lang, "zh"):
		return LangZhCN
	case strings.HasPrefix(lang, "ms"):
		return LangMS
	case strings.HasPrefix(lang, "en"):
		return LangEN
	default:
		return LangEN // 默认英语
	}
}

// Messages 多语言消息映射
var messages = map[string]map[string]string{
	// 简体中文
	LangZhCN: {
		"success":                    "操作成功",
		"error":                      "操作失败",
		"internal_error":             "服务器内部错误",
		"bad_request":                "请求参数错误",
		"unauthorized":               "未授权，请先登录",
		"forbidden":                  "没有权限访问",
		"not_found":                  "资源不存在",
		"validation_error":           "数据验证失败",
		"too_many_requests":          "请求过于频繁，请稍后再试",

		// 用户认证相关
		"register_success":           "注册成功",
		"login_success":              "登录成功",
		"logout_success":             "退出成功",
		"username_exists":            "用户名已存在",
		"invalid_credentials":        "用户名或密码错误",
		"account_disabled":           "账号已被禁用",
		"token_expired":              "登录已过期，请重新登录",
		"token_invalid":              "无效的登录凭证",
		"token_revoked":              "登录凭证已失效",
		"token_replaced":             "账号在其他地方登录",

		// 数据验证相关
		"username_invalid":           "用户名格式不正确（6-12位字母、数字或特殊字符）",
		"password_invalid":           "密码格式不正确（6-12位字母、数字或特殊字符）",
		"invalid_friend_code":        "无效的邀请码",
	},

	// English
	LangEN: {
		"success":                    "Success",
		"error":                      "Error",
		"internal_error":             "Internal server error",
		"bad_request":                "Bad request",
		"unauthorized":               "Unauthorized, please login",
		"forbidden":                  "Access forbidden",
		"not_found":                  "Resource not found",
		"validation_error":           "Validation error",
		"too_many_requests":          "Too many requests, please try again later",

		// Authentication
		"register_success":           "Registration successful",
		"login_success":              "Login successful",
		"logout_success":             "Logout successful",
		"username_exists":            "Username already exists",
		"invalid_credentials":        "Invalid username or password",
		"account_disabled":           "Account is disabled",
		"token_expired":              "Session expired, please login again",
		"token_invalid":              "Invalid token",
		"token_revoked":              "Token has been revoked",
		"token_replaced":             "Account logged in elsewhere",

		// Validation
		"username_invalid":           "Invalid username format (6-12 alphanumeric or special characters)",
		"password_invalid":           "Invalid password format (6-12 alphanumeric or special characters)",
		"invalid_friend_code":        "Invalid invite code",
	},

	// Bahasa Melayu (马来语)
	LangMS: {
		"success":                    "Berjaya",
		"error":                      "Ralat",
		"internal_error":             "Ralat pelayan dalaman",
		"bad_request":                "Permintaan tidak sah",
		"unauthorized":               "Tidak dibenarkan, sila log masuk",
		"forbidden":                  "Akses ditolak",
		"not_found":                  "Sumber tidak dijumpai",
		"validation_error":           "Ralat pengesahan",
		"too_many_requests":          "Terlalu banyak permintaan, sila cuba sebentar lagi",

		// Pengesahan
		"register_success":           "Pendaftaran berjaya",
		"login_success":              "Log masuk berjaya",
		"logout_success":             "Log keluar berjaya",
		"username_exists":            "Nama pengguna sudah wujud",
		"invalid_credentials":        "Nama pengguna atau kata laluan tidak sah",
		"account_disabled":           "Akaun telah dilumpuhkan",
		"token_expired":              "Sesi tamat tempoh, sila log masuk semula",
		"token_invalid":              "Token tidak sah",
		"token_revoked":              "Token telah dibatalkan",
		"token_replaced":             "Akaun log masuk di tempat lain",

		// Pengesahan
		"username_invalid":           "Format nama pengguna tidak sah (6-12 aksara alfanumerik atau khas)",
		"password_invalid":           "Format kata laluan tidak sah (6-12 aksara alfanumerik atau khas)",
		"invalid_friend_code":        "Kod jemputan tidak sah",
	},
}

// GetMessage 获取多语言消息
// 优先级：请求语言 -> 英语 -> 马来语 -> key本身
func GetMessage(lang, key string) string {
	if langMessages, ok := messages[lang]; ok {
		if msg, ok := langMessages[key]; ok {
			return msg
		}
	}

	// 如果找不到对应语言的消息，先尝试英文
	if langMessages, ok := messages[LangEN]; ok {
		if msg, ok := langMessages[key]; ok {
			return msg
		}
	}

	// 英文没有，尝试马来语
	if langMessages, ok := messages[LangMS]; ok {
		if msg, ok := langMessages[key]; ok {
			return msg
		}
	}

	// 都找不到，返回 key 本身
	return key
}
