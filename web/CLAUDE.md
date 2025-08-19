# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- **Start development server**: `pnpm dev` (runs on http://localhost:5173)
- **Build for production**: `pnpm build`
- **Preview production build**: `pnpm preview`
- **Format code**: `pnpm format` (uses Biome)
- **Lint and format**: `pnpm check` (uses Biome with auto-fix)

## Architecture Overview

This is a SolidJS frontend application built with Rsbuild that communicates with a Go backend API called "mindscape". The project structure follows a clear separation of concerns:

### Tech Stack
- **Frontend Framework**: SolidJS with TypeScript
- **Build Tool**: Rsbuild with Solid plugin
- **Styling**: TailwindCSS 4.x
- **Code Quality**: Biome (formatter + linter)
- **Package Manager**: pnpm

### API Integration
- Auto-generated TypeScript API client from OpenAPI/Swagger spec
- Located in `src/api/` with models and API classes
- Backend proxy configured to `http://0.0.0.0:60000` for `/api` and `/assets` routes
- Uses JWT authorization headers for authenticated endpoints

### Key API Features
- User authentication (login/signup)
- User management (CRUD operations, admin functions)
- Profile picture upload/download
- All API calls require authorization headers for authenticated endpoints

### Project Structure
```
src/
├── api/           # Auto-generated API client
│   ├── apis/      # API endpoint classes (UsersApi, etc.)
│   ├── models/    # TypeScript models/interfaces
│   └── runtime.ts # Base API runtime
├── App.tsx        # Main application component
└── index.tsx      # Application entry point
```

### Build Configuration
- Rsbuild with Solid plugin and Babel for JSX/TSX processing
- Path alias `@/*` maps to `./src/*`
- Development server with CORS enabled
- Proxy configuration for backend API integration

### Code Standards
- Single quotes for JavaScript/TypeScript
- Space-based indentation
- Biome handles all formatting and linting
- Import organization enabled
- CSS modules support configured


# TailwindCSS
In order to understand the documentation for TailwindCSS please check the [TailwindCSS documentation](https://tailwindcss.com/)

# Solid JS 
In order to understand the documentation for Solid JS please check the [Solid JS documentation](https://docs.solidjs.com/)

# Solid Router
In order to understand the documentation for Solid Router please check the [Solid Router documentation](https://docs.solidjs.com/solid-router/)
