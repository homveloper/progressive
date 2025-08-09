# Templ Syntax Reference Guide

Go용 템플릿 엔진인 templ의 문법과 사용법에 대한 종합 참조 문서입니다.

## 목차
1. [기본 구조](#기본-구조)
2. [표현식](#표현식)
3. [속성](#속성)
4. [조건문](#조건문)
5. [반복문](#반복문)
6. [컴포넌트 구성](#컴포넌트-구성)
7. [JavaScript 처리](#javascript-처리)

---

## 기본 구조

### 패키지 선언 및 Import
```templ
package main

import "fmt"
import "time"
import "strings"
```

### 기본 컴포넌트 구조
```templ
templ componentName(param1 string, param2 int) {
    <div>
        <h1>{ param1 }</h1>
        <p>Value: { param2 }</p>
    </div>
}
```

### Go 코드와 함께 사용
```templ
package main

// 일반 Go 코드
var greeting = "환영합니다!"

// templ 컴포넌트
templ headerTemplate(name string) {
    <header>
        <h1>{ name }</h1>
        <h2>{ greeting }</h2>
    </header>
}
```

---

## 표현식

### 지원되는 데이터 타입
- 문자열 (string)
- 숫자 (int, uint, float32, complex64)
- 불린 (boolean)
- 위 타입들을 기반으로 한 사용자 정의 타입

### 표현식 예제
```templ
templ expressionExamples(name string, age int) {
    <!-- 리터럴 -->
    <div>{ "안녕하세요" }</div>
    <div>{ 42 }</div>
    
    <!-- 변수 -->
    <div>{ name }</div>
    <div>나이: { age }세</div>
    
    <!-- 함수 호출 -->
    <div>{ strings.ToUpper("hello") }</div>
    <div>{ getString() }</div>
    
    <!-- 자동 HTML 이스케이핑 -->
    <div>{ `<script>alert('XSS')</script>` }</div>
}
```

### 보안 특징
- 모든 표현식은 XSS 방지를 위해 자동으로 HTML 이스케이프됩니다
- CSS 인젝션 공격도 방지됩니다

---

## 속성

### 정적 속성
```templ
<p data-testid="paragraph">텍스트</p>
<div class="container">내용</div>
```

### 동적 속성
```templ
templ dynamicAttributes(testID string, isActive bool) {
    <!-- 문자열 표현식 속성 -->
    <p data-testid={ testID }>텍스트</p>
    
    <!-- 불린 속성 -->
    <input disabled?={ isActive } />
    <hr noshade?={ false } />
}
```

### 조건부 속성
```templ
templ conditionalAttributes(showBorder bool) {
    <div 
        class="base-class"
        if showBorder {
            style="border: 1px solid #ccc"
        }
    >
        내용
    </div>
}
```

### 스프레드 속성
```templ
templ spreadAttributes(attrs map[string]any) {
    <p { attrs... }>동적 속성이 적용된 텍스트</p>
}
```

### 특별한 속성 처리
- URL 속성: 자동 검증 및 안전한 프로토콜 확인
- JavaScript 속성: 스크립트 템플릿 참조 필요
- JSON 속성: Go 데이터 구조에서 직렬화 가능

---

## 조건문

### if-else 구문
```templ
templ conditionalRendering(isLoggedIn bool, userType string) {
    if isLoggedIn {
        <div>환영합니다!</div>
        if userType == "admin" {
            <div>관리자 권한</div>
        } else {
            <div>일반 사용자</div>
        }
    } else {
        <input name="login" type="button" value="로그인"/>
    }
}
```

### 복잡한 조건
```templ
templ complexConditions(user User) {
    if user.IsActive && user.Role == "admin" {
        <div>활성 관리자</div>
    } else if user.IsActive {
        <div>활성 사용자</div>
    } else {
        <div>비활성 사용자</div>
    }
}
```

---

## 반복문

### 기본 for 반복문
```templ
templ listItems(items []Item) {
    <ul>
        for _, item := range items {
            <li>{ item.Name }</li>
        }
    </ul>
}
```

### 인덱스를 사용한 반복문
```templ
templ indexedList(items []string) {
    <ol>
        for i, item := range items {
            <li>{ i+1 }. { item }</li>
        }
    </ol>
}
```

### Map 반복문
```templ
templ mapIteration(data map[string]string) {
    <dl>
        for key, value := range data {
            <dt>{ key }</dt>
            <dd>{ value }</dd>
        }
    </dl>
}
```

### 조건과 결합한 반복문
```templ
templ conditionalLoop(users []User) {
    <ul>
        for _, user := range users {
            if user.IsActive {
                <li class="active">{ user.Name }</li>
            } else {
                <li class="inactive">{ user.Name }</li>
            }
        }
    </ul>
}
```

---

## 컴포넌트 구성

### 기본 컴포넌트 호출
```templ
templ layout() {
    @header()
    @main()
    @footer()
}
```

### 매개변수가 있는 컴포넌트 호출
```templ
templ page(title string, content string) {
    @headerWithTitle(title)
    @contentSection(content)
}

templ headerWithTitle(title string) {
    <header>
        <h1>{ title }</h1>
    </header>
}
```

### Children 전달
```templ
templ wrapper() {
    <div id="wrapper">
        { children... }
    </div>
}

templ pageWithWrapper() {
    @wrapper() {
        <p>이 내용이 wrapper의 children으로 전달됩니다</p>
    }
}
```

### 컴포넌트를 매개변수로 전달
```templ
templ layout(content templ.Component) {
    <html>
        <body>
            <div id="main">
                @content
            </div>
        </body>
    </html>
}
```

### 컴포넌트 조인
```templ
templ combinedComponents() {
    @templ.Join(hello(), world(), footer())
}
```

### 패키지 간 컴포넌트 공유
```templ
// 컴포넌트명을 대문자로 시작하여 export
templ SharedHeader(title string) {
    <header>
        <h1>{ title }</h1>
    </header>
}
```

---

## JavaScript 처리

### 기본 스크립트 포함
```templ
templ pageWithScript() {
    <html>
        <head>
            <script>
                console.log("페이지 로드됨");
            </script>
        </head>
    </html>
}
```

### Go 데이터를 JavaScript에 전달

#### 1. JSON 문자열 방식
```templ
templ dataToJS(user User) {
    <script>
        const userData = { templ.JSONString(user) };
        console.log(userData);
    </script>
}
```

#### 2. 스크립트 요소에 데이터 삽입
```templ
templ scriptData(config Config) {
    <script type="application/json" id="config">
        { templ.JSONString(config) }
    </script>
    <script>
        const config = JSON.parse(document.getElementById('config').textContent);
    </script>
}
```

#### 3. 함수 호출 방식
```templ
templ jsFunction(message string) {
    <script>
        function showMessage() {
            alert({ templ.JSONString(message) });
        }
        showMessage();
    </script>
}
```

### JavaScript 객체 생성
```templ
templ jsObjectCreation(items []Item) {
    <script>
        const items = {
            for i, item := range items {
                { item.ID }: {
                    name: { templ.JSONString(item.Name) },
                    value: { templ.JSONString(item.Value) }
                },
            }
        };
    </script>
}
```

---

## 실용적인 예제

### 폼 컴포넌트
```templ
templ form(action string, fields []FormField) {
    <form action={ action } method="post">
        for _, field := range fields {
            <div class="field">
                <label for={ field.ID }>{ field.Label }</label>
                <input 
                    type={ field.Type }
                    id={ field.ID }
                    name={ field.Name }
                    required?={ field.Required }
                />
            </div>
        }
        <button type="submit">제출</button>
    </form>
}
```

### 테이블 컴포넌트
```templ
templ dataTable(headers []string, rows [][]string) {
    <table class="data-table">
        <thead>
            <tr>
                for _, header := range headers {
                    <th>{ header }</th>
                }
            </tr>
        </thead>
        <tbody>
            for _, row := range rows {
                <tr>
                    for _, cell := range row {
                        <td>{ cell }</td>
                    }
                </tr>
            }
        </tbody>
    </table>
}
```

### 카드 리스트
```templ
templ cardList(items []Card) {
    <div class="card-grid">
        for _, item := range items {
            <div class="card">
                if item.ImageURL != "" {
                    <img src={ item.ImageURL } alt={ item.Title }/>
                }
                <div class="card-content">
                    <h3>{ item.Title }</h3>
                    <p>{ item.Description }</p>
                    if item.Price > 0 {
                        <div class="price">{ fmt.Sprintf("%.2f원", item.Price) }</div>
                    }
                </div>
            </div>
        }
    </div>
}
```

---

## 베스트 프랙티스

### 1. 컴포넌트 구조화
- 컴포넌트는 단일 책임을 가지도록 설계
- 재사용 가능한 작은 컴포넌트로 분할
- 매개변수는 명확하고 타입 안전하게 정의

### 2. 보안
- 모든 사용자 입력은 자동으로 이스케이프됨
- URL 속성은 안전한 프로토콜만 허용
- JavaScript에 전달하는 데이터는 templ.JSONString() 사용

### 3. 성능
- 불필요한 반복문 중첩 피하기
- 큰 데이터셋은 페이지네이션 고려
- CSS와 JavaScript는 외부 파일로 분리

### 4. 유지보수성
- 컴포넌트명은 명확하고 설명적으로
- 복잡한 로직은 Go 함수로 분리
- 템플릿 구조는 HTML 시맨틱을 따르기

---

## 빌드 및 사용

### 코드 생성
```bash
templ generate
```

### 파일 감시
```bash
templ generate --watch
```

### Go 코드에서 사용
```go
func handler(w http.ResponseWriter, r *http.Request) {
    component := myTemplate("Hello", 42)
    component.Render(r.Context(), w)
}
```

이 문서는 templ 사용 시 참조할 수 있는 종합적인 가이드입니다. 프로젝트 개발 중에 필요한 문법과 패턴을 빠르게 찾아볼 수 있습니다.