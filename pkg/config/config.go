package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Gin struct {
		Mode string `yaml:"mode"`
		Log  struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		} `yaml:"log"`
	} `yaml:"gin"`
	Application struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Environment string `yaml:"environment"`
		JWTSecret   string `yaml:"jwt_secret"`
	} `yaml:"application"`
	DB struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		User         string `yaml:"user"`
		Pass         string `yaml:"pass"`
		Name         string `yaml:"name"`
		InstanceConn string `yaml:"instance_connection_name"` // For Cloud SQL
	} `yaml:"db"`
	GCS struct {
		ProfileBucket string `yaml:"profile_bucket"`
		ObjectURL     string `yaml:"object_url"`
	} `yaml:"gcs"`
	GooglePlaces struct {
		APIKey      string `yaml:"api_key"`
		APIEndpoint string `yaml:"api_endpoint"`
	} `yaml:"google_places"`
	Vertex struct {
		ProjectID string `yaml:"project_id"`
		Location  string `yaml:"location"`
		Model     string `yaml:"model"`
	} `yaml:"vertex"`
}

func LoadConfig(ctx context.Context, configDir string) (*Config, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	path := fmt.Sprintf("./config/%s.yaml", appEnv)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if appEnv == "prod" {
		if err := resolveGCPSecrets(ctx, &cfg); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

func resolveGCPSecrets(ctx context.Context, cfg *Config) error {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	defer func() {
		cerr := client.Close()
		if cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close secretmanager client: %v\n", cerr)
		}
	}()

	resolve := func(val string) (string, error) {
		if !strings.HasPrefix(val, "gcp-secret://") {
			return val, nil
		}
		secretPath := strings.TrimPrefix(val, "gcp-secret://")
		accessReq := &secretmanagerpb.AccessSecretVersionRequest{Name: secretPath}
		result, err := client.AccessSecretVersion(ctx, accessReq)
		if err != nil {
			return "", err
		}
		return string(result.Payload.Data), nil
	}

	var errList []error

	cfg.Gin.Mode, err = resolve(cfg.Gin.Mode)
	if err != nil {
		errList = append(errList, fmt.Errorf("Gin.Mode: %w", err))
	}
	cfg.Gin.Log.Level, err = resolve(cfg.Gin.Log.Level)
	if err != nil {
		errList = append(errList, fmt.Errorf("Gin.Log.Level: %w", err))
	}
	cfg.Gin.Log.Format, err = resolve(cfg.Gin.Log.Format)
	if err != nil {
		errList = append(errList, fmt.Errorf("Gin.Log.Format: %w", err))
	}

	// Application
	cfg.Application.JWTSecret, err = resolve(cfg.Application.JWTSecret)
	if err != nil {
		errList = append(errList, fmt.Errorf("JWTSecret: %w", err))
	}
	// DB
	cfg.DB.Host, err = resolve(cfg.DB.Host)
	if err != nil {
		errList = append(errList, fmt.Errorf("DB.Host: %w", err))
	}
	cfg.DB.Port, err = resolve(cfg.DB.Port)
	if err != nil {
		errList = append(errList, fmt.Errorf("DB.Port: %w", err))
	}
	cfg.DB.User, err = resolve(cfg.DB.User)
	if err != nil {
		errList = append(errList, fmt.Errorf("DB.User: %w", err))
	}
	cfg.DB.Pass, err = resolve(cfg.DB.Pass)
	if err != nil {
		errList = append(errList, fmt.Errorf("DB.Pass: %w", err))
	}
	cfg.DB.Name, err = resolve(cfg.DB.Name)
	if err != nil {
		errList = append(errList, fmt.Errorf("DB.Name: %w", err))
	}
	cfg.DB.InstanceConn, err = resolve(cfg.DB.InstanceConn)
	if err != nil {
		errList = append(errList, fmt.Errorf("DB.InstanceConn: %w", err))
	}
	// GCS
	cfg.GCS.ProfileBucket, err = resolve(cfg.GCS.ProfileBucket)
	if err != nil {
		errList = append(errList, fmt.Errorf("GCS.ProfileBucket: %w", err))
	}
	// Google Places
	cfg.GooglePlaces.APIKey, err = resolve(cfg.GooglePlaces.APIKey)
	if err != nil {
		errList = append(errList, fmt.Errorf("GooglePlaces.APIKey: %w", err))
	}
	cfg.GooglePlaces.APIEndpoint, err = resolve(cfg.GooglePlaces.APIEndpoint)
	if err != nil {
		errList = append(errList, fmt.Errorf("GooglePlaces.APIEndpoint: %w", err))
	}
	// Vertex
	cfg.Vertex.ProjectID, err = resolve(cfg.Vertex.ProjectID)
	if err != nil {
		errList = append(errList, fmt.Errorf("Vertex.ProjectID: %w", err))
	}
	cfg.Vertex.Location, err = resolve(cfg.Vertex.Location)
	if err != nil {
		errList = append(errList, fmt.Errorf("Vertex.Location: %w", err))
	}
	cfg.Vertex.Model, err = resolve(cfg.Vertex.Model)
	if err != nil {
		errList = append(errList, fmt.Errorf("Vertex.Model: %w", err))
	}

	if len(errList) > 0 {
		return fmt.Errorf("config secret resolution errors: %v", errList)
	}
	return nil
}
