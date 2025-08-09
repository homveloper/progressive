# 기능 개발 워크플로우 가이드

## 📋 개요
개별 기능 개발 시 따라야 할 표준 워크플로우입니다. 각 기능은 1주 스프린트 단위로 진행됩니다.

## 🔄 기본 사이클 (1주 스프린트)

```
기획(1일) → UI/UX 설계(1일) → 개발(2-3일) → 테스트/피드백(0.5일) → 수정(0.5일) → 배포
```

---

## 📝 Phase 1: 기획 (0.5-1일)

### 체크리스트
- [ ] User Story 작성
- [ ] Acceptance Criteria 정의
- [ ] 기술 요구사항 정리
- [ ] 우선순위 설정

### User Story 템플릿
```markdown
As a [사용자 유형]
I want to [원하는 기능]
So that [달성하고자 하는 목표]
```

### Acceptance Criteria 템플릿
```markdown
Given [초기 상황]
When [사용자 행동]
Then [예상 결과]
```

### 산출물
- Feature Specification Document
- User Story & Acceptance Criteria
- 기술 스펙 초안

---

## 🎨 Phase 2: UI/UX 설계 (0.5-1일)

### 체크리스트
- [ ] 와이어프레임 스케치
- [ ] 인터랙션 플로우 정의
- [ ] 컴포넌트 명세 작성
- [ ] 프로토타입 제작 (필요시)

### 인터랙션 플로우 템플릿
```markdown
1. 초기 상태: [상태 설명]
2. 트리거: [사용자 액션]
3. 전환: [UI 변화]
4. 피드백: [시각적/청각적 피드백]
5. 완료 상태: [최종 상태]
```

### 컴포넌트 명세 템플릿
```markdown
## 컴포넌트명
- 용도: 
- Props: 
- States:
- Events:
- 디자인 토큰:
  - Colors:
  - Spacing:
  - Typography:
```

### 산출물
- 와이어프레임/스케치
- 인터랙션 명세서
- UI 컴포넌트 명세
- Figma/디자인 파일 (선택)

---

## 💻 Phase 3: 개발 (2-3일)

### Day 1: Backend 개발
```markdown
## 체크리스트
- [ ] API 엔드포인트 설계
- [ ] 데이터 모델 정의
- [ ] 비즈니스 로직 구현
- [ ] 유효성 검증 로직
- [ ] 단위 테스트 작성
```

### Day 2: Frontend 개발
```markdown
## 체크리스트
- [ ] UI 컴포넌트 구현
- [ ] 상태 관리 로직
- [ ] API 연동
- [ ] 에러 처리
- [ ] 로딩 상태 처리
```

### Day 3: 통합 및 마무리
```markdown
## 체크리스트
- [ ] Frontend-Backend 통합
- [ ] 통합 테스트
- [ ] 성능 최적화
- [ ] 코드 리뷰
- [ ] 문서 업데이트
```

### 코드 구조 템플릿

#### Backend (Go)
```go
// Handler
func (h *Handler) FeatureName(w http.ResponseWriter, r *http.Request) {
    // 1. 요청 파싱
    // 2. 유효성 검증
    // 3. 비즈니스 로직
    // 4. 응답 반환
}

// Service
func (s *Service) ProcessFeature(ctx context.Context, input Input) (Output, error) {
    // 비즈니스 로직
}

// Repository
func (r *Repository) SaveFeature(ctx context.Context, data Data) error {
    // 데이터 저장
}
```

#### Frontend (Templ)
```go
// Component
templ FeatureComponent(props Props) {
    <div class="feature-container">
        // UI 구현
    </div>
}

// Handler
func HandleFeature(data Data) {
    // 이벤트 처리
}
```

### 산출물
- 동작하는 기능 코드
- 테스트 코드
- API 문서
- 컴포넌트 문서

---

## 🧪 Phase 4: 테스트/피드백 (0.5일)

### 기능 테스트 체크리스트
```markdown
## 기본 기능
- [ ] 정상 동작 확인
- [ ] 엣지 케이스 처리
- [ ] 에러 상황 처리

## 사용성
- [ ] 직관적인 UI/UX
- [ ] 응답 시간 (<100ms)
- [ ] 피드백 메시지 적절성

## 호환성
- [ ] Chrome
- [ ] Firefox
- [ ] Safari
- [ ] Edge
- [ ] 모바일 반응형

## 접근성
- [ ] 키보드 네비게이션
- [ ] 스크린 리더 지원
- [ ] 색상 대비
```

### 피드백 수집 템플릿
```markdown
## 피드백 항목
- 기능: [작동 여부]
- 사용성: [개선 필요사항]
- 디자인: [UI 개선점]
- 성능: [속도 이슈]
- 버그: [발견된 버그]
```

### 산출물
- 테스트 결과 보고서
- 버그 리스트
- 개선 요구사항 목록

---

## 🔧 Phase 5: 수정 (0.5일)

### 우선순위 결정 매트릭스
```markdown
## Critical (즉시 수정)
- [ ] 기능 블로킹 버그
- [ ] 보안 이슈
- [ ] 데이터 손실 위험

## High (이번 스프린트)
- [ ] 주요 UX 이슈
- [ ] 성능 문제
- [ ] 접근성 문제

## Medium (다음 스프린트)
- [ ] 마이너 버그
- [ ] UI 개선사항
- [ ] 코드 리팩토링

## Low (백로그)
- [ ] 개선 제안
- [ ] 추가 기능
```

### 수정 작업 체크리스트
```markdown
- [ ] 버그 수정
- [ ] 테스트 업데이트
- [ ] 문서 업데이트
- [ ] 코드 리뷰 반영
- [ ] 최종 테스트
```

### 산출물
- 수정된 코드
- 업데이트된 테스트
- 수정 내역 문서

---

## 📅 주간 스프린트 일정

### 월요일: 기획 & 설계
```markdown
09:00-11:00  기획 회의, User Story 작성
11:00-12:00  Acceptance Criteria 정의
14:00-16:00  UI/UX 스케치
16:00-18:00  프로토타입/와이어프레임
```

### 화요일: Backend 개발
```markdown
09:00-12:00  API 설계 및 구현
14:00-17:00  비즈니스 로직 구현
17:00-18:00  단위 테스트 작성
```

### 수요일: Frontend 개발
```markdown
09:00-12:00  UI 컴포넌트 구현
14:00-17:00  상태 관리 및 API 연동
17:00-18:00  에러 처리
```

### 목요일: 통합
```markdown
09:00-11:00  Frontend-Backend 통합
11:00-12:00  통합 테스트
14:00-16:00  버그 수정
16:00-18:00  성능 최적화
```

### 금요일: 테스트 & 배포
```markdown
09:00-11:00  기능 테스트
11:00-12:00  피드백 수집
14:00-16:00  수정 작업
16:00-17:00  배포 준비
17:00-18:00  문서 정리
```

---

## 🚀 빠른 시작 체크리스트

### 새 기능 시작 시
```markdown
1. [ ] User Story 작성
2. [ ] 와이어프레임 스케치
3. [ ] API 스펙 정의
4. [ ] 개발 환경 준비
5. [ ] 테스트 계획 수립
```

### 기능 완료 조건
```markdown
1. [ ] 모든 AC 충족
2. [ ] 테스트 통과
3. [ ] 코드 리뷰 완료
4. [ ] 문서 업데이트
5. [ ] 배포 준비 완료
```

---

## 📚 참고 템플릿

### Feature Specification 템플릿
```markdown
# Feature: [기능명]

## 개요
[기능 설명]

## User Story
As a [user]
I want to [action]
So that [benefit]

## Acceptance Criteria
1. Given [context], When [action], Then [result]
2. ...

## 기술 요구사항
- Backend: 
- Frontend:
- Database:

## UI/UX 요구사항
- 디자인:
- 인터랙션:
- 반응형:

## 테스트 시나리오
1. 정상 케이스:
2. 엣지 케이스:
3. 에러 케이스:
```

### API 명세 템플릿
```markdown
## Endpoint: [METHOD] /api/[path]

### Request
```json
{
  "field1": "type",
  "field2": "type"
}
```

### Response
```json
{
  "status": "success",
  "data": {}
}
```

### Error Codes
- 400: Bad Request
- 401: Unauthorized
- 404: Not Found
- 500: Internal Server Error
```

---

## 💡 Best Practices

### 기획 단계
- 명확한 스코프 정의
- 측정 가능한 완료 조건
- 기술적 제약사항 사전 파악

### 설계 단계
- 기존 디자인 시스템 활용
- 재사용 가능한 컴포넌트 설계
- 모바일 우선 접근

### 개발 단계
- 작은 단위로 커밋
- 지속적인 테스트
- 코드 리뷰 문화

### 테스트 단계
- 실제 사용 시나리오 기반
- 다양한 환경에서 테스트
- 정량적 지표 측정

### 수정 단계
- 우선순위 기반 작업
- 근본 원인 해결
- 회귀 테스트 실시

---

## 📌 Quick Reference

### 시간 배분 가이드
- 기획: 15%
- 설계: 15%
- 개발: 50%
- 테스트: 10%
- 수정: 10%

### 커뮤니케이션 포인트
- 기획 완료 → 설계팀 리뷰
- 설계 완료 → 개발팀 리뷰
- 개발 완료 → QA 테스트
- 테스트 완료 → 이해관계자 리뷰
- 수정 완료 → 최종 승인

### 문서화 체크포인트
- [ ] User Story
- [ ] Technical Spec
- [ ] API Documentation
- [ ] Component Documentation
- [ ] Test Results
- [ ] Release Notes