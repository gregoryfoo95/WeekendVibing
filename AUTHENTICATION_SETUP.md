# FitHero Authentication Setup Guide

This guide explains how to set up and use the comprehensive Google OAuth authentication system implemented for FitHero.

## 🎯 Overview

The authentication system includes:
- **Backend**: Google OAuth 2.0 with JWT tokens stored in secure HTTP-only cookies
- **Frontend**: React authentication context with protected routes
- **Security**: User-based authorization ensuring users can only access their own data

## 🔧 Backend Setup

### 1. Environment Configuration

Copy the configuration template:
```bash
cp backend/config.template backend/.env
```

Fill in your Google OAuth credentials in `backend/.env`:
```env
# Google OAuth Configuration
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_at_least_32_characters_long
JWT_EXPIRATION_HOURS=24

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=fithero_db

# Cookie Configuration
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false  # Set to true in production with HTTPS
COOKIE_SAME_SITE=lax

# Server Configuration
PORT=8080
GIN_MODE=debug
```

### 2. Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing project
3. Enable Google+ API
4. Go to "Credentials" and create OAuth 2.0 Client ID
5. Set authorized redirect URI: `http://localhost:8080/api/auth/google/callback`
6. Copy Client ID and Client Secret to your `.env` file

### 3. Database Setup

Start PostgreSQL and create the database:
```sql
CREATE DATABASE fithero_db;
CREATE USER your_db_user WITH PASSWORD 'your_db_password';
GRANT ALL PRIVILEGES ON DATABASE fithero_db TO your_db_user;
```

### 4. Run Backend

```bash
cd backend
go mod tidy
go run main.go
```

The backend will:
- Auto-migrate database tables
- Start server on port 8080
- Set up authentication endpoints

## 🎨 Frontend Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Environment Configuration

Create `frontend/.env`:
```env
VITE_API_URL=http://localhost:8080
```

### 3. Run Frontend

```bash
npm run dev
```

The frontend will start on port 3000 with:
- Authentication-aware navigation
- Protected routes
- Google OAuth integration

## 🔐 Authentication Flow

### User Authentication Process

1. **Login Initiation**: User clicks "Sign In with Google"
2. **OAuth Redirect**: User is redirected to Google OAuth
3. **Authorization**: User authorizes the application
4. **Callback Handling**: Backend processes OAuth callback
5. **User Creation/Login**: Backend finds or creates user account
6. **JWT Generation**: Backend generates JWT token
7. **Cookie Setting**: Token is stored in secure HTTP-only cookie
8. **Frontend Update**: Frontend authentication context updates

### API Authentication

#### Protected Endpoints
All protected endpoints require authentication via:
- `Authorization: Bearer <token>` header, OR
- `auth_token` HTTP-only cookie

#### Authorization Checks
- Users can only access their own data
- Resource-level authorization in repositories and services
- Middleware validates JWT tokens and sets user context

## 🛡️ Security Features

### JWT Tokens
- Secure generation with configurable expiration
- HMAC-SHA256 signing
- Claims include user ID, email, username

### HTTP-Only Cookies
- `HttpOnly`: Prevents XSS attacks
- `Secure`: HTTPS only in production
- `SameSite`: CSRF protection

### Authorization
- User ownership validation for all protected resources
- Repository-level access control
- Service-layer business logic protection

### Input Validation
- Request validation using go-playground/validator
- Type-safe TypeScript interfaces
- Comprehensive error handling

## 📡 API Endpoints

### Authentication Endpoints
```
GET  /api/auth/google              - Initiate Google OAuth
GET  /api/auth/google/callback     - Handle OAuth callback
POST /api/auth/logout              - Logout user
GET  /api/auth/check               - Check authentication status
POST /api/auth/refresh             - Refresh JWT token
GET  /api/me                       - Get current user profile
```

### Protected Endpoints
```
GET  /api/profile                  - Get current user profile
PUT  /api/profile                  - Update current user profile
GET  /api/tasks/daily              - Get user's daily tasks
POST /api/tasks/daily/generate     - Generate daily tasks
POST /api/tasks/daily/:id/complete - Complete a task
GET  /api/achievements/user        - Get user's achievements
POST /api/achievements/:id/unlock  - Unlock achievement
```

### Public Endpoints
```
GET /api/public/tasks              - Get all available tasks
GET /api/public/achievements       - Get all achievements
```

## 🧪 Testing Authentication

### Backend Testing
```bash
# Check health
curl http://localhost:8080/health

# Test authentication (should return 401)
curl http://localhost:8080/api/profile

# Check public endpoint
curl http://localhost:8080/api/public/tasks
```

### Frontend Testing
1. Open http://localhost:3000
2. Click "Sign In with Google"
3. Complete OAuth flow
4. Verify protected routes are accessible
5. Check navigation shows user info

## 🔄 Development Workflow

### Making Changes

#### Backend Changes
1. Update Go code
2. Run `go run main.go` to restart server
3. Database migrations run automatically

#### Frontend Changes
1. Update React components
2. Development server auto-reloads
3. TypeScript compilation is automatic

### Adding New Protected Endpoints

1. **Backend**: Add to protected route group in `main.go`
2. **Service**: Implement authorization checks
3. **Frontend**: Add to API client with authentication

### Adding New Public Endpoints

1. **Backend**: Add to public route group in `main.go`
2. **Frontend**: Add to API client (no auth required)

## 🚀 Deployment Considerations

### Environment Variables
- Set `COOKIE_SECURE=true` for HTTPS
- Use strong `JWT_SECRET` (32+ characters)
- Configure proper `GOOGLE_REDIRECT_URL` for domain

### Database
- Use connection pooling for production
- Set up proper backup and migration strategy
- Consider read replicas for scaling

### Security
- Enable HTTPS/TLS
- Set up rate limiting
- Configure CORS for production domains
- Use environment-specific OAuth credentials

## 🎮 User Experience

### Authentication States
- **Unauthenticated**: Show login button, public content only
- **Authenticated**: Show user avatar, points, level, logout option
- **Loading**: Show loading indicators during auth checks

### Protected Features
- Dashboard with personalized content
- Daily task management
- Achievement unlocking
- Profile management
- User-specific data

### Error Handling
- Graceful fallbacks for authentication failures
- Clear error messages for users
- Automatic retry mechanisms where appropriate

## 🔍 Troubleshooting

### Common Issues

1. **OAuth Redirect Mismatch**
   - Verify `GOOGLE_REDIRECT_URL` matches Google Console
   - Check for trailing slashes

2. **Cookie Issues**
   - Verify `COOKIE_DOMAIN` matches your domain
   - Check browser developer tools for cookie settings

3. **CORS Errors**
   - Verify frontend URL in backend CORS configuration
   - Check `withCredentials: true` in API client

4. **Database Connection**
   - Verify PostgreSQL is running
   - Check database credentials and permissions

5. **JWT Validation Errors**
   - Verify `JWT_SECRET` consistency
   - Check token expiration settings

This authentication system provides a robust foundation for secure user management while maintaining excellent user experience. The modular architecture makes it easy to extend and maintain as your application grows. 