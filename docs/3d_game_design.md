# Go-App 3D 게임 구현 설계

## 아키텍처 개요

### 하이브리드 아키텍처: Go-App + Three.js

```
Go (WASM) ←→ JavaScript Bridge ←→ Three.js (WebGL)
```

## 구현 방법

### 1. JavaScript 브릿지 방식 (추천)

#### 장점
- Three.js의 강력한 3D 기능 활용
- Go의 게임 로직과 상태 관리
- 검증된 WebGL 라이브러리 사용
- 풍부한 3D 자산 생태계

#### 구조
```go
// main.go
type Game3D struct {
    app.Compo
    gameState *GameState
    renderer  *ThreeJSRenderer
}

type ThreeJSRenderer struct {
    scene    js.Value
    camera   js.Value
    renderer js.Value
}

func (g *Game3D) OnMount(ctx app.Context) {
    // Three.js 초기화
    g.initThreeJS()
    // 게임 루프 시작
    go g.gameLoop(ctx)
}

func (g *Game3D) initThreeJS() {
    // JavaScript에서 Three.js Scene 생성
    g.renderer.scene = js.Global().Get("THREE").Get("Scene").New()
    g.renderer.camera = js.Global().Get("THREE").Get("PerspectiveCamera").New(75, 1, 0.1, 1000)
    g.renderer.renderer = js.Global().Get("THREE").Get("WebGLRenderer").New()
}
```

#### HTML 템플릿
```html
<!DOCTYPE html>
<html>
<head>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/three.js/r128/three.min.js"></script>
</head>
<body>
    <div id="game-container">
        <canvas id="three-canvas"></canvas>
        <div id="go-app-hud"></div>
    </div>
</body>
</html>
```

### 2. WebGL 직접 제어 방식

#### 구조
```go
// webgl_renderer.go
type WebGLRenderer struct {
    gl       js.Value
    program  js.Value
    buffers  map[string]js.Value
    textures map[string]js.Value
}

func (r *WebGLRenderer) InitGL(canvasID string) error {
    canvas := js.Global().Get("document").Call("getElementById", canvasID)
    r.gl = canvas.Call("getContext", "webgl2")
    if r.gl.IsNull() {
        return errors.New("WebGL2 not supported")
    }
    return nil
}

func (r *WebGLRenderer) CreateShaderProgram(vertexSource, fragmentSource string) {
    // 셰이더 컴파일 및 프로그램 생성
}

func (r *WebGLRenderer) RenderFrame(gameState *GameState) {
    // 게임 상태를 기반으로 3D 렌더링
}
```

### 3. 게임 상태 관리

```go
// game_state.go
type GameState struct {
    Players   []*Player3D
    Objects   []*GameObject3D
    Camera    *Camera3D
    Lighting  *LightingState
    Physics   *PhysicsWorld
}

type Player3D struct {
    ID       string
    Position Vector3
    Rotation Vector3
    Velocity Vector3
    Health   float64
}

type GameObject3D struct {
    ID       string
    Position Vector3
    Rotation Vector3
    Scale    Vector3
    Model    string
    Texture  string
}

type Vector3 struct {
    X, Y, Z float64
}
```

### 4. 물리 엔진 (순수 Go 구현)

```go
// physics.go
type PhysicsWorld struct {
    Bodies    []*RigidBody
    Gravity   Vector3
    TimeStep  float64
}

type RigidBody struct {
    Position Vector3
    Velocity Vector3
    Mass     float64
    Shape    CollisionShape
}

func (p *PhysicsWorld) Update(deltaTime float64) {
    // 물리 시뮬레이션 업데이트
    for _, body := range p.Bodies {
        // 중력 적용
        body.Velocity.Y -= p.Gravity.Y * deltaTime
        
        // 위치 업데이트
        body.Position.X += body.Velocity.X * deltaTime
        body.Position.Y += body.Velocity.Y * deltaTime
        body.Position.Z += body.Velocity.Z * deltaTime
    }
    
    // 충돌 감지 및 해결
    p.detectCollisions()
}
```

## 구현 예제: 3D 큐브 게임

### main.go
```go
package main

import (
    "syscall/js"
    "time"
    
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Cube3DGame struct {
    app.Compo
    gameState *GameState
    renderer  *ThreeJSRenderer
    running   bool
}

func (c *Cube3DGame) OnMount(ctx app.Context) {
    c.gameState = NewGameState()
    c.renderer = NewThreeJSRenderer()
    c.renderer.Init("game-canvas")
    c.running = true
    
    go c.gameLoop(ctx)
}

func (c *Cube3DGame) gameLoop(ctx app.Context) {
    ticker := time.NewTicker(16 * time.Millisecond) // 60 FPS
    defer ticker.Stop()
    
    for c.running {
        select {
        case <-ticker.C:
            c.update()
            c.renderer.Render(c.gameState)
            ctx.Dispatch(func(ctx app.Context) {
                ctx.Update()
            })
        }
    }
}

func (c *Cube3DGame) update() {
    // 게임 로직 업데이트
    deltaTime := 0.016 // 60 FPS
    
    // 큐브 회전
    for _, obj := range c.gameState.Objects {
        obj.Rotation.Y += deltaTime
    }
}

func (c *Cube3DGame) Render() app.UI {
    return app.Div().Class("game-container").Body(
        // Three.js 캔버스는 JavaScript에서 관리
        app.Canvas().ID("game-canvas"),
        
        // Go-App으로 UI/HUD 관리
        app.Div().Class("game-hud").Body(
            app.H2().Text("3D Cube Game"),
            app.P().Text("FPS: 60"),
        ),
    )
}
```

### JavaScript 브릿지 (app.js)
```javascript
// Three.js 초기화 및 Go와의 브릿지
class ThreeJSBridge {
    constructor() {
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        this.renderer = new THREE.WebGLRenderer({ canvas: document.getElementById('game-canvas') });
        
        this.setupScene();
    }
    
    setupScene() {
        // 조명 설정
        const light = new THREE.DirectionalLight(0xffffff, 1);
        light.position.set(5, 5, 5);
        this.scene.add(light);
        
        this.camera.position.z = 5;
    }
    
    updateGameObjects(gameState) {
        // Go에서 전달받은 게임 상태로 3D 객체 업데이트
        this.scene.children = this.scene.children.filter(child => child.type === 'DirectionalLight');
        
        gameState.Objects.forEach(obj => {
            const geometry = new THREE.BoxGeometry();
            const material = new THREE.MeshPhongMaterial({ color: 0x00ff00 });
            const cube = new THREE.Mesh(geometry, material);
            
            cube.position.set(obj.Position.X, obj.Position.Y, obj.Position.Z);
            cube.rotation.set(obj.Rotation.X, obj.Rotation.Y, obj.Rotation.Z);
            
            this.scene.add(cube);
        });
    }
    
    render() {
        this.renderer.render(this.scene, this.camera);
    }
}

// Go-App과의 인터페이스
window.threeBridge = new ThreeJSBridge();
```

## 성능 최적화 전략

### 1. 렌더링 최적화
- Object Pooling으로 GC 압박 감소
- LOD (Level of Detail) 시스템
- Frustum Culling
- Batch Rendering

### 2. 메모리 관리
```go
type ObjectPool struct {
    objects chan *GameObject3D
}

func (p *ObjectPool) Get() *GameObject3D {
    select {
    case obj := <-p.objects:
        return obj
    default:
        return &GameObject3D{}
    }
}

func (p *ObjectPool) Put(obj *GameObject3D) {
    obj.Reset() // 객체 초기화
    select {
    case p.objects <- obj:
    default:
        // 풀이 가득참
    }
}
```

### 3. 비동기 리소스 로딩
```go
type AssetManager struct {
    models   map[string]*Model3D
    textures map[string]*Texture
    loading  sync.Map
}

func (a *AssetManager) LoadModelAsync(path string) <-chan *Model3D {
    ch := make(chan *Model3D, 1)
    
    go func() {
        if model, exists := a.models[path]; exists {
            ch <- model
            return
        }
        
        // 비동기 로딩
        model := a.loadModelFromPath(path)
        a.models[path] = model
        ch <- model
    }()
    
    return ch
}
```

## 권장 라이브러리

### Three.js 확장
- **Three.js**: 기본 3D 렌더링
- **Cannon.js**: 물리 엔진 (JavaScript 측)
- **GLTFLoader**: 3D 모델 로딩
- **OrbitControls**: 카메라 컨트롤

### Go 라이브러리
- **go-app**: WASM 프레임워크
- **mathgl**: 3D 수학 라이브러리
- **gorilla/websocket**: 멀티플레이어 지원

## 예제 게임 타입

### 1. 3D 플랫포머
- 플레이어 캐릭터 제어
- 중력과 점프 메커니즘
- 3D 환경 탐험

### 2. 3D 슈팅 게임
- FPS/TPS 카메라
- 프로젝타일 물리
- 적 AI

### 3. 3D 퍼즐 게임
- 객체 조작
- 공간 추론
- 물리 기반 퍼즐

## 결론

Go-App에서 3D 게임을 구현하는 가장 실용적인 방법은 **Three.js 브릿지 방식**입니다:

### 장점
✅ 검증된 3D 렌더링 파이프라인
✅ Go의 강력한 게임 로직 처리
✅ 풍부한 3D 자산 생태계
✅ 웹 표준 준수

### 단점
❌ JavaScript 의존성
❌ 브릿지 오버헤드
❌ 복잡한 디버깅

하지만 이 방식이 현재로서는 가장 균형 잡힌 접근법입니다.