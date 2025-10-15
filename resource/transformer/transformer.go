package transformer

import (
	"encoding/json"

	"github.com/relaunch-cot/lib-relaunch-cot/models"
	pbBaseModels "github.com/relaunch-cot/lib-relaunch-cot/proto/base_models"
)

func GetUserProfileToBaseModels(in *models.User) (*pbBaseModels.User, error) {
	baseModelsUser := &pbBaseModels.User{}

	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &baseModelsUser)
	if err != nil {
		return nil, err
	}

	return baseModelsUser, nil
}
