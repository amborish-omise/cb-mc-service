# System Architecture Diagram

This diagram shows the overall system architecture including all components, layers, and their relationships.

```mermaid
graph TB
    subgraph "External Systems"
        A[Ethoca Alerts API]
        B[Datadog APM]
        C[Datadog Logs]
    end
    
    subgraph "Load Balancer / API Gateway"
        D[Load Balancer]
    end
    
    subgraph "Application Layer"
        E[HTTP Server]
        F[gRPC Server]
        G[Combined Server]
    end
    
    subgraph "Handler Layer"
        H[Case Handlers]
        I[Document Handlers]
        J[Ethoca Webhook Handlers]
    end
    
    subgraph "Service Layer"
        K[Case Service]
        L[Document Service]
        M[Ethoca Webhook Service]
    end
    
    subgraph "Model Layer"
        N[Case Models]
        O[Document Models]
        P[Ethoca Webhook Models]
    end
    
    subgraph "Configuration"
        Q[App Config]
        R[Datadog Config]
        S[Ethoca Webhook Config]
    end
    
    subgraph "Middleware"
        T[CORS Middleware]
        U[Logging Middleware]
        V[Datadog Middleware]
        W[Recovery Middleware]
    end
    
    subgraph "Infrastructure"
        X[Docker Containers]
        Y[Buildkite CI/CD]
        Z[Environment Config]
    end
    
    A --> D
    D --> E
    D --> F
    D --> G
    
    E --> H
    E --> I
    E --> J
    F --> H
    F --> I
    F --> J
    G --> H
    G --> I
    G --> J
    
    H --> K
    I --> L
    J --> M
    
    K --> N
    L --> O
    M --> P
    
    Q --> E
    Q --> F
    Q --> G
    R --> V
    S --> M
    
    E --> T
    E --> U
    E --> V
    E --> W
    
    X --> E
    X --> F
    X --> G
    Y --> X
    Z --> Q
    Z --> R
    Z --> S
    
    B --> V
    C --> U
    
    style A fill:#e3f2fd
    style B fill:#e8f5e8
    style C fill:#e8f5e8
    style D fill:#fff3e0
    style E fill:#f3e5f5
    style F fill:#f3e5f5
    style G fill:#f3e5f5
    style J fill:#e8f5e8
    style M fill:#e8f5e8
    style P fill:#e8f5e8
    style S fill:#e8f5e8
```

## Architecture Layers

### **External Systems**
- **Ethoca Alerts API**: Source of webhook events
- **Datadog APM**: Application performance monitoring
- **Datadog Logs**: Centralized logging service

### **Load Balancer**
- Routes traffic to appropriate servers
- Handles SSL termination
- Provides health checks

### **Application Layer**
- **HTTP Server**: RESTful API endpoints
- **gRPC Server**: High-performance RPC calls
- **Combined Server**: Both HTTP and gRPC in single process

### **Handler Layer**
- **Case Handlers**: Case management operations
- **Document Handlers**: Document processing
- **Ethoca Webhook Handlers**: Webhook event processing

### **Service Layer**
- **Case Service**: Business logic for cases
- **Document Service**: Document processing logic
- **Ethoca Webhook Service**: Webhook processing logic

### **Model Layer**
- Data structures and validation rules
- JSON tags for serialization
- Validation tags for input validation

### **Configuration**
- Environment-specific settings
- Feature flags and toggles
- Service endpoints and credentials

### **Middleware**
- **CORS**: Cross-origin resource sharing
- **Logging**: Request/response logging
- **Datadog**: Tracing and metrics
- **Recovery**: Panic recovery

### **Infrastructure**
- **Docker**: Containerization
- **Buildkite**: CI/CD pipeline
- **Environment Config**: Deployment configuration
