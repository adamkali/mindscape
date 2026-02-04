# Widget System Architecture Planning


# Widget System Architecture Planning


## Overview

   The Mindscape application is being enhanced with a configurable widget/module system that will populate the right 2/3 of the screen with user-customizable components. This document outlines the complete architecture and implementation strategy developed through strategic analysis.

## Problem Statement

- Current UI has unused space (right 2/3 of screen) while left side contains folders/bookmarks navigation
- Need extensible system for adding functionality like SearXNG search, Plex integration, RSS notifications, PR notifications
- Widget configurations must be saved per user and be highly configurable

## Architecture Decision: S3/MinIO-Based Schema System


### Core Concept

    Widget schemas stored in S3 buckets define widget structure and capabilities, while user configurations are stored in individual user S3 buckets. This creates a schema-driven system that enables dynamic UI generation and maximum extensibility.

### Bucket Organization

```
    mindscape-widgets/
    ├── schemas/
    │   ├── plex/
    │   │   ├── schema.json
    │   │   └── config_template.json
    │   ├── rss/
    │   │   ├── schema.json
    │   │   └── config_template.json
    │   ├── searxng/
    │   │   ├── schema.json
    │   │   └── config_template.json
    │   └── github-pr/
    │       ├── schema.json
    │       └── config_template.json
    └── users/
        └── <user_id>/
            ├── widgets/
            │   ├── plex_config.json
            │   ├── rss_config.json
            │   └── searxng_config.json
            ├── layout.json
            └── widget_preferences.json
```

## Grid-Based Layout System


### Grid Structure

- 12-column CSS Grid system
- Configurable row height (default: 60px)
- 16px gap between widgets
- Widgets can span multiple columns and rows

### Widget Sizing Examples

- SearXNG: 12 columns × 2 rows (full width search bar)
- RSS Reader: 8 columns × 4 rows (wide article list)
- Plex Recently Added: 4 columns × 6 rows (compact movie tiles)
- GitHub PR Notifications: 6 columns × 3 rows (notification list)

### Layout Schema Format

```json
    {
      "gridSettings": {
        "columns": 12,
        "rowHeight": 60,
        "gap": 16
      },
      "widgets": [
        {
          "id": "widget_1",
          "type": "searxng",
          "position": {
            "x": 0,      // Column position
            "y": 0,      // Row position
            "width": 12, // Spans 12 columns (full width)
            "height": 2  // Spans 2 rows
          }
        }
      ]
    }
```

## Widget Schema Format


### Schema Structure (JSON Schema)

```json
    {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "type": "object",
      "title": "SearXNG Widget Schema",
      "description": "Configuration schema for SearXNG search integration",
      "layout": {
        "defaultSize": {
          "width": 12,     // Full width (12 columns)
          "height": 2      // 2 rows tall
        },
        "minSize": {
          "width": 6,      // Minimum 6 columns
          "height": 1      // Minimum 1 row
        },
        "maxSize": {
          "width": 12,     // Maximum full width
          "height": 4      // Maximum 4 rows
        },
        "resizable": true
      },
      "properties": {
        "serverUrl": {
          "type": "string",
          "format": "uri",
          "description": "SearXNG server URL (localhost, home server, tailscale DNS)"
        },
        "defaultEngine": {
          "type": "string",
          "description": "Default search engine",
          "enum": ["google", "duckduckgo", "bing"]
        },
        "categories": {
          "type": "array",
          "items": {"type": "string"},
	"description": "Enabled search categories"
        }
      },
      "required": ["serverUrl"]
    }
```

### Target Widget Schemas


#### Plex Widget

```json
     {
       "type": "object",
       "title": "Plex Widget",
       "layout": {
         "defaultSize": {"width": 4, "height": 6},
         "minSize": {"width": 3, "height": 4},
         "maxSize": {"width": 8, "height": 8}
       },
       "properties": {
         "serverUrl": {"type": "string", "format": "uri"},
         "apiToken": {"type": "string", "sensitive": true},
         "libraries": {"type": "array", "items": {"type": "string"}},
         "refreshInterval": {"type": "number", "default": 3600},
         "showRecentlyAdded": {"type": "boolean", "default": true},
         "maxItems": {"type": "number", "default": 10}
       }
     }
```

#### RSS Notification Widget

```json
     {
       "type": "object",
       "title": "RSS Notification Widget",
       "layout": {
         "defaultSize": {"width": 8, "height": 4},
         "minSize": {"width": 4, "height": 3},
         "maxSize": {"width": 12, "height": 8}
       },
       "properties": {
         "feedUrls": {"type": "array", "items": {"type": "string"}},
         "keywords": {"type": "array", "items": {"type": "string"}},
         "refreshInterval": {"type": "number", "default": 1800},
         "notificationChannels": {
           "type": "array", 
           "items": {
             "type": "string", 
             "enum": ["email", "telegram", "slack"]
           }
         },
         "maxArticles": {"type": "number", "default": 20}
       }
     }
```

#### GitHub PR Notification Widget

```json
     {
       "type": "object", 
       "title": "GitHub PR Notification Widget",
       "layout": {
         "defaultSize": {"width": 6, "height": 3},
         "minSize": {"width": 4, "height": 2},
         "maxSize": {"width": 8, "height": 6}
       },
       "properties": {
         "repositories": {"type": "array", "items": {"type": "string"}},
         "githubToken": {"type": "string", "sensitive": true},
         "notificationRules": {
           "type": "object",
           "properties": {
             "reviewRequested": {"type": "boolean"},
             "assignedPRs": {"type": "boolean"},
             "mentionedPRs": {"type": "boolean"},
             "authoredPRs": {"type": "boolean"}
           }
         },
         "refreshInterval": {"type": "number", "default": 300}
       }
     }
```

## Backend Architecture


### Database Changes (Minimal)

    Keep database lightweight - main data stored in S3:
```sql
    CREATE TABLE user_widgets (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        widget_type VARCHAR(50) NOT NULL,
        widget_instance_id UUID NOT NULL, -- Links to S3 config
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
```

### Service Layer

```go
    type IWidgetService interface {
        GetUserWidgets(userID string) ([]WidgetConfig, error)
        CreateWidget(userID, widgetType string, config map[string]interface{}) error
		UpdateWidget(userID, widgetID string, config map[string]interface{}) error
        DeleteWidget(userID, widgetID string) error
        GetWidgetSchema(widgetType string) (*WidgetSchema, error)
        ValidateWidgetConfig(widgetType string, config map[string]interface{}) error
    }

    type WidgetService struct {
        minioClient *minio.Client
        cache       *WidgetCache
        validator   *WidgetValidator
    }
```

### Caching Strategy

```go
    type WidgetCache struct {
        schemas     *cache.Cache  // Global widget schemas
        userConfigs *cache.Cache  // User-specific widget configurations
        layouts     *cache.Cache  // User layout configurations
        ttl         time.Duration // Cache time-to-live
    }
```

### API Endpoints

```
    GET    /api/widgets/schemas           # List available widget types
    GET    /api/widgets/schemas/:type     # Get specific widget schema
    GET    /api/widgets                  # Get user's widgets
    POST   /api/widgets                  # Create new widget
    PUT    /api/widgets/:id              # Update widget configuration
    DELETE /api/widgets/:id              # Delete widget
    GET    /api/widgets/layout           # Get user's layout
    PUT    /api/widgets/layout           # Update user's layout
```

## Frontend Architecture


### Grid Container Implementation

```typescript
    const WidgetGrid: Component = () => {
      const [layout, setLayout] = createSignal<LayoutConfig>();
      const [widgets, setWidgets] = createSignal<WidgetConfig[]>([]);
      
      return (
        <div 
          class="widget-grid"
          style={{
            display: 'grid',
            'grid-template-columns': 'repeat(12, 1fr)',
            'grid-auto-rows': '60px',
            gap: '16px'
          }}
        >
          <For each={widgets()}>
            {(widget) => (
              <WidgetWrapper
                widget={widget}
                style={{
                  'grid-column': `${widget.position.x + 1} / span ${widget.position.width}`,
                  'grid-row': `${widget.position.y + 1} / span ${widget.position.height}`
                }}
                onMove={handleWidgetMove}
                onResize={handleWidgetResize}
              />
            )}
          </For>
        </div>
      );
    };
```

### Dynamic UI Generation

    Schema-driven form generation for widget configuration:
```typescript
    function generateWidgetConfigForm(schema: WidgetSchema) {
        return schema.properties.map(field => {
            switch(field.type) {
                case 'string':
                    return <Input 
                        label={field.label} 
                        type={field.sensitive ? "password" : "text"}
                        description={field.description}
                    />;
                case 'boolean':
                    return <Checkbox 
                        label={field.label}
                        description={field.description}
                    />;
                case 'array':
                    return <MultiSelect
                        label={field.label}
                        options={field.items?.enum || []}
                    />;
            }
        });
    }
```

### Drag & Drop System

```typescript
    interface GridPosition {
      x: number;      // Column (0-11 for 12-column grid)
      y: number;      // Row (0+)
      width: number;  // Column span (1-12)
      height: number; // Row span (1+)
    }

    const handleWidgetMove = (widgetId: string, newPosition: GridPosition) => {
      // Collision detection
      if (hasCollision(newPosition, otherWidgets)) {
        return; // Prevent move
      }
      
      // Update layout
      updateWidgetPosition(widgetId, newPosition);
      
      // Save to S3
      saveLayoutToS3(currentLayout());
    };
```

## Security Considerations


### Encryption Strategy

- Sensitive fields (API tokens) encrypted at rest using AES-256
- Client-side encryption before S3 storage
- Separate encryption keys per user stored securely

### Access Control

- JWT authentication for all widget operations
- User isolation through S3 bucket namespacing
- Input sanitization and validation for all configurations
- Rate limiting on widget creation/updates

### Data Validation

- JSON Schema validation for all widget configurations
- SQL injection prevention (though minimal DB usage)
- XSS prevention in widget rendering
- CORS configuration for widget iframe content

## Performance Optimization


### Caching Layers

    1. Memory cache for frequently accessed schemas
    2. Redis cache for user configurations
    3. Browser localStorage for layout preferences
    4. CDN for static widget assets

### S3 Optimization

- Connection pooling for MinIO client
- Batch operations for multiple widget updates
- Compression for large configuration files
- Circuit breaker pattern for S3 failures

### Frontend Performance

- Lazy loading of widget components
- Virtual scrolling for large widget lists
- Debounced auto-save for layout changes
- Optimistic UI updates during drag operations

## Implementation Roadmap


### Phase 1: Foundation (Weeks 1-2)

- ( ) Design and implement S3 bucket structure
- ( ) Create schema validation service
- ( ) Develop caching mechanism
- ( ) Basic widget CRUD API endpoints
- ( ) Database table creation

### Phase 2: Core Widget System (Weeks 3-4)

- ( ) Grid layout container implementation
- ( ) Dynamic UI generation from schemas
- ( ) Widget configuration management
- ( ) Basic drag and drop functionality
- ( ) SearXNG widget implementation

### Phase 3: Advanced Features (Weeks 5-6)

- ( ) Collision detection system
- ( ) Widget resizing capabilities
- ( ) Plex widget implementation
- ( ) RSS notification widget
- ( ) GitHub PR notification widget

### Phase 4: Polish & Optimization (Weeks 7-8)

- ( ) Performance optimization
- ( ) Security hardening
- ( ) Responsive design implementation
- ( ) Error handling and recovery
- ( ) Documentation and testing

## Technical Risks & Mitigation


### Identified Risks

    1. *S3 Performance Bottlenecks*
- Mitigation: Multi-layer caching, connection pooling
    2. *Complex Grid Layout Edge Cases*
- Mitigation: Comprehensive collision detection, fallback positioning
    3. *Schema Evolution Management*
- Mitigation: Versioned schemas, backward compatibility layer
    4. *Widget Security Vulnerabilities*
- Mitigation: Sandboxed widget execution, strict CSP policies

### Alternative Approaches Considered

    1. Database-only storage (rejected: less flexible)
    2. Micro-frontend architecture (rejected: too complex)
    3. Static widget definitions (rejected: not extensible enough)

## Success Metrics


### Technical Metrics

- Widget load time < 200ms
- Layout save time < 100ms
- Support for 20+ concurrent widgets
- 99.9% uptime for widget system

### User Experience Metrics

- Widget configuration completion rate > 90%
- Average time to configure widget < 2 minutes
- User retention with widget usage
- Widget usage frequency per user

## Future Extensions


### Planned Enhancements

- Widget marketplace for community contributions
- Widget templates and presets
- Advanced widget interactions and data sharing
- Mobile-responsive widget layouts
- Widget performance analytics
- Import/export widget configurations

### Integration Opportunities

- Home Assistant integration
- Docker container monitoring
- System resource monitoring
- Calendar and task management widgets
- Social media feed aggregation

## References

- {/ JSON Schema Specification}[https://json-schema.org/]
- {/ CSS Grid Layout}[https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Grid_Layout]
- {/ MinIO Go SDK}[https://docs.min.io/docs/golang-client-quickstart-guide.html]
- {/ SolidJS Documentation}[https://www.solidjs.com/docs/latest]

## Related Documents

- {/ API Documentation}[./api-docs.norg]
