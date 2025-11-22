# Microservices Architecture for Abhi SDK Integration

## Executive Summary

This document outlines a comprehensive microservices architecture strategy for integrating the Abhi Go SDK into your backend services using API Gateway with RabbitMQ messaging. The architecture emphasizes security, scalability, and maintainability while leveraging all enhanced features of the Abhi SDK.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Service Decomposition](#service-decomposition)
3. [API Gateway Design](#api-gateway-design)
4. [RabbitMQ Integration](#rabbitmq-integration)
5. [Abhi Gateway Service](#abhi-gateway-service)
6. [Security Architecture](#security-architecture)
7. [Data Flow Patterns](#data-flow-patterns)
8. [Implementation Strategy](#implementation-strategy)
9. [Deployment Architecture](#deployment-architecture)
10. [Monitoring & Observability](#monitoring--observability)
11. [Best Practices & Guidelines](#best-practices--guidelines)

---

## Architecture Overview

### High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                 CLIENT LAYER                                     │
├─────────────────────────┬─────────────────────────┬─────────────────────────────┤
│     Web Dashboard       │      Mobile Apps        │      Third-Party APIs      │
└─────────────────────────┴─────────────────────────┴─────────────────────────────┘
                                          │
                                          ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                               API GATEWAY                                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                │
│  │  Authentication │  │  Rate Limiting  │  │  Load Balancer  │                │
│  │  & Authorization│  │  & Throttling   │  │  & Routing      │                │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                │
│  │  Request/Response│  │  Logging &      │  │  Circuit        │                │
│  │  Transformation │  │  Monitoring     │  │  Breaker        │                │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                │
└─────────────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              MESSAGE BROKER (RabbitMQ)                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                │
│  │   Work Queues   │  │  Publish/Subscribe │ │  RPC Queues   │                │
│  │   • Commands    │  │  • Events        │  │  • Sync Calls │                │
│  │   • Jobs        │  │  • Notifications │  │  • Queries    │                │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                │
└─────────────────────────────────────────────────────────────────────────────────┘
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    │                     │                     │
                    ▼                     ▼                     ▼
┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
│    EMPLOYEE SERVICE     │ │   TRANSACTION SERVICE   │ │  ORGANIZATION SERVICE  │
│                         │ │                         │ │                         │
│ • Employee CRUD         │ │ • Advance Requests      │ │ • Org Management        │
│ • Profile Management    │ │ • Repayment Processing  │ │ • Business Types        │
│ • Department Handling   │ │ • Balance Calculations  │ │ • Credit Limits         │
│ • Search & Filtering    │ │ • Transaction History   │ │ • Master Data           │
│                         │ │ • Validation Logic      │ │ • Hierarchy Management  │
└─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘
                    │                     │                     │
                    └─────────────────────┼─────────────────────┘
                                          │
                                          ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            ABHI GATEWAY SERVICE                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                │
│  │   Abhi SDK      │  │  Security       │  │  Circuit        │                │
│  │   Integration   │  │  Management     │  │  Breaker        │                │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                │
│  │  Rate Limiting  │  │  Token          │  │  Request        │                │
│  │  & Throttling   │  │  Management     │  │  Signing        │                │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                │
└─────────────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              ABHI OPEN API                                       │
│                          (External Service)                                      │
└─────────────────────────────────────────────────────────────────────────────────┘

                    SHARED INFRASTRUCTURE & CROSS-CUTTING CONCERNS
┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
│      AUTH SERVICE       │ │     CONFIG SERVICE      │ │   NOTIFICATION SERVICE  │
│                         │ │                         │ │                         │
│ • JWT Token Management  │ │ • Environment Config    │ │ • Email/SMS Sending     │
│ • Session Handling      │ │ • Feature Flags        │ │ • Push Notifications    │
│ • MFA Support          │ │ • Secret Management     │ │ • Event Broadcasting    │
│ • User Authorization    │ │ • Service Discovery     │ │ • Template Management   │
└─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘

┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
│    LOGGING SERVICE      │ │    MONITORING SERVICE   │ │     CACHE SERVICE       │
│                         │ │                         │ │                         │
│ • Centralized Logging   │ │ • Metrics Collection    │ │ • Redis Cluster         │
│ • Log Aggregation       │ │ • Health Checks        │ │ • Session Storage       │
│ • Search & Analytics    │ │ • Alerting Rules        │ │ • Distributed Cache     │
│ • Audit Trail          │ │ • Performance Tracking  │ │ • Rate Limit Storage    │
└─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘
```

---

## Service Decomposition

### Core Business Services

#### 1. Employee Service
```
Responsibilities:
├── Employee Lifecycle Management
│   ├── Create/Update/Delete employees
│   ├── Employee profile management
│   ├── Department and role assignments
│   └── Employment status tracking
├── Search & Discovery
│   ├── Employee search functionality
│   ├── Directory services
│   ├── Reporting and analytics
│   └── Bulk operations
└── Data Synchronization
    ├── Sync with Abhi API
    ├── Local cache management
    └── Event publishing
```

#### 2. Transaction Service
```
Responsibilities:
├── Transaction Processing
│   ├── Advance request creation
│   ├── Transaction validation
│   ├── Approval workflows
│   └── Status tracking
├── Financial Operations
│   ├── Balance calculations
│   ├── Repayment processing
│   ├── Interest calculations
│   └── Fee management
└── History & Reporting
    ├── Transaction history
    ├── Monthly statements
    ├── Financial reports
    └── Audit logs
```

#### 3. Organization Service
```
Responsibilities:
├── Organization Management
│   ├── Organization CRUD operations
│   ├── Hierarchy management
│   ├── Business type handling
│   └── Configuration management
├── Master Data
│   ├── Bank information
│   ├── Business types
│   ├── Industry classifications
│   └── Reference data
└── Credit & Limits
    ├── Credit limit management
    ├── Risk assessment
    ├── Policy enforcement
    └── Compliance tracking
```

### Infrastructure Services

#### 4. Abhi Gateway Service (Dedicated)
```
Responsibilities:
├── Abhi SDK Management
│   ├── SDK instance lifecycle
│   ├── Configuration management
│   ├── Connection pooling
│   └── Health monitoring
├── Security Implementation
│   ├── Request signing (HMAC-SHA256)
│   ├── Credential encryption (AES-GCM)
│   ├── Rate limiting (Token bucket)
│   └── Circuit breaker pattern
├── Message Processing
│   ├── RabbitMQ message handling
│   ├── Request/response correlation
│   ├── Error handling & retry
│   └── Dead letter queue management
└── Monitoring & Logging
    ├── Performance metrics
    ├── Security events
    ├── API usage tracking
    └── Error reporting
```

---

## API Gateway Design

### Gateway Architecture

```go
// api-gateway/internal/gateway/gateway.go
type Gateway struct {
    router       *gin.Engine
    authService  *AuthService
    rateLimiter  *RateLimiter
    publisher    *rabbitmq.Publisher
    circuitbreaker.CircuitBreakerManager
    metrics      *prometheus.Registry
}

type GatewayConfig struct {
    Port               string
    RabbitMQURL        string
    RedisURL           string
    JWTSecret          string
    RateLimitConfig    *RateLimitConfig
    CircuitBreakerConfig *CircuitBreakerConfig
    CORSConfig         *CORSConfig
}

func NewGateway(config *GatewayConfig) *Gateway {
    gateway := &Gateway{
        router:       gin.New(),
        publisher:    rabbitmq.NewPublisher(config.RabbitMQURL),
        rateLimiter:  NewRateLimiter(config.RateLimitConfig),
        metrics:      prometheus.NewRegistry(),
    }
    
    gateway.setupMiddleware()
    gateway.setupRoutes()
    return gateway
}
```

### Middleware Stack

```go
// api-gateway/internal/middleware/middleware.go

// 1. Request ID & Correlation
func RequestIDMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Header("X-Request-ID", requestID)
        c.Set("request_id", requestID)
        c.Next()
    })
}

// 2. CORS Middleware
func CORSMiddleware(config *CORSConfig) gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     config.AllowOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Request-ID"},
        ExposeHeaders:    []string{"X-Request-ID", "X-Rate-Limit-Remaining"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}

// 3. Rate Limiting Middleware
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        clientIP := c.ClientIP()
        userID := c.GetString("user_id")
        
        key := fmt.Sprintf("rate_limit:%s:%s", clientIP, userID)
        
        allowed, remaining, err := limiter.Allow(key)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
            c.Abort()
            return
        }
        
        c.Header("X-Rate-Limit-Remaining", strconv.Itoa(remaining))
        
        if !allowed {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": limiter.RetryAfter(key),
            })
            c.Abort()
            return
        }
        
        c.Next()
    })
}

// 4. Authentication Middleware
func AuthenticationMiddleware(authService *AuthService) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    })
}

// 5. Logging Middleware
func LoggingMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("[%s] %s %s %d %s %s %s\n",
            param.TimeStamp.Format("2006/01/02 15:04:05"),
            param.ClientIP,
            param.Method,
            param.StatusCode,
            param.Latency,
            param.Path,
            param.ErrorMessage,
        )
    })
}

// 6. Circuit Breaker Middleware
func CircuitBreakerMiddleware(manager *circuitbreaker.Manager) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        service := c.GetHeader("X-Target-Service")
        if service == "" {
            service = "default"
        }
        
        cb := manager.GetCircuitBreaker(service)
        if cb.State() == circuitbreaker.StateOpen {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "error": "Service temporarily unavailable",
                "service": service,
            })
            c.Abort()
            return
        }
        
        c.Next()
    })
}
```

### Route Handlers

```go
// api-gateway/internal/handlers/handlers.go

type Handlers struct {
    publisher *rabbitmq.Publisher
    consumer  *rabbitmq.Consumer
    redis     *redis.Client
}

// Employee Management Endpoints
func (h *Handlers) setupEmployeeRoutes(rg *gin.RouterGroup) {
    employees := rg.Group("/employees")
    {
        employees.GET("", h.handleListEmployees)
        employees.POST("", h.handleCreateEmployee)
        employees.GET("/:id", h.handleGetEmployee)
        employees.PUT("/:id", h.handleUpdateEmployee)
        employees.DELETE("/:id", h.handleDeleteEmployee)
        employees.GET("/search", h.handleSearchEmployees)
        employees.POST("/bulk", h.handleBulkEmployeeOperations)
    }
}

func (h *Handlers) handleCreateEmployee(c *gin.Context) {
    var req CreateEmployeeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Add metadata
    req.RequestID = c.GetString("request_id")
    req.UserID = c.GetString("user_id")
    req.Timestamp = time.Now()
    
    // Publish to RabbitMQ
    correlationID := uuid.New().String()
    message := &rabbitmq.Message{
        ID:            correlationID,
        Type:          "employee.create",
        Payload:       req,
        ReplyTo:       "api-gateway.responses",
        CorrelationID: correlationID,
    }
    
    err := h.publisher.PublishWithResponse(
        "employee.commands",
        message,
        30*time.Second, // timeout
    )
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
        return
    }
    
    // Wait for response
    response, err := h.consumer.WaitForResponse(correlationID, 30*time.Second)
    if err != nil {
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
        return
    }
    
    if response.Error != nil {
        c.JSON(response.StatusCode, gin.H{"error": response.Error.Message})
        return
    }
    
    c.JSON(http.StatusCreated, response.Data)
}

// Transaction Management Endpoints
func (h *Handlers) setupTransactionRoutes(rg *gin.RouterGroup) {
    transactions := rg.Group("/transactions")
    {
        transactions.POST("/advance", h.handleCreateAdvanceRequest)
        transactions.POST("/repayment", h.handleCreateRepayment)
        transactions.GET("/employee/:employeeId", h.handleGetEmployeeTransactions)
        transactions.GET("/:id", h.handleGetTransaction)
        transactions.PUT("/:id/status", h.handleUpdateTransactionStatus)
        transactions.POST("/validate", h.handleValidateTransaction)
        transactions.GET("/balance/:employeeId", h.handleGetBalance)
    }
}

func (h *Handlers) handleCreateAdvanceRequest(c *gin.Context) {
    var req CreateAdvanceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Validate permissions
    userRole := c.GetString("user_role")
    if !h.canCreateTransaction(userRole, req) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        return
    }
    
    req.RequestID = c.GetString("request_id")
    req.UserID = c.GetString("user_id")
    
    correlationID := uuid.New().String()
    message := &rabbitmq.Message{
        ID:            correlationID,
        Type:          "transaction.advance.create",
        Payload:       req,
        ReplyTo:       "api-gateway.responses",
        CorrelationID: correlationID,
        Priority:      2, // Higher priority for transactions
    }
    
    err := h.publisher.PublishWithResponse("transaction.commands", message, 45*time.Second)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
        return
    }
    
    response, err := h.consumer.WaitForResponse(correlationID, 45*time.Second)
    if err != nil {
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
        return
    }
    
    c.JSON(response.StatusCode, response.Data)
}

// Organization Management Endpoints
func (h *Handlers) setupOrganizationRoutes(rg *gin.RouterGroup) {
    organizations := rg.Group("/organizations")
    {
        organizations.GET("", h.handleListOrganizations)
        organizations.POST("", h.handleCreateOrganization)
        organizations.GET("/:id", h.handleGetOrganization)
        organizations.PUT("/:id", h.handleUpdateOrganization)
        organizations.GET("/:id/statistics", h.handleGetOrganizationStats)
        organizations.GET("/business-types", h.handleGetBusinessTypes)
        organizations.GET("/banks", h.handleGetBanks)
    }
}
```

---

## RabbitMQ Integration

### Message Broker Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              RABBITMQ BROKER                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐            │
│  │   EXCHANGES     │    │     QUEUES      │    │   DEAD LETTER   │            │
│  │                 │    │                 │    │     QUEUES      │            │
│  │ • Direct        │    │ • Work Queues   │    │ • Failed Msgs   │            │
│  │ • Fanout        │    │ • RPC Queues    │    │ • Retry Logic   │            │
│  │ • Topic         │    │ • Event Streams │    │ • Error Handling│            │
│  │ • Headers       │    │ • Priority Qs   │    │                 │            │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘            │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘

Message Flow Patterns:

1. Request/Response (RPC)
   API Gateway ───> Service Queue ───> Service ───> Response Queue ───> API Gateway

2. Event Publishing (Pub/Sub)
   Service ───> Event Exchange ───> Multiple Queues ───> Multiple Services

3. Work Distribution
   Producer ───> Work Queue ───> Worker Pool ───> Result Storage

4. Priority Processing
   High Priority ───> Priority Queue ───> Immediate Processing
   Low Priority  ───> Standard Queue ───> Batch Processing
```

### RabbitMQ Configuration

```go
// shared/messaging/rabbitmq/config.go
type RabbitMQConfig struct {
    URL                string        `json:"url" env:"RABBITMQ_URL"`
    ConnectionTimeout  time.Duration `json:"connection_timeout" env:"RABBITMQ_CONN_TIMEOUT" default:"30s"`
    HeartbeatInterval  time.Duration `json:"heartbeat_interval" env:"RABBITMQ_HEARTBEAT" default:"10s"`
    
    // Exchange configurations
    Exchanges []ExchangeConfig `json:"exchanges"`
    
    // Queue configurations  
    Queues []QueueConfig `json:"queues"`
    
    // Consumer settings
    ConsumerConfig ConsumerConfig `json:"consumer_config"`
    
    // Publisher settings
    PublisherConfig PublisherConfig `json:"publisher_config"`
}

type ExchangeConfig struct {
    Name        string            `json:"name"`
    Type        string            `json:"type"` // direct, fanout, topic, headers
    Durable     bool              `json:"durable"`
    AutoDelete  bool              `json:"auto_delete"`
    Arguments   map[string]interface{} `json:"arguments"`
}

type QueueConfig struct {
    Name        string                 `json:"name"`
    Durable     bool                   `json:"durable"`
    AutoDelete  bool                   `json:"auto_delete"`
    Exclusive   bool                   `json:"exclusive"`
    Arguments   map[string]interface{} `json:"arguments"`
    RoutingKeys []string               `json:"routing_keys"`
    Exchange    string                 `json:"exchange"`
}
```

### Message Patterns

```go
// shared/messaging/rabbitmq/publisher.go
type Publisher struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    config  *PublisherConfig
    logger  *logrus.Logger
}

type Message struct {
    ID            string                 `json:"id"`
    Type          string                 `json:"type"`
    Payload       interface{}            `json:"payload"`
    Headers       map[string]interface{} `json:"headers"`
    Timestamp     time.Time              `json:"timestamp"`
    CorrelationID string                 `json:"correlation_id"`
    ReplyTo       string                 `json:"reply_to"`
    Priority      uint8                  `json:"priority"`
    Expiration    string                 `json:"expiration"`
}

// 1. Fire-and-forget pattern
func (p *Publisher) Publish(exchange, routingKey string, message *Message) error {
    body, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }
    
    return p.channel.Publish(
        exchange,   // exchange
        routingKey, // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType:   "application/json",
            Body:          body,
            MessageId:     message.ID,
            CorrelationId: message.CorrelationID,
            ReplyTo:       message.ReplyTo,
            Priority:      message.Priority,
            Timestamp:     message.Timestamp,
        },
    )
}

// 2. Request/Response pattern
func (p *Publisher) PublishWithResponse(queue string, message *Message, timeout time.Duration) error {
    responseQueue, err := p.declareResponseQueue()
    if err != nil {
        return err
    }
    
    message.ReplyTo = responseQueue
    message.CorrelationID = uuid.New().String()
    
    return p.Publish("", queue, message)
}

// 3. Event publishing pattern
func (p *Publisher) PublishEvent(eventType string, payload interface{}) error {
    message := &Message{
        ID:        uuid.New().String(),
        Type:      eventType,
        Payload:   payload,
        Timestamp: time.Now(),
    }
    
    return p.Publish("events", eventType, message)
}

// shared/messaging/rabbitmq/consumer.go
type Consumer struct {
    conn     *amqp.Connection
    channel  *amqp.Channel
    config   *ConsumerConfig
    handlers map[string]MessageHandler
    logger   *logrus.Logger
}

type MessageHandler func(ctx context.Context, message *Message) (*Response, error)

type Response struct {
    Data       interface{} `json:"data"`
    Error      *Error      `json:"error,omitempty"`
    StatusCode int         `json:"status_code"`
}

func (c *Consumer) RegisterHandler(messageType string, handler MessageHandler) {
    c.handlers[messageType] = handler
}

func (c *Consumer) Start(queueName string) error {
    msgs, err := c.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return err
    }
    
    go c.processMessages(msgs)
    return nil
}

func (c *Consumer) processMessages(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        go c.handleMessage(d)
    }
}

func (c *Consumer) handleMessage(delivery amqp.Delivery) {
    var message Message
    err := json.Unmarshal(delivery.Body, &message)
    if err != nil {
        c.logger.Errorf("Failed to unmarshal message: %v", err)
        delivery.Nack(false, false)
        return
    }
    
    handler, exists := c.handlers[message.Type]
    if !exists {
        c.logger.Warnf("No handler found for message type: %s", message.Type)
        delivery.Nack(false, false)
        return
    }
    
    ctx := context.Background()
    response, err := handler(ctx, &message)
    
    if err != nil {
        c.handleError(delivery, message, err)
        return
    }
    
    // Send response if reply-to is specified
    if delivery.ReplyTo != "" {
        c.sendResponse(delivery, response)
    }
    
    delivery.Ack(false)
}

func (c *Consumer) sendResponse(delivery amqp.Delivery, response *Response) {
    body, err := json.Marshal(response)
    if err != nil {
        c.logger.Errorf("Failed to marshal response: %v", err)
        return
    }
    
    err = c.channel.Publish(
        "",              // exchange
        delivery.ReplyTo, // routing key
        false,           // mandatory
        false,           // immediate
        amqp.Publishing{
            ContentType:   "application/json",
            CorrelationId: delivery.CorrelationId,
            Body:          body,
        },
    )
    
    if err != nil {
        c.logger.Errorf("Failed to send response: %v", err)
    }
}
```

### Queue Topology

```yaml
# Infrastructure configuration for RabbitMQ exchanges and queues
exchanges:
  - name: "employee.commands"
    type: "direct"
    durable: true
    auto_delete: false
    
  - name: "transaction.commands" 
    type: "direct"
    durable: true
    auto_delete: false
    
  - name: "organization.commands"
    type: "direct" 
    durable: true
    auto_delete: false
    
  - name: "abhi.gateway.commands"
    type: "direct"
    durable: true
    auto_delete: false
    
  - name: "events"
    type: "topic"
    durable: true
    auto_delete: false

queues:
  # Command Queues (Request/Response)
  - name: "employee.create"
    durable: true
    exchange: "employee.commands"
    routing_keys: ["employee.create"]
    arguments:
      x-max-priority: 10
      
  - name: "employee.update"
    durable: true
    exchange: "employee.commands" 
    routing_keys: ["employee.update"]
    
  - name: "employee.delete"
    durable: true
    exchange: "employee.commands"
    routing_keys: ["employee.delete"]
    
  - name: "transaction.advance.create"
    durable: true
    exchange: "transaction.commands"
    routing_keys: ["transaction.advance.create"]
    arguments:
      x-max-priority: 10
      
  - name: "transaction.repayment.create"
    durable: true
    exchange: "transaction.commands"
    routing_keys: ["transaction.repayment.create"]
    
  - name: "organization.create"
    durable: true
    exchange: "organization.commands"
    routing_keys: ["organization.create"]
    
  # Abhi Gateway Queues
  - name: "abhi.employee.sync"
    durable: true
    exchange: "abhi.gateway.commands"
    routing_keys: ["abhi.employee.*"]
    
  - name: "abhi.transaction.process"
    durable: true
    exchange: "abhi.gateway.commands"
    routing_keys: ["abhi.transaction.*"]
    arguments:
      x-max-priority: 10
      
  - name: "abhi.organization.sync"
    durable: true
    exchange: "abhi.gateway.commands"
    routing_keys: ["abhi.organization.*"]
    
  # Event Queues (Pub/Sub)
  - name: "employee.events.notification"
    durable: true
    exchange: "events"
    routing_keys: ["employee.*"]
    
  - name: "transaction.events.audit"
    durable: true
    exchange: "events"
    routing_keys: ["transaction.*"]
    
  - name: "organization.events.sync"
    durable: true
    exchange: "events" 
    routing_keys: ["organization.*"]
    
  # Dead Letter Queues
  - name: "employee.dlq"
    durable: true
    arguments:
      x-message-ttl: 86400000  # 24 hours
      
  - name: "transaction.dlq"
    durable: true
    arguments:
      x-message-ttl: 86400000
      
  - name: "abhi.gateway.dlq"
    durable: true
    arguments:
      x-message-ttl: 86400000
```

---

## Abhi Gateway Service

### Service Architecture

```go
// abhi-gateway-service/internal/service/abhi_gateway.go
type AbhiGatewayService struct {
    abhiSDK         *abhi.SDK
    publisher       *rabbitmq.Publisher
    consumer        *rabbitmq.Consumer
    redis           *redis.Client
    circuitBreaker  *circuitbreaker.CircuitBreaker
    metrics         *AbhiMetrics
    config          *AbhiGatewayConfig
    logger          *logrus.Logger
}

type AbhiGatewayConfig struct {
    AbhiConfig      *AbhiConfig
    RabbitMQConfig  *RabbitMQConfig
    RedisConfig     *RedisConfig
    CircuitBreaker  *CircuitBreakerConfig
    Security        *SecurityConfig
}

func NewAbhiGatewayService(config *AbhiGatewayConfig) (*AbhiGatewayService, error) {
    // Initialize Abhi SDK with enhanced security
    var abhiSDK *abhi.SDK
    if config.AbhiConfig.Environment == "production" {
        abhiSDK = abhi.NewForProduction(
            config.AbhiConfig.Username, 
            config.AbhiConfig.Password,
        )
    } else {
        abhiSDK = abhi.NewForUAT(
            config.AbhiConfig.Username, 
            config.AbhiConfig.Password,
        )
    }
    
    // Configure security features
    abhiSDK.SetRateLimit(
        config.AbhiConfig.RateLimitRPS, 
        config.AbhiConfig.RateLimitBurst,
    )
    abhiSDK.EnableRequestSigning(config.Security.SigningSecret)
    abhiSDK.EnableCredentialEncryption(config.Security.EncryptionPassword)
    
    service := &AbhiGatewayService{
        abhiSDK: abhiSDK,
        config:  config,
        metrics: NewAbhiMetrics(),
        logger:  logrus.NewEntry(logrus.New()),
    }
    
    return service, nil
}

func (s *AbhiGatewayService) Start(ctx context.Context) error {
    // Initialize RabbitMQ connections
    if err := s.initializeMessaging(); err != nil {
        return fmt.Errorf("failed to initialize messaging: %w", err)
    }
    
    // Register message handlers
    s.registerMessageHandlers()
    
    // Start consuming messages
    if err := s.startConsumers(); err != nil {
        return fmt.Errorf("failed to start consumers: %w", err)
    }
    
    // Start health check routine
    go s.healthCheckRoutine(ctx)
    
    s.logger.Info("Abhi Gateway Service started successfully")
    return nil
}

func (s *AbhiGatewayService) registerMessageHandlers() {
    // Employee operations
    s.consumer.RegisterHandler("employee.create", s.handleEmployeeCreate)
    s.consumer.RegisterHandler("employee.update", s.handleEmployeeUpdate)
    s.consumer.RegisterHandler("employee.delete", s.handleEmployeeDelete)
    s.consumer.RegisterHandler("employee.get", s.handleEmployeeGet)
    s.consumer.RegisterHandler("employee.list", s.handleEmployeeList)
    s.consumer.RegisterHandler("employee.search", s.handleEmployeeSearch)
    
    // Transaction operations
    s.consumer.RegisterHandler("transaction.advance.create", s.handleTransactionAdvanceCreate)
    s.consumer.RegisterHandler("transaction.repayment.create", s.handleTransactionRepaymentCreate)
    s.consumer.RegisterHandler("transaction.validate", s.handleTransactionValidate)
    s.consumer.RegisterHandler("transaction.get", s.handleTransactionGet)
    s.consumer.RegisterHandler("transaction.list", s.handleTransactionList)
    s.consumer.RegisterHandler("transaction.balance", s.handleTransactionBalance)
    
    // Organization operations
    s.consumer.RegisterHandler("organization.create", s.handleOrganizationCreate)
    s.consumer.RegisterHandler("organization.update", s.handleOrganizationUpdate)
    s.consumer.RegisterHandler("organization.get", s.handleOrganizationGet)
    s.consumer.RegisterHandler("organization.list", s.handleOrganizationList)
    s.consumer.RegisterHandler("organization.banks", s.handleOrganizationBanks)
    s.consumer.RegisterHandler("organization.business-types", s.handleOrganizationBusinessTypes)
}
```

### Message Handlers

```go
// Employee Operations
func (s *AbhiGatewayService) handleEmployeeCreate(ctx context.Context, message *Message) (*Response, error) {
    var req CreateEmployeeRequest
    if err := s.unmarshalPayload(message.Payload, &req); err != nil {
        return s.errorResponse(http.StatusBadRequest, err), nil
    }
    
    // Convert to Abhi employee model
    employee := s.convertToAbhiEmployee(&req)
    
    // Execute with circuit breaker
    result, err := s.executeWithCircuitBreaker("employee.create", func() (interface{}, error) {
        return s.abhiSDK.Employee.CreateSingle(ctx, employee)
    })
    
    if err != nil {
        s.metrics.RecordError("employee.create")
        return s.errorResponse(http.StatusInternalServerError, err), nil
    }
    
    // Cache the created employee
    go s.cacheEmployee(ctx, result.(*models.Employee))
    
    // Publish event
    go s.publishEvent("employee.created", &EmployeeCreatedEvent{
        EmployeeID: result.(*models.Employee).ID,
        RequestID:  req.RequestID,
        Timestamp:  time.Now(),
    })
    
    s.metrics.RecordSuccess("employee.create")
    return &Response{
        Data:       result,
        StatusCode: http.StatusCreated,
    }, nil
}

func (s *AbhiGatewayService) handleEmployeeList(ctx context.Context, message *Message) (*Response, error) {
    var req ListEmployeeRequest
    if err := s.unmarshalPayload(message.Payload, &req); err != nil {
        return s.errorResponse(http.StatusBadRequest, err), nil
    }
    
    // Check cache first
    cacheKey := s.buildEmployeeListCacheKey(&req)
    if cached := s.getFromCache(ctx, cacheKey); cached != nil {
        return &Response{
            Data:       cached,
            StatusCode: http.StatusOK,
        }, nil
    }
    
    // Build Abhi query options
    opts := &models.EmployeeListOptions{
        Page:       req.Page,
        Limit:      req.Limit,
        Department: req.Department,
        Search:     req.Search,
        Status:     req.Status,
    }
    
    result, err := s.executeWithCircuitBreaker("employee.list", func() (interface{}, error) {
        return s.abhiSDK.Employee.List(ctx, opts)
    })
    
    if err != nil {
        s.metrics.RecordError("employee.list")
        return s.errorResponse(http.StatusInternalServerError, err), nil
    }
    
    // Cache the result
    go s.cacheResult(ctx, cacheKey, result, 5*time.Minute)
    
    s.metrics.RecordSuccess("employee.list")
    return &Response{
        Data:       result,
        StatusCode: http.StatusOK,
    }, nil
}

// Transaction Operations
func (s *AbhiGatewayService) handleTransactionAdvanceCreate(ctx context.Context, message *Message) (*Response, error) {
    var req CreateAdvanceRequest
    if err := s.unmarshalPayload(message.Payload, &req); err != nil {
        return s.errorResponse(http.StatusBadRequest, err), nil
    }
    
    // First validate the transaction
    validationReq := models.TransactionValidationRequest{
        EmployeeID: req.EmployeeID,
        Amount:     req.Amount,
    }
    
    validation, err := s.executeWithCircuitBreaker("transaction.validate", func() (interface{}, error) {
        return s.abhiSDK.Transaction.ValidateEmployeeTransaction(ctx, validationReq)
    })
    
    if err != nil {
        return s.errorResponse(http.StatusInternalServerError, err), nil
    }
    
    validationResult := validation.(*models.TransactionValidationResponse)
    if !validationResult.IsValid {
        return &Response{
            Error: &Error{
                Code:    "VALIDATION_FAILED",
                Message: validationResult.Message,
            },
            StatusCode: http.StatusBadRequest,
        }, nil
    }
    
    // Create the advance transaction
    result, err := s.executeWithCircuitBreaker("transaction.advance.create", func() (interface{}, error) {
        return s.abhiSDK.Transaction.CreateAdvanceTransaction(
            ctx, 
            req.EmployeeID, 
            req.Amount, 
            req.Description,
        )
    })
    
    if err != nil {
        s.metrics.RecordError("transaction.advance.create")
        return s.errorResponse(http.StatusInternalServerError, err), nil
    }
    
    transaction := result.(*models.Transaction)
    
    // Publish transaction created event
    go s.publishEvent("transaction.advance.created", &TransactionCreatedEvent{
        TransactionID: transaction.ID,
        EmployeeID:    req.EmployeeID,
        Amount:        req.Amount,
        Type:          "advance",
        RequestID:     req.RequestID,
        Timestamp:     time.Now(),
    })
    
    s.metrics.RecordSuccess("transaction.advance.create")
    return &Response{
        Data:       transaction,
        StatusCode: http.StatusCreated,
    }, nil
}

// Organization Operations
func (s *AbhiGatewayService) handleOrganizationCreate(ctx context.Context, message *Message) (*Response, error) {
    var req CreateOrganizationRequest
    if err := s.unmarshalPayload(message.Payload, &req); err != nil {
        return s.errorResponse(http.StatusBadRequest, err), nil
    }
    
    // Convert to Abhi organization model
    orgReq := s.convertToAbhiOrganization(&req)
    
    result, err := s.executeWithCircuitBreaker("organization.create", func() (interface{}, error) {
        return s.abhiSDK.Organization.Create(ctx, orgReq)
    })
    
    if err != nil {
        s.metrics.RecordError("organization.create")
        return s.errorResponse(http.StatusInternalServerError, err), nil
    }
    
    // Cache organization data
    go s.cacheOrganization(ctx, result.(*models.CreateOrganizationResponse))
    
    // Publish event
    go s.publishEvent("organization.created", &OrganizationCreatedEvent{
        OrganizationID: result.(*models.CreateOrganizationResponse).Data.OrganizationID,
        RequestID:      req.RequestID,
        Timestamp:      time.Now(),
    })
    
    s.metrics.RecordSuccess("organization.create")
    return &Response{
        Data:       result,
        StatusCode: http.StatusCreated,
    }, nil
}
```

### Circuit Breaker Implementation

```go
// abhi-gateway-service/internal/resilience/circuit_breaker.go
func (s *AbhiGatewayService) executeWithCircuitBreaker(operation string, fn func() (interface{}, error)) (interface{}, error) {
    cb := s.circuitBreaker.GetBreaker(operation)
    
    result, err := cb.Execute(func() (interface{}, error) {
        start := time.Now()
        defer func() {
            duration := time.Since(start)
            s.metrics.RecordRequestDuration(operation, duration)
        }()
        
        return fn()
    })
    
    if err != nil {
        s.logger.WithFields(logrus.Fields{
            "operation": operation,
            "error":     err,
        }).Error("Circuit breaker operation failed")
        
        // Check if circuit breaker is open
        if cb.State() == circuitbreaker.StateOpen {
            return nil, fmt.Errorf("service unavailable: circuit breaker is open for operation %s", operation)
        }
        
        return nil, err
    }
    
    return result, nil
}
```

### Caching Strategy

```go
// abhi-gateway-service/internal/cache/cache.go
func (s *AbhiGatewayService) cacheEmployee(ctx context.Context, employee *models.Employee) {
    key := fmt.Sprintf("employee:%s", employee.ID)
    data, err := json.Marshal(employee)
    if err != nil {
        s.logger.Errorf("Failed to marshal employee for caching: %v", err)
        return
    }
    
    err = s.redis.SetEX(ctx, key, data, 30*time.Minute).Err()
    if err != nil {
        s.logger.Errorf("Failed to cache employee: %v", err)
    }
}

func (s *AbhiGatewayService) getEmployeeFromCache(ctx context.Context, employeeID string) *models.Employee {
    key := fmt.Sprintf("employee:%s", employeeID)
    data, err := s.redis.Get(ctx, key).Result()
    if err != nil {
        return nil
    }
    
    var employee models.Employee
    if err := json.Unmarshal([]byte(data), &employee); err != nil {
        s.logger.Errorf("Failed to unmarshal cached employee: %v", err)
        return nil
    }
    
    return &employee
}

func (s *AbhiGatewayService) buildEmployeeListCacheKey(req *ListEmployeeRequest) string {
    hasher := sha256.New()
    hasher.Write([]byte(fmt.Sprintf("employee_list:%d:%d:%s:%s:%s", 
        req.Page, req.Limit, req.Department, req.Search, req.Status)))
    return fmt.Sprintf("cache:%x", hasher.Sum(nil))
}
```

---

## Security Architecture

### Authentication & Authorization Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          AUTHENTICATION FLOW                                     │
└─────────────────────────────────────────────────────────────────────────────────┘

1. Client Login Request
   ├── Client ────────────────> API Gateway ────────────> Auth Service
   │                                │                           │
   │                                │                           ▼
   │                                │                    Validate Credentials
   │                                │                           │
   │                                │                           ▼
   │                                │                    Generate JWT Token
   │                                │                           │
   │                                ▼                           │
   │                         Store Session Info                 │
   │                        (Redis/In-Memory)                   │
   │                                │                           │
   │                                ◄───────────────────────────┘
   │                                │
   └────◄───────────────────────────┘
      JWT Token + User Info

2. Subsequent API Requests
   ├── Client ────────────────> API Gateway
   │   (with JWT Token)               │
   │                                  ▼
   │                          Validate JWT Token
   │                                  │
   │                                  ▼
   │                          Check Permissions
   │                                  │
   │                                  ▼
   │                          Route to Service
   │                                  │
   └────◄─────────────────────────────┘
      API Response

3. Abhi API Integration
   ├── Service ───────────────> Abhi Gateway ───────────> Abhi API
   │   (Internal Token)               │               (Request Signed)
   │                                  │                      │
   │                                  ▼                      │
   │                       Get/Refresh Abhi Token           │
   │                                  │                      │
   │                                  ▼                      │
   │                        Sign Request (HMAC-SHA256)      │
   │                                  │                      │
   │                                  ▼                      │
   │                           Apply Rate Limiting           │
   │                                  │                      │
   │                                  ├─────────────────────►│
   │                                  │                      │
   └────◄─────────────────────────────◄──────────────────────┘
      Abhi API Response
```

### JWT Token Management

```go
// auth-service/internal/auth/jwt_manager.go
type JWTManager struct {
    secretKey     []byte
    tokenExpiry   time.Duration
    refreshExpiry time.Duration
    redis         *redis.Client
}

type TokenClaims struct {
    UserID       string   `json:"user_id"`
    Username     string   `json:"username"`
    Role         string   `json:"role"`
    Permissions  []string `json:"permissions"`
    SessionID    string   `json:"session_id"`
    jwt.RegisteredClaims
}

func (jm *JWTManager) GenerateToken(userID, role string, permissions []string) (*TokenPair, error) {
    sessionID := uuid.New().String()
    
    // Access token (short-lived)
    accessClaims := &TokenClaims{
        UserID:      userID,
        Role:        role,
        Permissions: permissions,
        SessionID:   sessionID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.tokenExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "abhi-auth-service",
            Subject:   userID,
        },
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(jm.secretKey)
    if err != nil {
        return nil, err
    }
    
    // Refresh token (long-lived)
    refreshClaims := &TokenClaims{
        UserID:    userID,
        SessionID: sessionID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.refreshExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "abhi-auth-service",
            Subject:   userID,
        },
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(jm.secretKey)
    if err != nil {
        return nil, err
    }
    
    // Store session in Redis
    sessionData := &SessionData{
        UserID:       userID,
        Role:         role,
        Permissions:  permissions,
        CreatedAt:    time.Now(),
        LastActivity: time.Now(),
        Active:       true,
    }
    
    err = jm.storeSession(sessionID, sessionData)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresIn:    int64(jm.tokenExpiry.Seconds()),
        TokenType:    "Bearer",
    }, nil
}

func (jm *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jm.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
        // Validate session in Redis
        sessionData, err := jm.getSession(claims.SessionID)
        if err != nil || !sessionData.Active {
            return nil, fmt.Errorf("invalid session")
        }
        
        // Update last activity
        go jm.updateLastActivity(claims.SessionID)
        
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

### Request Signing & Encryption

```go
// shared/security/request_signing.go
type RequestSecurityManager struct {
    signingSecret    string
    encryptionKey    []byte
    rateLimitManager *RateLimitManager
}

func (rsm *RequestSecurityManager) SignRequest(req *http.Request, body []byte) error {
    // Generate timestamp
    timestamp := time.Now().Unix()
    req.Header.Set("X-Timestamp", strconv.FormatInt(timestamp, 10))
    
    // Create canonical string
    canonicalString := rsm.createCanonicalString(req, body, timestamp)
    
    // Generate HMAC-SHA256 signature
    h := hmac.New(sha256.New, []byte(rsm.signingSecret))
    h.Write([]byte(canonicalString))
    signature := hex.EncodeToString(h.Sum(nil))
    
    req.Header.Set("X-Signature", signature)
    return nil
}

func (rsm *RequestSecurityManager) VerifySignature(req *http.Request, body []byte) bool {
    timestampStr := req.Header.Get("X-Timestamp")
    signature := req.Header.Get("X-Signature")
    
    if timestampStr == "" || signature == "" {
        return false
    }
    
    timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return false
    }
    
    // Check timestamp (5 minute window)
    now := time.Now().Unix()
    if abs(now-timestamp) > 300 {
        return false
    }
    
    // Verify signature
    canonicalString := rsm.createCanonicalString(req, body, timestamp)
    h := hmac.New(sha256.New, []byte(rsm.signingSecret))
    h.Write([]byte(canonicalString))
    expectedSignature := hex.EncodeToString(h.Sum(nil))
    
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (rsm *RequestSecurityManager) createCanonicalString(req *http.Request, body []byte, timestamp int64) string {
    var parts []string
    
    // HTTP method
    parts = append(parts, req.Method)
    
    // Path
    parts = append(parts, req.URL.Path)
    
    // Query parameters (sorted)
    if req.URL.RawQuery != "" {
        params := strings.Split(req.URL.RawQuery, "&")
        sort.Strings(params)
        parts = append(parts, strings.Join(params, "&"))
    } else {
        parts = append(parts, "")
    }
    
    // Headers (specific headers only, sorted)
    var headerParts []string
    headersToSign := []string{"authorization", "content-type", "x-timestamp"}
    
    for _, headerName := range headersToSign {
        value := req.Header.Get(headerName)
        if value != "" {
            headerParts = append(headerParts, fmt.Sprintf("%s:%s", headerName, strings.TrimSpace(value)))
        }
    }
    parts = append(parts, strings.Join(headerParts, "\n"))
    
    // Body hash
    bodyHash := sha256.Sum256(body)
    parts = append(parts, hex.EncodeToString(bodyHash[:]))
    
    // Timestamp
    parts = append(parts, strconv.FormatInt(timestamp, 10))
    
    return strings.Join(parts, "\n")
}

func abs(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}
```

### Rate Limiting Implementation

```go
// shared/security/rate_limiting.go
type RateLimitManager struct {
    redis  *redis.Client
    config *RateLimitConfig
}

type RateLimitConfig struct {
    DefaultRPS    int           `json:"default_rps"`
    DefaultBurst  int           `json:"default_burst"`
    WindowSize    time.Duration `json:"window_size"`
    Rules         []RateLimitRule `json:"rules"`
}

type RateLimitRule struct {
    Pattern string `json:"pattern"` // URL pattern or user role
    RPS     int    `json:"rps"`
    Burst   int    `json:"burst"`
}

func (rlm *RateLimitManager) CheckRateLimit(key string, rps, burst int) (bool, int, error) {
    now := time.Now()
    windowKey := fmt.Sprintf("rate_limit:%s:%d", key, now.Unix()/60) // 1-minute window
    
    pipe := rlm.redis.Pipeline()
    
    // Get current count
    countCmd := pipe.Get(context.Background(), windowKey)
    
    // Increment count
    incrCmd := pipe.Incr(context.Background(), windowKey)
    
    // Set expiry
    pipe.Expire(context.Background(), windowKey, time.Minute)
    
    _, err := pipe.Exec(context.Background())
    if err != nil && err != redis.Nil {
        return false, 0, err
    }
    
    currentCount := int(incrCmd.Val())
    
    // Check if within limits
    if currentCount > rps {
        remaining := 0
        if currentCount < burst {
            remaining = burst - currentCount
        }
        return false, remaining, nil
    }
    
    remaining := rps - currentCount
    if remaining < 0 {
        remaining = 0
    }
    
    return true, remaining, nil
}

func (rlm *RateLimitManager) GetRateLimitRule(path, userRole string) (int, int) {
    for _, rule := range rlm.config.Rules {
        if matched, _ := filepath.Match(rule.Pattern, path); matched {
            return rule.RPS, rule.Burst
        }
        if rule.Pattern == userRole {
            return rule.RPS, rule.Burst
        }
    }
    return rlm.config.DefaultRPS, rlm.config.DefaultBurst
}
```

---

## Data Flow Patterns

### 1. Employee Creation Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        EMPLOYEE CREATION FLOW                                    │
└─────────────────────────────────────────────────────────────────────────────────┘

Step 1: Request Initiation
┌─────────────┐    POST /employees    ┌─────────────┐
│   Client    │ ────────────────────> │ API Gateway │
└─────────────┘    (Employee Data)    └─────────────┘
                                             │
                                             ▼
Step 2: Authentication & Validation          
                                      ┌─────────────┐
                                      │ Middleware  │
                                      │ Stack       │
                                      │ • Auth      │
                                      │ • Rate Limit│
                                      │ • Validation│
                                      └─────────────┘
                                             │
                                             ▼
Step 3: Message Publishing
                                      ┌─────────────┐
                                      │  RabbitMQ   │
                                      │  Publisher  │
                                      └─────────────┘
                                             │
                                             ▼
                                   employee.create queue
                                             │
                                             ▼
Step 4: Service Processing
                                      ┌─────────────┐
                                      │  Employee   │
                                      │  Service    │
                                      │ Consumer    │
                                      └─────────────┘
                                             │
                                             ▼
Step 5: Business Logic & Validation
                                      ┌─────────────┐
                                      │ Business    │
                                      │ Logic       │
                                      │ • Validate  │
                                      │ • Transform │
                                      │ • Enrich    │
                                      └─────────────┘
                                             │
                                             ▼
Step 6: Abhi API Integration
                                      ┌─────────────┐    abhi.employee.create    ┌─────────────┐
                                      │  Employee   │ ─────────────────────────> │    Abhi     │
                                      │  Service    │          queue             │  Gateway    │
                                      └─────────────┘                           └─────────────┘
                                             │                                          │
                                             │                                          ▼
                                             │                                   ┌─────────────┐
                                             │                                   │ Abhi SDK    │
                                             │                                   │ • Security  │
                                             │                                   │ • Circuit   │
                                             │                                   │   Breaker   │
                                             │                                   │ • Retry     │
                                             │                                   └─────────────┘
                                             │                                          │
                                             │                                          ▼
                                             │                                   ┌─────────────┐
                                             │                                   │ Abhi Open   │
                                             │                                   │ API         │
                                             │                                   │ (External)  │
                                             │                                   └─────────────┘
                                             │                                          │
                                             ◄──────────────────────────────────────────┘
                                             │                Response
                                             ▼
Step 7: Event Publishing & Caching
                                      ┌─────────────┐
                                      │ Event Bus   │
                                      │ • employee  │
                                      │   .created  │
                                      │ • Cache     │
                                      │   Update    │
                                      └─────────────┘
                                             │
                                             ▼
Step 8: Response
                                      ┌─────────────┐
                                      │ Response    │
                                      │ Queue       │
                                      └─────────────┘
                                             │
                                             ▼
                                      ┌─────────────┐
                                      │ API Gateway │
                                      │ Response    │
                                      │ Handler     │
                                      └─────────────┘
                                             │
                                             ▼
┌─────────────┐      200 Created      ┌─────────────┐
│   Client    │ ◄──────────────────── │ API Gateway │
└─────────────┘    (Employee Data)    └─────────────┘
```

### 2. Transaction Processing Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                       TRANSACTION PROCESSING FLOW                                │
└─────────────────────────────────────────────────────────────────────────────────┘

Step 1: Advance Request
┌─────────────┐ POST /transactions/advance ┌─────────────┐
│   Client    │ ─────────────────────────> │ API Gateway │
└─────────────┘     (Transaction Data)     └─────────────┘
                                                  │
                                                  ▼
Step 2: Pre-validation & Routing
                                           ┌─────────────┐
                                           │ Gateway     │
                                           │ Validation  │
                                           │ • Amount > 0│
                                           │ • Employee  │
                                           │   exists    │
                                           └─────────────┘
                                                  │
                                                  ▼
                                           ┌─────────────┐
                                           │  RabbitMQ   │
                                           │  High       │
                                           │  Priority   │
                                           │  Queue      │
                                           └─────────────┘
                                                  │
                                                  ▼
Step 3: Transaction Service Processing
                                           ┌─────────────┐
                                           │Transaction  │
                                           │Service      │
                                           │Consumer     │
                                           └─────────────┘
                                                  │
                                                  ▼
Step 4: Business Rules Engine
                                           ┌─────────────┐
                                           │ Business    │
                                           │ Rules       │
                                           │ • Eligibility│
                                           │ • Limits    │
                                           │ • Policy    │
                                           └─────────────┘
                                                  │
                                                  ▼
Step 5: Abhi Validation (Parallel Calls)
     ┌─────────────────────┬──────────────────────┬─────────────────────┐
     │                     │                      │                     │
     ▼                     ▼                      ▼                     ▼
┌─────────────┐    ┌─────────────┐      ┌─────────────┐    ┌─────────────┐
│   Validate  │    │Get Employee │      │Get Balance  │    │Check Limits │
│ Transaction │    │   Details   │      │Information  │    │& Policies   │
└─────────────┘    └─────────────┘      └─────────────┘    └─────────────┘
     │                     │                      │                     │
     └─────────────────────┴──────────────────────┴─────────────────────┘
                                    │
                                    ▼
Step 6: Create Transaction (If Valid)
                             ┌─────────────┐
                             │ Abhi SDK    │
                             │ Transaction │
                             │ Creation    │
                             └─────────────┘
                                    │
                                    ▼
Step 7: Event Broadcasting
                             ┌─────────────┐
                             │Event        │
                             │Broadcasting │
                             │• Audit Log  │
                             │• Notification│
                             │• Analytics  │
                             └─────────────┘
                                    │
                                    ▼ 
Step 8: Response with Details
┌─────────────┐      201 Created       ┌─────────────┐
│   Client    │ ◄─────────────────────── │ API Gateway │
└─────────────┘   (Transaction Details) └─────────────┘
                    • Transaction ID
                    • Status
                    • Amount  
                    • Fees
                    • Due Date
```

### 3. Real-time Balance Update Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        REAL-TIME BALANCE UPDATE FLOW                             │
└─────────────────────────────────────────────────────────────────────────────────┘

Trigger: Transaction Status Change
                             ┌─────────────┐
                             │ Abhi API    │
                             │ Webhook     │
                             │ (Optional)  │
                             └─────────────┘
                                    │
                                    ▼
                             ┌─────────────┐
                             │ Webhook     │
                             │ Handler     │
                             │ Service     │
                             └─────────────┘
                                    │
                                    ▼
Alternative: Polling Strategy
                             ┌─────────────┐
                             │ Scheduled   │
                             │ Job         │
                             │ (Cron)      │
                             └─────────────┘
                                    │
                                    ▼
                             ┌─────────────┐
                             │ Balance     │
                             │ Sync        │
                             │ Service     │
                             └─────────────┘
                                    │
                                    ▼
Parallel Balance Updates
     ┌─────────────────────┬──────────────────────┐
     │                     │                      │
     ▼                     ▼                      ▼
┌─────────────┐    ┌─────────────┐      ┌─────────────┐
│Fetch Latest │    │Update Cache │      │Broadcast    │
│from Abhi    │    │(Redis)      │      │Event        │
│API          │    │             │      │             │
└─────────────┘    └─────────────┘      └─────────────┘
     │                     │                      │
     └─────────────────────┼──────────────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │ Update      │
                    │ Local DB    │
                    │ (Optional)  │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │ Notify      │
                    │ Subscribers │
                    │ • Frontend  │
                    │ • Mobile    │
                    │ • Analytics │
                    └─────────────┘
```

---

## Implementation Strategy

### Phase 1: Foundation Setup (Week 1-2)

#### Infrastructure Components
```bash
# 1. RabbitMQ Setup
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin123 \
  rabbitmq:3-management

# 2. Redis Setup  
docker run -d --name redis \
  -p 6379:6379 \
  redis:7-alpine

# 3. PostgreSQL Setup
docker run -d --name postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=admin123 \
  -e POSTGRES_DB=abhi_microservices \
  postgres:15-alpine
```

#### Project Structure
```
abhi-microservices/
├── api-gateway/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── server/
│   ├── pkg/
│   └── deployments/
├── employee-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── domain/
│   │   ├── service/
│   │   ├── repository/
│   │   └── handlers/
│   └── pkg/
├── transaction-service/
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── organization-service/
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── abhi-gateway-service/
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── shared/
│   ├── messaging/
│   ├── security/
│   ├── monitoring/
│   ├── database/
│   └── config/
└── scripts/
    ├── setup/
    ├── migration/
    └── deployment/
```

### Phase 2: Core Services Development (Week 3-6)

#### Development Priority
1. **Shared Libraries** (Week 3)
   - RabbitMQ messaging framework
   - Security components (JWT, encryption, signing)
   - Configuration management
   - Database utilities

2. **Abhi Gateway Service** (Week 4)
   - Abhi SDK integration with enhanced security
   - Message processing handlers
   - Circuit breaker implementation
   - Caching layer

3. **API Gateway** (Week 5)
   - Route configuration
   - Middleware stack implementation
   - Request/response handling
   - Load balancing

4. **Business Services** (Week 6)
   - Employee Service
   - Transaction Service  
   - Organization Service

### Phase 3: Integration & Testing (Week 7-8)

#### Integration Testing Strategy
```go
// integration-tests/employee_test.go
func TestEmployeeCreationFlow(t *testing.T) {
    // Setup test environment
    testEnv := setupTestEnvironment(t)
    defer testEnv.Cleanup()
    
    // Test API Gateway -> Employee Service -> Abhi Gateway flow
    client := testEnv.APIGatewayClient()
    
    employee := &CreateEmployeeRequest{
        FirstName:    "John",
        LastName:     "Doe",
        Email:        "john.doe@test.com",
        Department:   "Engineering",
        // ... other fields
    }
    
    // Send request
    resp, err := client.CreateEmployee(context.Background(), employee)
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, resp.StatusCode)
    
    // Verify employee was created in Abhi API (mock)
    testEnv.AbhiMock.AssertEmployeeCreated(employee.Email)
    
    // Verify event was published
    testEnv.EventStore.AssertEventPublished("employee.created")
    
    // Verify caching
    cachedEmployee := testEnv.Cache.GetEmployee(resp.Data.ID)
    require.NotNil(t, cachedEmployee)
}

func TestTransactionValidationAndCreation(t *testing.T) {
    testEnv := setupTestEnvironment(t)
    defer testEnv.Cleanup()
    
    // Create employee first
    employee := testEnv.CreateTestEmployee()
    
    // Test transaction validation
    client := testEnv.APIGatewayClient()
    
    transaction := &CreateAdvanceRequest{
        EmployeeID:  employee.ID,
        Amount:      1000.0,
        Description: "Medical emergency",
    }
    
    // Validate first
    validation, err := client.ValidateTransaction(context.Background(), &ValidateTransactionRequest{
        EmployeeID: employee.ID,
        Amount:     transaction.Amount,
    })
    require.NoError(t, err)
    require.True(t, validation.IsValid)
    
    // Create transaction
    resp, err := client.CreateAdvanceTransaction(context.Background(), transaction)
    require.NoError(t, err)
    require.Equal(t, http.StatusCreated, resp.StatusCode)
    
    // Verify transaction creation
    testEnv.AbhiMock.AssertTransactionCreated(employee.ID, transaction.Amount)
}
```

### Phase 4: Security Implementation (Week 9)

#### Security Checklist
- [ ] JWT token management with refresh mechanism
- [ ] Request signing (HMAC-SHA256) on all Abhi API calls
- [ ] Credential encryption (AES-GCM) for sensitive data
- [ ] Rate limiting on all endpoints
- [ ] Input validation on all requests
- [ ] HTTPS/TLS termination at API Gateway
- [ ] Secrets management with Kubernetes secrets
- [ ] Database connection encryption
- [ ] Inter-service communication security

### Phase 5: Production Deployment (Week 10-12)

#### Kubernetes Deployment Configuration

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: abhi-microservices
  labels:
    name: abhi-microservices

---
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: abhi-credentials
  namespace: abhi-microservices
type: Opaque
data:
  username: <base64-encoded-username>
  password: <base64-encoded-password>
  signing-secret: <base64-encoded-signing-secret>
  encryption-password: <base64-encoded-encryption-password>

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: abhi-config
  namespace: abhi-microservices
data:
  environment: "uat"
  rabbitmq-url: "amqp://admin:admin123@rabbitmq-service:5672/"
  redis-url: "redis://redis-service:6379"
  rate-limit-rps: "10"
  rate-limit-burst: "20"
  circuit-breaker-threshold: "5"
  circuit-breaker-timeout: "60"

---
# k8s/api-gateway-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: abhi-microservices
  labels:
    app: api-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: your-registry/api-gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: RABBITMQ_URL
          valueFrom:
            configMapKeyRef:
              name: abhi-config
              key: rabbitmq-url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: abhi-config
              key: redis-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: abhi-credentials
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"  
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
# k8s/abhi-gateway-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: abhi-gateway
  namespace: abhi-microservices
  labels:
    app: abhi-gateway
spec:
  replicas: 2
  selector:
    matchLabels:
      app: abhi-gateway
  template:
    metadata:
      labels:
        app: abhi-gateway
    spec:
      containers:
      - name: abhi-gateway
        image: your-registry/abhi-gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: ABHI_ENV
          valueFrom:
            configMapKeyRef:
              name: abhi-config
              key: environment
        - name: ABHI_USERNAME
          valueFrom:
            secretKeyRef:
              name: abhi-credentials
              key: username
        - name: ABHI_PASSWORD
          valueFrom:
            secretKeyRef:
              name: abhi-credentials
              key: password
        - name: ABHI_SIGNING_SECRET
          valueFrom:
            secretKeyRef:
              name: abhi-credentials
              key: signing-secret
        - name: ABHI_ENCRYPTION_PASS
          valueFrom:
            secretKeyRef:
              name: abhi-credentials
              key: encryption-password
        - name: RABBITMQ_URL
          valueFrom:
            configMapKeyRef:
              name: abhi-config
              key: rabbitmq-url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: abhi-config
              key: redis-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
          limits:
            memory: "512Mi"
            cpu: "400m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

---

## Monitoring & Observability

### Prometheus Metrics Configuration

```yaml
# monitoring/prometheus-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: abhi-microservices
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      evaluation_interval: 15s
    
    rule_files:
      - "/etc/prometheus/rules/*.yml"
    
    scrape_configs:
      - job_name: 'api-gateway'
        kubernetes_sd_configs:
          - role: endpoints
            namespaces:
              names:
                - abhi-microservices
        relabel_configs:
          - source_labels: [__meta_kubernetes_service_name]
            action: keep
            regex: api-gateway
          - source_labels: [__meta_kubernetes_endpoint_port_name]
            action: keep
            regex: metrics
    
      - job_name: 'abhi-gateway'
        kubernetes_sd_configs:
          - role: endpoints
            namespaces:
              names:
                - abhi-microservices
        relabel_configs:
          - source_labels: [__meta_kubernetes_service_name]
            action: keep
            regex: abhi-gateway
    
      - job_name: 'employee-service'
        kubernetes_sd_configs:
          - role: endpoints
            namespaces:
              names:
                - abhi-microservices
        relabel_configs:
          - source_labels: [__meta_kubernetes_service_name]
            action: keep
            regex: employee-service
    
      - job_name: 'transaction-service'
        kubernetes_sd_configs:
          - role: endpoints
            namespaces:
              names:
                - abhi-microservices
        relabel_configs:
          - source_labels: [__meta_kubernetes_service_name]
            action: keep
            regex: transaction-service

---
# monitoring/alert-rules.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-rules
  namespace: abhi-microservices
data:
  abhi-alerts.yml: |
    groups:
    - name: abhi-microservices
      rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is above 10% for {{ $labels.service }}"
    
      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is above 2s for {{ $labels.service }}"
    
      - alert: AbhiAPIDown
        expr: up{job="abhi-gateway"} == 0
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Abhi Gateway service is down"
          description: "Abhi Gateway service has been down for more than 2 minutes"
    
      - alert: RabbitMQQueueBacklog
        expr: rabbitmq_queue_messages > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "RabbitMQ queue backlog"
          description: "Queue {{ $labels.queue }} has more than 1000 messages"
    
      - alert: HighMemoryUsage
        expr: container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Container {{ $labels.container }} memory usage is above 90%"
```

### Grafana Dashboard Configuration

```json
{
  "dashboard": {
    "id": null,
    "title": "Abhi Microservices Dashboard",
    "tags": ["abhi", "microservices"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{service}} - {{method}} {{endpoint}}"
          }
        ],
        "yAxes": [
          {
            "label": "Requests/sec"
          }
        ]
      },
      {
        "id": 2,
        "title": "Error Rate",
        "type": "graph", 
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "{{service}} - Errors"
          }
        ]
      },
      {
        "id": 3,
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{service}} - 95th percentile"
          },
          {
            "expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{service}} - Median"
          }
        ]
      },
      {
        "id": 4,
        "title": "Abhi API Metrics",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(abhi_requests_total[5m])",
            "legendFormat": "{{method}} - {{status}}"
          }
        ]
      },
      {
        "id": 5,
        "title": "RabbitMQ Queue Depth",
        "type": "graph",
        "targets": [
          {
            "expr": "rabbitmq_queue_messages",
            "legendFormat": "{{queue}}"
          }
        ]
      },
      {
        "id": 6,
        "title": "Circuit Breaker Status",
        "type": "stat",
        "targets": [
          {
            "expr": "circuit_breaker_state",
            "legendFormat": "{{name}}"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "mappings": [
              {
                "options": {
                  "0": {
                    "text": "CLOSED",
                    "color": "green"
                  },
                  "1": {
                    "text": "HALF_OPEN", 
                    "color": "yellow"
                  },
                  "2": {
                    "text": "OPEN",
                    "color": "red"
                  }
                }
              }
            ]
          }
        }
      },
      {
        "id": 7,
        "title": "Rate Limiting",
        "type": "graph",
        "targets": [
          {
            "expr": "rate_limit_hits_total",
            "legendFormat": "{{service}} - Rate Limited"
          }
        ]
      }
    ],
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "refresh": "30s"
  }
}
```

### Centralized Logging with ELK Stack

```yaml
# monitoring/elasticsearch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: elasticsearch
  namespace: abhi-microservices
spec:
  replicas: 1
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      containers:
      - name: elasticsearch
        image: docker.elastic.co/elasticsearch/elasticsearch:8.5.0
        env:
        - name: discovery.type
          value: single-node
        - name: ES_JAVA_OPTS
          value: "-Xms512m -Xmx512m"
        - name: xpack.security.enabled
          value: "false"
        ports:
        - containerPort: 9200
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"

---
# monitoring/logstash.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: logstash-config
  namespace: abhi-microservices
data:
  logstash.conf: |
    input {
      beats {
        port => 5044
      }
    }
    
    filter {
      if [fields][service] {
        mutate {
          add_field => { "service" => "%{[fields][service]}" }
        }
      }
      
      # Parse JSON logs
      if [message] =~ /^\{.*\}$/ {
        json {
          source => "message"
        }
      }
      
      # Parse timestamp
      if [timestamp] {
        date {
          match => [ "timestamp", "ISO8601" ]
        }
      }
      
      # Add Abhi-specific fields
      if [service] == "abhi-gateway" {
        if [abhi_request_id] {
          mutate {
            add_field => { "abhi_request" => "true" }
          }
        }
      }
    }
    
    output {
      elasticsearch {
        hosts => ["elasticsearch-service:9200"]
        index => "abhi-microservices-%{+YYYY.MM.dd}"
      }
      
      stdout {
        codec => rubydebug
      }
    }

---
# monitoring/filebeat.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: filebeat-config
  namespace: abhi-microservices
data:
  filebeat.yml: |
    filebeat.inputs:
    - type: container
      paths:
        - /var/log/containers/*abhi*.log
      processors:
        - add_kubernetes_metadata:
            host: ${NODE_NAME}
            matchers:
            - logs_path:
                logs_path: "/var/log/containers/"
    
    output.logstash:
      hosts: ["logstash-service:5044"]
    
    logging.level: info
```

---

## Best Practices & Guidelines

### Code Organization Standards

#### 1. Project Structure
```
service-name/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/                     # Private code
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── domain/
│   │   ├── entities/            # Business entities
│   │   ├── repositories/        # Repository interfaces
│   │   └── services/            # Business logic
│   ├── infrastructure/
│   │   ├── database/            # Database implementations
│   │   ├── messaging/           # Message queue implementations
│   │   └── external/            # External API clients
│   ├── handlers/
│   │   ├── http/                # HTTP handlers
│   │   └── messaging/           # Message handlers
│   └── middleware/
│       └── auth.go              # Authentication middleware
├── pkg/                         # Public code
│   └── api/
│       └── models.go            # API models
├── deployments/
│   ├── docker/
│   │   └── Dockerfile
│   └── k8s/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── tests/
│   ├── integration/
│   └── unit/
├── go.mod
├── go.sum
└── README.md
```

#### 2. Error Handling Patterns

```go
// shared/errors/errors.go
type ServiceError struct {
    Code       string                 `json:"code"`
    Message    string                 `json:"message"`
    Details    map[string]interface{} `json:"details,omitempty"`
    Cause      error                  `json:"-"`
    Service    string                 `json:"service"`
    Timestamp  time.Time              `json:"timestamp"`
    TraceID    string                 `json:"trace_id,omitempty"`
}

func (e *ServiceError) Error() string {
    return fmt.Sprintf("[%s] %s: %s", e.Service, e.Code, e.Message)
}

func NewServiceError(service, code, message string) *ServiceError {
    return &ServiceError{
        Code:      code,
        Message:   message,
        Service:   service,
        Timestamp: time.Now(),
    }
}

func (e *ServiceError) WithDetails(details map[string]interface{}) *ServiceError {
    e.Details = details
    return e
}

func (e *ServiceError) WithCause(err error) *ServiceError {
    e.Cause = err
    return e
}

func (e *ServiceError) WithTraceID(traceID string) *ServiceError {
    e.TraceID = traceID
    return e
}

// Common error codes
const (
    ErrCodeValidation    = "VALIDATION_ERROR"
    ErrCodeNotFound      = "NOT_FOUND"
    ErrCodeUnauthorized  = "UNAUTHORIZED"
    ErrCodeForbidden     = "FORBIDDEN"
    ErrCodeInternal      = "INTERNAL_ERROR"
    ErrCodeExternalAPI   = "EXTERNAL_API_ERROR"
    ErrCodeRateLimit     = "RATE_LIMIT_EXCEEDED"
    ErrCodeCircuitOpen   = "CIRCUIT_BREAKER_OPEN"
)
```

#### 3. Configuration Management

```go
// shared/config/config.go
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    RabbitMQ RabbitMQConfig `json:"rabbitmq"`
    Redis    RedisConfig    `json:"redis"`
    Security SecurityConfig `json:"security"`
    Logging  LoggingConfig  `json:"logging"`
    Metrics  MetricsConfig  `json:"metrics"`
}

type ServerConfig struct {
    Host            string        `json:"host" env:"SERVER_HOST" default:"0.0.0.0"`
    Port            int           `json:"port" env:"SERVER_PORT" default:"8080"`
    ReadTimeout     time.Duration `json:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30s"`
    WriteTimeout    time.Duration `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30s"`
    ShutdownTimeout time.Duration `json:"shutdown_timeout" env:"SERVER_SHUTDOWN_TIMEOUT" default:"30s"`
}

func LoadConfig() (*Config, error) {
    config := &Config{}
    
    // Load from environment variables
    if err := env.Parse(config); err != nil {
        return nil, fmt.Errorf("failed to parse environment variables: %w", err)
    }
    
    // Load from config file if exists
    configFile := os.Getenv("CONFIG_FILE")
    if configFile != "" {
        if err := loadConfigFromFile(configFile, config); err != nil {
            return nil, fmt.Errorf("failed to load config from file: %w", err)
        }
    }
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("configuration validation failed: %w", err)
    }
    
    return config, nil
}
```

#### 4. Structured Logging

```go
// shared/logging/logger.go
type Logger struct {
    *logrus.Entry
    service string
}

func NewLogger(service string) *Logger {
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "level",
            logrus.FieldKeyMsg:   "message",
        },
    })
    
    return &Logger{
        Entry:   log.WithField("service", service),
        service: service,
    }
}

func (l *Logger) WithRequestID(requestID string) *Logger {
    return &Logger{
        Entry:   l.Entry.WithField("request_id", requestID),
        service: l.service,
    }
}

func (l *Logger) WithUserID(userID string) *Logger {
    return &Logger{
        Entry:   l.Entry.WithField("user_id", userID),
        service: l.service,
    }
}

func (l *Logger) WithAbhiRequestID(abhiRequestID string) *Logger {
    return &Logger{
        Entry:   l.Entry.WithField("abhi_request_id", abhiRequestID),
        service: l.service,
    }
}

func (l *Logger) LogServiceError(err *ServiceError) {
    fields := logrus.Fields{
        "error_code": err.Code,
        "service":    err.Service,
    }
    
    if err.TraceID != "" {
        fields["trace_id"] = err.TraceID
    }
    
    if err.Details != nil {
        fields["error_details"] = err.Details
    }
    
    if err.Cause != nil {
        fields["cause"] = err.Cause.Error()
    }
    
    l.WithFields(fields).Error(err.Message)
}
```

### Security Guidelines

#### 1. Authentication Flow
- Use JWT tokens with short expiry (15 minutes)
- Implement refresh token rotation
- Store session data in Redis with TTL
- Validate tokens on every request
- Implement proper logout (token invalidation)

#### 2. Authorization Patterns
```go
// shared/auth/rbac.go
type Permission string

const (
    PermissionEmployeeRead   Permission = "employee:read"
    PermissionEmployeeWrite  Permission = "employee:write"
    PermissionEmployeeDelete Permission = "employee:delete"
    
    PermissionTransactionRead   Permission = "transaction:read"
    PermissionTransactionWrite  Permission = "transaction:write"
    PermissionTransactionApprove Permission = "transaction:approve"
    
    PermissionOrganizationRead  Permission = "organization:read"
    PermissionOrganizationWrite Permission = "organization:write"
)

type Role string

const (
    RoleEmployee    Role = "employee"
    RoleManager     Role = "manager"
    RoleAdmin       Role = "admin"
    RoleSuperAdmin  Role = "super_admin"
    RoleAPIClient   Role = "api_client"
)

var RolePermissions = map[Role][]Permission{
    RoleEmployee: {
        PermissionEmployeeRead,
        PermissionTransactionRead,
    },
    RoleManager: {
        PermissionEmployeeRead,
        PermissionEmployeeWrite,
        PermissionTransactionRead,
        PermissionTransactionWrite,
        PermissionTransactionApprove,
    },
    RoleAdmin: {
        PermissionEmployeeRead,
        PermissionEmployeeWrite,
        PermissionEmployeeDelete,
        PermissionTransactionRead,
        PermissionTransactionWrite,
        PermissionTransactionApprove,
        PermissionOrganizationRead,
        PermissionOrganizationWrite,
    },
    // ... other roles
}

func (r Role) HasPermission(permission Permission) bool {
    permissions, exists := RolePermissions[r]
    if !exists {
        return false
    }
    
    for _, p := range permissions {
        if p == permission {
            return true
        }
    }
    return false
}
```

#### 3. Input Validation
```go
// shared/validation/validator.go
type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    v := validator.New()
    
    // Register custom validators
    v.RegisterValidation("emirates_id", validateEmiratesID)
    v.RegisterValidation("phone_uae", validateUAEPhone)
    
    return &Validator{validate: v}
}

func validateEmiratesID(fl validator.FieldLevel) bool {
    emiratesID := fl.Field().String()
    
    // UAE Emirates ID format: 784-YYYY-XXXXXXX-X (15 digits total)
    pattern := `^784-\d{4}-\d{7}-\d{1}$`
    matched, _ := regexp.MatchString(pattern, emiratesID)
    return matched
}

func validateUAEPhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    
    // UAE phone format: +971XXXXXXXXX
    pattern := `^\+971[5][0-9]{8}$`
    matched, _ := regexp.MatchString(pattern, phone)
    return matched
}
```

### Performance Guidelines

#### 1. Database Optimization
- Use connection pooling
- Implement proper indexing strategy
- Use read replicas for queries
- Implement database-level caching
- Use database migrations for schema changes

#### 2. Caching Strategy
```go
// shared/cache/cache.go
type CacheManager struct {
    redis *redis.Client
}

type CacheOptions struct {
    TTL        time.Duration
    Namespace  string
    Serialize  func(interface{}) ([]byte, error)
    Deserialize func([]byte, interface{}) error
}

func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, opts *CacheOptions) error {
    if opts == nil {
        opts = &CacheOptions{
            TTL: 5 * time.Minute,
            Serialize: json.Marshal,
        }
    }
    
    fullKey := fmt.Sprintf("%s:%s", opts.Namespace, key)
    
    data, err := opts.Serialize(value)
    if err != nil {
        return err
    }
    
    return cm.redis.SetEX(ctx, fullKey, data, opts.TTL).Err()
}

func (cm *CacheManager) Get(ctx context.Context, key string, dest interface{}, opts *CacheOptions) error {
    if opts == nil {
        opts = &CacheOptions{
            Deserialize: json.Unmarshal,
        }
    }
    
    fullKey := fmt.Sprintf("%s:%s", opts.Namespace, key)
    
    data, err := cm.redis.Get(ctx, fullKey).Result()
    if err != nil {
        return err
    }
    
    return opts.Deserialize([]byte(data), dest)
}

// Cache patterns for different data types
func (cm *CacheManager) CacheEmployee(ctx context.Context, employee *models.Employee) error {
    return cm.Set(ctx, employee.ID, employee, &CacheOptions{
        TTL:       30 * time.Minute,
        Namespace: "employee",
    })
}

func (cm *CacheManager) CacheTransactionList(ctx context.Context, employeeID string, page int, transactions []models.Transaction) error {
    key := fmt.Sprintf("%s:page:%d", employeeID, page)
    return cm.Set(ctx, key, transactions, &CacheOptions{
        TTL:       5 * time.Minute,
        Namespace: "transaction_list",
    })
}
```

#### 3. Message Queue Optimization
- Use message priorities for critical operations
- Implement dead letter queues for failed messages
- Use message batching for bulk operations
- Configure appropriate prefetch counts
- Monitor queue depths and processing times

### Testing Strategy

#### 1. Unit Testing
```go
// employee-service/internal/service/employee_service_test.go
func TestEmployeeService_CreateEmployee(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateEmployeeRequest
        setup   func(*testing.T) (*EmployeeService, *mocks.AbhiClient, *mocks.Repository)
        assert  func(*testing.T, *models.Employee, error)
    }{
        {
            name: "successful creation",
            input: &CreateEmployeeRequest{
                FirstName: "John",
                LastName:  "Doe",
                Email:     "john.doe@test.com",
                // ... other fields
            },
            setup: func(t *testing.T) (*EmployeeService, *mocks.AbhiClient, *mocks.Repository) {
                abhiClient := mocks.NewAbhiClient(t)
                repo := mocks.NewRepository(t)
                
                abhiClient.EXPECT().CreateEmployee(mock.Anything, mock.Anything).
                    Return(&models.Employee{ID: "emp-123"}, nil)
                
                repo.EXPECT().CreateEmployee(mock.Anything, mock.Anything).
                    Return(nil)
                
                service := NewEmployeeService(abhiClient, repo)
                return service, abhiClient, repo
            },
            assert: func(t *testing.T, result *models.Employee, err error) {
                require.NoError(t, err)
                require.NotNil(t, result)
                require.Equal(t, "emp-123", result.ID)
            },
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service, abhiClient, repo := tt.setup(t)
            
            result, err := service.CreateEmployee(context.Background(), tt.input)
            
            tt.assert(t, result, err)
            abhiClient.AssertExpectations(t)
            repo.AssertExpectations(t)
        })
    }
}
```

#### 2. Integration Testing
```go
// tests/integration/employee_integration_test.go
func TestEmployeeIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup test environment
    testEnv := setupIntegrationTest(t)
    defer testEnv.Cleanup()
    
    // Test complete flow
    t.Run("CreateEmployee", func(t *testing.T) {
        employee := &CreateEmployeeRequest{
            FirstName: "Integration",
            LastName:  "Test",
            Email:     "integration.test@example.com",
        }
        
        // Send via API Gateway
        result, err := testEnv.APIClient.CreateEmployee(context.Background(), employee)
        require.NoError(t, err)
        require.NotNil(t, result)
        
        // Verify in Abhi mock
        testEnv.AbhiMock.AssertEmployeeExists("integration.test@example.com")
        
        // Verify event was published
        testEnv.EventAsserter.AssertEventPublished("employee.created", 5*time.Second)
        
        // Verify caching
        cached := testEnv.Cache.GetEmployee(result.ID)
        require.NotNil(t, cached)
    })
}
```

---

## Conclusion

This microservices architecture provides a robust, scalable, and secure foundation for integrating the enhanced Abhi Go SDK into your backend services. The design emphasizes:

### Key Benefits

1. **Separation of Concerns**: Each service handles a specific domain with clear boundaries
2. **Scalability**: Services can be scaled independently based on demand
3. **Security**: Multiple layers of security including JWT, request signing, encryption, and rate limiting
4. **Resilience**: Circuit breakers, retry logic, and graceful degradation
5. **Observability**: Comprehensive monitoring, logging, and alerting
6. **Maintainability**: Clean architecture with established patterns and conventions

### Enhanced Security Features Utilized

- ✅ **Request Signing**: All Abhi API calls signed with HMAC-SHA256
- ✅ **Credential Encryption**: AES-GCM encryption for sensitive data
- ✅ **Rate Limiting**: Token bucket algorithm preventing API abuse
- ✅ **Circuit Breakers**: Protection against cascade failures
- ✅ **JWT Authentication**: Secure token-based authentication
- ✅ **Input Validation**: Comprehensive validation at multiple layers

### Implementation Timeline

- **Phase 1-2**: Foundation (4 weeks)
- **Phase 3-4**: Core Services (4 weeks)
- **Phase 5**: Production Ready (4 weeks)
- **Total**: 12 weeks to production deployment

This architecture leverages all the enhanced security and performance features of your Abhi SDK while providing a production-ready microservices foundation that can grow with your business needs.

The combination of API Gateway + RabbitMQ + dedicated Abhi Gateway service provides the optimal balance of performance, security, and maintainability for your use case.