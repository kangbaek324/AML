// Package userinfo는 유저의 자산/위험 등급처럼 주기적으로 갱신되어야 하는
// 유저 정보를 관리하는 워커들을 묶는다.
package userinfo

import "context"

type Manager struct {
	assetTier *AssetTierWorker
	userRisk  *UserRiskWorker
}

func New(assetTier *AssetTierWorker, userRisk *UserRiskWorker) *Manager {
	return &Manager{assetTier: assetTier, userRisk: userRisk}
}

func (m *Manager) Start(ctx context.Context) {
	go m.assetTier.Start(ctx)
	go m.userRisk.Start(ctx)
}

// RunNow는 스케줄과 무관하게 즉시 전체 유저 정보를 갱신한다 (수동 갱신 API용).
func (m *Manager) RunNow(ctx context.Context) {
	m.assetTier.run(ctx)
	m.userRisk.run(ctx)
}
