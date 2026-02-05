# Athema: Decision Matrix & Quick Reference

## The 4 Critical Decisions

### Decision 1: Rendering Technology

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        RENDERING TECHNOLOGY                                 │
├─────────────────────────┬─────────────────────┬─────────────────────────────┤
│      LIVE2D (2D)        │   SPINE (2D)        │       VRM (3D)              │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ✅ Authentic anime feel │ ✅ Smooth animation │ ✅ Maximum flexibility      │
│ ✅ Industry standard    │ ✅ Lower cost       │ ✅ Physics-based            │
│ ✅ Established tooling  │ ✅ Web-optimized    │ ✅ Multi-angle shots        │
│ ✅ VTuber ecosystem     │ ✅ Game-ready       │ ✅ AR/VR ready              │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ❌ License fees         │ ❌ Less anime feel  │ ❌ Steeper learning curve   │
│ ❌ 2D only              │ ❌ Smaller anime    │ ❌ More complex art         │
│ ❌ Performance cost     │    community        │    pipeline                 │
│                         │                     │ ❌ Can look "uncanny"       │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ BEST FOR:               │ BEST FOR:           │ BEST FOR:                   │
│ Authentic VTuber-style  │ Games, web apps     │ Future-proof, premium       │
│ companion apps          │ on budget           │ experience                  │
└─────────────────────────┴─────────────────────┴─────────────────────────────┘
```

**RECOMMENDATION**: Live2D for authentic anime companion. VRM if you might expand to VR/AR later.

---

### Decision 2: Platform Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        PLATFORM ARCHITECTURE                                │
├─────────────────────────┬─────────────────────┬─────────────────────────────┤
│   WEB (Electron/Tauri)  │  DESKTOP (Godot)    │    HYBRID                   │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ✅ Fastest dev cycle    │ ✅ Best performance │ ✅ Flexible                 │
│ ✅ Easy AI integration  │ ✅ True offline     │ ✅ Can pivot later          │
│ ✅ One codebase         │ ✅ Better audio     │ ✅ Leverage both            │
│ ✅ Easy updates         │ ✅ System tray      │    ecosystems               │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ❌ Memory (Electron)    │ ❌ Longer dev time  │ ❌ More complex             │
│ ❌ Performance ceiling  │ ❌ Larger download  │ ❌ Integration overhead     │
│ ❌ Limited offline      │ ❌ Build per OS     │                             │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ BEST FOR:               │ BEST FOR:           │ BEST FOR:                   │
│ MVPs, rapid iteration   │ Premium experience  │ Uncertain requirements      │
└─────────────────────────┴─────────────────────┴─────────────────────────────┘
```

**RECOMMENDATION**: Tauri (not Electron) for performance + web ecosystem. Godot if premium feel is priority.

---

### Decision 3: AI Backend

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           AI BACKEND                                        │
├─────────────────────────┬─────────────────────┬─────────────────────────────┤
│   CLOUD API             │   LOCAL LLM         │    HYBRID                   │
│   (Claude/OpenAI)       │   (llama.cpp)       │    (Smart fallback)         │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ✅ Best quality         │ ✅ Privacy          │ ✅ Balanced                 │
│ ✅ Easy setup           │ ✅ No API costs     │ ✅ Graceful degradation     │
│ ✅ Fast responses       │ ✅ Offline capable  │ ✅ User choice              │
│ ✅ Always improving     │ ✅ No rate limits   │                             │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ❌ Ongoing costs        │ ❌ Hardware reqs    │ ❌ More complex             │
│ ❌ Requires internet    │ ❌ Lower quality    │ ❌ More testing             │
│ ❌ Privacy concerns     │ ❌ Slower on CPU    │                             │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ BEST FOR:               │ BEST FOR:           │ BEST FOR:                   │
│ Start fast, iterate     │ Privacy-focused     │ Production apps             │
└─────────────────────────┴─────────────────────┴─────────────────────────────┘
```

**RECOMMENDATION**: Start with Claude API, abstract the interface so local LLM can be added later.

---

### Decision 4: Interaction Mode

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        INTERACTION MODE                                     │
├─────────────────────────┬─────────────────────┬─────────────────────────────┤
│   TEXT-ONLY             │   TURN-BASED VOICE  │    REAL-TIME CONVERSATION   │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ✅ Simplest to build    │ ✅ Natural feel     │ ✅ Most engaging            │
│ ✅ No audio complexity  │ ✅ Lip sync works   │ ✅ Interruptible            │
│ ✅ Lower costs          │ ✅ Predictable flow │ ✅ Human-like               │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ ❌ Less engaging        │ ❌ Wait for TTS     │ ❌ Complex audio pipeline   │
│                         │ ❌ Not interruptible│ ❌ Harder to debug          │
├─────────────────────────┼─────────────────────┼─────────────────────────────┤
│ BEST FOR:               │ BEST FOR:           │ BEST FOR:                   │
│ MVP, testing AI         │ Balanced approach   │ Premium experience          │
└─────────────────────────┴─────────────────────┴─────────────────────────────┘
```

**RECOMMENDATION**: Start text-only, add turn-based voice in Phase 3, consider real-time for v2.0.

---

## Quick Start Combinations

### Combo A: "Fast to Market" ⭐ RECOMMENDED
```
Platform:     Tauri (Rust + WebView)
Animation:    Live2D Cubism
AI:           Claude API
Interaction:  Text → Turn-based Voice
Timeline:     6-8 weeks to MVP
```
**Best for**: Validating the concept, building audience

---

### Combo B: "Premium Experience"
```
Platform:     Godot 4
Animation:    VRM (3D)
AI:           Claude API + Local LLM option
Interaction:  Turn-based Voice
Timeline:     12-16 weeks to MVP
```
**Best for**: Paid product, offline users, VR expansion plans

---

### Combo C: "Web-First"
```
Platform:     Next.js + WebGL
Animation:    Spine or WebGL shaders
AI:           Claude API
Interaction:  Text-only initially
Timeline:     4-6 weeks to MVP
```
**Best for**: Browser access, easy sharing, lower barrier

---

## My Recommendation for Athema

Given your requirements:
- ✅ Sophisticated anime style → **Live2D**
- ✅ AI conversation → **Claude API** with abstraction
- ✅ Graphically engaging → **Tauri + GPU acceleration**
- ✅ Clean architecture → **Layered architecture with DI**

### Suggested Stack:
```
┌─────────────────────────────────────┐
│  Tauri (Rust backend)              │
│  + React/TypeScript frontend       │
├─────────────────────────────────────┤
│  Live2D Cubism SDK                 │
│  Zustand for state                 │
├─────────────────────────────────────┤
│  Claude API (conversation)         │
│  ElevenLabs/Azure TTS              │
│  Whisper.cpp (local STT)           │
├─────────────────────────────────────┤
│  SQLite (local storage)            │
└─────────────────────────────────────┘
```

### Why This Stack?
1. **Tauri** = Better performance than Electron, still web tech
2. **Live2D** = Industry standard for anime VTubers
3. **Claude** = Best conversation quality for companions
4. **Layered architecture** = Can swap any component later

### First 2 Weeks Priority:
1. Set up Tauri + React project structure
2. Integrate Live2D SDK, display a basic model
3. Connect to Claude API, echo responses
4. Build the state management core

This gives you a working foundation to iterate on.
