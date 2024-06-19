package news

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Izumra/2Handlers/domain/dto/requests"
	"github.com/Izumra/2Handlers/domain/dto/responses"
	"github.com/Izumra/2Handlers/internal/service/news"
	"github.com/gofiber/fiber/v2"
)

type GroupHandlers struct {
	service *news.Service
}

func MountNewsHandlers(router fiber.Router, service *news.Service) {
	gh := &GroupHandlers{
		service,
	}

	router.Put("/edit/:Id", gh.EditNews)
	router.Get("/list/:offset/:count", gh.List)
}

// @Summary Редактировать новость
// @Description Метод API, позволяющий авторизированному пользователю отредактировать новость по id. Учтены пожелания и предусмотрено сохранение полей от перезаписи,  если значение не было передано в запросе.
// @Tags News
// @Accept json
// @Produce  json
// @Param Id path int true "Идентификатор новости" default(0)
// @Param UpdateNews body requests.NewsData true "Тело запроса формата 'application/json', содержащее информацию для изменения записи"
// @Success 200 {object} responses.Pattern{data=responses.ListResBody,error=nil} "Структура успешного ответа на запрос изменения записи"
// @Failure 500 {object} responses.Pattern{data=nil,error=string} "Структура неудачного ответа на запрос изменения записи"
// @Security Authorization
// @Router /edit/{Id}/ [put]
func (gh *GroupHandlers) EditNews(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		c.Status(fiber.StatusForbidden)
		return c.JSON("Доступ запрещен")
	}

	idParam := c.Params("Id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON("Формат идентификатора статьи должен быть числовым значением")
	}

	var body requests.NewsData
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON("Тело запроса не соответсвует формату")
	}

	err = gh.service.EditNews(c.Context(), accessToken, id, body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON("Запись обновлена")
}

// @Summary Список новостей
// @Description Метод API, позволяющий авторизированному пользователю получить список новостей в заданном диапазоне параметроми 'offset' и 'count'
// @Tags News
// @Accept json
// @Produce  json
// @Param offset path int true "Шаг смещения" default(0)
// @Param count path int true "Количество" default(0)
// @Success 200 {object} responses.Pattern{data=responses.ListResBody,error=nil} "Структура успешного ответа на запрос получения списка новостей"
// @Failure 500 {object} responses.Pattern{data=nil,error=string} "Структура неудачного ответа на запрос получения списка новостей"
// @Security Authorization
// @Router /list/{offset}/{count} [get]
func (gh *GroupHandlers) List(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		c.Status(fiber.StatusForbidden)
		return c.JSON(responses.Pattern{
			Error: fmt.Errorf("Доступ запрещен"),
		})
	}

	offsetParam := c.Params("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(responses.Pattern{
			Error: fmt.Errorf("Шаг смещения статей должен быть числовым значением"),
		})
	}

	countParam := c.Params("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(responses.Pattern{
			Error: fmt.Errorf("Количество статей должно быть числовым значением"),
		})
	}

	news, err := gh.service.ListNews(c.Context(), accessToken, offset, count)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(responses.Pattern{
			Error: err,
		})
	}

	return c.JSON(responses.Pattern{
		Data: responses.ListResBody{
			Success: true,
			News:    news,
		},
	})
}
