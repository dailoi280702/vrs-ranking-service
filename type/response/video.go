package response

import "github.com/dailoi280702/vrs-ranking-service/type/model"

type Video struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func FormVideoResponse(video model.Video) Video {
	return Video{
		ID:   video.ID,
		Name: video.Name,
	}
}

func FormVideosResponse(videos ...model.Video) []Video {
	res := make([]Video, len(videos))

	for i := range videos {
		res[i] = FormVideoResponse(videos[i])
	}

	return res
}
