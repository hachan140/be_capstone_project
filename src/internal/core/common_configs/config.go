package common_configs

// Store defines repository relevant info
type Store struct {
	Driver          string `envconfig:"DB_DRIVER" default:"postgresql" `
	Host            string `envconfig:"DB_HOST"`
	Port            int64  `envconfig:"DB_PORT"`
	User            string `envconfig:"DB_USERNAME"`
	Password        string `envconfig:"DB_PASSWORD"`
	Db              string `envconfig:"DB_SCHEMA_NAME"`
	IsDebug         bool   `envconfig:"DB_IS_DEBUG"`
	MaxOpenConns    int64  `envconfig:"DB_MAX_OPEN_CONNS" default:"500" `
	MaxIdleConns    int64  `envconfig:"DB_MAX_IDLE_CONNS" envconfig:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifeTime int64  `envconfig:"DB_CONN_MAX_LIFE_TIME" default:"5" `
}

type App struct {
	Name        string `envconfig:"APP_NAME"`
	Host        string ` default:"127.0.0.1" envconfig:"APP_HOST"`
	Port        int64  `default:"8080" envconfig:"APP_PORT"`
	Environment string `default:"dev" envconfig:"APP_ENV"`
}

type Redis struct {
	Host string `yaml:"host" envconfig:"REDIS_URL"`
}

type KafkaSaramaConfig struct {
	TopicUserAuthV2                       string `envconfig:"KAFKA_TOPIC_USER_AUTHV2" required:"false"`
	TopicUserAttributeSetupRecoveryMethod string `envconfig:"KAFKA_TOPIC_USER_ATTRIBUTE_SETUP_RECOVERY_METHOD" required:"false"`
	TopicUserSignInSignUp                 string `envconfig:"KAFKA_TOPIC_USER_SIGNIN_SIGNUP" required:"false"`
	TopicUserProfileUpdated               string `envconfig:"KAFKA_TOPIC_USER_PROFILE_UPDATED" required:"false"`
	TopicUserProfilePhoneUpdated          string `envconfig:"KAFKA_TOPIC_USER_PROFILE_PHONE_NUMBER_UPDATED" required:"false"`
	TopicUserAuth                         string `envconfig:"KAFKA_TOPIC_USER_AUTH" required:"false"`

	GroupID            string   `envconfig:"KAFKA_GROUP_ID" required:"false"`
	Brokers            []string `envconfig:"KAFKA_BROKERS" required:"true"`
	MaxRetry           int      `envconfig:"KAFKA_MAX_RETRY" required:"true"`
	EnableTLS          bool     `envconfig:"KAFKA_ENABLE_TLS" required:"true"`
	InsecureSkipVerify bool     `envconfig:"KAFKA_INSECURE_SKIP_VERIFY" required:"true"`
	AuthSDKCACertFile  string   `envconfig:"KAFKA_AUTHEN_SDK_CA_CERT_FILE"`
	AuthSDKCertFile    string   `envconfig:"KAFKA_AUTHEN_SDK_CERT_FILE"`
	AuthSDKKeyFile     string   `envconfig:"KAFKA_AUTHEN_SDK_KEY_FILE"`
}
