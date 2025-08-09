# Progressive Utility CSS System Design

## 2. 유틸리티 클래스 시스템 설계

### Spacing System (Tailwind 스타일)
```css
/* 간격 시스템 - 8px 기반 */
.p-0 { padding: 0; }
.p-1 { padding: 0.25rem; }    /* 4px */
.p-2 { padding: 0.5rem; }     /* 8px */
.p-3 { padding: 0.75rem; }    /* 12px */
.p-4 { padding: 1rem; }       /* 16px */
.p-5 { padding: 1.25rem; }    /* 20px */
.p-6 { padding: 1.5rem; }     /* 24px */

/* 방향별 패딩 */
.px-4 { padding-left: 1rem; padding-right: 1rem; }
.py-2 { padding-top: 0.5rem; padding-bottom: 0.5rem; }

/* 마진도 동일한 시스템 */
.m-0, .m-1, .m-2, .mx-4, .my-2 등...
```

### Color System
```css
/* 색상 시스템 - Excel 테마 기반 */
.text-excel-primary { color: var(--excel-primary); }
.text-excel-secondary { color: var(--excel-secondary); }
.bg-excel-primary { background-color: var(--excel-primary); }
.bg-grid-header { background-color: var(--grid-header-bg); }

/* 상태 색상 */
.text-success { color: var(--status-success); }
.text-warning { color: var(--status-warning); }
.text-danger { color: var(--status-danger); }
.text-info { color: var(--status-info); }
```

### Layout Utilities
```css
/* Flexbox */
.flex { display: flex; }
.flex-col { flex-direction: column; }
.flex-row { flex-direction: row; }
.items-center { align-items: center; }
.justify-between { justify-content: space-between; }
.justify-center { justify-content: center; }

/* Grid */
.grid { display: grid; }
.grid-cols-5 { grid-template-columns: repeat(5, minmax(0, 1fr)); }
.col-span-2 { grid-column: span 2 / span 2; }

/* Position */
.relative { position: relative; }
.absolute { position: absolute; }
.fixed { position: fixed; }
```

### Typography System
```css
/* 폰트 크기 */
.text-xs { font-size: var(--font-size-xs); }
.text-sm { font-size: var(--font-size-sm); }
.text-base { font-size: var(--font-size-base); }
.text-lg { font-size: var(--font-size-lg); }
.text-xl { font-size: var(--font-size-xl); }

/* 폰트 굵기 */
.font-normal { font-weight: 400; }
.font-medium { font-weight: 500; }
.font-semibold { font-weight: 600; }
.font-bold { font-weight: 700; }
```

## 3. 포터블 CSS 시스템

### Standalone CSS 생성
```css
/* portable.css - 외부 의존성 없는 독립형 CSS */
:root {
  /* 모든 CSS 변수를 하드코딩된 값으로 변환 */
  --p-primary: #217346;
  --p-spacing-1: 0.25rem;
  --p-spacing-2: 0.5rem;
  /* ... */
}

.p-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.5rem 1rem;
  background: #217346;
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.p-btn:hover {
  background: #1e6b40;
  transform: translateY(-1px);
}
```

## 4. 빌드 시스템 설계

### 자동화된 CSS 생성
```javascript
// build.js
const fs = require('fs');
const postcss = require('postcss');
const autoprefixer = require('autoprefixer');
const cssnano = require('cssnano');

const buildConfig = {
  // 유틸리티 생성 규칙
  spacing: {
    values: [0, 1, 2, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24, 32],
    properties: ['margin', 'padding'],
    directions: ['', 'x', 'y', 't', 'r', 'b', 'l']
  },
  
  colors: {
    'excel-primary': '#217346',
    'excel-secondary': '#f2f8f4',
    'success': '#28a745',
    'warning': '#ffc107',
    'danger': '#dc3545'
  },
  
  // 반응형 브레이크포인트
  breakpoints: {
    'sm': '640px',
    'md': '768px',
    'lg': '1024px',
    'xl': '1280px'
  }
};

function generateUtilities() {
  let css = '';
  
  // 간격 유틸리티 생성
  buildConfig.spacing.values.forEach(value => {
    const rem = value * 0.25;
    css += `.p-${value} { padding: ${rem}rem; }\n`;
    css += `.m-${value} { margin: ${rem}rem; }\n`;
  });
  
  return css;
}
```

## 5. 현재 시스템 → 유틸리티 마이그레이션 예시

### Before (Component CSS)
```css
.menu-button {
  background: var(--grid-cell-bg);
  border: 1px solid var(--grid-border);
  border-radius: 4px;
  padding: var(--spacing-xs) var(--spacing-sm);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 4px;
}
```

### After (Utility Classes)
```html
<button class="bg-white border border-gray-300 rounded px-2 py-1 text-sm cursor-pointer transition-all duration-200 flex items-center gap-1 hover:bg-gray-50 hover:border-blue-500">
  📄 New
</button>
```

## 6. 포터블 CSS 배포 전략

### CDN 스타일 배포
```html
<!-- 포터블 버전 - 외부 의존성 없음 -->
<link href="https://cdn.progressive-css.com/v1/portable.min.css" rel="stylesheet">

<!-- 커스터마이즈 가능한 버전 -->
<link href="https://cdn.progressive-css.com/v1/progressive.min.css" rel="stylesheet">
<script>
  ProgressiveCSS.config({
    theme: {
      primary: '#217346',
      spacing: '8px'
    }
  });
</script>
```

### npm 패키지 형태
```json
{
  "name": "@progressive/ui-css",
  "version": "1.0.0",
  "files": [
    "dist/progressive.css",
    "dist/progressive.min.css", 
    "dist/portable.css",
    "dist/portable.min.css"
  ],
  "exports": {
    "./portable": "./dist/portable.css",
    "./full": "./dist/progressive.css"
  }
}
```

## 7. 성능 최적화

### CSS 크기 비교 (예상)
- **현재 시스템**: 38KB (압축 전)
- **유틸리티 시스템**: 
  - 전체: ~25KB
  - 사용된 클래스만: ~8KB (PurgeCSS 적용)
  - 포터블 버전: ~12KB

### Tree Shaking & PurgeCSS
```javascript
// 사용된 클래스만 추출
const purgecss = require('@fullhuman/postcss-purgecss');

module.exports = {
  plugins: [
    purgecss({
      content: ['./spreadsheet/**/*.go', './web/**/*.html'],
      extractors: [
        {
          extractor: /class="([^"]+)"/g,
          extensions: ['go', 'html']
        }
      ]
    })
  ]
};
```

## 8. Go 템플릿과의 통합

### Go 코드에서 유틸리티 클래스 사용
```go
// Before
app.Button().Class("menu-button").Text("📄 New")

// After - 유틸리티 클래스 사용
app.Button().
    Class("bg-white border border-gray-300 rounded px-2 py-1 text-sm flex items-center gap-1 hover:bg-gray-50").
    Text("📄 New")

// Helper 함수 생성
func MenuButton(text string) app.HTMLButton {
    return app.Button().
        Class("btn btn-menu"). // 조합된 유틸리티 클래스
        Text(text)
}
```

## 9. 장점과 단점

### ✅ 장점
- **크기 감소**: PurgeCSS로 사용하지 않는 CSS 제거
- **일관성**: 디자인 시스템 기반 일관된 스타일
- **개발 속도**: 미리 정의된 클래스로 빠른 개발
- **포터블리티**: 독립형 CSS 파일 제공 가능
- **유지보수**: 중앙화된 디자인 토큰

### ❌ 단점
- **HTML 복잡성**: 클래스 이름이 길어질 수 있음
- **초기 학습**: 새로운 클래스 체계 학습 필요
- **빌드 도구**: 빌드 시스템 구축 필요

## 10. 구현 로드맵

1. **Phase 1**: 유틸리티 시스템 구축 (2주)
2. **Phase 2**: 현재 컴포넌트 마이그레이션 (3주)
3. **Phase 3**: 포터블 CSS 빌드 시스템 (1주)
4. **Phase 4**: 성능 최적화 및 문서화 (1주)