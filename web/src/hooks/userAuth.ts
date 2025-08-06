import { useState, useEffect } from "react";
import api from "../lib/api";

interface User {
    username: string;
    sessionToken: string;
    csrfToken: string;
}

interface UserAuthResponse {
    user: User | null;
    isAuthenticated: boolean;
    loading: boolean;
    error: string | null;
}

// is user authenticated?
export function userAuth(): UserAuthResponse {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await api.post('/auth');
                setUser(response.data);
                setIsAuthenticated(true);
            } catch (err) {
                setUser(null);
                setIsAuthenticated(false);
                setError(err instanceof Error ? err.message : 'An error occurred');
            } finally {
                setLoading(false);
            }
        }
        fetchUser();
    }, []);

    return { user, isAuthenticated, loading, error };
}