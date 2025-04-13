import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../../config/api';

function Login({ onSwitch }) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        if (!username || !password) {
            setError('Please fill in all fields');
            setLoading(false);
            return;
        }

        try {
            console.log('Sending login:', { username, password });
            const response = await axios.post(`${API_BASE_URL}/api/v1/auth/login`, {
                username,
                password,
            });

            console.log('Login success:', response.data);
            localStorage.setItem('access_token', response.data.tokens.access_token);
            localStorage.setItem('refresh_token', response.data.tokens.refresh_token);
            localStorage.setItem('profile', JSON.stringify(response.data.profile));

            navigate('/profile');
        } catch (err) {
            setLoading(false);
            console.error('Login error:', err);
            setError(err.response?.data?.message || 'Invalid username or password');
        }
    };

    return (
        <div className="auth-form">
            <h2>Log In</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>Username</label>
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        placeholder="Enter your username"
                        disabled={loading}
                    />
                </div>
                <div className="form-group">
                    <label>Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Enter your password"
                        disabled={loading}
                    />
                </div>
                {error && <p className="error">{error}</p>}
                <button type="submit" disabled={loading}>
                    {loading ? 'Logging in...' : 'Log In'}
                </button>
            </form>
            <p>
                Don’t have an account?{' '}
                <span className="switch-link" onClick={() => onSwitch('register')}>
          Sign Up
        </span>
            </p>
        </div>
    );
}

export default Login;