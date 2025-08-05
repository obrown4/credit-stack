# Credit Stack

A browser extension that helps optimize credit card rewards for online purchases.

## Features

- **Authentication System**: Secure login/registration with Go backend
- **Shopping Site Detection**: Automatically detects when you're on shopping websites
- **Credit Card Optimization**: Analyzes purchases and recommends optimal credit cards
- **Floating UI**: Non-intrusive floating button for easy access

## Setup

### Backend (Go)

1. Navigate to the backend directory:
   ```bash
   cd /path/to/credit-stack
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up your environment variables in a `.env` file:
   ```
   MONGODB_URI=your_mongodb_connection_string
   ```

4. Run the backend server:
   ```bash
   go run cmd/creditStack/main.go
   ```

   The server will start on `http://localhost:8080`

### Frontend (Browser Extension)

1. Navigate to the web directory:
   ```bash
   cd web
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Build the extension:
   ```bash
   npm run build
   ```

4. Load the extension in Chrome:
   - Open Chrome and go to `chrome://extensions/`
   - Enable "Developer mode"
   - Click "Load unpacked"
   - Select the `web/dist` directory

## Testing the Authentication

1. **Start the backend server** (see setup instructions above)

2. **Load the extension** in Chrome

3. **Navigate to any shopping website** (e.g., amazon.com, walmart.com)

4. **You should see an authentication modal** with:
   - Username and password fields
   - "Sign In" button
   - "Register here" link

5. **Test Registration**:
   - Click "Register here"
   - Enter a username (at least 8 characters)
   - Enter a password (at least 8 characters)
   - Click "Create Account"
   - You should see a success message

6. **Test Login**:
   - Use the credentials you just created
   - Click "Sign In"
   - You should see a success message and the modal should close

7. **Test the Optimization Feature**:
   - After successful login, you should see a floating "💳 [Username]" button
   - Click it to see the optimization modal with credit card recommendations
   - The modal should show your username and a logout button

8. **Test Logout**:
   - Click the logout button in the optimization modal
   - You should be logged out and the floating button should disappear
   - Navigate to another shopping site to see the auth modal again

## API Endpoints

- `POST /register` - User registration
- `POST /login` - User login
- `POST /logout` - User logout
- `POST /protected` - Protected endpoint (requires authentication)

## Architecture

- **Backend**: Go with MongoDB for data storage
- **Frontend**: TypeScript browser extension
- **Authentication**: Session-based with CSRF protection
- **CORS**: Configured to allow cross-origin requests from the extension

## Security Features

- Password hashing with bcrypt
- Session tokens with expiration
- CSRF protection
- Secure cookie handling
- Input validation and sanitization
