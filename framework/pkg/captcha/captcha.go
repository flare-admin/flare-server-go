package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store = base64Captcha.DefaultMemStore

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	// 通用配置
	Width      int64       // 宽度
	Height     int64       // 高度
	NoiseCount int         // 噪声数量
	ShowLine   bool        // 是否显示干扰线
	BgColor    *color.RGBA // 背景颜色
	FontSize   int         // 字体大小
	FontStyle  string      // 字体样式

	// 数字验证码配置
	DigitLength int     // 数字长度
	DigitNoise  float64 // 数字噪声程度

	// 算术验证码配置
	MathLength int    // 算术表达式长度
	MathType   string // 算术类型：add(加法), sub(减法), mul(乘法), div(除法)

	// 字符串验证码配置
	StringLength  int    // 字符串长度
	StringSource  string // 字符串来源：number(数字), letter(字母), mixed(混合), alphanumeric(字母数字混合)
	CaseSensitive bool   // 是否区分大小写
}

// DefaultConfig 返回默认配置
func DefaultConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Width:         200,
		Height:        60,
		NoiseCount:    2,
		ShowLine:      true,
		BgColor:       &color.RGBA{R: 99, G: 253, B: 124, A: 100},
		FontSize:      40,
		FontStyle:     "default",
		DigitLength:   6,
		DigitNoise:    0.7,
		MathLength:    2,
		MathType:      "add",
		StringLength:  4,
		StringSource:  "alphanumeric",
		CaseSensitive: true,
	}
}

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, string, string, error) {
	// Configure captcha parameters
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)

	// Generate the captcha
	id, b64s, a, err := c.Generate()
	if err != nil {
		return "", "", "", err
	}
	return id, b64s, a, nil
}

// GetMathCaptcha create return id, b64s, err
func GetMathCaptcha(width, height int64) (string, string, string, error) {
	if width <= 0 {
		width = 200
	}
	if height <= 0 {
		height = 60
	}
	// 配置算术验证码
	driver := base64Captcha.NewDriverMath(
		int(height),                        // 高度
		int(width),                         // 宽度
		2,                                  // 噪声数量
		base64Captcha.OptionShowHollowLine, // 干扰线选项
		&color.RGBA{R: 99, G: 253, B: 124, A: 100}, // 背景颜色
		nil, // 使用默认字体存储
		nil, // 使用默认字体
	)
	// 生成验证码实例
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验码
	return captcha.Generate()
}

// GetDigitCaptcha create return id, b64s, err
func GetDigitCaptcha(width, height, size int64) (string, string, string, error) {
	if width <= 0 {
		width = 200
	}
	if height <= 0 {
		height = 60
	}
	if width < 120 {
		width = 120
	}
	if height < 32 {
		height = 32
	}
	// 配置算术验证码
	driver := base64Captcha.NewDriverDigit(int(height), int(width), int(size), 0.1, 80)
	// 生成验证码实例
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验码
	return captcha.Generate()
}

// GenerateCustomCaptcha 生成自定义验证码
func GenerateCustomCaptcha(config *CaptchaConfig) (string, string, string, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if config.Width <= 0 || config.Height <= 0 {
		def := DefaultConfig()
		if config.Width <= 0 {
			config.Width = def.Width
		}
		if config.Height <= 0 {
			config.Height = def.Height
		}
	}
	var driver base64Captcha.Driver
	switch config.MathType {
	case "add", "sub", "mul", "div":
		// 算术验证码
		options := base64Captcha.OptionShowHollowLine
		if !config.ShowLine {
			options = 0
		}
		driver = base64Captcha.NewDriverMath(
			int(config.Height),
			int(config.Width),
			config.NoiseCount,
			options,
			config.BgColor,
			nil,
			nil,
		)
	case "digit":
		// 数字验证码
		driver = base64Captcha.NewDriverDigit(
			int(config.Height),
			int(config.Width),
			config.DigitLength,
			config.DigitNoise,
			config.FontSize,
		)
	default:
		// 字符串验证码
		var source string
		switch config.StringSource {
		case "number":
			source = "0123456789"
		case "letter":
			if config.CaseSensitive {
				source = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			} else {
				source = "abcdefghijklmnopqrstuvwxyz"
			}
		case "alphanumeric":
			if config.CaseSensitive {
				source = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			} else {
				source = "0123456789abcdefghijklmnopqrstuvwxyz"
			}
		default:
			if config.CaseSensitive {
				source = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			} else {
				source = "0123456789abcdefghijklmnopqrstuvwxyz"
			}
		}
		showLineOpt := 0
		if config.ShowLine {
			showLineOpt = base64Captcha.OptionShowHollowLine
		}
		driver = base64Captcha.NewDriverString(
			int(config.Height),
			int(config.Width),
			config.NoiseCount,
			showLineOpt,
			config.StringLength,
			source,
			config.BgColor,
			nil,
			nil,
		)
	}

	captcha := base64Captcha.NewCaptcha(driver, store)
	return captcha.Generate()
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(captchaId, value string) bool {
	// Verify the captcha
	if base64Captcha.DefaultMemStore.Verify(captchaId, value, true) {
		return true
	}
	return false
}
