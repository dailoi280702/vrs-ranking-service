package request

import "context"

type VideoInteraction string

const (
	VideoInteractionView    = "view"
	VideoInteractionLike    = "like"
	VideoInteractionShare   = "share"
	VideoInteractionComment = "comment"
	VideoInteractionWatch   = "watch"
)

type UpdateInteraction struct {
	VideoId          int64            `json:"-" validate:"required,gt=0"`
	Type             VideoInteraction `json:"type" validate:"required,oneof=view like share comment watch" example:"view" description:"Interaction type (view, like, share, comment, watch)"`
	WatchTimeSeconds int64            `json:"watch_time" validate:"gte=0" example:"300" description:"Watch time in seconds"`
}

func (r *UpdateInteraction) Validate() error {
	if err := transformer.Struct(context.Background(), r); err != nil {
		return err
	}

	if err := validator.Struct(r); err != nil {
		return err
	}

	return nil
}

type GetTopVideos struct {
	UserId *int64 `query:"user_id" validate:"omitempty,gt=0"`
}

func (r *GetTopVideos) Validate() error {
	if err := transformer.Struct(context.Background(), r); err != nil {
		return err
	}

	if err := validator.Struct(r); err != nil {
		return err
	}

	return nil
}
