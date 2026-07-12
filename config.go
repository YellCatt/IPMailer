package main



type Config struct {
	SMTPHost      string
	SMTPPort      int
	FromEmail     string
	AuthCode      string
	ToEmail       string
	TLSSkipVerify bool
}

func LoadConfig() *Config {
	cfg := &Config{
		SMTPHost:      "smtp.qq.com",
		SMTPPort:      465,
		FromEmail:     "768305875@qq.com",
		AuthCode:      "gpfruabgjebubdad",
		ToEmail:       "768305875@qq.com",
		TLSSkipVerify: true,
	}

	return cfg
}