package request

import "devdiaries/models"

type UserQuery struct {
	models.User
	IncludeComments         bool `schema:"include_comments"`
	IncludeBlogs            bool `schema:"include_blogs"`
	IncludeCommentReactions bool `schema:"include_comment_reactions"`
	IncludeBlogReactions    bool `schema:"include_blog_reactions"`
	IncludeFollowers        bool `schema:"include_followers"`
	IncludeFollowing        bool `schema:"include_following"`
}
