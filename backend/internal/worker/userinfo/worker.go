// Package userinfo는 유저의 자산/위험 등급처럼 주기적으로 갱신되어야 하는
// 유저 정보를 관리하는 워커들을 묶는다.
package userinfo

import "context"

type Worker struct {
	assetTier *AssetTierWorker
	// userRisk *UserRiskWorker // TODO: User Risk 워커 추가 예정
}

func New(assetTier *AssetTierWorker) *Worker {
	return &Worker{assetTier: assetTier}
}

func (w *Worker) Start(ctx context.Context) {
	go w.assetTier.Start(ctx)
}
