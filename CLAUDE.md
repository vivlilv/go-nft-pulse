# go_nft_bot

## Project
Pet project to practice Golang and build a useful NFT trading bot.
Architecture: event-driven (events → reducer → state → strategy → execution layers).

## Current Status
Implemented:
- **events layer** — OpenSea WebSocket events, parsed into domain types (`internal/events/opensea/`)
- **reducer layer** — pure function, handles collection/trait/item offers + cancels + sales (`internal/reducer/`)
- **offer expiration** — synthetic `ExpiredOfferEvent` via `ExpiryScheduler` + `time.AfterFunc` (`internal/events/opensea/expiry.go`) — implemented, needs integration testing
- **strategy layer** - currently on hold, implemented partially

## Architecture Notes (for AI to load fast)
- `domain.Event` is an empty interface; reducer type-switches on **value types** (e.g. `domain.CollectionOfferEvent`, not pointer)
- `State.Items` is keyed by `domain.NftID` (full string: `"chain/contract/tokenID"`); items must have `Slug` set in `items.json` for slug-based filtering to work
- Reducer helpers: `selectAllItemsForSlug`, `filterItemsByTraits` (filters by slug first, then traits), `applyOffer`, `clearOffer`

## Key decisions
- *big.Int fields in ItemState are treated as immutable - never mutated in-place. New values are always created with new(big.Int). This makes Snapshot() safe without deep-copying them.
- The reason *state is passed to Reduce is: you want to pass state by value (a copy) so the reducer can't accidentally mutate the live state. The * is just dereferencing the pointer to get a value copy. That way old state remains unchanged until the new state atomically updated is replacing old one in stateManager.

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
