package helper

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Retryable func(error) (bool, error)

func RetryWhen[T any](ctx context.Context, timeout time.Duration, f func() (T, error), retryable Retryable) (T, error) {
	var output T

	// リトライ可能なエラーの場合はf()を再試行させる
	// リトライ不可能なエラーの場合はそのままエラーを送出する
	// エラーが発生せず実行に成功した場合は終了して, outputにある結果を利用する
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var err error
		var again bool

		output, err = f()

		// retryableでリトライ可能か判定する
		// またエラーの有無も判定する
		again, err = retryable(err)

		if again {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})

	// 途中でタイムアウトとなった場合
	// 最後にもう一度実行し, 再試行せずその結果をそのまま返す
	if Timedout(err) {
		output, err = f()
	}

	if err != nil {
		// 最終結果にnilが指定できないので, 最後の実行結果を返す
		return output, err
	}

	return output, nil
}

func RetryWhenNotFoundFromClient[T any](ctx context.Context, timeout time.Duration, f func() (T, error)) (T, error) {
	return RetryWhen[T](ctx, timeout, f, func(err error) (bool, error) {
		if NotFoundFromClient(err) {
			return true, err
		}
		return false, err
	})
}
