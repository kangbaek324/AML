# AML Frontend

AML 담당자용 대시보드 프론트엔드입니다. [Rule-based Detection](../backend/README.md#rule-based-detection)으로 발생한 Alert를 담당자가 확인하고 처리하는 화면을 제공합니다.

## Content

- [Structure](#structure)
- [Implement](#implement)
  - [Dashboard](#dashboard)
  - [Alerts](#alerts)
- [How to Run](#how-to-run)
- [Deploy](#deploy)

## Structure

```
src/
├── components/   # 재사용 UI 컴포넌트 (layout, 도메인별 하위 폴더)
├── pages/        # 라우트 단위 화면
├── services/api/ # 백엔드 API 클라이언트
├── hooks/        # 커스텀 훅
├── types/        # 공용 타입
├── utils/        # 유틸 함수
├── constants/    # 상수
└── router.tsx    # 라우트 정의
```

## Implement

### Dashboard

유저별 평균 자산, 위험도, 자산 등급을 조회합니다. 행을 펼치면 계좌 목록이, 계좌를 펼치면 보유 종목까지 드릴다운됩니다.

<img width="1000" alt="Dashboard 목록" src="../docs/image/fe_1.png" />
<img width="1000" alt="Dashboard 드릴다운" src="../docs/image/fe_2.png" />

### Alerts

Rule-based Detection이 발생시킨 Alert 목록을 상태(PENDING/NORMAL/ABNORMAL)와 함께 보여줍니다. 행을 펼치면 Alert 판단 근거가 된 거래/이체 내역을 확인할 수 있습니다.

PENDING 상태일 경우 클릭해 상태를 바꿀 수 있습니다.

<img width="1000" alt="Alerts 목록" src="../docs/image/fe_3.png" />
<img width="1000" alt="Alert 상세" src="../docs/image/fe_4.png" />

## How to Run

### Requirements

| 항목 | 버전 / 값 |
| ---- | --------- |
| Node | `22+`     |

### Environment Variables

`.env.example`를 참고하여 `.env` 파일을 생성하세요.

```
VITE_API_BASE_URL=http://localhost:8080
```

### Run

```bash
npm install
npm run dev
```

## Deploy

이미지는 환경에 상관없이 한 번만 빌드하고, API 주소는 컨테이너 실행 시점에 `API_BASE_URL` 환경변수로 주입합니다.

```bash
docker build -t aml-frontend .
docker run -p 81:81 -e API_BASE_URL=http://api.example.com aml-frontend
```

nginx가 `API_BASE_URL`을 채운 `/config.js`를 서빙하고, `index.html`이 이를 앱 번들보다 먼저 로드합니다. `public/config.js`는 `npm run dev`에서 같은 경로가 404가 나지 않도록 두는 빈 스텁이며, 개발 환경에서는 `.env`의 `VITE_API_BASE_URL`이 사용됩니다.
