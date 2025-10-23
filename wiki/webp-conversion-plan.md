# WebP Image Conversion Implementation Plan


# WebP Image Conversion Implementation Plan


## Problem Statement

   The current 4K background images stored in the S3 bucket are causing significant loading delays
   during background image selection. Users experience slow preview loading times that impact
   the overall user experience.

## Current Architecture Analysis


### Upload Flow

- Images uploaded via `UploadProfilePictureHandler` and `UploadBackgroundHandler`
- Files go directly to MinIO S3-compatible storage via `MinioService`
- No image processing or optimization occurs during upload

### Storage Structure

- **Background Images**: Large 4K images stored in "mindspace" bucket
- **Profile Pictures**: User-specific uploads stored in user UUID buckets
- Current formats: JPEG, PNG (unoptimized)

### Performance Issues

- 4K images causing slow loading during background selection
- No compression optimization
- Large file sizes impacting network transfer

## Proposed Solution: FFmpeg WebP Conversion


### Benefits

- **File Size Reduction**: 25-35% smaller than JPEG, 80% smaller than PNG
- **Loading Speed**: Significantly faster background image previews
- **Quality**: Better compression with same visual quality
- **Compatibility**: 97%+ modern browser support

## Implementation Plan


### 1. WebP Conversion Service

1. Location: `services/ConversionService.go`

#### Interface Design

```go
     type IConversionService interface {
         ConvertToWebP(input io.Reader, quality int) (io.Reader, error)
         ValidateImageFormat(filename string) bool
         GetOptimizedFilename(originalFilename string) string
     }
```

#### Implementation Details

- Use Go's `os/exec` to call FFmpeg binary
- Support quality settings (75-90 for optimal balance)
- Handle multiple input formats: JPEG, PNG, BMP, TIFF
- Automatic file extension handling (.webp)

#### Error Handling

- FFmpeg binary availability check
- Conversion failure fallback to original
- Input validation for supported formats
- Memory management for large files

### 2. FFmpeg Integration Points


#### Profile Picture Upload

1. File: `models/handlers/user_handlers/UploadProfilePictureHandler.go`
1. Line: 111 (before `ms.Upload()`)
     
- Convert image to WebP before MinIO storage
- Update filename to use `.webp` extension
- Maintain error handling chain

#### Background Image Upload

1. File: `models/handlers/user_handlers/UploadBackgroundHandler.go` 
1. Line: 84 (before storage)
     
- Convert 4K images to optimized WebP
- Update Redis caching with new filenames
- Ensure presigned URL generation works with WebP

### 3. Configuration Management


#### Config File Updates

1. Files: `config/development.yaml`, `config/production.yaml`
     
```yaml
     conversion:
       ffmpeg_path: "/usr/bin/ffmpeg"  # System path to FFmpeg binary
       webp_quality: 80                # Quality setting (1-100)
       max_file_size: 50485760        # 50MB limit
       supported_formats: ["jpg", "jpeg", "png", "bmp", "tiff"]
```

#### Configuration Service Integration

1. File: `cmd/configuration/configuration.go`
     
- Add conversion settings to Configuration struct
- Environment variable overrides
- Validation for FFmpeg binary existence

### 4. Handler Updates


#### Upload Profile Picture Handler

```go
     // Add conversion step before MinIO upload
     convertedFile, err := h.cs.ConvertToWebP(src, config.Conversion.WebPQuality)
     if err != nil {
         // Fallback to original file
         convertedFile = src
         filename = file.Filename
     } else {
         filename = h.cs.GetOptimizedFilename(file.Filename)
     }
     
     err = h.ms.Upload(userId, filename, convertedFile, file.Size)
```

#### Upload Background Handler

```go
     // Convert before storage and update filename
     convertedFile, err := h.cs.ConvertToWebP(file, config.Conversion.WebPQuality)
     if err == nil {
         request.File.Filename = h.cs.GetOptimizedFilename(request.File.Filename)
     } else {
         convertedFile = file
     }
```

### 5. Service Layer Integration


#### Controller Updates

1. File: `controllers/user_controller.go`
     
- Add ConversionService to UserController struct
- Inject service in BuildUserController function
- Pass to handlers requiring conversion

#### Service Registration

1. File: Controller registration location
     
- Create ConversionService instance
- Configure with FFmpeg path and settings
- Register in dependency injection

### 6. Migration Strategy


#### Existing Image Conversion

- Background job to convert existing 4K images
- Batch processing with rate limiting
- Backup original files before conversion
- Update database references to new filenames

#### Backward Compatibility

- `GetBackgroundChoices()` prioritizes WebP files
- Fallback to original format if WebP not available
- Gradual migration without breaking existing functionality

### 7. Docker & Deployment Updates


#### Dockerfile Modifications

```dockerfile
     # Add FFmpeg installation
     RUN apt-get update && apt-get install -y ffmpeg
     
     # Verify installation
     RUN ffmpeg -version
```

#### Docker Compose Updates

- Environment variables for FFmpeg path
- Volume mounts if using external FFmpeg
- Health checks for FFmpeg availability

### 8. Frontend Considerations


#### API Client Updates

- No significant changes needed
- Browsers handle WebP transparently
- Update file extension expectations in TypeScript models

#### Browser Compatibility

- WebP support: 97%+ modern browsers
- Automatic fallback handled by server
- No client-side changes required

## Testing Strategy


### Unit Tests

1. Location: `services/ConversionService_test.go`
    
#### Test Cases

- WebP conversion with various quality settings
- Input format validation
- Error handling for unsupported formats
- FFmpeg binary availability checks
- Memory management with large files

### Integration Tests

1. Location: `models/handlers/user_handlers/*_test.go`
    
#### Handler Testing

- Upload with conversion enabled
- Fallback to original on conversion failure
- File size validation post-conversion
- MinIO storage with WebP files

### Performance Benchmarks

- File size reduction measurements
- Conversion time benchmarks
- Memory usage during conversion
- Network transfer speed improvements

## Implementation Phases


### Phase 1: Core Service (Week 1)

- ( ) Implement ConversionService interface
- ( ) Add configuration management
- ( ) Create unit tests
- ( ) FFmpeg integration and error handling

### Phase 2: Handler Integration (Week 1)

- ( ) Update UploadProfilePictureHandler
- ( ) Update UploadBackgroundHandler
- ( ) Integration testing
- ( ) Controller dependency injection

### Phase 3: Deployment Preparation (Week 2)

- ( ) Docker configuration updates
- ( ) Environment-specific testing
- ( ) Migration strategy for existing images
- ( ) Performance benchmarking

### Phase 4: Migration & Monitoring (Week 2)

- ( ) Deploy to staging environment
- ( ) Background conversion of existing images
- ( ) Production deployment
- ( ) Performance monitoring and optimization

## Risk Assessment


### Technical Risks

- FFmpeg binary availability in production
- Conversion failures causing upload errors
- Memory usage with large file processing
- Storage space during migration period

### Mitigation Strategies

- Comprehensive fallback mechanisms
- Memory-efficient streaming conversion
- Gradual migration with rollback capability
- Monitoring and alerting for conversion failures

## Success Metrics


### Performance Improvements

- Background image loading time reduction (target: 60-80%)
- File size reduction measurements
- User experience improvements in image selection

### Technical Metrics

- Conversion success rate (target: >95%)
- Error handling effectiveness
- System resource usage optimization

## Future Enhancements


### Advanced Features

- Multiple quality tiers for different use cases
- Progressive JPEG fallback for older browsers
- Real-time conversion status updates
- Batch conversion API for admin users

### Optimization Opportunities

- CDN integration for WebP delivery
- Client-side format detection
- Adaptive quality based on network conditions
- Image resizing for different viewport sizes
