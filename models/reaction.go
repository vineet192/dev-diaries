package models

import "time"

type reaction_type string

type Reaction struct {
	id            uint
	comment_id    uint
	blog_id       uint
	user_id       uint
	reaction_type reaction_type
	posted_on     time.Time
}
