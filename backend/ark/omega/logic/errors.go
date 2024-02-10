package logic

// マイクロサービスとして分離を考えると、モジュール毎に定義するのが無難?

import (
	"github.com/morikuni/failure"
)

// レスポンスに含めることを考えると、ドメインサービスで定義したエラーとは異なるステータスコードを意識したエラーを定義
var (
	InvalidArgument     failure.StringCode = "InvalidArgument"
	NotFound            failure.StringCode = "NotFound"
	Forbidden           failure.StringCode = "Forbidden"
	IntervalServerError failure.StringCode = "IntervalServerError"
)
