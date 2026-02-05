# Athema: Animated AI Companion - Architecture Proposal

## Executive Summary

Athema is an animated AI companion featuring a sophisticated anime-style character with natural conversation capabilities. This architecture proposal outlines a modular, extensible system built on SOLID principles.

---

## 1. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PRESENTATION LAYER                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Render    │  │   Input     │  │   Audio     │  │   UI/Overlay        │ │
│  │   Engine    │  │   Handler   │  │   System    │  │   System            │ │
│  │  (2D/3D)    │  │             │  │             │  │                     │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
│         │                │                │                    │            │
│         └────────────────┴────────────────┴────────────────────┘            │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     CHARACTER CONTROLLER                            │   │
│  │         (Animation State Machine + Behavior Orchestrator)           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CORE SERVICES LAYER                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────┐  │
│  │  Conversation   │  │   Emotion       │  │      Memory Manager         │  │
│  │     Engine      │  │   Engine        │  │  (Short & Long-term)        │  │
│  │                 │  │                 │  │                             │  │
│  │ • LLM Interface │  │ • Mood States   │  │ • Conversation History      │  │
│  │ • Intent Parser │  │ • Expression    │  │ • User Preferences          │  │
│  │ • Response Gen  │  │   Mapper        │  │ • Relationship Model        │  │
│  └────────┬────────┘  └────────┬────────┘  └─────────────┬───────────────┘  │
│           │                    │                         │                  │
│           └────────────────────┴─────────────────────────┘                  │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      STATE MANAGEMENT CORE                          │   │
│  │              (Single Source of Truth - Observable)                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      INFRASTRUCTURE LAYER                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   AI/LLM    │  │   Asset     │  │ Persistence │  │   Platform Abstr.   │ │
│  │   Adapter   │  │   Pipeline  │  │   Layer     │  │      Layer          │ │
│  │             │  │             │  │             │  │                     │ │
│  │ • OpenAI    │  │ • L2D/Spine │  │ • Local DB  │  │ • Web/Desktop       │ │
│  │ • Claude    │  │ • VRM       │  │ • Cloud Sync│  │ • Mobile            │ │
│  │ • Local LLM │  │ • Sprite    │  │ • Config    │  │ • Event Bridge      │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Major Components & Responsibilities

### 2.1 Presentation Layer

#### Render Engine
| Aspect | Details |
|--------|---------|
| **Responsibility** | Display character, handle animations, visual effects |
| **Interface** | `IRenderer` - abstraction for 2D/3D implementations |
| **Key Methods** | `render()`, `playAnimation()`, `setExpression()`, `transition()` |
| **SOLID** | OCP: New renderers (WebGL, Unity, Godot) without changing core |

#### Input Handler
| Aspect | Details |
|--------|---------|
| **Responsibility** | Capture and normalize user input (text, voice, clicks) |
| **Interface** | `IInputSource` - text, voice, gesture inputs |
| **Key Methods** | `onInput()`, `registerCallback()`, `setMode()` |
| **SOLID** | ISP: Separate interfaces for different input types |

#### Audio System
| Aspect | Details |
|--------|---------|
| **Responsibility** | Text-to-speech, speech-to-text, ambient audio, lip sync |
| **Interface** | `IAudioEngine`, `ITTSProvider`, `ISTTProvider` |
| **Key Methods** | `speak()`, `listen()`, `lipSyncData()`, `setVoice()` |
| **SOLID** | DIP: TTS depends on abstract provider, not specific service |

### 2.2 Character Controller (The Heart)

```typescript
// Core abstraction - Character Controller coordinates everything
interface ICharacterController {
  // Receives processed input + emotional context
  processInteraction(input: ProcessedInput): Promise<InteractionResult>;
  
  // Updates animation based on emotion + response
  updateState(emotion: EmotionState, response: AIResponse): void;
  
  // Observable state for UI components
  getCurrentState(): CharacterState;
}

// Animation State Machine
interface IAnimationStateMachine {
  transition(to: AnimationState, blending: BlendMode): void;
  getCurrentAnimation(): AnimationClip;
  triggerExpression(expression: ExpressionType, intensity: number): void;
}
```

**Responsibility**: Orchestrates between AI responses, emotions, and visual representation. The ONLY component that knows about animation-emotion mapping.

### 2.3 Core Services Layer

#### Conversation Engine
| Aspect | Details |
|--------|---------|
| **Responsibility** | LLM communication, prompt management, context handling |
| **Interface** | `IConversationProvider`, `IPromptBuilder` |
| **Key Classes** | `ConversationManager`, `ContextBuilder`, `ResponseParser` |
| **SOLID** | SRP: Only handles conversation logic, not rendering or emotions |

#### Emotion Engine
| Aspect | Details |
|--------|---------|
| **Responsibility** | Analyze sentiment, maintain mood state, map to expressions |
| **Interface** | `IEmotionAnalyzer`, `IExpressionMapper` |
| **Key Classes** | `EmotionAnalyzer`, `MoodStateMachine`, `ExpressionMapper` |
| **SOLID** | OCP: New emotion models without changing animation code |

#### Memory Manager
| Aspect | Details |
|--------|---------|
| **Responsibility** | Store/retrieve conversation history, user preferences, relationship data |
| **Interface** | `IMemoryStore`, `IUserProfileManager` |
| **Key Classes** | `ConversationMemory`, `UserProfile`, `RelationshipTracker` |
| **SOLID** | DIP: Core logic depends on `IMemoryStore`, not specific DB |

### 2.4 State Management Core

**Pattern**: Observable Store (Redux-like or Signals-based)

```typescript
// Single source of truth
interface AppState {
  character: CharacterState;
  conversation: ConversationState;
  emotion: EmotionState;
  user: UserState;
  settings: SettingsState;
}

// Character State
interface CharacterState {
  currentAnimation: AnimationId;
  currentExpression: ExpressionId;
  position: { x: number; y: number };
  isSpeaking: boolean;
  isListening: boolean;
  mood: MoodVector; // Multi-dimensional emotion
}
```

**Responsibility**: Centralized, immutable state. All components react to state changes.

---

## 3. Tech Stack Options

### Option A: Web-Based (Electron/Tauri) + 2D Live2D
**Best for**: Cross-platform, easier distribution, web-native AI integration

```
┌─────────────────────────────────────────┐
│  Frontend: React/Vue/Svelte            │
│  Animation: Live2D Cubism SDK          │
│  Rendering: WebGL/Canvas               │
│  State: Zustand/Jotai/Redux            │
├─────────────────────────────────────────┤
│  Backend (main process):               │
│  Runtime: Node.js / Rust (Tauri)       │
│  AI: OpenAI/Claude API                 │
│  TTS: ElevenLabs / Azure / Local       │
│  DB: SQLite / IndexedDB                │
└─────────────────────────────────────────┘
```

| Pros | Cons |
|------|------|
| ✅ Single codebase, all platforms | ❌ Performance ceiling for complex animations |
| ✅ Easy AI API integration | ❌ Live2D license costs for commercial use |
| ✅ Rich web ecosystem | ❌ Limited offline capability |
| ✅ Easy updates | ❌ Memory usage (Electron) |
| ✅ Excellent for 2D anime style | ❌ Web audio latency |

**Recommended if**: You want fastest time-to-market, don't need heavy 3D, plan to be online-first.

---

### Option B: Desktop Native (Unity/Unreal/Godot) + 3D/2.5D
**Best for**: Premium visual quality, complex animations, offline-first

```
┌─────────────────────────────────────────┐
│  Engine: Unity / Godot / Unreal        │
│  Character: VRM (3D) or Spine (2D)     │
│  UI: Unity UI / Godot UI / UMG         │
│  Scripting: C# / GDScript / C++        │
├─────────────────────────────────────────┤
│  Backend Integration:                  │
│  Local LLM: llama.cpp / Ollama         │
│  Cloud AI: REST API calls              │
│  TTS: Runtime voice synthesis          │
│  DB: Local SQLite / JSON               │
└─────────────────────────────────────────┘
```

| Pros | Cons |
|------|------|
| ✅ Superior animation quality | ❌ Steeper learning curve |
| ✅ True offline capability | ❌ Larger download size |
| ✅ VRM ecosystem (anime-ready) | ❌ Platform-specific builds |
| ✅ Better audio/visual sync | ❌ More complex deployment |
| ✅ Physics, lighting, effects | ❌ Harder to iterate quickly |

**Recommended if**: You want premium experience, offline capability, complex interactions, or plan to go 3D eventually.

---

### Option C: Hybrid (Godot/WebView) - "Best of Both"
**Best for**: Flexibility, future-proofing

```
┌─────────────────────────────────────────┐
│  Rendering: Godot (2D/3D capable)      │
│  Character: VRM or custom 2D           │
│  Logic Layer: Godot C# / GDScript      │
├─────────────────────────────────────────┤
│  Embedded WebView (optional):          │
│  AI Interface via JS bridge            │
│  Web-based UI overlays                 │
├─────────────────────────────────────────┤
│  Backend: Same as Option B             │
└─────────────────────────────────────────┘
```

| Pros | Cons |
|------|------|
| ✅ Can leverage web AI libraries | ❌ More complex architecture |
| ✅ Native performance where needed | ❌ WebView integration overhead |
| ✅ Gradual migration path | ❌ Potential sync issues |
| ✅ Flexible deployment | ❌ More maintenance surface |

**Recommended if**: You're uncertain about 2D vs 3D, want maximum flexibility, or plan to iterate on the tech stack.

---

## 4. KEY DECISIONS (Must Answer Before Implementation)

### Decision 1: Online vs Offline-First
| Online-First | Offline-First |
|--------------|---------------|
| Cloud LLM APIs (OpenAI, Claude) | Local LLM (llama.cpp, Ollama) |
| Better conversation quality | Privacy-focused |
| Requires subscription/auth | Higher hardware requirements |
| Easier to implement | More complex setup |

**Recommendation**: Start online with OpenAI/Claude, architect for offline capability later.

### Decision 2: 2D Live2D vs 3D VRM
| Live2D (2D) | VRM (3D) |
|-------------|----------|
| Authentic anime aesthetic | More expressive range |
| Lower performance cost | Physics-based movement |
| Easier asset creation | Better camera angles |
| Established ecosystem | Can do 2.5D anime rendering |

**Recommendation**: Live2D for authentic anime feel, VRM if you want flexibility to go 3D later.

### Decision 3: Desktop vs Web Deployment
| Desktop App | Web App |
|-------------|---------|
| Better system integration | Instant access, no install |
| Can run local LLM | Easier sharing/distribution |
| File system access | Platform agnostic |
| Systray/background | Lower barrier to entry |

**Recommendation**: Desktop (Electron/Tauri) for companion apps (stays in background), Web for accessibility.

### Decision 4: Real-time vs Turn-based Interaction
| Real-time | Turn-based |
|-----------|------------|
| Interruptible speech | Simpler state management |
| More natural feel | Easier to implement |
| Complex audio handling | Better for text-first |
| Lip sync challenges | Clearer interaction flow |

**Recommendation**: Start turn-based, add real-time as enhancement.

---

## 5. MVP Roadmap (Build in This Order)

### Phase 1: Foundation (Weeks 1-2)
```
Priority: CRITICAL
Components:
├── State Management Core
│   └── Observable store with character state
├── Basic Render Engine
│   └── Simple 2D sprite display (or Live2D basic integration)
└── Input Handler
    └── Text input only
```
**Goal**: Character appears on screen, can update state programmatically.

### Phase 2: Conversation (Weeks 3-4)
```
Priority: HIGH
Components:
├── Conversation Engine
│   └── OpenAI/Claude integration
│   └── Basic prompt template
├── Character Controller
│   └── Map AI response to basic emotions
└── Simple Expression System
    └── 5-10 basic expressions
```
**Goal**: Can chat with character, see basic emotional reactions.

### Phase 3: Animation & Polish (Weeks 5-6)
```
Priority: MEDIUM
Components:
├── Animation State Machine
│   └── Idle, talking, reacting states
├── Audio System
│   └── Basic TTS integration
│   └── Lip sync (basic)
└── Memory Manager
    └── Conversation history
```
**Goal**: Character speaks, remembers conversation, smooth animations.

### Phase 4: Interaction Layer (Weeks 7-8)
```
Priority: MEDIUM
Components:
├── Voice Input (STT)
├── Click/Touch interactions
├── Customization (appearance, voice)
└── Settings persistence
```
**Goal**: Multi-modal input, personalization, settings.

### Phase 5: Advanced Features (Ongoing)
```
Priority: LOW
Components:
├── Advanced emotions (mood that persists)
├── Mini-games/activities
├── Calendar/reminders integration
├── Multi-language support
└── Mobile support
```

---

## 6. SOLID Principles Application

### Single Responsibility Principle (SRP)
```typescript
// ❌ BAD: One class doing everything
class Character {
  render() { /* drawing */ }
  think() { /* AI logic */ }
  save() { /* database */ }
  speak() { /* audio */ }
}

// ✅ GOOD: Each class has one reason to change
class CharacterRenderer implements IRenderer { }
class ConversationEngine implements IConversationProvider { }
class MemoryStore implements IMemoryStore { }
class AudioEngine implements IAudioEngine { }
class CharacterController {
  // Orchestrates, doesn't implement
  constructor(
    private renderer: IRenderer,
    private conversation: IConversationProvider,
    private memory: IMemoryStore,
    private audio: IAudioEngine
  ) {}
}
```

### Open/Closed Principle (OCP)
```typescript
// ✅ Extend behavior without modifying
interface IAnimationProvider {
  loadModel(path: string): Promise<Model>;
  playAnimation(name: string): void;
}

class Live2DProvider implements IAnimationProvider { }
class SpineProvider implements IAnimationProvider { }
class VRMProvider implements IAnimationProvider { }

// CharacterController works with any provider
// Add new animation tech without changing core code
```

### Liskov Substitution Principle (LSP)
```typescript
// ✅ Any TTS provider can be substituted
interface ITTSProvider {
  speak(text: string, options: TTSOptions): Promise<AudioBuffer>;
  getVoices(): Voice[];
}

class ElevenLabsTTS implements ITTSProvider { }
class AzureTTS implements ITTSProvider { }
class LocalTTS implements ITTSProvider { }

// All work interchangeably
const tts: ITTSProvider = config.useLocal ? new LocalTTS() : new ElevenLabsTTS();
```

### Interface Segregation Principle (ISP)
```typescript
// ❌ BAD: Fat interface
interface ICharacter {
  render(): void;
  think(): void;
  save(): void;
  loadAsset(): void;
  processAI(): void;
}

// ✅ GOOD: Split by role
interface IRenderable { render(): void; }
interface IThinkable { think(input: string): Promise<Response>; }
interface IPersistable { save(): void; load(): void; }
interface IAssetLoader { loadAsset(path: string): void; }

class Character implements IRenderable, IThinkable, IPersistable { }
```

### Dependency Inversion Principle (DIP)
```typescript
// ❌ BAD: Depends on concrete implementation
class ConversationEngine {
  private llm = new OpenAIClient(); // Hard dependency
}

// ✅ GOOD: Depends on abstraction
interface ILLMClient {
  complete(prompt: string): Promise<string>;
}

class ConversationEngine {
  constructor(private llm: ILLMClient) {}
}

// Can inject OpenAI, Claude, local LLM, mock for testing
```

---

## 7. Project Structure (Recommended)

```
athema/
├── src/
│   ├── core/                      # Domain logic, pure TypeScript
│   │   ├── models/               # Data models (Character, Emotion, etc.)
│   │   ├── services/             # Business logic interfaces
│   │   └── events/               # Domain events
│   │
│   ├── presentation/             # UI & Rendering
│   │   ├── renderer/            # IRenderer implementations
│   │   │   ├── live2d/
│   │   │   └── webgl/
│   │   ├── components/          # UI components
│   │   └── hooks/               # State connection hooks
│   │
│   ├── application/             # Application services
│   │   ├── character-controller/
│   │   ├── conversation/
│   │   ├── emotion/
│   │   └── memory/
│   │
│   ├── infrastructure/          # External concerns
│   │   ├── ai/                 # LLM adapters
│   │   ├── audio/              # TTS/STT implementations
│   │   ├── persistence/        # Database implementations
│   │   └── platform/           # Platform-specific code
│   │
│   └── shared/                 # Utilities, types
│       ├── types/
│       └── utils/
│
├── assets/                     # Art assets
│   ├── models/                # Live2D/VRM files
│   ├── animations/            # Animation data
│   └── audio/                 # Voice samples, SFX
│
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
└── docs/
    ├── architecture/
    └── api/
```

---

## 8. Risk Mitigation

| Risk | Mitigation |
|------|------------|
| AI API costs | Abstract LLM interface, support local models |
| Animation complexity | Start with sprite-based, migrate to Live2D |
| Performance issues | Profile early, keep rendering isolated |
| Scope creep | Strict MVP boundaries, modular design |
| Asset creation | Start with free/open models, customize later |

---

## 9. Recommended Path Forward

Based on requirements (sophisticated anime, AI companion, clean architecture):

**Primary Recommendation: Option A (Web/Desktop Hybrid)**
- **Platform**: Tauri (Rust backend) or Electron
- **Animation**: Live2D Cubism (industry standard for anime)
- **AI**: Start with Claude API, architect for local LLM
- **State**: Zustand or Jotai (simple, effective)

**Why**: 
- Fastest iteration for MVPin the anime VTuber space
- Clean separation enables tech swaps later
- Balanced performance vs development speed

**Alternative**: If premium feel is paramount, choose **Godot + VRM** (Option B).

---

## 10. Immediate Next Steps

1. **Decide**: 2D Live2D vs 3D VRM
2. **Prototype**: Basic character display with chosen tech
3. **Validate**: AI conversation flow with simple UI
4. **Iterate**: Add animation layer once basics work

---

*Document Version: 1.0*
*Created: Architecture design phase*
*Next Review: After tech stack decision*
