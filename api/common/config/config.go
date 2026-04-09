// 鏂囦欢閰嶇疆,瑙ｆ瀽yaml閰嶇鏂囦欢
// author xiaoRui

package config

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// 鎬婚厤鏂囦欢
type config struct {
	Server        server        `yaml:"server"`
	Db            db            `yaml:"db"`
	Redis         redis         `yaml:"redis"`
	ImageSettings imageSettings `yaml:"imageSettings"`
	Log           log           `yaml:"log"`
	Monitor       monitor       `yaml:"monitor"`
	AI            aiRuntime     `yaml:"ai"`
}

// 鐩戞帶閰嶇疆
type monitor struct {
	Prometheus  prometheus  `yaml:"prometheus"`
	Pushgateway pushgateway `yaml:"pushgateway"`
	Agent       agent       `yaml:"agent"`
	Webhook     webhook     `yaml:"webhook"`
}

// Pushgateway閰嶇疆
type pushgateway struct {
	URL string `yaml:"url"`
}

// Agent閰嶇疆
type agent struct {
	HeartbeatServerURL string `yaml:"heartbeat_server_url"`
	HeartbeatToken     string `yaml:"heartbeat_token"`
}

// Webhook閰嶇疆
type webhook struct {
	Token string `yaml:"token"`
}

// Prometheus閰嶇疆
type prometheus struct {
	URL string `yaml:"url"`
}

// 椤圭洰绔彛閰嶇疆
type server struct {
	Address       string `yaml:"address"`
	Model         string `yaml:"model"`
	EnableSwagger bool   `yaml:"enableSwagger"` // 鏄惁鍚敤Swagger鏂囨。
}

// 鏁版嵁搴撻厤缃?
type db struct {
	Dialects string `yaml:"dialects"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

// redis閰嶇疆
type redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

// imageSettings鍥剧墖涓婁紶閰嶇疆
type imageSettings struct {
	UploadDir string `yaml:"uploadDir"`
	ImageHost string `yaml:"imageHost"`
}

// log鏃ュ織閰嶇疆
type log struct {
	Path  string `yaml:"path"`
	Name  string `yaml:"name"`
	Model string `yaml:"model"`
}

type aiRuntime struct {
	Enabled         bool   `yaml:"enabled"`
	Provider        string `yaml:"provider"`
	BaseURL         string `yaml:"baseUrl"`
	APIKey          string `yaml:"apiKey"`
	Model           string `yaml:"model"`
	ReasoningEffort string `yaml:"reasoningEffort"`
	TimeoutSeconds  int    `yaml:"timeoutSeconds"`
}

var Config *config

// 閰嶇疆鍒濆鍖?
func init() {
	// 鍒濆鍖栨椂鍏堜笉鍔犺浇閰嶇疆鏂囦欢锛岀瓑寰匧oadConfig()琚皟鐢?
}

// LoadConfig 浠庢寚瀹氳矾寰勫姞杞介厤缃枃浠?
func LoadConfig(configPath string) error {
	if configPath == "" {
		configPath = "./config.yaml" // 榛樿閰嶇疆鏂囦欢璺緞
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// 缁戝畾鍊?
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		return err
	}

	overrideFromEnv(Config)
	return nil
}

// GetConfig 鑾峰彇鏁版嵁搴撻厤缃?
func GetConfig() *db {
	if Config == nil {
		panic("Config is not initialized")
	}
	return &Config.Db
}

// GetRedisConfig 鑾峰彇Redis閰嶇疆
func GetRedisConfig() *redis {
	if Config == nil {
		panic("Config is not initialized")
	}
	return &Config.Redis
}

// Setup 鍒濆鍖栭厤缃紙涓轰簡鍏煎migrate.go鐨勮皟鐢級
func Setup() {
	// 閰嶇疆宸茬粡鍦╥nit()鏂规硶涓垵濮嬪寲浜嗭紝杩欓噷鍙槸鎻愪緵涓€涓吋瀹规€ф柟娉?
	if Config == nil {
		panic("Config initialization failed")
	}
}

func overrideFromEnv(cfg *config) {
	if cfg == nil {
		return
	}

	overrideString(&cfg.Server.Address, "SERVER_ADDRESS")
	overrideString(&cfg.Server.Model, "SERVER_MODEL")
	overrideBool(&cfg.Server.EnableSwagger, "SERVER_ENABLE_SWAGGER")

	overrideString(&cfg.Db.Dialects, "DB_DIALECTS")
	overrideString(&cfg.Db.Host, "DB_HOST")
	overrideInt(&cfg.Db.Port, "DB_PORT")
	overrideString(&cfg.Db.Db, "DB_NAME")
	overrideString(&cfg.Db.Username, "DB_USER")
	overrideString(&cfg.Db.Password, "DB_PASSWORD")
	overrideString(&cfg.Db.Charset, "DB_CHARSET")
	overrideInt(&cfg.Db.MaxIdle, "DB_MAX_IDLE")
	overrideInt(&cfg.Db.MaxOpen, "DB_MAX_OPEN")

	overrideString(&cfg.Redis.Address, "REDIS_ADDR")
	overrideString(&cfg.Redis.Password, "REDIS_PASSWORD")

	overrideString(&cfg.ImageSettings.UploadDir, "UPLOAD_DIR")
	overrideString(&cfg.ImageSettings.ImageHost, "IMAGE_HOST")

	overrideString(&cfg.Log.Path, "LOG_PATH")
	overrideString(&cfg.Log.Name, "LOG_NAME")
	overrideString(&cfg.Log.Model, "LOG_MODEL")

	overrideString(&cfg.Monitor.Prometheus.URL, "MONITOR_PROMETHEUS_URL")
	overrideString(&cfg.Monitor.Pushgateway.URL, "MONITOR_PUSHGATEWAY_URL")
	overrideString(&cfg.Monitor.Agent.HeartbeatServerURL, "MONITOR_AGENT_HEARTBEAT_SERVER_URL")
	overrideString(&cfg.Monitor.Agent.HeartbeatToken, "MONITOR_AGENT_HEARTBEAT_TOKEN")
	overrideString(&cfg.Monitor.Webhook.Token, "MONITOR_WEBHOOK_TOKEN")

	overrideBool(&cfg.AI.Enabled, "AI_ENABLED")
	overrideString(&cfg.AI.Provider, "AI_PROVIDER")
	overrideString(&cfg.AI.BaseURL, "AI_BASE_URL")
	overrideString(&cfg.AI.APIKey, "AI_API_KEY")
	overrideString(&cfg.AI.Model, "AI_MODEL")
	overrideString(&cfg.AI.ReasoningEffort, "AI_REASONING_EFFORT")
	overrideInt(&cfg.AI.TimeoutSeconds, "AI_TIMEOUT_SECONDS")
}

func overrideString(target *string, envKey string) {
	if target == nil {
		return
	}

	if value, ok := os.LookupEnv(envKey); ok && strings.TrimSpace(value) != "" {
		*target = value
	}
}

func overrideInt(target *int, envKey string) {
	if target == nil {
		return
	}

	value, ok := os.LookupEnv(envKey)
	if !ok {
		return
	}

	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err == nil {
		*target = parsed
	}
}

func overrideBool(target *bool, envKey string) {
	if target == nil {
		return
	}

	value, ok := os.LookupEnv(envKey)
	if !ok {
		return
	}

	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err == nil {
		*target = parsed
	}
}
