import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

interface User {
    username: string;
    sessionToken: string;
    csrfToken: string;
}

interface AuthContextType {
    user: User | null;
    login: (username: string, password: string) => Promise<boolean>;
    logout: () => void;
    register: (username: string, password: string) => Promise<boolean>;
    isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Check for existing session on app load
    useEffect(() => {
        const checkAuthStatus = async () => {
            const sessionToken = localStorage.getItem('sessionToken');
            const csrfToken = localStorage.getItem('csrfToken');
            const username = localStorage.getItem('username');

            if (sessionToken && csrfToken && username) {
                // You can add session validation here if needed
                setUser({ username, sessionToken, csrfToken });
            }
            setIsLoading(false);
        };

        checkAuthStatus();
    }, []);

    const login = async (username: string, password: string): Promise<boolean> => {
        try {
            const formData = new FormData();
            formData.append('username', username);
            formData.append('password', password);

            const response = await fetch('/login', {
                method: 'POST',
                body: formData,
            });

            if (response.ok) {
                const data = await response.json();
                
                // Store tokens in localStorage for persistence
                localStorage.setItem('sessionToken', data.session_token);
                localStorage.setItem('csrfToken', data.csrf_token);
                localStorage.setItem('username', data.username);

                setUser({
                    username: data.username,
                    sessionToken: data.session_token,
                    csrfToken: data.csrf_token,
                });
                return true;
            } else {
                const errorData = await response.text();
                console.error('Login failed:', errorData);
                return false;
            }
        } catch (error) {
            console.error('Login error:', error);
            return false;
        }
    };

    const register = async (username: string, password: string): Promise<boolean> => {
        try {
            const formData = new FormData();
            formData.append('username', username);
            formData.append('password', password);

            const response = await fetch('/register', {
                method: 'POST',
                body: formData,
            });

            if (response.ok) {
                return true;
            } else {
                const errorData = await response.text();
                console.error('Registration failed:', errorData);
                return false;
            }
        } catch (error) {
            console.error('Registration error:', error);
            return false;
        }
    };

    const logout = async () => {
        if (user) {
            try {
                const formData = new FormData();
                formData.append('username', user.username);
                formData.append('session_token', user.sessionToken);

                await fetch('/logout', {
                    method: 'POST',
                    body: formData,
                });
            } catch (error) {
                console.error('Logout error:', error);
            }
        }

        // Clear local storage and state
        localStorage.removeItem('sessionToken');
        localStorage.removeItem('csrfToken');
        localStorage.removeItem('username');
        setUser(null);
    };

    const value = {
        user,
        login,
        logout,
        register,
        isLoading,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};