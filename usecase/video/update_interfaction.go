package video

import (
	"context"

	"github.com/dailoi280702/vrs-ranking-service/type/request"
)

func (u *Usecase) UpdateInteraction(ctx context.Context, req request.UpdateInteraction) error {
	if err := req.Validate(); err != nil {
		return err
	}

	return nil
}
