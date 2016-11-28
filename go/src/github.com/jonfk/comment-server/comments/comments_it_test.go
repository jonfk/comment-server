package comments

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"

	"github.com/jonfk/comment-server/accounts"
)

var (
	DBUser, DBName, DBPassword string
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBUser = os.Getenv("DATABASE_USER")
	DBName = os.Getenv("DATABASE_NAME")
	DBPassword = os.Getenv("DATABASE_PASSWORD")

	os.Exit(m.Run())
}

func TestCommentThreads(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DBUser, DBName, DBPassword))
	if err != nil {
		t.Fatalf("sqlx.Connect failed : %v\n", err)
	}

	comments := &Comments{DB: db}

	expectedCommentThread := CommentThread{
		CreatedOn: time.Now().UTC().Round(time.Second),
		PageUrl:   "pageUrl",
		Title:     "title",
	}

	createdCommentThread, err := comments.CreateNewThread(
		expectedCommentThread.PageUrl,
		expectedCommentThread.Title,
		expectedCommentThread.CreatedOn)

	if err != nil {
		t.Fatalf("comments.CreateNewThread failed : %v\n", err)
	}

	if createdCommentThread.Title != expectedCommentThread.Title ||
		createdCommentThread.PageUrl != expectedCommentThread.PageUrl ||
		!createdCommentThread.CreatedOn.Equal(expectedCommentThread.CreatedOn) {
		t.Fatalf("createdCommentThread != expectedCommentThread\n (createdCommentThread) %v != (expectedCommentThread) %v",
			createdCommentThread, expectedCommentThread)
	}

	expectedCommentThread.CommentThreadId = createdCommentThread.CommentThreadId

	fetchedCommentThread, err := comments.GetThreadByThreadId(createdCommentThread.CommentThreadId)
	if err != nil {
		t.Fatalf("comments.GetThreadByThreadId failed : %v\n", err)
	}

	if !fetchedCommentThread.Equal(expectedCommentThread) {
		t.Fatalf("fetchedCommentThread != expectedCommentThread\n (fetchedCommentThread) %v != (expectedCommentThread) %v)",
			fetchedCommentThread, expectedCommentThread)
	}

	_, err = comments.DeleteThreadById(createdCommentThread.CommentThreadId)
	if err != nil {
		t.Fatalf("comments.DeleteThread failed : %v\n", err)
	}
}

func TestCreateNewComment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DBUser, DBName, DBPassword))
	if err != nil {
		t.Fatalf("sqlx.Connect failed : %v\n", err)
	}

	accountsModule := &accounts.Accounts{DB: db}
	commentsModule := &Comments{DB: db}

	account, commentThread, err := CreateAccountAndCommentThread(accountsModule, commentsModule)

	cleanUp := func(commentId uuid.UUID) {
		if !uuid.Equal(commentId, uuid.Nil) {
			_, err := commentsModule.DeleteCommentById(commentId)
			if err != nil {
				t.Fatalf("commentsModule.DeleteCommentById failed : %v\n", err)
			}
		}
		_, err := commentsModule.DeleteThreadById(commentThread.CommentThreadId)
		if err != nil {
			t.Fatalf("commentsModule.DeleteThreadById failed : %v\n", err)
		}
		_, err = accountsModule.DeleteById(account.AccountId)
		if err != nil {
			t.Fatalf("AccountsModule.DeleteById failed : %v\n", err)
		}
	}

	expectedComment := Comment{
		Timestamp:       time.Now().UTC().Round(time.Second),
		Data:            "this is a comment",
		CommentThreadId: commentThread.CommentThreadId,
		AccountId:       account.AccountId,
	}

	createdComment, err := commentsModule.CreateNewComment(expectedComment)
	if err != nil {
		cleanUp(uuid.Nil)
		t.Fatalf("commentsModule.CreateNewComment failed : %v\n", err)
	}

	if !createdComment.Timestamp.Equal(expectedComment.Timestamp) ||
		createdComment.Data != expectedComment.Data ||
		!uuid.Equal(createdComment.AccountId, expectedComment.AccountId) ||
		!uuid.Equal(createdComment.CommentThreadId, expectedComment.CommentThreadId) {
		cleanUp(createdComment.CommentId)
		t.Fatalf("createdComment != expectedComment\n (createdComment) %v != (expectedComment) %v", createdComment, expectedComment)
	}

	cleanUp(createdComment.CommentId)
}

func CreateAccountAndCommentThread(accountsModule *accounts.Accounts, commentsModule *Comments) (accounts.Account, CommentThread, error) {
	var (
		account       accounts.Account
		commentThread CommentThread
		err           error
	)

	commentThread, err = commentsModule.CreateNewThread("newpageurl", "newtitle", time.Now().UTC().Round(time.Second))
	if err != nil {
		return account, commentThread, err
	}

	account, err = accountsModule.CreateNewAccount(accounts.Account{
		Username:  "CreateAccountAndCommentThreadusername",
		Email:     "CreateAccountAndCommentThreademail",
		CreatedOn: time.Now().UTC().Round(time.Second),
	}, "password")
	if err != nil {
		return account, commentThread, err
	}

	return account, commentThread, err
}
