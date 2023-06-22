// Copyright Contributors to the Open Cluster Management project

package config

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strconv"

	klog "k8s.io/klog/v2"
)

var Cfg = new()

// Define a Config type to hold our config properties.
type Config struct {
	API_SERVER_URL      string // address for Kubernetes API Server
	AuthCacheTTL        int    // Time-to-live (milliseconds) of Authentication (TokenReview) cache.
	SharedCacheTTL      int    // Time-to-live (milliseconds) of common resources (shared across users) cache.
	UserCacheTTL        int    // Time-to-live (milliseconds) of namespaced resources (specifc to users) cache.
	ContextPath         string
	DBHost              string
	DBMinConns          int // Overrides pgxpool.Config{ MinConns } Default: 0
	DBMaxConns          int // Overrides pgxpool.Config{ MaxConns } Default: 10
	DBMaxConnIdleTime   int // Overrides pgxpool.Config{ MaxConnIdleTime } Default: 30 min
	DBMaxConnLifeTime   int // Overrides pgxpool.Config{ MaxConnLifetime } Default: 60 min
	DBMaxConnLifeJitter int // Overrides pgxpool.Config{ MaxConnLifetimeJitter } Default: 2 min
	DBName              string
	DBPass              string
	DBPort              int
	DBUser              string
	HttpPort            int
	PlaygroundMode      bool // Enable the GraphQL Playground client.
	QueryLimit          int  // The default LIMIT to use on queries. Client can override.
	RelationLevel       int  // The number of levels/hops for finding relationships for a particular resource
	SlowLog             int  // Logs when queries are slower than the specified time duration in ms. Default 300ms
	// Placeholder for future use.
	// QueryLoopLimit          int // number of queries handled at a time
	// RBAC_INACTIVITY_TIMEOUT int
}

func new() *Config {
	// If environment variables are set, use default values
	// Simply put, the order of preference is env -> default values (from left to right)
	conf := &Config{
		API_SERVER_URL: getEnv("API_SERVER_URL", "https://kubernetes.default.svc"),
		AuthCacheTTL:   getEnvAsInt("AUTH_CACHE_TTL", 60000),    // 1 minute
		SharedCacheTTL: getEnvAsInt("SHARED_CACHE_TTL", 300000), // 5 min (increase to 10min after implementation)
		UserCacheTTL:   getEnvAsInt("USER_CACHE_TTL", 300000),   // 5 min (increase to 10min after implementation)
		ContextPath:    getEnv("CONTEXT_PATH", "/searchapi"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		// Postgres has 100 conns by default. Using 20 allows scaling indexer and api.
		DBMaxConns:          getEnvAsInt("DB_MAX_CONNS", 10),                   // 10 - Overrides pgxpool default
		DBMaxConnIdleTime:   getEnvAsInt("DB_MAX_CONN_IDLE_TIME", 30*60*1000),  // 30 min - Default for pgxpool.Config
		DBMaxConnLifeJitter: getEnvAsInt("DB_MAX_CONN_LIFE_JITTER", 2*60*1000), // 2 min - Overrides pgxpool default
		DBMaxConnLifeTime:   getEnvAsInt("DB_MAX_CONN_LIFE_TIME", 60*60*1000),  // 60 min - Default for pgxpool.Config
		DBMinConns:          getEnvAsInt("DB_MIN_CONNS", 0),                    // Default for pgxpool.Config
		DBName:              getEnv("DB_NAME", ""),
		DBPass:              getEnv("DB_PASS", ""),
		DBPort:              getEnvAsInt("DB_PORT", 5432),
		DBUser:              getEnv("DB_USER", ""),
		HttpPort:            getEnvAsInt("HTTP_PORT", 4010),
		PlaygroundMode:      getEnvAsBool("PLAYGROUND_MODE", false),
		QueryLimit:          getEnvAsInt("QUERY_LIMIT", 1000),
		SlowLog:             getEnvAsInt("SLOW_LOG", 300),
		// Setting default level to 0 to check if user has explicitly set this variable
		// This will be updated to 1 for default searches and 3 for applications - unless set by the user
		RelationLevel: getEnvAsInt("RELATION_LEVEL", 0),
		// Placeholder for future use.
		// QueryLoopLimit:          getEnvAsInt("QUERY_LOOP_LIMIT", 5000),
		// RBAC_INACTIVITY_TIMEOUT: getEnvAsInt("RBAC_INACTIVITY_TIMEOUT", 600000), // 10 minutes
	}
	conf.DBPass = url.QueryEscape(conf.DBPass)
	return conf
}

// Format and print environment to logger.
func (cfg *Config) PrintConfig() {
	// Make a copy to redact secrets and sensitive information.
	tmp := *cfg
	tmp.DBPass = "[REDACTED]"

	// Convert to JSON for nicer formatting.
	cfgJSON, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		klog.Warning("Encountered a problem formatting configuration. ", err)
		klog.Infof("Configuration %#v\n", tmp)
	}
	klog.Infof("Using configuration:\n%s\n", string(cfgJSON))
}

// Validate required configuration.
func (cfg *Config) Validate() error {
	if cfg.DBName == "" {
		return errors.New("required environment DB_NAME is not set")
	}
	if cfg.DBUser == "" {
		return errors.New("required environment DB_USER is not set")
	}
	if cfg.DBPass == "" {
		return errors.New("required environment DB_PASS is not set")
	}
	return nil
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
