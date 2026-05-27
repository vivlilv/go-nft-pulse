# go-nft-pulse

> Real-time NFT trading bot and analytics platform built in Go.

**go-nft-pulse** ingests live NFT market events from OpenSea via WebSocket — planned to extend to Kafka streams and raw on-chain data — and processes them through a pure event-driven pipeline. Trading strategies run continuously against a live state snapshot, while the web dashboard lets users add or remove target items as they use it.

---

## Architecture

The system is built around a unidirectional data flow:

```
OpenSea WebSocket
       │
       ▼
  Event Layer          — parses raw WS events into typed domain events
       │
       ▼
   Reducer             — pure function; folds events into state
       │
       ▼
  State Manager        — concurrent-safe state with RWMutex snapshots
       │
       ▼
   Strategy            — evaluates state, emits trading intents
       │
       ▼
  Execution Layer      — acts on intents (place/cancel offers)  [planned]
       │
       ▼Lipgloss
  Web Dashboard        — live UI for settings, insights, control  [in progress]
```

Each layer is decoupled — the reducer is a pure function with no side effects, the strategy layer only reads snapshots, and execution is fully isolated from market data ingestion.

---

## Features

### Implemented
- **OpenSea WebSocket listener** — subscribes to collection, trait, and item offer events in real time
- **Event-driven reducer** — handles collection offers, trait offers, item offers, cancellations, and sales
- **Offer expiration** — synthetic expiry events via `ExpiryScheduler` using `time.AfterFunc`; keeps state clean without polling
- **Analytics TUI** — live terminal table (Bubbletea), 200ms tick-based state snapshots
- **Strategy layer** — evaluates current market state and emits intents for collection-wide offers

### In Progress
- **Web dashboard** — live HTTP server; user-facing interface for configuring trading parameters and viewing real-time analytics

### Planned
- **Execution layer** — place and cancel offers on OpenSea; graceful bot shutdown with notification when a target item is purchased
- **Richer analytics** — offer depth, trait rarity overlays

---

## Tech Stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 |
| Market data | OpenSea WebSocket API |
| Terminal UI | [Bubbletea](https://github.com/charmbracelet/bubbletea) + Lipgloss |
| Web server | Go `net/http` |
| Config | `envconfig` + `.env` |

---

## Getting Started

### Prerequisites
- Go 1.21+
- OpenSea API key
- A wallet address (can be a read-only address for observation mode)

### Setup

```bash
git clone https://github.com/vivlilv/go-nft-pulse.git
cd go-nft-pulse
cp .env.example .env   # fill in your OpenSea API key and wallet address
```

Configure your target NFT collections and items in `items.json`:

```json
[
  {
    "id": "ethereum/0xContractAddress/tokenId",
    "slug": "collection-slug"
  }
]
```
Example of real data:
```
[
    {
        "nftID": "ethereum/0xbd3531da5cf5857e7cfaa92426877b022e612cf8/920",
        "Slug": "pudgypenguins",
        "isPending": false,
        "tokenID": 6,
        "maxBidAllowedWei": 1000000000000000000,
        "offerStepWei": 100000000000000000,
        "traits": [
            {
                "trait_type": "Background",
                "trait_name": "Tangerine"
            },
            {
                "trait_type": "Skin",
                "trait_name": "Olive Green"
            },
            {
                "trait_type": "Body",
                "trait_name": "Lei Pink"
            },
            {
                "trait_type": "Face",
                "trait_name": "Beard"
            },
            {
                "trait_type": "Head",
                "trait_name": "Jester's Hat"
            }
        ],
        "ImgURL": "https://gateway.pinata.cloud/ipfs/QmNf1UsmdGaMbpatQ6toXSkzDpizaGmC9zfunCyoz1enD5/penguin/6.png"
    }
]
```

### Run

```bash
go run ./cmd/nft_trader/main.go
```

the analytics TUI launches in the terminal. The web dashboard is available at `http://localhost:8080` (port configurable via env).

---

## Screenshots

> *Coming soon — TUI analytics and web dashboard previews.*

---

## Project Status

active development. Core pipeline (events → reducer → state → strategy) is functional. Web dashboard and execution layer are the current focus.

---

*Built as a learning project — practicing production-like Go patterns: event-driven architecture, pure reducer functions, concurrent state management, and real-time data pipelines.*
