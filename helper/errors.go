package helper

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/rollbar/terraform-provider-rollbar/client"
)

func Timedout(err error) bool {
	timedoutErr, ok := err.(*retry.TimeoutError)
	return ok && timedoutErr.LastError == nil
}

func NotFound(err error) bool {
	var e *retry.NotFoundError
	return errors.As(err, &e)
}

// Rollbar APIで結果が見つからなかったときのエラーかどうか判定
func NotFoundFromClient(err error) bool {
	return errors.Is(err, client.ErrNotFound)
	//return err == client.ErrNotFound
}

// Rollbar APIで認証に失敗した時のエラーかどうか判定
func UnauthorizedFromClient(err error) bool {
	return errors.Is(err, client.ErrUnauthorized)
}
