package auth

import (
	"encoding/json"
	"fmt"

	"github.com/Izumra/2Handlers/domain/dto/requests"
	"github.com/Izumra/2Handlers/domain/dto/responses"
	"github.com/Izumra/2Handlers/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

type GroupHandlers struct {
	service *auth.Service
}

func MountAuthHandlers(router fiber.Router, service *auth.Service) {
	gh := &GroupHandlers{
		service,
	}

	router.Post("/login", gh.Login)
}

// @Summary Авторизация
// @Description Метод API, позволяющий пользователю произвести авторизацию в системе. После успешного прохождения авторизации, в ответе выдастся JWT токен доступа, требуемый для остальных обработчиков в качестве заголовка 'Authorization'
// @Tags Auth
// @Accept json
// @Produce  json
// @Param AuthBody body requests.LoginBody true "Тело запроса формата 'application/json', содержащее информацию для авторизации"
// @Success 200 {object} responses.Pattern{data=string,error=nil} "Структура успешного ответа на запрос авторизации"
// @Failure 500 {object} responses.Pattern{data=nil,error=string} "Структура неудачного ответа на запрос авторизации"
// @Router /login [post]
func (gh *GroupHandlers) Login(c *fiber.Ctx) error {

	var body requests.LoginBody
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(responses.Pattern{
			Error: fmt.Errorf("Тело запроса не соответсвует формату"),
		})
	}

	token, err := gh.service.Login(c.Context(), body.Username, body.Password)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(responses.Pattern{
			Error: err,
		})
	}

	return c.JSON(responses.Pattern{
		Data: token,
	})
}
