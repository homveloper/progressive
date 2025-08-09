# Portable CSS Build System - 설계 문서

## 🎯 설계 목표
- Tailwind CSS 스타일의 유틸리티 CSS 시스템 구축
- 포터블한 독립형 CSS 파일 생성 (외부 의존성 없음)
- 현재 38KB CSS를 ~8KB로 최적화 (PurgeCSS 적용)
- Go 프로젝트와 원활한 통합

## 🏗️ 시스템 아키텍처

### 디렉토리 구조
```
progressive/
├── css-system/                 # 새로운 CSS 시스템
│   ├── src/                   # 소스 CSS 파일들
│   │   ├── core/
│   │   │   ├── reset.css      # 브라우저 초기화
│   │   │   ├── variables.css  # CSS 변수 (현재 app.css에서 추출)
│   │   │   └── base.css       # 기본 요소 스타일
│   │   ├── utilities/
│   │   │   ├── spacing.css    # p-1, m-2, px-4 등
│   │   │   ├── typography.css # text-sm, font-bold 등  
│   │   │   ├── layout.css     # flex, grid, items-center 등
│   │   │   ├── colors.css     # text-excel-primary, bg-success 등
│   │   │   ├── borders.css    # border, rounded 등
│   │   │   └── effects.css    # shadow, transition 등
│   │   └── components/
│   │       ├── buttons.css    # .btn-* 컴포넌트 클래스
│   │       ├── forms.css      # .form-* 컴포넌트 클래스
│   │       └── grids.css      # .grid-cell-* 전용 클래스
│   ├── build/                 # 빌드된 CSS 파일들
│   │   ├── progressive.css    # 개발용 (주석 포함)
│   │   ├── progressive.min.css # 프로덕션용 (압축)
│   │   ├── portable.css       # 포터블용 (CSS 변수 하드코딩)
│   │   └── portable.min.css   # 포터블 압축버전
│   ├── tools/                 # 빌드 도구들
│   │   ├── build.js          # 메인 빌드 스크립트
│   │   ├── generate-utils.js # 유틸리티 클래스 생성기
│   │   ├── purge.js          # 미사용 CSS 제거
│   │   └── portable.js       # 포터블 CSS 생성기
│   └── config/
│       ├── tailwind.config.js # Tailwind 호환 설정
│       └── design-tokens.json # 디자인 토큰 정의
└── web/                       # 기존 웹 자산
    ├── app.css               # 기존 파일 (단계적 교체 예정)
    └── progressive.min.css   # 새 시스템으로 교체
```

## 📐 설계 원칙

### 1. Utility-First 원칙
```css
/* ❌ 기존 방식 - 컴포넌트 기반 */
.menu-button {
  background: var(--grid-cell-bg);
  border: 1px solid var(--grid-border);
  padding: var(--spacing-xs) var(--spacing-sm);
  font-size: var(--font-size-sm);
}

/* ✅ 새로운 방식 - 유틸리티 조합 */
.bg-white .border .border-gray-300 .px-2 .py-1 .text-sm
```

### 2. 포터블 CSS 설계
```css
/* 일반 버전 - CSS 변수 사용 */
.text-excel-primary { color: var(--excel-primary); }

/* 포터블 버전 - 하드코딩된 값 */
.text-excel-primary { color: #217346; }
```

### 3. 반응형 디자인 시스템
```css
/* 모바일 퍼스트 */
.p-2 { padding: 0.5rem; }

/* 반응형 변형 */
@media (min-width: 768px) {
  .md:p-4 { padding: 1rem; }
}
@media (min-width: 1024px) {
  .lg:p-6 { padding: 1.5rem; }
}
```

## 🎨 디자인 토큰 시스템

### design-tokens.json
```json
{
  "colors": {
    "excel": {
      "primary": "#217346",
      "primaryHover": "#1e6b40",
      "secondary": "#f2f8f4",
      "accent": "#0078d4",
      "accentHover": "#106ebe"
    },
    "status": {
      "success": "#28a745",
      "warning": "#ffc107", 
      "danger": "#dc3545",
      "info": "#17a2b8"
    },
    "grid": {
      "border": "#d0d7de",
      "headerBg": "#f6f8fa",
      "headerText": "#24292f",
      "cellBg": "#ffffff",
      "cellSelected": "#cce7ff",
      "cellEditing": "#fff4ce",
      "cellHover": "#f6f8fa"
    }
  },
  "spacing": {
    "scale": [0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80, 96],
    "unit": "px"
  },
  "typography": {
    "fontSizes": {
      "xs": "11px",
      "sm": "12px", 
      "base": "14px",
      "lg": "16px",
      "xl": "18px"
    },
    "fontWeights": {
      "normal": 400,
      "medium": 500,
      "semibold": 600,
      "bold": 700
    },
    "lineHeights": {
      "none": 1,
      "tight": 1.25,
      "normal": 1.5,
      "relaxed": 1.75
    }
  },
  "borderRadius": {
    "sm": "3px",
    "md": "6px", 
    "lg": "8px",
    "xl": "12px",
    "full": "9999px"
  },
  "shadows": {
    "sm": "0 1px 3px rgba(0, 0, 0, 0.1)",
    "md": "0 4px 6px rgba(0, 0, 0, 0.1)",
    "lg": "0 10px 15px rgba(0, 0, 0, 0.1)"
  }
}
```

## 🛠️ 빌드 시스템 설계

### 빌드 파이프라인
```
1. [디자인 토큰] → generate-utils.js → [유틸리티 CSS 생성]
2. [소스 CSS들] → build.js → [통합 CSS 파일]
3. [통합 CSS] → purge.js → [미사용 CSS 제거]
4. [정리된 CSS] → portable.js → [포터블 CSS 생성]
5. [최종 CSS] → 압축 → [배포용 파일들]
```

### build.js 설계
```javascript
// 의사코드 - 실제 구현은 하지 않음
const BuildSystem = {
  // 1. 유틸리티 클래스 자동 생성
  generateUtilities() {
    // spacing: .p-0 ~ .p-24, .m-0 ~ .m-24, .px-*, .py-*, .pt-* 등
    // colors: .text-*, .bg-*, .border-*
    // typography: .text-xs ~ .text-xl, .font-*, .leading-*
    // layout: .flex, .grid, .items-*, .justify-*
    // borders: .border, .border-*, .rounded-*
    // effects: .shadow-*, .transition-*, .opacity-*
  },

  // 2. CSS 파일 병합 및 최적화
  bundleCSS() {
    // reset.css + variables.css + base.css + utilities/* + components/*
    // PostCSS로 처리 (autoprefixer, cssnano)
  },

  // 3. 포터블 버전 생성
  createPortable() {
    // CSS 변수를 실제 값으로 치환
    // 외부 의존성 제거
    // 독립 실행 가능한 CSS 생성
  },

  // 4. 미사용 CSS 제거
  purgeUnusedCSS() {
    // Go 파일들 스캔해서 사용된 클래스만 추출
    // PurgeCSS 또는 커스텀 파서 사용
  }
};
```

## 🔄 마이그레이션 전략

### Phase 1: 기반 구축 (설계만)
```
1. css-system/ 디렉토리 구조 생성
2. 현재 app.css에서 CSS 변수들 추출 → variables.css
3. 디자인 토큰 JSON 파일 생성
4. 유틸리티 생성기 로직 설계
```

### Phase 2: 유틸리티 시스템 (설계만)
```
1. spacing.css - .p-*, .m-* 클래스들
2. typography.css - .text-*, .font-* 클래스들  
3. layout.css - .flex, .grid 관련 클래스들
4. colors.css - 색상 유틸리티들
5. 빌드 스크립트 작성
```

### Phase 3: 컴포넌트 마이그레이션 (설계만)
```
1. 기존 .menu-button → 유틸리티 조합으로 변경
2. .grid-cell → .grid-cell-* 유틸리티로 분해
3. .modal-* → 유틸리티 조합으로 변경
4. Go 템플릿에서 클래스 이름 업데이트
```

### Phase 4: 최적화 및 배포 (설계만)
```
1. PurgeCSS 설정으로 사용하지 않는 CSS 제거
2. 포터블 CSS 빌드 시스템 완성
3. CDN 배포를 위한 파일 준비
4. npm 패키지 형태 배포 준비
```

## 📊 예상 성능 개선

### 파일 크기 비교
```
현재 app.css: 38KB (1,810 라인)
↓
새 시스템:
- progressive.css: ~25KB (유틸리티 + 컴포넌트)
- progressive.min.css: ~15KB (압축 후)
- 실제 사용: ~8KB (PurgeCSS 적용 후)
- portable.css: ~12KB (독립형 버전)
```

### 개발 경험 개선
```
❌ 현재: 새 스타일 → CSS 파일 수정 → 클래스 생성
✅ 미래: 유틸리티 조합으로 즉시 스타일링
```

## 🎁 포터블 CSS 배포 형태

### CDN 배포 (설계)
```html
<!-- 기본 버전 -->
<link href="https://cdn.progressive-ui.com/v1.0/progressive.min.css" rel="stylesheet">

<!-- 포터블 버전 (의존성 없음) -->
<link href="https://cdn.progressive-ui.com/v1.0/portable.min.css" rel="stylesheet">

<!-- 커스텀 테마 버전 -->
<link href="https://cdn.progressive-ui.com/v1.0/progressive.min.css?theme=dark" rel="stylesheet">
```

### npm 패키지 (설계)
```json
{
  "name": "@progressive/ui-css",
  "version": "1.0.0",
  "description": "Progressive Spreadsheet UI CSS Framework",
  "main": "dist/progressive.css",
  "files": ["dist/"],
  "exports": {
    "./portable": "./dist/portable.css",
    "./full": "./dist/progressive.css",
    "./utilities": "./dist/utilities.css"
  }
}
```

## 🔧 Go 통합 설계

### 유틸리티 헬퍼 함수 (설계)
```go
// 설계만 - 실제 구현 안함
package ui

// CSS 클래스 조합 헬퍼
func Classes(classes ...string) string {
    return strings.Join(classes, " ")
}

// 자주 사용하는 조합들을 미리 정의
var (
    ButtonPrimary = Classes("bg-excel-primary", "text-white", "px-4", "py-2", "rounded", "hover:bg-excel-primary-hover")
    ButtonSecondary = Classes("bg-gray-200", "text-gray-800", "px-4", "py-2", "rounded", "hover:bg-gray-300")
    GridCell = Classes("border", "border-gray-300", "p-2", "text-sm", "hover:bg-gray-50")
)

// 조건부 클래스 적용
func ConditionalClass(condition bool, ifTrue, ifFalse string) string {
    if condition {
        return ifTrue
    }
    return ifFalse
}
```

### Go 템플릿 사용 예시 (설계)
```go
// Before (기존 방식)
app.Button().Class("menu-button").Text("📄 New")

// After (유틸리티 방식)  
app.Button().
    Class("bg-white border border-gray-300 rounded px-2 py-1 text-sm flex items-center gap-1 hover:bg-gray-50 hover:border-blue-500 transition-all duration-200").
    Text("📄 New")

// 또는 헬퍼 사용
app.Button().Class(ui.ButtonSecondary).Text("📄 New")
```

## 📝 결론

이 설계는 다음과 같은 이점을 제공합니다:

1. **개발 효율성**: 미리 정의된 유틸리티 클래스로 빠른 스타일링
2. **일관성**: 디자인 토큰 기반의 일관된 디자인 시스템
3. **성능**: PurgeCSS로 75% 크기 감소 (38KB → ~8KB)
4. **포터블리티**: 외부 의존성 없는 독립형 CSS 제공
5. **유지보수성**: 중앙화된 디자인 토큰 관리

이 설계 문서를 기반으로 단계적으로 구현할 수 있으며, 각 단계별로 점진적인 마이그레이션이 가능합니다.