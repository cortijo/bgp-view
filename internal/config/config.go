package config

import "os"

// Config carrega as variáveis de ambiente obrigatórias e com default.
// Segredos (BGP_API_KEY, BGP_MD5) não têm default e devem vir do .env/compose.
type Config struct {
	APIPort             string
	APIKey              string
	LocalASN            string
	Neighbor            string
	NeighborASN         string
	MD5                 string
	MaxPrefixes         string
	MaxProfiles         string
	DBPath              string
	AuditRetentionDays  string
	AppLogRetentionDays string
}

// Load monta a config a partir do ambiente.
func Load() Config {
	return Config{
		APIPort:             envOr("BGP_API_PORT", "8080"),
		APIKey:              os.Getenv("BGP_API_KEY"),
		LocalASN:            os.Getenv("BGP_LOCAL_ASN"),
		Neighbor:            os.Getenv("BGP_NEIGHBOR_ADDR"),
		NeighborASN:         os.Getenv("BGP_NEIGHBOR_ASN"),
		MD5:                 os.Getenv("BGP_MD5"),
		MaxPrefixes:         envOr("BGP_MAX_PREFIXES", "500"),
		MaxProfiles:         envOr("BGP_MAX_PROFILES", "50"),
		DBPath:              envOr("BGP_DB_PATH", "/data/bgpadvert.db"),
		AuditRetentionDays:  envOr("BGP_AUDIT_RETENTION_DAYS", "90"),
		AppLogRetentionDays: envOr("BGP_APPLOG_RETENTION_DAYS", "14"),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
