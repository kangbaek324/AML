// Package worker는 유저의 자산/위험 등급처럼 주기적으로 갱신되어야 하는
// 유저 정보를 관리하는 워커들을 묶는다.
package worker

import "context"

type Worker struct {
	assetTier *AssetTierWorker
	userRisk  *UserRiskWorker
}

func New(assetTier *AssetTierWorker, userRisk *UserRiskWorker) *Worker {
	return &Worker{assetTier: assetTier, userRisk: userRisk}
}

func (w *Worker) Start(ctx context.Context) {
	go w.assetTier.Start(ctx)
	go w.userRisk.Start(ctx)
}
