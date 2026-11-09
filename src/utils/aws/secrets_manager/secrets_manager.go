package secrets_manager

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"y-ogura/echo_lambda_sample/src/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerService interface {
	GetRDS(string, string) (*Values, error)
	Create(string, string) error
	Delete(string) error
}

type secretsManagerService struct {
	Client secretsmanager.Client
}

type Values struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSecretsManager(cfg aws.Config) SecretsManagerService {
	var client *secretsmanager.Client
	endpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	if endpoint != "" {
		client = secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = secretsmanager.NewFromConfig(cfg)
	}
	return &secretsManagerService{Client: *client}
}

// Secrets Managerの値を取得
func (s *secretsManagerService) GetRDS(secretName string, region string) (*Values, error) {
	secretString, err := s.getSecret(secretName, region)
	if err != nil {
		return nil, err
	}

	var secret Values
	if err := json.Unmarshal([]byte(secretString), &secret); err != nil {
		logger.Log.Error("Failed to unmarshal Secrets Manager value", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("Complete Secrets Manager fetch")

	return &secret, nil
}

func (s *secretsManagerService) Create(name string, secret string) error {
	_, err := s.Client.CreateSecret(context.Background(),
		&secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String(secret),
		})
	if err != nil {
		return err
	}
	return nil
}

func (s *secretsManagerService) Delete(name string) error {
	_, err := s.Client.DeleteSecret(context.Background(),
		&secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(name),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
	if err != nil {
		return err
	}
	return nil
}

func (s *secretsManagerService) getSecret(secretName string, region string) (string, error) {
	logger.Log.Info(
		"Start Secrets Manager fetch",
		slog.String("secret_name", secretName),
		slog.String("region", region),
	)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := s.Client.GetSecretValue(context.Background(), input)
	if err != nil {
		logger.Log.Error(
			"Failed to get Secrets Manager value",
			slog.Any("error", err),
		)
		return "", err
	}

	logger.Log.Info("Complete Secrets Manager fetch")

	return *result.SecretString, nil
}
