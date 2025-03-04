package friends

import "github.com/gofiber/fiber/v2"

// Configuration
const MaximumFriends = 100

func Authorized(router fiber.Router) {
	router.Post("/add", addFriend)
	router.Post("/remove", removeFriend)
	router.Post("/update", updateFriend)
	router.Post("/list", listFriends)
	router.Post("/update_receive_date", updateReceiveDate)
	router.Post("/get_receive_date", getReceiveDate)
}
