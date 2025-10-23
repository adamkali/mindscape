# Background Context System Implementation Plan


# Background Context System Implementation Plan


## Project Overview

   This document outlines the implementation of a centralized background image management system
   for the Mindscape frontend. The system will replace individual background fetching on each page
   with a unified context that automatically manages user background preferences.

## Current State Analysis


### Existing Implementation

    Currently, each page in the application individually handles background management:

```typescript
    // Current pattern in each page (Home.tsx, Login.tsx, etc.)
    const [defaultBackground] = createResource(async () => {
        const api = new BackgroundApi();
        const response = await api.getDefaultBackground();
        if (response.success && response.data) {
            return response.data;
        } else {
            throw new Error('Failed to fetch default background: ' + response.message);
        }
    });

    // Applied via inline styles
    style={{ 'background-image': `url(${defaultBackground()})` }}
```

### New Background Selection Feature

    The EditProfile.tsx now includes `handleBackgroundSelect` functionality:

```typescript
    const handleBackgroundSelect = async (backgroundUrl: string) => {
        setSelectedBackground(backgroundUrl);
        setSuccess('Background updated! Changes are applied immediately.');

        const userApi = new UserApi();
        try {
            const result = await userApi.setUserBackground({
                authorization: `Bearer ${auth.token()}`,
                background: backgroundUrl
            });
            if (!result.success) {
                setError(result.message || 'Failed to update background');
            } else {
                console.log('Background updated successfully:', result.data);
            }
        } catch (error: any) {
            setError(error.message || 'Failed to update background');
            if ((error as ResponseError).response.status === 401) {
                auth.logout();
            }
        }
        setError('');
    };
```

### Current Issues

- ( ) Duplicate API calls across pages for default background
- ( ) No centralized user background preference management
- ( ) Background changes in EditProfile don't affect other pages immediately
- ( ) Inconsistent background loading patterns
- ( ) No caching mechanism for background data

## Proposed Architecture


### Background Context Structure

```typescript
    // web/src/contexts/BackgroundContext.tsx
    interface BackgroundContextValue {
        currentBackground: () => string | undefined;
        setUserBackground: (backgroundUrl: string) => Promise<void>;
        isLoading: () => boolean;
        error: () => string | null;
        refreshBackground: () => Promise<void>;
        backgroundChoices: () => string[] | undefined;
        isLoadingChoices: () => boolean;
    }

    interface BackgroundProviderProps {
        children: JSX.Element;
    }

    const BackgroundProvider = (props: BackgroundProviderProps) => {
        // Centralized background state management
        // Auto-fetches user's selected background on authentication
        // Falls back to default background when no user preference
        // Provides methods for updating and refreshing background data
        // Caches background choices and default background
    };
```

### Context Features

#### Automatic Background Loading

- ( ) Loads user's selected background on authentication
- ( ) Falls back to default background if no user preference exists
- ( ) Caches background data to avoid repeated API calls
- ( ) Automatically refreshes on user login/logout events
- ( ) Handles loading and error states gracefully

#### Background Management Methods

- ( ) `setUserBackground()` - Updates user preference and persists to backend
- ( ) `refreshBackground()` - Manually refresh background data
- ( ) `currentBackground()` - Get current background URL (user or default)
- ( ) Background choices loading and caching
- ( ) Error handling and retry mechanisms

#### Performance Optimizations

- ( ) Single API call per session for user background
- ( ) Caching mechanism for background URLs and choices
- ( ) Lazy loading of background images
- ( ) Debounced background updates to prevent API spam
- ( ) Memory management for cached data

## Implementation Plan


### Phase 1: Create Background Context Infrastructure


#### Background Context Implementation

```typescript
     // web/src/contexts/BackgroundContext.tsx
     import { createContext, createSignal, createEffect, createResource, useContext } from 'solid-js';
     import { useAuth } from '@/contexts/AuthContext';
     import { BackgroundApi, UserApi } from '@/api';

     const BackgroundContext = createContext<BackgroundContextValue>();

     export const BackgroundProvider = (props: BackgroundProviderProps) => {
         const auth = useAuth();
         const [userBackground, setUserBackground] = createSignal<string>('');
         const [error, setError] = createSignal<string | null>(null);

         // Load default background once
         const [defaultBackground] = createResource(async () => {
             const api = new BackgroundApi();
             const response = await api.getDefaultBackground();
             if (response.success && response.data) {
                 return response.data;
             }
             throw new Error('Failed to fetch default background');
         });

         // Load background choices once
         const [backgroundChoices] = createResource(async () => {
             const backgroundApi = new BackgroundApi();
             const userApi = new UserApi();
             
             // Implementation similar to EditProfile.tsx background choices logic
             // Combine global and user-specific choices
         });

         // Auto-fetch user background on auth change
         createEffect(async () => {
             if (auth.isAuthenticated() && auth.token()) {
                 await loadUserBackground();
             } else {
                 setUserBackground('');
                 setError(null);
             }
         });

         const loadUserBackground = async () => {
             try {
                 const userApi = new UserApi();
                 const response = await userApi.getUserBackground({
                     authorization: `Bearer ${auth.token()}`
                 });
                 if (response.success && response.data) {
                     setUserBackground(response.data);
                 }
             } catch (err: any) {
                 console.warn('Failed to load user background:', err);
                 setError(err.message);
             }
         };

         const updateUserBackground = async (backgroundUrl: string) => {
             try {
                 const userApi = new UserApi();
                 const result = await userApi.setUserBackground({
                     authorization: `Bearer ${auth.token()}`,
                     background: backgroundUrl
                 });
                 
                 if (result.success) {
                     setUserBackground(backgroundUrl);
                     setError(null);
                 } else {
                     throw new Error(result.message || 'Failed to update background');
                 }
             } catch (err: any) {
                 setError(err.message);
                 throw err;
             }
         };

         const contextValue: BackgroundContextValue = {
             currentBackground: () => userBackground() || defaultBackground(),
             setUserBackground: updateUserBackground,
             isLoading: () => defaultBackground.loading,
             error,
             refreshBackground: loadUserBackground,
             backgroundChoices,
             isLoadingChoices: () => backgroundChoices.loading
         };

         return (
             <BackgroundContext.Provider value={contextValue}>
                 {props.children}
             </BackgroundContext.Provider>
         );
     };

     export const useBackground = () => {
         const context = useContext(BackgroundContext);
         if (!context) {
             throw new Error('useBackground must be used within BackgroundProvider');
         }
         return context;
     };
```

#### Background Style Hook

```typescript
     // web/src/hooks/useBackground.ts
     import { createMemo } from 'solid-js';
     import { useBackground } from '@/contexts/BackgroundContext';

     export const useBackgroundStyle = () => {
         const { currentBackground } = useBackground();
         
         return createMemo(() => ({
             'background-image': `url(${currentBackground()})`,
             'background-size': 'cover',
             'background-position': 'center center',
             'background-repeat': 'no-repeat',
             'background-attachment': 'fixed'
         }));
     };

     export const useBackgroundClass = () => {
         const { currentBackground } = useBackground();
         
         return createMemo(() => 
             currentBackground() ? 'min-h-screen bg-background' : 'min-h-screen bg-background bg-gray-900'
         );
     };
```

### Phase 2: App Integration


#### Update App.tsx

```typescript
     // web/src/App.tsx
     import { BackgroundProvider } from '@/contexts/BackgroundContext';

     function App() {
         return (
             <AuthProvider>
                 <BackgroundProvider>
                     <Router>
                         <Routes>
                             {/* existing routes */}
                         </Routes>
                     </Router>
                 </BackgroundProvider>
             </AuthProvider>
         );
     }
```

#### Provider Hierarchy

     ```
     App
     └── AuthProvider
         └── BackgroundProvider
             └── Router
                 └── Routes
                     └── Individual Pages
     ```

### Phase 3: Page Updates


#### Home Page Update

```typescript
     // web/src/pages/Home.tsx
     import { useBackgroundStyle } from '@/hooks/useBackground';

     const Home = () => {
         const auth = useAuth();
         const backgroundStyle = useBackgroundStyle();
         
         // Remove existing defaultBackground createResource
         // Remove existing background fetching logic

         return (
             <div 
                 class="min-h-screen bg-background"
                 style={backgroundStyle()}
                 onClick={(e) => {
                     // existing click handler
                 }}
             >
                 <Header />
                 {/* rest of existing content */}
             </div>
         );
     };
```

#### Login Page Update

```typescript
     // web/src/pages/Login.tsx
     import { useBackgroundStyle } from '@/hooks/useBackground';

     const Login = () => {
         const backgroundStyle = useBackgroundStyle();
         
         return (
             <div 
                 class="min-h-screen flex items-center justify-center bg-background"
                 style={backgroundStyle()}
             >
                 {/* existing login form content */}
             </div>
         );
     };
```

#### Signup Page Update

```typescript
     // web/src/pages/Signup.tsx  
     import { useBackgroundStyle } from '@/hooks/useBackground';

     const Signup = () => {
         const backgroundStyle = useBackgroundStyle();
         
         return (
             <div 
                 class="min-h-screen flex items-center justify-center bg-background"
                 style={backgroundStyle()}
             >
                 {/* existing signup form content */}
             </div>
         );
     };
```

### Phase 4: EditProfile Integration


#### Update EditProfile Background Selection

```typescript
     // web/src/pages/EditProfile.tsx
     import { useBackground, useBackgroundStyle } from '@/hooks/useBackground';

     const EditProfile = () => {
         const { 
             setUserBackground, 
             currentBackground, 
             backgroundChoices, 
             isLoadingChoices 
         } = useBackground();
         const backgroundStyle = useBackgroundStyle();
         
         // Remove existing background fetching logic
         // Remove selectedBackground signal (use currentBackground from context)

         const handleBackgroundSelect = async (backgroundUrl: string) => {
             try {
                 await setUserBackground(backgroundUrl);
                 setSuccess('Background updated! Changes are applied immediately.');
                 setError('');
             } catch (error: any) {
                 setError(error.message || 'Failed to update background');
                 if ((error as ResponseError).response.status === 401) {
                     auth.logout();
                 }
             }
         };

         return (
             <div 
                 class="min-h-screen bg-background"
                 style={backgroundStyle()}
             >
                 {/* existing profile editing content */}
                 
                 {/* Background Selection Section */}
                 <Card variant="glass">
                     <CardHeader title="Background Settings" subtitle="Choose from available backgrounds or upload a custom one" />
                     <CardBody padding="lg">
                         <Show when={isLoadingChoices()}>
                             <div class="text-white/70">Loading backgrounds...</div>
                         </Show>
                         <Show when={backgroundChoices()}>
                             <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
                                 <For each={backgroundChoices()}>
                                     {(backgroundUrl) => (
                                         <div 
                                             class={`relative aspect-video rounded-lg overflow-hidden cursor-pointer border-2 transition-all duration-300 hover:scale-105 ${
                                                 currentBackground() === backgroundUrl 
                                                     ? 'border-white/70 ring-2 ring-white/50' 
                                                     : 'border-white/20 hover:border-white/40'
                                             }`}
                                             onClick={() => handleBackgroundSelect(backgroundUrl)}
                                         >
                                             {/* existing background option rendering */}
                                         </div>
                                     )}
                                 </For>
                             </div>
                         </Show>
                     </CardBody>
                 </Card>
             </div>
         );
     };
```

## File Structure Changes


### New Files

    ```
    web/src/
    ├── contexts/
    │   ├── AuthContext.tsx          # Existing
    │   └── BackgroundContext.tsx    # New - Core background management
    ├── hooks/
    │   └── useBackground.ts         # New - Convenience hooks for styling
    ```

### Modified Files

    ```
    web/src/
    ├── App.tsx                      # Add BackgroundProvider
    ├── pages/
    │   ├── Home.tsx                 # Remove individual background fetching
    │   ├── EditProfile.tsx          # Use context for background management
    │   ├── Login.tsx                # Use background context
    │   ├── Signup.tsx               # Use background context
    │   └── admin/
    │       ├── index.tsx            # Use background context (if needed)
    │       └── Showcase.tsx         # Use background context (if needed)
    ```

## Implementation Benefits


### Performance Improvements

- ( ) Eliminates duplicate API calls for default background across pages
- ( ) Single user background fetch per session
- ( ) Cached background choices and data
- ( ) Reduced network requests and improved page load times
- ( ) Efficient memory usage with proper cleanup

### User Experience Enhancements

- ( ) Instant background changes across all pages when updated in EditProfile
- ( ) Consistent background experience during navigation
- ( ) Smooth loading states with proper fallbacks
- ( ) Immediate visual feedback for background changes
- ( ) Graceful error handling and recovery

### Code Quality Improvements

- ( ) Centralized background state management
- ( ) Elimination of duplicate background fetching logic
- ( ) Clean separation of concerns
- ( ) Consistent API usage patterns
- ( ) Easier testing and maintenance

### Maintainability Benefits

- ( ) Single source of truth for background management
- ( ) Easy to add new background-related features
- ( ) Simplified page components with removed boilerplate
- ( ) Better error handling and debugging capabilities
- ( ) Future-proof architecture for background enhancements

## API Integration


### No Backend Changes Required

    The implementation uses existing API endpoints:
- ( ) `backgroundApi.getDefaultBackground()` - Fetches system default background
- ( ) `userApi.setUserBackground()` - Saves user's background preference  
- ( ) `userApi.getUserBackground()` - Retrieves user's saved background
- ( ) `backgroundApi.getBackgroundChoices()` - Gets available background options
- ( ) `userApi.getUserBackgroundChoices()` - Gets user-specific backgrounds

### Authentication Integration

- ( ) Automatically responds to auth state changes
- ( ) Clears user background on logout
- ( ) Loads user background on login
- ( ) Handles authentication errors gracefully
- ( ) Integrates seamlessly with existing AuthContext

## Implementation Timeline


### Day 1: Core Infrastructure

- ( ) Create BackgroundContext with basic functionality
- ( ) Implement useBackground and useBackgroundStyle hooks
- ( ) Add BackgroundProvider to App.tsx
- ( ) Test basic context functionality
- ( ) Verify authentication integration

### Day 2: Page Integration  

- ( ) Update Home.tsx to use background context
- ( ) Update Login.tsx and Signup.tsx
- ( ) Update admin pages if applicable
- ( ) Remove duplicate background fetching logic
- ( ) Test page navigation and background consistency

### Day 3: EditProfile Integration & Testing

- ( ) Update EditProfile.tsx background selection
- ( ) Test immediate background updates across pages
- ( ) Verify background persistence and loading
- ( ) Perform cross-browser testing
- ( ) Document usage patterns and API

## Testing Strategy


### Unit Testing

- ( ) Test BackgroundContext state management
- ( ) Test useBackground hook functionality
- ( ) Test background style generation
- ( ) Test error handling scenarios
- ( ) Test authentication integration

### Integration Testing

- ( ) Test context provider hierarchy
- ( ) Test page-to-page background consistency
- ( ) Test EditProfile background selection flow
- ( ) Test login/logout background behavior
- ( ) Test API error handling and recovery

### User Experience Testing

- ( ) Verify immediate background updates
- ( ) Test loading states and fallbacks
- ( ) Verify background persistence across sessions
- ( ) Test responsive background behavior
- ( ) Validate accessibility considerations

## Future Enhancements


### Advanced Features

- ( ) Background caching with expiration
- ( ) Progressive image loading
- ( ) Background animation support
- ( ) Multiple background themes
- ( ) Background scheduling (time-based changes)
- ( ) Background sharing between users
- ( ) Background categories and filtering

### Performance Optimizations

- ( ) Image optimization and compression
- ( ) Lazy loading for background choices
- ( ) Service worker caching
- ( ) CDN integration for backgrounds
- ( ) Background preloading strategies

## Conclusion

This Background Context System provides a robust, centralized solution for managing background images across the Mindscape application. The implementation eliminates code duplication, improves performance, and enhances user experience while maintaining compatibility with existing APIs.

The phased approach ensures stable incremental delivery with thorough testing at each step. The architecture is designed to be extensible for future background-related features while keeping the current implementation clean and maintainable.
