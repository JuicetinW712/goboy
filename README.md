# GoBoy - Game Boy Emulator

Welcome to **GoBoy**, a Game Boy (DMG) emulator written in Go and Ebitengine

## Getting Started

### Prerequisites
- **Go 1.25+**
- **C Compiler** (GCC or Clang) for Ebitengine's dependencies (OpenGL/ALSA on Linux).

### Build and Run
```bash
> go build -o goboy app/main.go
> ./goboy [flags] roms/YourGame.gb
```

**Flags:**
- `-scale <int>`: Window scale (default: 5)
- `-debug`: Enable FPS/TPS overlays
- `-nosync`: Disable frame synchronization (unlocks speed)

### Controls
| Game Boy | Keyboard |
| :--- | :--- |
| **D-Pad** | `W`, `A`, `S`, `D` |
| **A / B** | `Z` / `X` |
| **Start / Select** | `Enter` / `Shift` |

---

## TODO

Several hardware features and optimizations are still in development:

- [ ] **PPU Enhancements**:
  - Enforce the hardware limit of 10 sprites per scanline.
  - Implement full sprite-to-background priority logic.
- [ ] **Timing & Accuracy**:
  - Add the required cycle penalty for DMA transfers (CPU should be stalled).
  - Refine timer behavior to ensure they continue ticking during `HALT` but pause during `STOP`. (low priority / not important)
- [ ] **Input**:
  - Fully implement the Joypad interrupt (low priority / not important).
- [] **Audio (APU)**:
  - Implementation of Pulse, Wave, and Noise channels for full sound support. (low priority)
