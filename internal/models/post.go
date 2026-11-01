package models

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Version   int      `json:"version"`
}

type PostWithComments struct {
	Post
	Comments []CommentWithUser `json:"comments"`
	User     User              `json:"user"`
}

type PostWithMetadata struct {
	PostWithComments
	CommentsCount int `json:"comments_count"`
}
