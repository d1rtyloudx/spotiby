import { useState } from 'react';
import Login from './Login';
import Register from './Register';
import '../../styles/Auth.css';

function AuthContainer() {
    const [view, setView] = useState('login');

    return (
        <div className="auth-container">
            {view === 'login' ? (
                <Login onSwitch={setView} />
            ) : (
                <Register onSwitch={setView} />
            )}
        </div>
    );
}

export default AuthContainer;