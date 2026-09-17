// Package rule은 거래 데이터를 주기적으로 조회해 이상거래를 탐지하는 룰들을 관리한다.
// 룰마다 거래를 조회하는 주기가 다르기 때문에 각 룰은 독립된 goroutine으로 실행된다.
package rule

import (
	"context"
	"time"
)

// Rule은 하나의 이상거래 탐지 규칙이다. since 이후 데이터를 조회/처리하고, 다음번에
// 이어서 조회할 기준 시각을 반환한다. 워터마크 저장/복원은 Engine이 Name()을 키로
// cursors 테이블에 대신 처리하므로, Rule 구현체는 그 부분을 신경 쓸 필요가 없다.
// 새 룰을 추가하려면 이 인터페이스를 구현하고 engine 생성 시점(main.go)에 목록에
// 추가하기만 하면 된다.
type Rule interface {
	Name() string
	Interval() time.Duration
	Run(ctx context.Context, since time.Time) (next time.Time, err error)
}
