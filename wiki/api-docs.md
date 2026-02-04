# Mindscape API Documentation


# Mindscape API Documentation


## Overview

   This document provides comprehensive documentation for the Mindscape REST API endpoints. The API follows RESTful conventions and uses JWT-based authentication for protected endpoints.

- *Base Path*: `/api`
- *Version*: 0.0.1
- *Schemes*: HTTP
- *Authentication*: Bearer JWT tokens

## Authentication

   Most endpoints require JWT authentication via the Authorization header:
```
   Authorization: Bearer <jwt_token>
```

   Admin-only endpoints require tokens with elevated privileges.

## API Endpoints


### Background Management


#### Get Default Background

- *Endpoint*: `GET /api/background`
- *Summary*: Get Default Background
- *Authentication*: Not required
- *Response*: StringResponse with background data
- *Status Codes*:
  - 200: OK
  - 404: Not Found
  - 500: Internal Server Error

#### Get Background Choices

- *Endpoint*: `GET /api/background/choices`
- *Summary*: Get Background Choices
- *Authentication*: Not required
- *Response*: StringResponse with available background options
- *Status Codes*:
  - 200: OK
  - 404: Not Found
  - 500: Internal Server Error

#### Upload User Background

- *Endpoint*: `PUT /api/user/background`
- *Summary*: Upload custom background for user
- *Authentication*: Not required (needs to be fixed)
- *Response*: StringResponse
- *Status Codes*:
  - 200: OK
  - 404: Not Found
  - 500: Internal Server Error

### Bookmark Management


#### Create Bookmark

- *Endpoint*: `POST /api/bookmarks`
- *Summary*: Create a new Bookmark
- *Authentication*: Required (`Authorization: Bearer token`)
- *Request Body*: CreateBookmarkParams
```json
       {
         "folder_id": "string",
         "link": "string", 
         "name": "string",
         "user_id": "string"
       }
```
- *Response*: BookmarkResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

#### Move Bookmark

- *Endpoint*: `PATCH /api/bookmarks`
- *Summary*: Move a Bookmark to different folder
- *Authentication*: Required (`Authorization: Bearer token`)
- *Request Body*: MoveBookmarkParams
```json
       {
         "folder_id": "string",
         "id": "string"
       }
```
- *Response*: BookmarksResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

#### Get Bookmarks by Folder

- *Endpoint*: `GET /api/bookmarks/folder/{parent_id}`
- *Summary*: Get all Bookmarks in a specific folder
- *Authentication*: Required (`Authorization: Bearer token`)
- *Path Parameters*:
  - `parent_id` (string, required): Parent Folder ID
- *Response*: BookmarksResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

#### Delete Bookmark

- *Endpoint*: `DELETE /api/bookmarks/folder/{parent_id}`
- *Summary*: Delete a Bookmark from folder
- *Authentication*: Required (`Authorization: Bearer token`)
- *Path Parameters*:
  - `parent_id` (string, required): Parent Folder ID
- *Response*: BookmarksResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

### Folder Management


#### Get Root Folders

- *Endpoint*: `GET /api/folders`
- *Summary*: Get the Root Folders associated with the user
- *Description*: Returns root folders and their children
- *Authentication*: Required (`Authorization: Bearer token`)
- *Response*: FoldersResponse (array of FolderData)
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

#### Create Folder

- *Endpoint*: `POST /api/folders`
- *Summary*: Create a new Folder
- *Authentication*: Required (`Authorization: Bearer token`)
- *Request Body*: CreateFolderParams
```json
       {
         "description": "string",
         "name": "string",
         "parent_id": "string",
         "user_id": "string"
       }
```
- *Response*: FolderResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

#### Move Folder

- *Endpoint*: `PATCH /api/folders`
- *Summary*: Move a Folder to different parent
- *Authentication*: Required (`Authorization: Bearer token`)
- *Request Body*: MoveFolderRequest
```json
       {
         "folderId": "string",
         "newParentId": "string",
         "userId": "string"
       }
```
- *Response*: FolderResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 403: Forbidden
  - 404: Not Found
  - 500: Internal Server Error

#### Get Folder by ID

- *Endpoint*: `GET /api/folders/{folder_id}`
- *Summary*: Get Folder and its children by ID
- *Authentication*: Required (`Authorization: Bearer token`)
- *Path Parameters*:
  - `folder_id` (string, required): Folder ID
- *Response*: FolderResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 404: Not Found
  - 500: Internal Server Error

#### Delete Folder

- *Endpoint*: `DELETE /api/folders/{folder_id}`
- *Summary*: Delete a Folder (cascade delete)
- *Authentication*: Required (`Authorization: Bearer token`)
- *Path Parameters*:
  - `folder_id` (string, required): Folder ID
- *Response*: FolderResponse
- *Status Codes*:
  - 200: OK
  - 404: Not Found
  - 500: Internal Server Error

### User Management


#### User Signup

- *Endpoint*: `POST /api/users/signup`
- *Summary*: Signup to the app
- *Authentication*: Not required
- *Request Body*: NewUserRequest
```json
       {
         "email": "string",
         "isAdmin": false,
         "password": "string",
         "username": "string"
       }
```
- *Response*: LoginResponse (includes JWT token)
- *Status Codes*:
  - 200: OK
  - 400: Bad Request
  - 500: Internal Server Error

#### User Login

- *Endpoint*: `POST /api/users/login`
- *Summary*: Login with email/username and password
- *Authentication*: Not required
- *Request Body*: LoginRequest
```json
       {
         "email": "string",
         "password": "string",
         "username": "string"
       }
```
- *Response*: LoginResponse (includes JWT token)
- *Status Codes*:
  - 200: OK
  - 400: Bad Request
  - 401: Unauthorized
  - 500: Internal Server Error

#### Get Current User

- *Endpoint*: `GET /api/users/current`
- *Summary*: Get Current User from JWT claims
- *Authentication*: Required (`authorization: Bearer token`)
- *Response*: UserResponse
- *Status Codes*:
  - 200: OK
  - 400: Bad Request
  - 401: Unauthorized

#### Update User Credentials

- *Endpoint*: `POST /api/users/creds`
- *Summary*: Update User Credentials (must be same user)
- *Authentication*: Required (`Authorization: Bearer token`)
- *Request Body*: UpdateCredentialsRequest
```json
       {
         "email": "string",
         "id": "string",
         "oldPassword": "string",
         "password": "string",
         "username": "string"
       }
```
- *Response*: UpdateUserResponse (includes new JWT)
- *Status Codes*:
  - 200: OK
  - 400: Bad Request
  - 401: Unauthorized
  - 500: Internal Server Error

#### Get All Users (Admin)

- *Endpoint*: `GET /api/users/`
- *Summary*: Get All Users (Admin only)
- *Authentication*: Required (`Authorization: Bearer token` - Admin privileges)
- *Response*: UsersResponse (array of UserData)
- *Status Codes*:
  - 200: OK
  - 403: Forbidden
  - 404: Not Found
  - 500: Internal Server Error

#### Delete User (Admin)

- *Endpoint*: `DELETE /api/users/{user_id}`
- *Summary*: Delete User by UUID (Admin only)
- *Authentication*: Required (`Authorization: Bearer token` - Admin privileges)
- *Path Parameters*:
  - `user_id` (string, required): User ID (UUID)
- *Response*: DeleteUserResponse
- *Status Codes*:
  - 200: OK

### Profile Management


#### Upload Profile Picture

- *Endpoint*: `POST /api/users/profile`
- *Summary*: Upload profile picture file
- *Authentication*: Required (`authorization: Bearer token`)
- *Content Type*: `multipart/form-data`
- *Form Parameters*:
  - `file` (file, required): Profile picture file
- *Response*: String (file URL)
- *Status Codes*:
  - 200: OK
  - 400: Bad Request
  - 404: Not Found

#### Get Profile Picture

- *Endpoint*: `GET /api/users/profile`
- *Summary*: Get User Profile Picture URL
- *Authentication*: Required (`Authorization: Bearer token`)
- *Response*: StringResponse
- *Status Codes*:
  - 200: OK
  - 401: Unauthorized
  - 403: Forbidden
  - 500: Internal Server Error

## Data Models


### Request Models


#### NewUserRequest

```json
     {
       "email": "string",
       "isAdmin": "boolean",
       "password": "string", 
       "username": "string"
     }
```

#### LoginRequest

```json
     {
       "email": "string",
       "password": "string",
       "username": "string"
     }
```

#### UpdateCredentialsRequest

```json
     {
       "email": "string",
       "id": "string",
       "oldPassword": "string",
       "password": "string",
       "username": "string"
     }
```

#### CreateFolderParams

```json
     {
       "description": "string",
       "name": "string",
       "parent_id": "string",
       "user_id": "string"
     }
```

#### MoveFolderRequest

```json
     {
       "folderId": "string",
       "newParentId": "string",
       "userId": "string"
     }
```

#### CreateBookmarkParams

```json
     {
       "folder_id": "string",
       "link": "string",
       "name": "string",
       "user_id": "string"
     }
```

#### MoveBookmarkParams

```json
     {
       "folder_id": "string",
       "id": "string"
     }
```


### Response Models


#### LoginResponse

```json
     {
       "data": "UserData",
       "jwt": "string",
       "message": "string", 
       "success": "boolean"
     }
```

#### UserResponse

```json
     {
       "data": "UserData",
       "message": "string",
       "success": "boolean"
     }
```

#### UsersResponse

```json
     {
       "data": ["UserData"],
       "message": "string",
       "success": "boolean"
     }
```

#### UpdateUserResponse

```json
     {
       "data": "UserData",
       "jwt": "string",
       "message": "string",
       "success": "boolean"
     }
```

#### DeleteUserResponse

```json
     {
       "data": "string",
       "message": "string",
       "success": "boolean"
     }
```

#### FolderResponse

```json
     {
       "data": "FolderData",
       "message": "string",
       "success": "boolean"
     }
```

#### FoldersResponse

```json
     {
       "data": ["FolderData"],
       "message": "string",
       "success": "boolean"
     }
```

#### BookmarkResponse

```json
     {
       "data": "Bookmark",
       "message": "string",
       "success": "boolean"
     }
```

#### BookmarksResponse

```json
     {
       "data": ["Bookmark"],
       "message": "string",
       "success": "boolean"
     }
```

#### StringResponse

```json
     {
       "data": "string",
       "message": "string",
       "success": "boolean"
     }
```

### Entity Models


#### UserData

```json
     {
       "admin": "boolean",
       "created_datetime": "string",
       "email": "string",
       "id": "string",
       "profile_pic_url": "string",
       "updated_datetime": "string",
       "username": "string"
     }
```

#### Folder

```json
     {
       "created_datetime": "string",
       "description": "string",
       "id": "string",
       "name": "string",
       "parent_id": "string",
       "updated_datetime": "string",
       "user_id": "string"
     }
```

#### FolderData (Extended)

```json
     {
       "bookmarks": ["Bookmark"],
       "children": ["Folder"],
       "created_datetime": "string",
       "description": "string",
       "id": "string",
       "name": "string",
       "notes": ["Note"],
       "parent_id": "string",
       "updated_datetime": "string",
       "user_id": "string"
     }
```

#### Bookmark

```json
     {
       "created_datetime": "string",
       "description": "string",
       "folder_id": "string",
       "icon": "string",
       "id": "string",
       "link": "string",
       "name": "string",
       "updated_datetime": "string",
       "user_id": "string"
     }
```

#### Note

```json
     {
       "content": "string",
       "created_datetime": "string", 
       "description": "string",
       "folder_id": "string",
       "id": "string",
       "name": "string",
       "updated_datetime": "string",
       "user_id": "string"
     }
```


## Authentication Flow


### User Registration

    1. `POST /api/users/signup` with NewUserRequest
    2. Receive LoginResponse with JWT token
    3. Use JWT token for authenticated requests

### User Login

    1. `POST /api/users/login` with LoginRequest (email/username + password)
    2. Receive LoginResponse with JWT token
    3. Use JWT token for authenticated requests

### Authenticated Requests

    Include JWT token in Authorization header:
```
    Authorization: Bearer <jwt_token>
```

## Error Handling


### Common HTTP Status Codes

- *200 OK*: Request successful
- *400 Bad Request*: Invalid request data
- *401 Unauthorized*: Missing or invalid JWT token
- *403 Forbidden*: Insufficient privileges (admin required)
- *404 Not Found*: Resource not found
- *500 Internal Server Error*: Server error

### Error Response Format

    All error responses follow the standard response format:
```json
    {
      "data": null,
      "message": "Error description",
      "success": false
    }
```

## Usage Examples


### Authentication Example

```javascript
    // Signup
    const signupResponse = await fetch('/api/users/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: 'user@example.com',
        username: 'username',
        password: 'password123',
        isAdmin: false
      })
    });
    const { jwt } = await signupResponse.json();

    // Use JWT for authenticated requests
    const foldersResponse = await fetch('/api/folders', {
      headers: { 'Authorization': `Bearer ${jwt}` }
    });
```

### Folder Operations Example

```javascript
    // Create folder
    const createFolderResponse = await fetch('/api/folders', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${jwt}`
      },
      body: JSON.stringify({
        name: 'My Folder',
        description: 'Personal bookmarks',
        parent_id: null,
        user_id: 'user-uuid'
      })
    });

    // Get root folders
    const foldersResponse = await fetch('/api/folders', {
      headers: { 'Authorization': `Bearer ${jwt}` }
    });
```

### Bookmark Operations Example

```javascript
    // Create bookmark
    const createBookmarkResponse = await fetch('/api/bookmarks', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${jwt}`
      },
      body: JSON.stringify({
        name: 'Example Site',
        link: 'https://example.com',
        folder_id: 'folder-uuid',
        user_id: 'user-uuid'
      })
    });

    // Get bookmarks in folder
    const bookmarksResponse = await fetch(`/api/bookmarks/folder/${folderId}`, {
      headers: { 'Authorization': `Bearer ${jwt}` }
    });
```

## Notes


### API Design Patterns

- All endpoints follow RESTful conventions
- Consistent response format with `data`, `message`, `success` fields
- JWT-based authentication for protected endpoints
- UUID identifiers for all entities
- Timestamps in ISO format

### Known Issues

- Background upload endpoint (`PUT /api/user/background`) lacks authentication
- Some endpoint descriptions could be more detailed
- Missing pagination for list endpoints

### Future Enhancements

- Pagination support for large data sets
- Rate limiting documentation
- WebSocket endpoints for real-time updates
- API versioning strategy
- Comprehensive error code documentation
