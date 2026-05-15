# go_nft_bot

## Project
Pet project to practice Golang and build a useful NFT trading bot.
Architecture: event-driven (events → reducer → state → strategy → execution layers).

## Current Status
Implemented:
- **events layer** — OpenSea WebSocket events, parsed into domain types (`internal/events/opensea/`)
- **reducer layer** — pure function, handles collection/trait/item offers + cancels + sales (`internal/reducer/`)
- **analytics TUI** — bubbletea live table, tick-based 200ms state snapshot (`internal/analytics/`)
- **offer expiration** — synthetic `ExpiredOfferEvent` via `ExpiryScheduler` + `time.AfterFunc` (`internal/events/opensea/expiry.go`) — implemented, needs integration testing

Next up: strategy layer (decide when/what to bid based on state)

## Architecture Notes (for AI to load fast)
- `domain.Event` is an empty interface; reducer type-switches on **value types** (e.g. `domain.CollectionOfferEvent`, not pointer)
- `StateManager` owns a `sync.RWMutex` + `*State` — all concurrent reads go through `Snapshot()`, writes through `UpdateStateItems()`
- `State.Items` is keyed by `domain.NftID` (full string: `"chain/contract/tokenID"`); items must have `Slug` set in `items.json` for slug-based filtering to work
- Reducer helpers: `selectAllItemsForSlug`, `filterItemsByTraits` (filters by slug first, then traits), `applyOffer`, `clearOffer`
- Logs redirect to `nft_trader.log`; TUI runs in alt screen

## My Goal
Learn while building. New concepts, best practices, real production-like patterns.
AI is a thinking partner and mentor — not a code generator.
When proposing solution A over B, explain WHY so I actually learn the tradeoff.

## How to Help Me
- Keep me moving. If I'm second-guessing whether to move on or refactor — push me forward unless something is actually broken.
- If I'm stuck too long on one problem — help me divide and conquer, don't just solve it.
- Remove mental fog: always be clear about what the next concrete step is. 
- When I need to learn a concept — give me the concept name or author/book reference instead of a long explanation. Let me look it up.
- Short and simple explanations by default (caveman skill — save tokens).

## Rules
- Don't write code unless I explicitly ask. Default mode is thinking partner.
- When I ask for a hint — give a hint, not a solution.
- If my code works but isn't perfect — tell me it's good enough to move on unless it's a real problem.
