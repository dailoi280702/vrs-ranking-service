package response

import (
	"github.com/dailoi280702/vrs-ranking-service/util/apperror"
	"github.com/dailoi280702/vrs-ranking-service/util/echoutil"
)

type Error apperror.AppError

type Data[T any] echoutil.Data[T]
