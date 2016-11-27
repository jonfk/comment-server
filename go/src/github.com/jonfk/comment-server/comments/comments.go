package comments

import (
	"time"

	"github.com/satori/go.uuid"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Comment struct {
	CommentId       uuid.UUID     `db:"comment_id"`
	Timestamp       time.Time     `db:"timestamp"`
	Data            string        `db:"data"`
	ParentId        uuid.NullUUID `db:"parent_id"`
	CommentThreadId uuid.UUID     `db:"comment_thread_id"`
	AccountId       uuid.UUID     `db:"account_id"`
}

type CommentThread struct {
	CommentThreadId uuid.UUID `db:"comment_thread_id"`
	CreatedOn       time.Time `db:"created_on"`
	PageUrl         string    `db:"page_url"`
	Title           string    `db:"title"`
}

func (a CommentThread) Equal(b CommentThread) bool {
	if !uuid.Equal(a.CommentThreadId, b.CommentThreadId) ||
		!a.CreatedOn.Equal(b.CreatedOn) ||
		a.PageUrl != b.PageUrl ||
		a.Title != b.Title {
		return false
	}
	return true
}

type Comments struct {
	DB *sqlx.DB
}

func (t *Comments) CreateNewThread(pageUrl, title string, createdOn time.Time) (CommentThread, error) {
	var (
		newThread CommentThread
	)

	err := t.DB.QueryRowx("INSERT INTO comment_threads (comment_thread_id,created_on,page_url,title) VALUES ($1,$2,$3,$4) RETURNING comment_thread_id,created_on,page_url,title",
		uuid.NewV4().String(), createdOn, pageUrl, title).StructScan(&newThread)

	return newThread, err
}

func (t *Comments) DeleteThreadById(commentThreadId uuid.UUID) (uuid.UUID, error) {
	var deletedId uuid.UUID
	err := t.DB.QueryRowx("DELETE FROM comment_threads where comment_thread_id = $1 RETURNING comment_thread_id", commentThreadId).Scan(&deletedId)
	return deletedId, err
}

func (t *Comments) GetThreadByThreadId(commentThreadId uuid.UUID) (CommentThread, error) {
	var thread CommentThread
	err := t.DB.Get(&thread, "SELECT comment_thread_id,created_on,page_url,title FROM comment_threads where comment_thread_id = $1",
		commentThreadId)
	return thread, err
}

func (t *Comments) CreateNewComment(comment Comment) (Comment, error) {
	var (
		newComment Comment
	)

	err := t.DB.QueryRowx("INSERT INTO comments (comment_id,timestamp,data,parent_id,comment_thread_id,account_id) VALUES ($1,$2,$3,$4,$5,$6) RETURNING comment_id,timestamp,data,parent_id,comment_thread_id,account_id",
		uuid.NewV4().String(),
		comment.Timestamp,
		comment.Data,
		comment.ParentId,
		comment.CommentThreadId,
		comment.AccountId).StructScan(&newComment)

	return newComment, err
}

func (t *Comments) DeleteCommentById(commentId uuid.UUID) (uuid.UUID, error) {
	var deletedId uuid.UUID
	err := t.DB.QueryRowx("DELETE FROM comments where comment_id = $1 RETURNING comment_id", commentId).Scan(&deletedId)
	return deletedId, err
}
