# Tailwind CSS CLI Integration Design

## 🎯 목표
기존에 잘 만들어진 Tailwind CSS CLI를 활용하여 Progressive Spreadsheet 프로젝트에 유틸리티-first CSS 시스템 도입

## 🛠️ Tailwind CSS CLI 통합 설계

### 1. 프로젝트 구조 설계
```
progressive/
├── package.json              # Node.js 의존성 관리
├── tailwind.config.js        # Tailwind 설정 (커스텀 테마)
├── input.css                 # Tailwind 진입점 CSS
├── web/
│   ├── tailwind.css          # Tailwind 빌드 결과물
│   ├── tailwind.min.css      # 압축 버전
│   └── app.css               # 기존 CSS (단계적 교체)
├── Makefile                  # 빌드 프로세스 통합
└── .github/workflows/        # CI/CD 자동화
    └── build-css.yml
```

### 2. Tailwind CLI 설치 및 설정

#### package.json 설계
```json
{
  "name": "progressive-spreadsheet",
  "version": "1.0.0",
  "description": "Progressive Web Application for Data Management",
  "scripts": {
    "css:dev": "tailwindcss -i input.css -o web/tailwind.css --watch",
    "css:build": "tailwindcss -i input.css -o web/tailwind.css --minify",
    "css:portable": "tailwindcss -i input.css -o web/tailwind.css --minify --purge"
  },
  "devDependencies": {
    "tailwindcss": "^3.4.0",
    "@tailwindcss/forms": "^0.5.7",
    "@tailwindcss/typography": "^0.5.10"
  }
}
```

#### tailwind.config.js 설계 - Excel 테마 적용
```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./spreadsheet/**/*.go",    // Go 템플릿에서 클래스 추출
    "./web/**/*.html",
    "./web/**/*.js"
  ],
  theme: {
    extend: {
      // 기존 CSS 변수들을 Tailwind 테마로 변환
      colors: {
        excel: {
          primary: '#217346',
          'primary-hover': '#1e6b40',
          secondary: '#f2f8f4',
          accent: '#0078d4',
          'accent-hover': '#106ebe',
        },
        grid: {
          border: '#d0d7de',
          'header-bg': '#f6f8fa',
          'header-text': '#24292f',
          'cell-bg': '#ffffff',
          'cell-selected': '#cce7ff',
          'cell-editing': '#fff4ce',
          'cell-hover': '#f6f8fa',
        },
        status: {
          success: '#28a745',
          warning: '#ffc107',
          danger: '#dc3545',
          info: '#17a2b8',
        }
      },
      spacing: {
        // 기존 spacing 변수들
        'header': '60px',
        'toolbar': '48px',
        'statusbar': '32px',
        'sidebar': '250px',
      },
      fontSize: {
        // 기존 font-size 변수들  
        'xs': '11px',
        'sm': '12px',
        'base': '14px',
        'lg': '16px',
        'xl': '18px',
      },
      fontFamily: {
        'main': ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial', 'sans-serif'],
        'mono': ['SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', 'Consolas', 'Courier New', 'monospace'],
      },
      boxShadow: {
        'sm': '0 1px 3px rgba(0, 0, 0, 0.1)',
        'md': '0 4px 6px rgba(0, 0, 0, 0.1)',
        'lg': '0 10px 15px rgba(0, 0, 0, 0.1)',
      },
      borderRadius: {
        'sm': '3px',
        'md': '6px', 
        'lg': '8px',
      },
      transitionDuration: {
        'fast': '150ms',
        'base': '250ms',
        'slow': '350ms',
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),        // 폼 스타일링
    require('@tailwindcss/typography'),   // 타이포그래피
    
    // 커스텀 플러그인 - Excel 전용 유틸리티
    function({ addUtilities }) {
      const newUtilities = {
        '.grid-cell': {
          display: 'flex',
          alignItems: 'center',
          padding: '4px 8px',
          borderRight: '1px solid #d0d7de',
          minWidth: '120px',
          maxWidth: '200px',
          fontSize: '12px',
          lineHeight: '1.4',
          textOverflow: 'ellipsis',
          overflow: 'hidden',
          whiteSpace: 'nowrap',
          cursor: 'cell',
          transition: 'all 150ms',
          userSelect: 'none',
        },
        '.grid-cell-selected': {
          backgroundColor: '#cce7ff !important',
          border: '2px solid #0078d4 !important',
          boxShadow: 'inset 0 0 0 1px #0078d4',
          zIndex: '10',
          position: 'relative',
        },
        '.excel-btn': {
          display: 'inline-flex',
          alignItems: 'center',
          gap: '4px',
          padding: '4px 8px',
          backgroundColor: '#ffffff',
          border: '1px solid #d0d7de',
          borderRadius: '4px',
          fontSize: '12px',
          cursor: 'pointer',
          transition: 'all 200ms',
          fontFamily: 'inherit',
        }
      }
      addUtilities(newUtilities)
    }
  ],
}
```

#### input.css 설계
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

/* 커스텀 CSS 변수 (포터블리티를 위해) */
@layer base {
  :root {
    --excel-primary: theme('colors.excel.primary');
    --excel-accent: theme('colors.excel.accent');
    --grid-border: theme('colors.grid.border');
    /* 기존 CSS 변수들을 Tailwind 테마와 연결 */
  }
}

/* 기존 컴포넌트 중 꼭 필요한 것들만 유지 */
@layer components {
  .schema-banner {
    @apply flex items-center gap-3 p-3 mb-4 rounded-lg text-sm;
  }
  
  .schema-banner.default {
    @apply bg-blue-50 border border-blue-200 text-blue-700;
  }
  
  .schema-banner.loaded {
    @apply bg-green-50 border border-green-200 text-green-700;
  }
}
```

### 3. 빌드 시스템 통합

#### Makefile 수정
```makefile
# 기존 Makefile에 CSS 빌드 추가
.PHONY: css-dev css-build css-watch

# CSS 개발 모드 (watch)
css-dev:
	npm run css:dev

# CSS 프로덕션 빌드
css-build:
	npm run css:build

# CSS watch 모드
css-watch:
	npm run css:dev &

# 전체 개발 환경 실행 (CSS watch + Go server)
dev: css-watch
	go run *.go

# 프로덕션 빌드 (CSS + WASM)
build: css-build
	GOOS=js GOARCH=wasm go build -o web/app.wasm *.go
	@echo "Build complete: CSS + WebAssembly"

# 클린
clean:
	rm -f web/app.wasm web/tailwind.css web/tailwind.min.css
```

### 4. 기존 CSS → Tailwind 마이그레이션 전략

#### 4.1. 클래스 매핑 테이블
```
기존 CSS                    →  Tailwind 클래스
.menu-button               →  bg-white border border-gray-300 rounded px-2 py-1 text-sm hover:bg-gray-50
.grid-cell                 →  grid-cell (커스텀 유틸리티) 또는 flex items-center p-1 border-r text-xs
.modal-overlay             →  fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50
.status-indicator.saved    →  text-green-600 bg-green-100 px-2 py-1 rounded text-xs
.excel-editor              →  flex flex-col h-screen bg-white font-main text-sm text-gray-800
```

#### 4.2. Go 템플릿 변경 예시
```go
// Before
app.Button().Class("menu-button").Text("📄 New")

// After  
app.Button().
    Class("bg-white border border-gray-300 rounded px-2 py-1 text-sm cursor-pointer transition-all duration-200 flex items-center gap-1 hover:bg-gray-50 hover:border-blue-500").
    Text("📄 New")

// 또는 @layer components에 정의한 클래스 사용
app.Button().Class("excel-btn excel-btn-primary").Text("📄 New")
```

### 5. 포터블 CSS 생성 전략

#### 5.1. Standalone Tailwind Build
```bash
# 포터블 CSS 빌드 (CDN 없이 사용 가능)
npx tailwindcss -i input.css -o web/portable.css --minify

# 사용하지 않는 CSS 완전 제거
npx tailwindcss -i input.css -o web/optimized.css --minify --purge
```

#### 5.2. 포터블 HTML 템플릿
```html
<!DOCTYPE html>
<html>
<head>
    <!-- 외부 CDN 없이 로컬 CSS 사용 -->
    <link href="./portable.css" rel="stylesheet">
</head>
<body class="bg-gray-50 font-main">
    <div class="excel-editor">
        <!-- Tailwind 클래스로 스타일링 -->
    </div>
</body>
</html>
```

### 6. 성능 최적화

#### 6.1. PurgeCSS 설정 (자동 포함)
```javascript
// tailwind.config.js에서 content 경로 지정으로 자동 purge
module.exports = {
  content: [
    "./spreadsheet/**/*.go",    // Go 파일에서 클래스 추출
    "./web/**/*.html",
  ],
  // ...
}
```

#### 6.2. 예상 파일 크기
```
현재 app.css: 38KB
↓
Tailwind 개발용: ~3.5MB (모든 유틸리티)
Tailwind 프로덕션: ~15KB (사용된 클래스만)
최적화 후: ~8KB (실제 사용 클래스 + 커스텀)
```

### 7. CI/CD 자동화

#### .github/workflows/build-css.yml
```yaml
name: Build CSS

on:
  push:
    paths:
      - 'spreadsheet/**/*.go'
      - 'input.css'
      - 'tailwind.config.js'

jobs:
  build-css:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Install dependencies
        run: npm install
        
      - name: Build Tailwind CSS
        run: npm run css:build
        
      - name: Commit generated CSS
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add web/tailwind.css
          git diff --staged --quiet || git commit -m "Auto-build Tailwind CSS"
          git push
```

### 8. 장점 정리

#### ✅ Tailwind CLI 사용의 장점
1. **검증된 프레임워크**: 수백만 개발자가 사용하는 안정적인 도구
2. **자동 최적화**: PurgeCSS, 압축, prefixing 자동 처리
3. **풍부한 생태계**: 플러그인, 템플릿, 커뮤니티 지원
4. **지속적인 업데이트**: 최신 CSS 기능 및 성능 최적화
5. **툴링**: VS Code 확장, IntelliSense, 자동완성
6. **문서화**: 완벽한 공식 문서 및 학습 자료

#### ✅ Progressive 프로젝트 맞춤화
1. **Excel 테마**: 기존 CSS 변수들을 Tailwind 테마로 완벽 변환
2. **Go 통합**: Go 템플릿에서 Tailwind 클래스 자동 추출
3. **포터블 빌드**: 외부 의존성 없는 독립형 CSS 생성
4. **점진적 마이그레이션**: 기존 CSS와 병행 사용 가능

### 9. 구현 로드맵

#### Week 1: 설정 및 기반 구축
- package.json, tailwind.config.js 생성
- input.css 작성 및 기본 빌드 확인
- Makefile CSS 빌드 통합

#### Week 2: 테마 및 커스텀 유틸리티
- Excel 색상 테마 Tailwind로 변환
- grid-cell, excel-btn 등 커스텀 컴포넌트 정의
- 기존 CSS 변수와 Tailwind 연동

#### Week 3-4: 컴포넌트 마이그레이션  
- 주요 컴포넌트들 Tailwind 클래스로 변환
- Go 템플릿 클래스명 업데이트
- 기능별 테스트 및 검증

#### Week 5: 최적화 및 배포
- PurgeCSS 최적화
- 포터블 CSS 빌드 스크립트
- CI/CD 자동화 설정

## 결론

커스텀 프레임워크 대신 **Tailwind CSS CLI를 활용하는 것이 훨씬 현명한 선택**입니다:

- ⚡ **빠른 구현**: 기존 도구 활용으로 개발 시간 단축
- 🛡️ **안정성**: 검증된 프레임워크로 버그 위험 최소화  
- 📈 **확장성**: 풍부한 플러그인 생태계 활용
- 🔧 **유지보수**: 커뮤니티 지원 및 지속적인 업데이트
- 💰 **비용 효율**: 개발 리소스를 핵심 기능에 집중