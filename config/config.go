package config

type ExternalConfig struct {
}

type ServiceConfig struct {
	Gcp GcpDetails `mapstructure:"gcp" validate:"required"`
	//External *ExternalConfig `mapstructure:"service" validate:"required"`
}

type GcpDetails struct {
	ProjectID         string `mapstructure:"project_id" validate:"required"`
	BigQueryProjectId string `mapstructure:"bigquery_project_id" validate:"required"`
}
