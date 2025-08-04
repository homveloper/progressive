// =============================================================================
// Three.js Bridge for Qube 3D Game
// =============================================================================

class ThreeJSBridge {
    constructor() {
        this.scene = null;
        this.camera = null;
        this.renderer = null;
        this.canvas = null;
        this.gameObjects = new Map();
        this.animationId = null;
        this.initialized = false;
        
        // Game state
        this.currentGameState = null;
        
        // Performance tracking
        this.lastFrameTime = 0;
        this.frameCount = 0;
        this.fps = 60;
        
        // Debug timing
        this.lastDebugTime = 0;
        this.debugInterval = 1000; // 1 second
    }

    // Initialize Three.js scene
    init(canvasId) {
        try {
            this.canvas = document.getElementById(canvasId);
            if (!this.canvas) {
                console.error('Canvas not found:', canvasId);
                return false;
            }

            // Create scene
            this.scene = new THREE.Scene();
            this.scene.background = new THREE.Color(0x87CEEB); // Sky blue

            // Create camera
            this.camera = new THREE.PerspectiveCamera(
                75, // FOV
                this.canvas.clientWidth / this.canvas.clientHeight, // Aspect ratio
                0.1, // Near plane
                1000 // Far plane
            );

            // Create renderer
            this.renderer = new THREE.WebGLRenderer({ 
                canvas: this.canvas,
                antialias: true,
                alpha: true
            });
            this.renderer.setSize(this.canvas.clientWidth, this.canvas.clientHeight);
            this.renderer.shadowMap.enabled = true;
            this.renderer.shadowMap.type = THREE.PCFSoftShadowMap;

            // Setup lighting
            this.setupLighting();

            // Setup initial scene
            this.setupScene();

            // Handle window resize
            window.addEventListener('resize', () => this.handleResize());
            
            // Handle pointer lock
            this.setupPointerLock();

            this.initialized = true;
            console.log('Three.js initialized successfully');

            // Start render loop
            this.startRenderLoop();

            return true;
        } catch (error) {
            console.error('Failed to initialize Three.js:', error);
            return false;
        }
    }

    setupLighting() {
        // Ambient light
        const ambientLight = new THREE.AmbientLight(0x404040, 0.4);
        this.scene.add(ambientLight);

        // Directional light (sun)
        const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight.position.set(10, 10, 5);
        directionalLight.castShadow = true;
        
        // Shadow camera settings
        directionalLight.shadow.camera.near = 0.1;
        directionalLight.shadow.camera.far = 50;
        directionalLight.shadow.camera.left = -20;
        directionalLight.shadow.camera.right = 20;
        directionalLight.shadow.camera.top = 20;
        directionalLight.shadow.camera.bottom = -20;
        directionalLight.shadow.mapSize.width = 2048;
        directionalLight.shadow.mapSize.height = 2048;
        
        this.scene.add(directionalLight);
    }

    setupScene() {
        // Create ground grid helper
        const gridHelper = new THREE.GridHelper(20, 20, 0x888888, 0xcccccc);
        gridHelper.position.y = -0.5;
        this.scene.add(gridHelper);

        // Create sky dome
        const skyGeometry = new THREE.SphereGeometry(500, 32, 32);
        const skyMaterial = new THREE.MeshBasicMaterial({ 
            color: 0x87CEEB,
            side: THREE.BackSide
        });
        const sky = new THREE.Mesh(skyGeometry, skyMaterial);
        this.scene.add(sky);
    }

    setupPointerLock() {
        this.canvas.addEventListener('click', () => {
            this.canvas.requestPointerLock();
        });

        document.addEventListener('pointerlockchange', () => {
            if (document.pointerLockElement === this.canvas) {
                console.log('Pointer locked');
            } else {
                console.log('Pointer unlocked');
            }
        });
    }

    // Update game state from Go
    updateGameState(gameState) {
        if (!this.initialized || !gameState) {
            return;
        }

        this.currentGameState = gameState;

        // Update camera
        if (gameState.camera) {
            this.updateCamera(gameState.camera);
        }

        // Update lighting
        if (gameState.lighting) {
            this.updateLighting(gameState.lighting);
        }

        // Update players
        if (gameState.players) {
            this.updatePlayers(gameState.players);
        }

        // Update objects
        if (gameState.objects) {
            this.updateObjects(gameState.objects);
        }
    }

    updateCamera(cameraData) {
        if (!this.camera) return;

        const now = Date.now();
        const shouldDebug = now - this.lastDebugTime > this.debugInterval;

        if (shouldDebug) {
            console.log('Updating camera:', cameraData);
        }

        // Safely update camera position
        if (cameraData.position) {
            const newPos = {
                x: cameraData.position.x || 0,
                y: cameraData.position.y || 0,
                z: cameraData.position.z || 0
            };
            if (shouldDebug) {
                console.log('Setting camera position to:', newPos);
            }
            this.camera.position.set(newPos.x, newPos.y, newPos.z);
        }

        // Safely update camera target using Three.js Vector3
        if (cameraData.target) {
            const targetVector = new THREE.Vector3(
                cameraData.target.x || 0,
                cameraData.target.y || 0,
                cameraData.target.z || 0
            );
            if (shouldDebug) {
                console.log('Setting camera target to:', targetVector);
            }
            this.camera.lookAt(targetVector);
        }

        this.camera.fov = cameraData.fov || 75;
        this.camera.updateProjectionMatrix();
    }

    updateLighting(lightingData) {
        // Find and update lights
        this.scene.traverse((child) => {
            if (child instanceof THREE.AmbientLight) {
                if (lightingData.ambientColor) {
                    child.color.setHex(parseInt(lightingData.ambientColor.replace('#', '0x')));
                }
                if (typeof lightingData.ambientIntensity === 'number') {
                    child.intensity = lightingData.ambientIntensity;
                }
            }
            
            if (child instanceof THREE.DirectionalLight) {
                const dir = lightingData.directionalLight;
                if (dir && dir.direction) {
                    child.position.set(
                        (dir.direction.x || 0) * 10, 
                        (dir.direction.y || 0) * 10, 
                        (dir.direction.z || 0) * 10
                    );
                }
                if (dir && dir.color) {
                    child.color.setHex(parseInt(dir.color.replace('#', '0x')));
                }
                if (dir && typeof dir.intensity === 'number') {
                    child.intensity = dir.intensity;
                }
            }
        });
    }

    updatePlayers(players) {
        const now = Date.now();
        const shouldDebug = now - this.lastDebugTime > this.debugInterval;
        
        if (shouldDebug) {
            console.log('Updating players:', players);
        }
        
        players.forEach(player => {
            if (!player.isActive) {
                if (shouldDebug) console.log('Player not active:', player.id);
                return;
            }

            let playerMesh = this.gameObjects.get(player.id);
            
            if (!playerMesh) {
                console.log('Creating new player mesh for:', player.id); // Always log creation
                // Create new player mesh using CylinderGeometry (compatible with r128)
                const geometry = new THREE.CylinderGeometry(0.5, 0.5, 1.5, 8);
                const material = new THREE.MeshLambertMaterial({ color: 0x00ff00 });
                playerMesh = new THREE.Mesh(geometry, material);
                playerMesh.castShadow = true;
                playerMesh.receiveShadow = true;
                
                this.scene.add(playerMesh);
                this.gameObjects.set(player.id, playerMesh);
                console.log('Player mesh created and added to scene');
            }

            // Update position and rotation safely
            if (player.position) {
                const newPos = {
                    x: player.position.x || 0,
                    y: (player.position.y || 0) + 1, // Offset for cylinder center
                    z: player.position.z || 0
                };
                if (shouldDebug) {
                    console.log('Updating player position to:', newPos);
                }
                playerMesh.position.set(newPos.x, newPos.y, newPos.z);
            }

            if (player.rotation) {
                playerMesh.rotation.set(
                    player.rotation.x || 0,
                    player.rotation.y || 0,
                    player.rotation.z || 0
                );
            }
        });
    }

    updateObjects(objects) {
        objects.forEach(obj => {
            if (!obj.isActive) {
                // Remove inactive objects
                const mesh = this.gameObjects.get(obj.id);
                if (mesh) {
                    this.scene.remove(mesh);
                    this.gameObjects.delete(obj.id);
                }
                return;
            }

            let mesh = this.gameObjects.get(obj.id);
            
            if (!mesh) {
                // Create new object mesh
                mesh = this.createObjectMesh(obj);
                if (mesh) {
                    this.scene.add(mesh);
                    this.gameObjects.set(obj.id, mesh);
                }
            }

            if (mesh) {
                // Update position safely
                if (obj.position) {
                    mesh.position.set(
                        obj.position.x || 0,
                        obj.position.y || 0,
                        obj.position.z || 0
                    );
                }

                // Update rotation safely
                if (obj.rotation) {
                    mesh.rotation.set(
                        obj.rotation.x || 0,
                        obj.rotation.y || 0,
                        obj.rotation.z || 0
                    );
                }

                // Update scale safely
                if (obj.scale) {
                    mesh.scale.set(
                        obj.scale.x || 1,
                        obj.scale.y || 1,
                        obj.scale.z || 1
                    );
                }
            }
        });
    }

    createObjectMesh(obj) {
        let geometry;
        let material;

        // Create geometry based on model type
        switch (obj.model) {
            case 'box':
                geometry = new THREE.BoxGeometry(1, 1, 1);
                break;
            case 'sphere':
                geometry = new THREE.SphereGeometry(0.5, 16, 16);
                break;
            case 'cylinder':
                geometry = new THREE.CylinderGeometry(0.5, 0.5, 1, 16);
                break;
            default:
                geometry = new THREE.BoxGeometry(1, 1, 1);
        }

        // Create material with safe color parsing
        let color = 0xcccccc; // Default gray color
        if (obj.color && typeof obj.color === 'string') {
            try {
                color = parseInt(obj.color.replace('#', '0x'));
            } catch (e) {
                console.warn('Invalid color format:', obj.color);
            }
        }
        
        if (obj.id === 'ground') {
            material = new THREE.MeshLambertMaterial({ 
                color: color,
                transparent: true,
                opacity: 0.8
            });
        } else {
            material = new THREE.MeshLambertMaterial({ color: color });
        }

        const mesh = new THREE.Mesh(geometry, material);
        
        // Enable shadows
        mesh.castShadow = true;
        mesh.receiveShadow = true;

        // Add special effects for collectibles
        if (obj.id.includes('cube_') && obj.id !== 'ground') {
            // Add a subtle glow effect
            const glowMaterial = new THREE.MeshBasicMaterial({
                color: color,
                transparent: true,
                opacity: 0.3
            });
            const glowMesh = new THREE.Mesh(
                new THREE.BoxGeometry(1.2, 1.2, 1.2),
                glowMaterial
            );
            mesh.add(glowMesh);
        }

        return mesh;
    }

    // Render loop
    startRenderLoop() {
        const animate = (currentTime) => {
            if (!this.initialized) return;

            this.animationId = requestAnimationFrame(animate);

            // Calculate FPS
            if (currentTime - this.lastFrameTime >= 1000) {
                this.fps = this.frameCount;
                this.frameCount = 0;
                this.lastFrameTime = currentTime;
            }
            this.frameCount++;

            // Render scene
            this.render();
        };

        animate(0);
    }

    render() {
        if (!this.renderer || !this.scene || !this.camera) return;

        // Add some dynamic effects
        if (this.currentGameState && this.currentGameState.time) {
            // Rotate collectible cubes (visual enhancement)
            this.gameObjects.forEach((mesh, id) => {
                if (id.includes('cube_') && id !== 'ground') {
                    // Additional rotation for visual appeal
                    mesh.rotation.x += 0.01;
                    mesh.rotation.z += 0.005;
                    
                    // Floating animation
                    mesh.position.y += Math.sin(this.currentGameState.time * 2 + mesh.position.x) * 0.001;
                }
            });
        }

        this.renderer.render(this.scene, this.camera);
    }

    handleResize() {
        if (!this.camera || !this.renderer || !this.canvas) return;

        const width = this.canvas.clientWidth;
        const height = this.canvas.clientHeight;

        this.camera.aspect = width / height;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(width, height);
    }

    // Cleanup
    dispose() {
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
        }

        if (this.renderer) {
            this.renderer.dispose();
        }

        this.gameObjects.clear();
        this.initialized = false;
    }

    // Get performance info
    getPerformanceInfo() {
        return {
            fps: this.fps,
            objects: this.gameObjects.size,
            triangles: this.renderer?.info?.render?.triangles || 0,
            calls: this.renderer?.info?.render?.calls || 0
        };
    }
}

// Global instance
let threeBridge = null;

// Global functions for Go to call
window.initThreeJSBridge = function(canvasId) {
    console.log('Initializing Three.js Bridge for:', canvasId);
    
    try {
        // Wait for Three.js to load
        if (typeof THREE === 'undefined') {
            console.error('Three.js not loaded!');
            return false;
        }

        // Dispose existing bridge if any
        if (threeBridge) {
            threeBridge.dispose();
        }

        threeBridge = new ThreeJSBridge();
        const result = threeBridge.init(canvasId);
        console.log('Three.js bridge initialization result:', result);
        return result;
    } catch (error) {
        console.error('Error initializing Three.js bridge:', error);
        return false;
    }
};

window.updateGameState = function(gameStateJson) {
    try {
        if (threeBridge && threeBridge.initialized) {
            // Parse JSON string from Go
            let gameState;
            if (typeof gameStateJson === 'string') {
                gameState = JSON.parse(gameStateJson);
            } else {
                gameState = gameStateJson; // Fallback for direct object
            }
            
            // Debug: Log game state to console only once per second
            const now = Date.now();
            if (now - threeBridge.lastDebugTime > threeBridge.debugInterval) {
                console.log('Updating game state:', gameState);
                threeBridge.lastDebugTime = now;
            }
            
            threeBridge.updateGameState(gameState);
        } else {
            console.warn('ThreeBridge not initialized or not available');
        }
    } catch (error) {
        console.error('Error updating game state:', error);
    }
};

window.getThreeJSPerformance = function() {
    if (threeBridge && threeBridge.initialized) {
        return threeBridge.getPerformanceInfo();
    }
    return { fps: 0, objects: 0, triangles: 0, calls: 0 };
};

window.disposeThreeJS = function() {
    try {
        if (threeBridge) {
            console.log('Disposing Three.js bridge...');
            threeBridge.dispose();
            threeBridge = null;
        }
    } catch (error) {
        console.error('Error disposing Three.js bridge:', error);
    }
};

// Debug function
window.debugThreeJS = function() {
    if (threeBridge) {
        console.log('Three.js Bridge Debug Info:', {
            initialized: threeBridge.initialized,
            gameObjects: threeBridge.gameObjects.size,
            performance: threeBridge.getPerformanceInfo(),
            scene: threeBridge.scene,
            camera: threeBridge.camera,
            renderer: threeBridge.renderer
        });
    }
};

console.log('Three.js Bridge loaded successfully');