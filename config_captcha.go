package config

type CaptchaConfig struct {
	KeyPrefix          string `mapstructure:"key_prefix" description:"键前缀"`
	ImgWidth           int    `mapstructure:"img_width" description:"验证码宽度"`
	ImgHeight          int    `mapstructure:"img_height" description:"验证码高度"`
	OpenCaptcha        int    `mapstructure:"open_captcha" description:"防爆破验证码开启次数(0表示每次都需)"`
	OpenCaptchaTimeOut int    `mapstructure:"open_captcha_timeout" description:"防爆破验证码超时时间(秒)"`
}
