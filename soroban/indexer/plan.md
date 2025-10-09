# Standalone Indexer Architecture Plan

## Overview

Build a standalone indexer service that polls Stellar RPC API for events and transactions, then writes directly to PostgreSQL. This replaces the fork-based architecture with a clean, maintainable solution.

## Architecture

### Simplified Data Flow (No Redis)

```
┌─────────────────────────────┐
│  Upstream Stellar RPC       │
│  (v23.0.4 - Unmodified)     │
│                             │
│  SQLite (7-day retention)   │
│  Port 8000 (RPC API)        │
└─────────────┬───────────────┘
              │
              │ Poll every 5s
              │ getEvents() / getTransactions()
              ▼
┌─────────────────────────────┐
│  Standalone Indexer Service │
│  (Single Go Binary)         │
│                             │
│  1. Poll RPC API            │
│  2. Parse events/txs        │
│  3. Write to PostgreSQL     │
│  4. Track cursor in DB      │
└─────────────┬───────────────┘
              │
              │ Direct DB writes
              ▼
       ┌──────────────┐
       │ PostgreSQL   │
       │              │
       │  - events    │
       │  - txs       │
       │  - tokens    │
       │  - metadata  │
       │  - cursor    │
       └──────────────┘
```

### Why No Redis?

**Old Architecture Issues:**
- Redis = Extra infrastructure to manage
- Redis = Single point of failure
- Redis = Network hop overhead
- Queue consumer = Separate process to maintain

**New Approach:**
- Direct PostgreSQL writes (simpler)
- PostgreSQL handles concurrency
- Cursor stored in DB (no state loss)
- Single binary (easier deployment)

## Directory Structure

```
indexer/
├── cmd/
│   └── indexer/
│       └── main.go              # Main entry point
├── pkg/
│   ├── client/
│   │   └── rpc.go               # Stellar RPC JSON-RPC client
│   ├── poller/
│   │   ├── poller.go            # Main polling loop
│   │   ├── events.go            # Event polling logic
│   │   └── transactions.go      # Transaction polling logic
│   ├── parser/
│   │   ├── event.go             # Parse RPC events → DB model
│   │   ├── transaction.go       # Parse RPC txs → DB model
│   │   └── token.go             # Parse token operations
│   ├── models/
│   │   ├── event.go             # Event GORM model
│   │   ├── transaction.go       # Transaction GORM model
│   │   ├── token_metadata.go    # Token metadata model
│   │   ├── token_operation.go   # Token operation model
│   │   └── cursor.go            # Cursor tracking model
│   └── db/
│       └── postgres.go          # PostgreSQL connection & helpers
├── go.mod
├── go.sum
├── Dockerfile
└── plan.md                      # This file
```

## Component Design

### 1. RPC Client (`pkg/client/rpc.go`)

**Purpose:** HTTP client for Stellar RPC JSON-RPC API

```go
type Client struct {
    endpoint   string
    httpClient *http.Client
}

// Core methods
func (c *Client) GetLatestLedger(ctx context.Context) (uint32, error)
func (c *Client) GetEvents(ctx context.Context, req GetEventsRequest) (*GetEventsResponse, error)
func (c *Client) GetTransactions(ctx context.Context, startLedger uint32) ([]Transaction, error)
func (c *Client) GetNetwork(ctx context.Context) (*NetworkInfo, error)

// Request/Response types match stellar-rpc protocol
type GetEventsRequest struct {
    StartLedger uint32
    Filters     EventFilters
    Pagination  Pagination
}

type GetEventsResponse struct {
    Events        []Event
    LatestLedger  uint32
}
```

**RPC Methods Used:**
- `getLatestLedger` - Get current ledger number
- `getEvents` - Fetch contract events with filters
- `getTransactions` - Fetch transaction details
- `getNetwork` - Network info (passphrase, protocol)

### 2. Poller (`pkg/poller/`)

**Purpose:** Orchestrate polling, parsing, and writing

```go
type Poller struct {
    rpcClient *client.Client
    db        *gorm.DB
    logger    *log.Entry
    interval  time.Duration
}

func (p *Poller) Start(ctx context.Context) error {
    ticker := time.NewTicker(p.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            if err := p.poll(ctx); err != nil {
                p.logger.WithError(err).Error("poll failed")
            }
        }
    }
}

func (p *Poller) poll(ctx context.Context) error {
    // 1. Get last indexed ledger from DB
    cursor := p.getCursor()

    // 2. Fetch new events from RPC
    events, err := p.pollEvents(ctx, cursor)
    if err != nil {
        return err
    }

    // 3. Fetch new transactions
    txs, err := p.pollTransactions(ctx, cursor)
    if err != nil {
        return err
    }

    // 4. Parse and write to DB (in transaction)
    return p.processBatch(ctx, events, txs)
}
```

**Polling Logic:**
- Poll every 5 seconds (configurable)
- Track cursor (last processed ledger) in DB
- Handle RPC errors with retry/backoff
- Batch writes to PostgreSQL

### 3. Parser (`pkg/parser/`)

**Purpose:** Transform RPC responses to database models

```go
// Event parser
func ParseEvent(rpcEvent protocol.Event) (*models.Event, error) {
    // Parse XDR topics and values
    topics := parseTopics(rpcEvent.Topic)
    value := parseValue(rpcEvent.Value)

    return &models.Event{
        ID:             rpcEvent.ID,
        ContractID:     rpcEvent.ContractID,
        EventType:      rpcEvent.Type,
        Ledger:         rpcEvent.Ledger,
        Topic:          topics,  // JSON
        Value:          value,   // JSON
        // ... other fields
    }, nil
}

// Transaction parser
func ParseTransaction(rpcTx protocol.Transaction) (*models.Transaction, error) {
    // Parse transaction envelope, signatures, etc.
    envelope := parseEnvelope(rpcTx.EnvelopeXdr)

    return &models.Transaction{
        ID:            rpcTx.Hash,
        Status:        rpcTx.Status,
        Ledger:        rpcTx.Ledger,
        SourceAccount: envelope.SourceAccount,
        // ... other fields
    }, nil
}
```

**Reuse from old indexer:**
- XDR parsing logic
- JSON conversion
- Token operation detection

### 4. Database Models (`pkg/models/`)

**Reuse existing GORM models** from old indexer:

```go
// Event model
type Event struct {
    ID                       string      `gorm:"primaryKey"`
    TxIndex                  int32
    EventType                string
    Ledger                   int32
    LedgerClosedAt           string
    ContractID               string
    PagingToken              string
    Topic                    interface{} `gorm:"type:jsonb"`
    Value                    interface{} `gorm:"type:jsonb"`
    InSuccessfulContractCall bool
    LastModifiedLedgerSeq    uint32
    CreatedAt                time.Time
    UpdatedAt                time.Time
}

// Transaction model
type Transaction struct {
    ID               string `gorm:"primaryKey"`
    Status           string
    Ledger           uint32
    ApplicationOrder int32
    SourceAccount    string
    Fee              int32
    FeeCharged       int32
    // ... (keep all fields from old indexer)
}

// Cursor model (NEW)
type Cursor struct {
    ID           int    `gorm:"primaryKey"`
    LastLedger   uint32 `gorm:"not null"`
    UpdatedAt    time.Time
}
```

**Database Operations:**
```go
func UpsertEvent(db *gorm.DB, event *Event) error {
    return db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        DoUpdates: clause.AssignmentColumns([]string{...}),
    }).Create(event).Error
}

func GetCursor(db *gorm.DB) (uint32, error)
func UpdateCursor(db *gorm.DB, ledger uint32) error
```

### 5. Database Layer (`pkg/db/postgres.go`)

**Purpose:** PostgreSQL connection and helpers

```go
type DB struct {
    *gorm.DB
}

func Connect(dsn string) (*DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Error),
    })
    if err != nil {
        return nil, err
    }

    // Auto-migrate models
    db.AutoMigrate(&models.Event{}, &models.Transaction{}, &models.Cursor{})

    return &DB{DB: db}, nil
}

func (db *DB) WithTransaction(fn func(*gorm.DB) error) error {
    return db.Transaction(fn)
}
```

## Data Flow

### Event Processing Flow

```
1. RPC Request:
   POST http://stellar-rpc:8000
   {"jsonrpc":"2.0","method":"getEvents","params":{
     "startLedger": 12345,
     "filters": {...},
     "pagination": {"limit": 1000}
   }}

2. RPC Response:
   {
     "events": [
       {
         "id": "0000123456-0000000001",
         "type": "contract",
         "ledger": 123456,
         "contractId": "CXXX...",
         "topic": ["AAAADwAAAAdkZXBvc2l0"],
         "value": "AAAAAQAAAA8..."
       }
     ],
     "latestLedger": 123500
   }

3. Parse:
   Event{
     ID: "0000123456-0000000001",
     ContractID: "CXXX...",
     Topic: ["deposit"],  // Decoded
     Value: {...},        // Decoded JSON
   }

4. Write to PostgreSQL:
   INSERT INTO events (...) VALUES (...)
   ON CONFLICT (id) DO UPDATE SET ...

5. Update Cursor:
   UPDATE cursor SET last_ledger = 123500
```

### Transaction Processing Flow

Similar to events, but uses `getTransactions` method.

## Configuration

### Environment Variables

```bash
# Stellar RPC
STELLAR_RPC_URL=http://stellar-rpc:8000

# PostgreSQL
POSTGRES_DSN=postgresql://user:pass@postgres:5432/stellar_indexer

# Polling
POLL_INTERVAL=5s            # How often to poll
BATCH_SIZE=1000             # Events per request
START_LEDGER=0              # Initial ledger (0 = latest)

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Configuration Struct

```go
type Config struct {
    RPCURL       string
    PostgresDSN  string
    PollInterval time.Duration
    BatchSize    int
    StartLedger  uint32
    LogLevel     string
}

func LoadConfig() (*Config, error) {
    return &Config{
        RPCURL:       getEnv("STELLAR_RPC_URL", "http://localhost:8000"),
        PostgresDSN:  getEnv("POSTGRES_DSN", ""),
        PollInterval: getDuration("POLL_INTERVAL", 5*time.Second),
        BatchSize:    getInt("BATCH_SIZE", 1000),
        StartLedger:  getUint32("START_LEDGER", 0),
        LogLevel:     getEnv("LOG_LEVEL", "info"),
    }, nil
}
```

## Docker Deployment

### Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o indexer cmd/indexer/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/indexer /app/indexer
WORKDIR /app
CMD ["/app/indexer"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  stellar-rpc:
    image: stellar/stellar-rpc:23.0.4
    ports: ["8000:8000"]
    volumes:
      - ./deploy/config/stellar-rpc.toml:/config/stellar-rpc.toml
      - ./deploy/data/stellar-rpc:/data
      - ./deploy/data/captive-core:/captive-core

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: stellar_indexer
      POSTGRES_USER: stellar
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports: ["5432:5432"]
    volumes:
      - postgres-data:/var/lib/postgresql/data

  indexer:
    build: ./indexer
    depends_on:
      - stellar-rpc
      - postgres
    environment:
      STELLAR_RPC_URL: http://stellar-rpc:8000
      POSTGRES_DSN: postgresql://stellar:${POSTGRES_PASSWORD}@postgres:5432/stellar_indexer
      POLL_INTERVAL: 5s
      LOG_LEVEL: info
    restart: unless-stopped

volumes:
  postgres-data:
```

## Error Handling

### Retry Strategy

```go
type RetryConfig struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
    Multiplier  float64
}

func (p *Poller) pollWithRetry(ctx context.Context) error {
    var lastErr error
    wait := p.retryConfig.InitialWait

    for i := 0; i < p.retryConfig.MaxRetries; i++ {
        if err := p.poll(ctx); err == nil {
            return nil
        } else {
            lastErr = err
            p.logger.WithError(err).Warnf("poll failed, retry %d/%d", i+1, p.retryConfig.MaxRetries)
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(wait):
            wait = time.Duration(float64(wait) * p.retryConfig.Multiplier)
            if wait > p.retryConfig.MaxWait {
                wait = p.retryConfig.MaxWait
            }
        }
    }

    return fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

### Error Types

- **RPC Errors**: Retry with backoff
- **Parse Errors**: Log and skip (don't block pipeline)
- **DB Errors**: Retry with backoff, alert if persistent
- **Network Errors**: Retry immediately (RPC may be restarting)

## Monitoring

### Metrics to Track

```go
type Metrics struct {
    LastProcessedLedger  uint32
    EventsProcessed      int64
    TransactionsProcessed int64
    ErrorCount           int64
    LastPollDuration     time.Duration
    LastPollTime         time.Time
}

// Expose via /metrics endpoint (Prometheus format)
func (p *Poller) exposeMetrics() {
    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "indexer_last_ledger %d\n", p.metrics.LastProcessedLedger)
        fmt.Fprintf(w, "indexer_events_total %d\n", p.metrics.EventsProcessed)
        fmt.Fprintf(w, "indexer_errors_total %d\n", p.metrics.ErrorCount)
    })
    http.ListenAndServe(":9090", nil)
}
```

### Health Check

```go
func (p *Poller) healthCheck() error {
    // Check RPC connectivity
    if _, err := p.rpcClient.GetLatestLedger(context.Background()); err != nil {
        return fmt.Errorf("rpc unhealthy: %w", err)
    }

    // Check DB connectivity
    if err := p.db.Raw("SELECT 1").Error; err != nil {
        return fmt.Errorf("db unhealthy: %w", err)
    }

    // Check lag (< 1 minute behind)
    lag := time.Since(p.metrics.LastPollTime)
    if lag > time.Minute {
        return fmt.Errorf("indexer lagging: %s", lag)
    }

    return nil
}
```

## Migration from Old Indexer

### Step 1: Database Schema

**Reuse existing PostgreSQL tables** - no schema changes needed!

```sql
-- Tables already exist from old indexer
SELECT table_name FROM information_schema.tables
WHERE table_schema = 'public';
-- events
-- transactions
-- token_metadata
-- token_operations
```

### Step 2: Initial Cursor

```sql
-- Set starting point to latest event in DB
INSERT INTO cursor (id, last_ledger, updated_at)
SELECT 1, MAX(ledger), NOW() FROM events;
```

### Step 3: Deployment

1. Stop old indexer (fork-based)
2. Deploy upstream stellar-rpc (v23.0.4)
3. Deploy new standalone indexer
4. Verify events are being indexed

### Step 4: Backfill (Optional)

If there's a gap between old and new indexer:

```bash
# Run indexer with specific start ledger
START_LEDGER=12345000 ./indexer
```

## Advantages Over Old Architecture

| Feature | Old (Fork + Redis) | New (Standalone) |
|---------|-------------------|------------------|
| **Components** | RPC fork + Redis + Consumer | RPC upstream + Indexer |
| **Maintenance** | Very High (merge conflicts) | Low (no fork) |
| **Upgrades** | Very Hard (v20→v23) | Easy (docker pull) |
| **Infrastructure** | RPC + Redis + PostgreSQL | RPC + PostgreSQL |
| **Latency** | 0ms (inline) | ~5s (polling) |
| **Reliability** | Coupled (fork breaks = no indexing) | Independent |
| **Testability** | Hard (need full RPC) | Easy (mock RPC client) |
| **State Management** | Redis queue | PostgreSQL cursor |
| **Recovery** | Complex (replay Redis) | Simple (replay from cursor) |

## Implementation Phases

### Phase 1: Core Infrastructure ✅
- [x] Create directory structure
- [ ] Setup go.mod with dependencies
- [ ] Create PostgreSQL models (copy from old)
- [ ] Setup docker-compose

### Phase 2: RPC Client
- [ ] Implement JSON-RPC client
- [ ] Add getEvents method
- [ ] Add getTransactions method
- [ ] Add getLatestLedger method
- [ ] Add error handling & retries

### Phase 3: Poller & Parser
- [ ] Implement main polling loop
- [ ] Implement event parser
- [ ] Implement transaction parser
- [ ] Add cursor tracking
- [ ] Add batch processing

### Phase 4: Testing
- [ ] Unit tests for parser
- [ ] Integration tests with mock RPC
- [ ] End-to-end test with local stellar-rpc
- [ ] Load testing

### Phase 5: Production
- [ ] Add metrics endpoint
- [ ] Add health check endpoint
- [ ] Add graceful shutdown
- [ ] Production deployment
- [ ] Monitoring & alerts

## Performance Considerations

### Polling Interval

- **5 seconds**: Good balance (12 polls/minute)
- **Too fast (<2s)**: Wastes resources, hits RPC limits
- **Too slow (>10s)**: Increases lag

### Batch Size

- **1000 events**: Good default
- RPC SQLite has 7-day retention (~120K ledgers)
- At ~10 events/ledger = ~1.2M events max
- Pagination handles large batches

### Database Optimization

```sql
-- Indexes for fast queries
CREATE INDEX idx_events_ledger ON events(ledger);
CREATE INDEX idx_events_contract_id ON events(contract_id);
CREATE INDEX idx_events_ledger_closed_at ON events(ledger_closed_at);
CREATE INDEX idx_txs_ledger ON transactions(ledger);
```

### Resource Usage

- **Memory**: ~100MB (minimal)
- **CPU**: <5% (mostly idle)
- **Network**: ~1KB/s polling overhead
- **Database**: Same as old indexer

## Future Enhancements

1. **Multi-threaded polling**: Parallel event/tx fetching
2. **GraphQL API**: Query indexed data
3. **Webhook support**: Real-time notifications
4. **Historical backfill**: Index from genesis
5. **Multiple networks**: Index testnet + mainnet
6. **Event filtering**: Only index specific contracts

## Conclusion

This architecture is **simpler, more maintainable, and easier to operate** than the fork-based approach, with acceptable 5-second latency for analytics use cases.

**Key Benefits:**
- No fork to maintain
- Easy upgrades (docker pull)
- Single binary deployment
- PostgreSQL as single source of truth
- Clear separation of concerns
