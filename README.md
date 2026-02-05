# Athema 🌟

An animated AI companion with sophisticated anime aesthetics and natural conversation capabilities.

> *"A digital companion that feels alive, understands you, and grows alongside you."*

---

## Project Status

🚧 **Architecture Phase** - Designing the foundation for a modular, extensible AI companion.

---

## Quick Links

| Document | Purpose |
|----------|---------|
| [Architecture Overview](docs/architecture/ATHEMA_ARCHITECTURE.md) | Complete system design & component breakdown |
| [Decision Matrix](docs/architecture/DECISION_MATRIX.md) | Tech stack comparisons & recommendations |
| [Diagrams](docs/architecture/DIAGRAMS.md) | Visual flow charts & state machines |

---

## The Vision

Athema is designed to be:

- **🎨 Visually Engaging**: Sophisticated anime character with smooth, expressive animations
- **🧠 Conversational**: Natural AI dialogue that remembers context and builds rapport
- **🏗️ Architecturally Sound**: Clean separation of concerns, SOLID principles, extensible design
- **🔧 Customizable**: Pluggable AI providers, animation systems, and interaction modes

---

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────────────────────────┐
│  PRESENTATION        Render │ Input │ Audio │ UI                       │
├─────────────────────────────────────────────────────────────────────────┤
│  APPLICATION         Character Controller │ Conversation │ Emotion       │
├─────────────────────────────────────────────────────────────────────────┤
│  DOMAIN              Interfaces: IRenderer │ IAIProvider │ IAudioEngine │
├─────────────────────────────────────────────────────────────────────────┤
│  INFRASTRUCTURE      Live2D │ Claude API │ TTS │ SQLite                │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Recommended Tech Stack

Based on the architecture analysis, the recommended path forward:

| Component | Recommendation | Rationale |
|-----------|---------------|-----------|
| **Platform** | Tauri (Rust + WebView) | Performance + web ecosystem |
| **Frontend** | React + TypeScript | Developer experience, tooling |
| **Animation** | Live2D Cubism | Industry standard for anime VTubers |
| **State** | Zustand | Simple, effective, TypeScript-native |
| **AI** | Claude API | Best conversation quality |
| **TTS** | ElevenLabs/Azure | Natural voice output |
| **Storage** | SQLite | Local, reliable, portable |

---

## Key Decisions Needed

Before implementation begins, please confirm:

1. **Rendering**: Live2D (2D, authentic anime) or VRM (3D, more flexible)?
2. **Platform**: Desktop app (Tauri/Electron) or Web-first?
3. **AI Strategy**: Cloud-only initially, or architect for local LLM from start?
4. **Interaction**: Text-first MVP, or voice from day one?

See [Decision Matrix](docs/architecture/DECISION_MATRIX.md) for detailed trade-offs.

---

## MVP Roadmap

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| **1. Foundation** | 2 weeks | Character on screen, state management |
| **2. Conversation** | 2 weeks | AI integration, basic emotions |
| **3. Animation** | 2 weeks | TTS, lip sync, smooth transitions |
| **4. Interaction** | 2 weeks | Voice input, customization, settings |
| **5. Polish** | Ongoing | Advanced features, optimization |

---

## Project Structure (Proposed)

```
athema/
├── src/
│   ├── core/              # Domain models & interfaces
│   ├── application/       # Business logic services
│   ├── presentation/      # UI components & renderers
│   ├── infrastructure/    # External API implementations
│   └── shared/           # Utilities & types
├── assets/               # Live2D models, audio, images
├── tests/               # Unit, integration, e2e
└── docs/                # Architecture & API docs
```

---

## SOLID Principles

The architecture is designed around:

- **S**ingle Responsibility: Each component has one job
- **O**pen/Closed: Extend behavior without modifying existing code
- **L**iskov Substitution: Swappable AI providers, renderers, audio engines
- **I**nterface Segregation: Lean interfaces, no fat contracts
- **D**ependency Inversion: Core logic depends on abstractions

---

## Next Steps

1. Review the [architecture documentation](docs/architecture/ATHEMA_ARCHITECTURE.md)
2. Make the [key decisions](docs/architecture/DECISION_MATRIX.md#the-4-critical-decisions)
3. Validate the recommended stack
4. Begin Phase 1: Foundation

---

## Development Philosophy

> *"Build the core right, then make it beautiful."*

- Start with clean interfaces
- Implement one provider at a time
- Test components in isolation
- Iterate based on feel, not features

---

*Architecture proposal v1.0 - Ready for review*
