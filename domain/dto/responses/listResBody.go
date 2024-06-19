package responses

import "github.com/Izumra/2Handlers/domain/entity"

type ListResBody struct {
	Success bool
	News    []entity.News
}
