package transformer

import (
	"encoding/json"

	"github.com/relaunch-cot/lib-relaunch-cot/models"
	pbBaseModels "github.com/relaunch-cot/lib-relaunch-cot/proto/base_models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUserProfileToBaseModels(in *models.User) (*pbBaseModels.User, error) {
	baseModelsUser := &pbBaseModels.User{}

	b, err := json.Marshal(in)
	if err != nil {
		return nil, status.Error(codes.Internal, "error marshalling user model. Details: "+err.Error())
	}

	err = json.Unmarshal(b, &baseModelsUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "error unmarshalling user model. Details: "+err.Error())
	}

	return baseModelsUser, nil
}
