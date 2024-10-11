package post_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/posts/list_after
func listAfter(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Creator string   // The account id of the creator
		Key     string   // Required to access posts to friends (locked with stored action key)
		Groups  []string `json:"groups"` // Required to access posts to groups
		After   int64    // The time after
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Stuff to search for in the visibility table
	var types []string
	var identifiers []string

	// Check if public posts can be revealed
	if req.Creator != "" {
		var key database.StoredActionKey
		if err := database.DBConn.Where("id = ?", req.Creator).Take(&key).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorServer, nil)
		}

		// Add the friend visibility if the key is correct
		if key.Key == req.Key {
			types = append(types, database.VisibilityFriends)
			identifiers = append(identifiers, "-")
		}
	}

	// Add all the group identifiers
	if len(req.Groups) > 0 {
		types = append(types, database.VisibilityConversation)
		identifiers = append(identifiers, req.Groups...)
	}

	// Get all the visibilities that match from the database
	var unfilteredVis []database.PostVisibility
	if err := database.DBConn.Where("creator = ? AND type IN ? AND identifier IN ? AND creation > ?", req.Creator, types, identifiers, req.After).Limit(60).Find(&unfilteredVis).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the ids of the posts we have to grab
	postMap := map[string][]database.PostVisibility{}
	postsToGrab := []string{}
	for _, vis := range unfilteredVis {
		_, valid := postMap[vis.Post.String()]

		// If the post is not found yet, add the post id and the visibility
		if !valid {
			postMap[vis.Post.String()] = []database.PostVisibility{vis}
			postsToGrab = append(postsToGrab, vis.Post.String())
		} else {

			// Only add the visibility if the post has already been added
			postMap[vis.Post.String()] = append(postMap[vis.Post.String()], vis)
		}
	}

	// Grab all the posts from the database
	var posts []database.Post
	if err := database.DBConn.Where("id IN ?", postsToGrab).Limit(30).Find(&posts).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Convert all the data to sendable posts
	sendables := []database.SentPost{}
	for _, post := range posts {
		sendables = append(sendables, post.ToSent(postMap[post.ID.String()]))
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"posts":   sendables,
	})
}
