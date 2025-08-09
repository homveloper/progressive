# Progressive 도메인 언어 (Domain Language)

JSON Schema 기반 테이블 에디터의 핵심 도메인 용어와 개념 정의

## 📊 계층 구조

```
Team (팀)
├── Individual Account (개인 계정)
└── Organization (조직)
    └── Workspace (워크스페이스)
        └── Table (테이블)
            ├── Schema (스키마)
            ├── Record (레코드)
            └── Field (필드)
```

---

## 🏢 Team Level

### Team (팀)
```yaml
정의: Progressive 서비스의 최상위 단위
특징:
  - 개인 계정 또는 조직으로 구분
  - 결제 및 구독 관리 단위
  - 사용자 인증의 기본 단위
예시: "김개발의 개인 팀", "스타트업 XYZ"
```

#### Individual Account (개인 계정)
```yaml
정의: 개인 사용자의 팀
특징:
  - 한 명의 사용자만 소속
  - 개인 프로젝트 관리
  - 무료 플랜 기본 제공
용량: 워크스페이스 5개, 테이블 50개 제한
```

#### Organization (조직)
```yaml
정의: 여러 사용자가 협업하는 팀
특징:
  - 다중 사용자 지원
  - 권한 관리 (Admin, Member, Viewer)
  - 팀 협업 기능
  - 고급 보안 설정
용량: 사용자 수에 따른 차등 요금제
```

---

## 🗂️ Workspace Level

### Workspace (워크스페이스)
```yaml
정의: 관련된 테이블들을 그룹화하는 논리적 공간
목적: 프로젝트, 부서, 용도별로 테이블 분류
특징:
  - 각 워크스페이스는 독립적인 권한 설정
  - 색상 테마 및 아이콘 커스터마이징
  - 워크스페이스별 멤버 초대 가능
예시:
  - "마케팅 캠페인"
  - "고객 관리"
  - "재고 관리"
  - "프로젝트 추적"
```

#### Workspace Types (워크스페이스 유형)
```yaml
Personal (개인):
  - 개인 작업용
  - 비공개 기본 설정

Shared (공유):
  - 팀 협업용
  - 멤버 권한 관리

Template (템플릿):
  - 재사용 가능한 워크스페이스
  - 새 워크스페이스 생성 시 복사 가능
```

---

## 📋 Table Level

### Table (테이블)
```yaml
정의: 구조화된 데이터를 저장하는 기본 단위
특징:
  - JSON Schema 기반 데이터 검증
  - 엑셀 형태의 편집 인터페이스
  - 실시간 협업 지원
  - 버전 관리 및 히스토리 추적
구성요소:
  - Schema (스키마)
  - Records (레코드들)
  - Views (뷰들)
```

### Schema (스키마)
```yaml
정의: 테이블의 데이터 구조와 검증 규칙을 정의하는 JSON Schema
역할:
  - 필드 타입 정의
  - 유효성 검증 규칙
  - 기본값 설정
  - 필수/선택 필드 지정
지원 타입:
  - string (문자열)
  - number (숫자)
  - integer (정수)
  - boolean (불린)
  - date (날짜)
  - email (이메일)
  - url (URL)
  - enum (선택 목록)
  - array (배열)
  - object (객체)
```

### Record (레코드)
```yaml
정의: 테이블의 각 행(row)에 해당하는 개별 데이터 항목
특징:
  - 스키마에 정의된 구조를 따름
  - 고유한 ID를 가짐
  - 생성/수정 시간 자동 기록
  - 실시간 동기화
예시: 고객 테이블의 "김고객" 레코드
```

### Field (필드)
```yaml
정의: 테이블의 각 열(column)에 해당하는 데이터 속성
구성요소:
  - name: 필드명
  - type: 데이터 타입
  - required: 필수 여부
  - default: 기본값
  - validation: 검증 규칙
예시: "이름", "이메일", "가입일", "활성 상태"
```

---

## 🔍 View & Edit Level

### View (뷰)
```yaml
정의: 테이블 데이터를 특정 방식으로 표시하는 인터페이스
종류:
  Grid View: 스프레드시트 형태
  Card View: 카드 레이아웃
  Calendar View: 달력 형태
  Kanban View: 칸반 보드
기능:
  - 필터링
  - 정렬
  - 그룹화
  - 컬럼 숨김/표시
```

### Edit Mode (편집 모드)
```yaml
정의: 레코드의 데이터를 수정할 수 있는 상태
종류:
  Inline Edit: 셀 직접 편집
  Modal Edit: 팝업 창 편집
  Bulk Edit: 일괄 편집
특징:
  - 실시간 검증
  - 자동 저장
  - 충돌 해결
  - 변경 이력 추적
```

---

## 🔐 Permission & Access

### Role (역할)
```yaml
Owner (소유자):
  - 모든 권한
  - 팀/조직 관리
  - 결제 관리

Admin (관리자):
  - 워크스페이스 관리
  - 사용자 초대/제거
  - 테이블 생성/삭제

Editor (편집자):
  - 테이블 편집
  - 레코드 CRUD
  - 뷰 생성/수정

Viewer (보기 전용):
  - 데이터 조회만 가능
  - 내보내기 권한
```

### Permission Level (권한 레벨)
```yaml
Team Level:
  - 팀 설정 관리
  - 구독 관리
  - 전체 사용자 관리

Workspace Level:
  - 워크스페이스 설정
  - 멤버 관리
  - 테이블 생성/삭제

Table Level:
  - 스키마 수정
  - 데이터 편집
  - 뷰 관리

Record Level:
  - 개별 레코드 편집
  - 필드별 권한 설정
```

---

## 📈 Data & Analytics

### Import/Export (가져오기/내보내기)
```yaml
지원 형식:
  - CSV (Comma-Separated Values)
  - JSON (JavaScript Object Notation)
  - Excel (.xlsx, .xls)
  - Google Sheets (API 연동)

특징:
  - 스키마 자동 감지
  - 데이터 타입 변환
  - 오류 레포트 생성
  - 대용량 파일 처리 지원
```

### History & Version (히스토리 & 버전)
```yaml
Change Log (변경 로그):
  - 모든 데이터 변경 추적
  - 사용자별 변경 이력
  - 시간 기반 필터링

Revision (리비전):
  - 스키마 변경 버전
  - 롤백 기능
  - 변경점 비교

Backup (백업):
  - 자동 백업 (매일)
  - 수동 백업
  - 백업본 복원
```

---

## 🔄 Collaboration & Sync

### Real-time Collaboration (실시간 협업)
```yaml
기능:
  - 동시 편집 지원
  - 사용자 커서 표시
  - 변경사항 실시간 동기화
  - 충돌 자동 해결

Presence (현재 상태):
  - 온라인 사용자 표시
  - 편집 중인 셀 표시
  - 마지막 활동 시간
```

### Comments & Discussion (댓글 & 토론)
```yaml
Comment System:
  - 셀별 댓글
  - 레코드별 댓글
  - 멘션 기능 (@사용자)
  - 댓글 알림

Thread (스레드):
  - 댓글 답글
  - 해결됨 표시
  - 댓글 히스토리
```

---

## 🚀 Integration & API

### Webhook & Automation (웹훅 & 자동화)
```yaml
Trigger Events:
  - 레코드 생성/수정/삭제
  - 스키마 변경
  - 사용자 권한 변경

Actions:
  - 외부 서비스 알림
  - 이메일 발송
  - 데이터 동기화
```

### API Access (API 접근)
```yaml
REST API:
  - 모든 CRUD 작업
  - 인증: JWT 토큰
  - 속도 제한: 시간당 1000 요청

GraphQL API:
  - 필요한 데이터만 조회
  - 실시간 구독
  - 스키마 기반 타입 안전성
```

---

## 📱 User Interface

### Navigation (내비게이션)
```yaml
Sidebar:
  - 워크스페이스 선택기
  - 테이블 목록
  - 즐겨찾기
  - 최근 작업

Breadcrumb:
  - 현재 위치 표시
  - 상위 레벨로 이동
  - 빠른 네비게이션
```

### Context Menu (컨텍스트 메뉴)
```yaml
Table Context:
  - 복제
  - 내보내기
  - 설정
  - 삭제

Record Context:
  - 편집
  - 복제
  - 삭제
  - 댓글 추가

Cell Context:
  - 복사
  - 붙여넣기
  - 셀 형식
  - 유효성 검사 정보
```

---

## 💡 Usage Examples

### 실제 사용 시나리오

#### 스타트업 고객 관리
```yaml
Team: "스타트업 ABC"
Workspace: "고객 관리"
Table: "리드 고객"
Schema:
  - 이름 (string, required)
  - 이메일 (email, required)
  - 회사명 (string)
  - 관심도 (enum: [높음, 중간, 낮음])
  - 등록일 (date, default: today)
Records: 각각의 잠재 고객 정보
```

#### 프로젝트 태스크 관리
```yaml
Team: "개발팀"
Workspace: "프로젝트 관리"
Table: "백로그"
Schema:
  - 제목 (string, required)
  - 설명 (string)
  - 담당자 (string)
  - 상태 (enum: [TODO, 진행중, 완료])
  - 우선순위 (integer, 1-5)
  - 마감일 (date)
Records: 각각의 태스크
Views:
  - 칸반 뷰 (상태별)
  - 달력 뷰 (마감일별)
```

---

## 🔍 Glossary (용어집)

| 한국어 | 영어 | 정의 |
|--------|------|------|
| 팀 | Team | 최상위 조직 단위 |
| 조직 | Organization | 다중 사용자 팀 |
| 워크스페이스 | Workspace | 테이블 그룹화 공간 |
| 테이블 | Table | 데이터 저장 기본 단위 |
| 스키마 | Schema | 데이터 구조 정의 |
| 레코드 | Record | 개별 데이터 항목 |
| 필드 | Field | 데이터 속성/컬럼 |
| 뷰 | View | 데이터 표시 방식 |
| 권한 | Permission | 접근 제어 |
| 역할 | Role | 사용자 권한 그룹 |

---

## 📝 Naming Conventions

### 파일/폴더 명명 규칙
```yaml
Database Tables:
  - teams
  - organizations  
  - workspaces
  - tables
  - records
  - fields

API Endpoints:
  - /api/teams
  - /api/workspaces/{id}/tables
  - /api/tables/{id}/records

Component Names:
  - WorkspaceSelector
  - TableGrid
  - RecordEditor
  - SchemaBuilder
```

### URL 구조
```yaml
Path Structure:
  - /team/{team-id}
  - /workspace/{workspace-id}
  - /table/{table-id}
  - /table/{table-id}/record/{record-id}

Query Parameters:
  - ?view=grid|card|calendar
  - ?filter=field:value
  - ?sort=field:asc|desc
```

---

이 도메인 언어는 Progressive 서비스 전반에서 일관되게 사용되어야 하며, 새로운 기능 개발 시 참고 기준으로 활용됩니다.