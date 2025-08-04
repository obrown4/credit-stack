import React, { useState } from 'react';
import { useAuth } from '../AuthContext';

const LoginForm: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isRegistering, setIsRegistering] = useState(false);
  const { login, register } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    const success = isRegistering 
      ? await register(username, password)
      : await login(username, password);
      
    if (!success) {
      setError(isRegistering 
        ? 'Registration failed. Username might already exist or password is too short.'
        : 'Invalid username or password'
      );
    }
    setIsLoading(false);
  };

  return (
    <div className="login-container">
      <h2>{isRegistering ? 'Create Account' : 'Sign In'}</h2>
      <form onSubmit={handleSubmit} className="login-form">
        <div className="form-group">
          <label htmlFor="username">Username:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            minLength={8}
            placeholder="Enter username (min 8 characters)"
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            placeholder="Enter password (min 8 characters)"
          />
        </div>
        {error && <div className="error-message">{error}</div>}
        <button type="submit" disabled={isLoading}>
          {isLoading 
            ? (isRegistering ? 'Creating account...' : 'Signing in...') 
            : (isRegistering ? 'Create Account' : 'Sign In')
          }
        </button>
      </form>
      <div className="form-toggle">
        <button 
          type="button" 
          onClick={() => {
            setIsRegistering(!isRegistering);
            setError('');
          }}
          className="toggle-btn"
        >
          {isRegistering 
            ? 'Already have an account? Sign In' 
            : 'Need an account? Register'
          }
        </button>
      </div>
    </div>
  );
};

export default LoginForm; 