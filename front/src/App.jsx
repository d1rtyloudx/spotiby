import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AuthContainer from './components/Auth/AuthContainer';
import Profile from './components/Profile';

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/auth" element={<AuthContainer />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/" element={<AuthContainer />} />
            </Routes>
        </Router>
    );
}

export default App;