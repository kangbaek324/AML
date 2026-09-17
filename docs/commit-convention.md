# 커밋 컨벤션

커밋 메시지는 `type(scope): message` 형식을 따른다.

- **type**: `feat`, `fix`, `refactor`, `docs`, `chore`
- **scope**: `backend`, `frontend`

예시:

```
feat(backend): 유저 기능 추가
fix(frontend): 로그인 버튼 클릭 오류 수정
docs(backend): README 갱신
```

## 커밋 범위

한 커밋에는 Claude가 해당 세션에서 직접 작업한 변경사항만 포함한다. 사용자가 직접 수정했거나 다른 도구가 만든 변경사항이 워킹 트리에 같이 있어도 임의로 묶어서 커밋하지 않는다. 관련 없는 변경사항이 있으면 커밋하지 않고 사용자에게 알린다.
