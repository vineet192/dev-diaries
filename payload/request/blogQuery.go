package request

import "strings"

type BlogQuery struct {
	Title            *string `schema:"title"`
	Content          *string `schema:"content"`
	DisableComments  bool    `schema:"disable_comments"`
	DisableReactions bool    `schema:"disable_reactions"`
	Page             uint    `schema:"page"`
	PageSize         uint    `schema:"page_size"`
}

func (bq *BlogQuery) DBQuery() (string, []interface{}) {

	query := make([]string, 0)
	values := make([]interface{}, 0)

	if bq.Title != nil {
		query = append(query, "title LIKE ?")
		values = append(values, "%"+*bq.Title+"%")
	}

	if bq.Content != nil {
		query = append(query, "content LIKE ?")
		values = append(values, "%"+*bq.Content+"%")
	}

	return strings.Join(query, "AND"), values
}
