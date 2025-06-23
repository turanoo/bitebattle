package config

type Config struct {
	Gin          Gin          `yaml:"gin"`
	Auth0        Auth0        `yaml:"auth0"`
	Application  Application  `yaml:"application"`
	DB           Db           `yaml:"db"`
	GCS          Gcs          `yaml:"gcs"`
	GooglePlaces GooglePlaces `yaml:"google_places"`
	Vertex       Vertex       `yaml:"vertex"`
}

type Gin struct {
	Mode string `yaml:"mode"`
	Log  Log    `yaml:"log"`
}

type Log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Auth0 struct {
	Domain       string `yaml:"domain"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Application struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	JWTSecret   string `yaml:"jwt_secret"`
}

type Db struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Pass         string `yaml:"pass"`
	Name         string `yaml:"name"`
	InstanceConn string `yaml:"instance_connection_name"`
}

type Gcs struct {
	ProfileBucket string `yaml:"profile_bucket"`
	ObjectURL     string `yaml:"object_url"`
}

type GooglePlaces struct {
	APIKey      string `yaml:"api_key"`
	APIEndpoint string `yaml:"api_endpoint"`
}

type Vertex struct {
	ProjectID string `yaml:"project_id"`
	Location  string `yaml:"location"`
	Model     string `yaml:"model"`
	AuthToken string `yaml:"auth_token"`
}
