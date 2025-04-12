package response

import "github.com/dailoi280702/vrs-ranking-service/type/model"

type Video struct {
	ID     int64  `json:"id" example:"5" description:"Video ID"`
	Name   string `json:"name" example:"Funny Cat Video" description:"Video Name"`
	Likes  int64  `json:"likes" example:"1000" description:"Number of likes"`
	Views  int64  `json:"view" example:"10000" description:"Number of views"`
	Shares int64  `json:"share" example:"500" description:"Number of shares"`
	Length int64  `json:"length" example:"600" description:"Video length in seconds"`
}

func FormVideoResponse(video model.Video) Video {
	return Video{
		ID:     video.ID,
		Name:   video.Name,
		Likes:  video.Likes,
		Views:  video.Views,
		Shares: video.Shares,
		Length: video.Length,
	}
}

func FormVideosResponse(videos ...model.Video) []Video {
	res := make([]Video, len(videos))

	for i := range videos {
		res[i] = FormVideoResponse(videos[i])
	}

	return res
}
