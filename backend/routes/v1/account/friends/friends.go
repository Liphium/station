package friends

import "github.com/gofiber/fiber/v2"

// Configuration
const MaximumFriends = 100

func Authorized(router fiber.Router) {
	router.Post("/add", addFriend)
	router.Post("/remove", removeFriend)
	router.Post("/list", listFriends)
	router.Post("/exists", existsFriend)
	router.Post("/update_receive_date", updateReceiveDate)
	router.Post("/get_receive_date", getReceiveDate)
}
