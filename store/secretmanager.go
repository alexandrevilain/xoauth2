package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GCPSecretManager struct {
	log logr.Logger
	client *secretmanager.Client
	projectID string
	secretName string
}

func NewGCPSecretManager(ctx context.Context, client *secretmanager.Client, projectID, secretName string) (TokenStore, error) {
	s := &GCPSecretManager{
		client: client,
		projectID: projectID,
		secretName: secretName,
	}

	s.log = logr.FromContext(ctx)
	if s.log == nil {
		var err error
		s.log, err = s.defaultLogger()
		if err != nil {
			return nil, err
		}
	}
	
	err := s.ensureSecretExists()
	return s, err
}

func (s *GCPSecretManager) defaultLogger() (logr.Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return zapr.NewLogger(log), nil
}

func (s *GCPSecretManager) ensureSecretExists() error {
	ctx := context.Background()
	req := &secretmanagerpb.GetSecretRequest{
		Name: s.getSecretName(),
	}

	result, err := s.client.GetSecret(ctx, req)
	if err != nil {
		grpcError, ok := status.FromError(err); 
		if !ok {
			return err
		}

		if grpcError.Code() != codes.NotFound {
			return err
		}

		req := &secretmanagerpb.CreateSecretRequest{
			Parent:   fmt.Sprintf("projects/%s", s.projectID),
			SecretId: s.secretName,
			Secret: &secretmanagerpb.Secret{
					Replication: &secretmanagerpb.Replication{
							Replication: &secretmanagerpb.Replication_Automatic_{
									Automatic: &secretmanagerpb.Replication_Automatic{},
							},
					},
			},
		}

		result, err := s.client.CreateSecret(ctx, req)
        if err != nil {
                return fmt.Errorf("failed to create secret: %v", err)
        }
		s.log.V(1).Info("created secret: %s", result.Name)
        return nil
	}
	s.log.V(1).Info("secret already exists and was created at: %v", result.CreateTime)
	return nil
}

func (s *GCPSecretManager) getSecretName() string {
	return fmt.Sprintf("projects/%s/secrets/%s", s.projectID, s.secretName)
}


func (s *GCPSecretManager) Get() (*oauth2.Token, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("%s/versions/latest", s.getSecretName()), 
	}

	result, err := s.client.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %v", err)
	}
	
	token := &oauth2.Token{}
	err = json.Unmarshal(result.Payload.Data, token)
	return token, err
}

func (s *GCPSecretManager) Save(token *oauth2.Token) error {
	ctx := context.Background()

	payload, err := json.Marshal(token)
	if err != nil {
		return err
	}

	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: s.getSecretName(),
		Payload: &secretmanagerpb.SecretPayload{
				Data: payload,
		},
	}

	result, err := s.client.AddSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %v", err)
	}

	parts := strings.Split(result.Name, "/")

	latestVersion := parts[len(parts)-1]
	latestVersionInt, err  := strconv.Atoi(latestVersion)
	if err != nil {
		return err
	}

	previousVersion := latestVersionInt - 1
	if previousVersion >= 1 {
		req := &secretmanagerpb.DestroySecretVersionRequest{
			Name: fmt.Sprintf("%s/versions/%d", s.getSecretName(), previousVersion), 
		}

		_, err := s.client.DestroySecretVersion(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to destroy secret version: %v", err)
		}
	}
	s.log.V(1).Info("created new secret version: %v", result.Name)
	return nil
}