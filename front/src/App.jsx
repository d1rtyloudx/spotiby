import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AuthContainer from './components/Auth/AuthContainer';
import Profile from './components/Profile';
import ProfilesList from './components/ProfilesList';
import ProfileById from "./components/ProfileById.jsx";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/auth" element={<AuthContainer />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/profiles" element={<ProfilesList />} />
                <Route path="/" element={<AuthContainer />} />
                <Route path="/profile/:id" element={<ProfileById />} />
            </Routes>
        </Router>
    );
}

export default App;