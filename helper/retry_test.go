package helper

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rollbar/terraform-provider-rollbar/client"
)

func TestRetryWhenNotFoundFromClient(t *testing.T) {
	t.Parallel()

	const RETRY_TIMEOUT = 5 * time.Second
	const RETRY_COUNT_IN_TEST int32 = 5

	var currentRetryCount = RETRY_COUNT_IN_TEST

	testCases := []struct {
		Name        string
		F           func() (any, error)
		ExpectError bool
	}{
		{
			// 渡したハンドラ内でエラーの発生しない
			Name: "no error",
			F: func() (any, error) {
				return nil, nil
			},
			ExpectError: false,
		},
		{
			// APIのレスポンスでリソースがない場合にはリトライを繰り返す
			// ただし指定したタイムアウト値以上リソースが確認できない場合はエラーとする
			Name: "retryable NotFound timeout",
			F: func() (any, error) {
				return nil, client.ErrNotFound
			},
			ExpectError: true,
		},
		{
			Name: "non-retryable Error",
			F: func() (any, error) {
				return nil, client.ErrUnauthorized
			},
			ExpectError: true,
		},
		{
			Name: "retryable NotFoundError success",
			F: func() (any, error) {
				// 上記で指定した再試行回数だけAPIレスポンスに起因するNotFoundエラーを返す
				// ここだけパラレルに実行されないようにCAS関数を使用する
				// 実行にかかる時間はタイムアウト値未満なので全体としてはエラーとならない
				if atomic.CompareAndSwapInt32(&currentRetryCount, 0, 1) {
					return nil, client.ErrNotFound
				}
				return nil, nil
			},
			ExpectError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			defer func() {
				// 1つのtestCaseが完了したら、再試行回数のカウンターは元の値に戻しておく
				currentRetryCount = RETRY_COUNT_IN_TEST
			}()

			_, err := RetryWhenNotFoundFromClient(context.Background(), RETRY_TIMEOUT, testCase.F)
			if testCase.ExpectError && err == nil {
				t.Fatal("expected error")
			} else if !testCase.ExpectError && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}
