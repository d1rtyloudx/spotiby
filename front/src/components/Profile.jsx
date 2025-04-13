import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/Profile.css';

function Profile() {
    const navigate = useNavigate();
    const [profile, setProfile] = useState(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [isEditing, setIsEditing] = useState(false);
    const [editForm, setEditForm] = useState({
        display_name: '',
        first_name: '',
        last_name: '',
    });
    const [editError, setEditError] = useState('');

    useEffect(() => {
        const storedProfile = localStorage.getItem('profile');
        if (storedProfile) {
            const parsedProfile = JSON.parse(storedProfile);
            setProfile(parsedProfile);
            setEditForm({
                display_name: parsedProfile.display_name || '',
                first_name: parsedProfile.first_name || '',
                last_name: parsedProfile.last_name || '',
            });
        } else {
            navigate('/auth');
        }
    }, [navigate]);

    const handleLogout = () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('profile');
        navigate('/auth');
    };

    const handleSearch = (e) => {
        e.preventDefault();
        console.log('Search query:', searchQuery);
    };

    const openEditModal = () => {
        setIsEditing(true);
        setEditError('');
    };

    const closeEditModal = () => {
        setIsEditing(false);
        setEditError('');
    };

    const handleEditChange = (e) => {
        setEditForm({ ...editForm, [e.target.name]: e.target.value });
    };

    const handleAvatarChange = async (e) => {
        const file = e.target.files[0];
        if (file) {
            try {
                const formData = new FormData();
                formData.append('file', file);
                const uploadResponse = await axios.post(
                    `${API_BASE_URL}/api/v1/image/upload/profile/${profile.id}`,
                    formData,
                    {
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                            'Content-Type': 'multipart/form-data',
                        },
                    }
                );

                const newAvatarUrl = uploadResponse.data.avatar_url;
                const updatedProfile = { ...profile, avatar_url: newAvatarUrl };
                setProfile(updatedProfile);
                localStorage.setItem('profile', JSON.stringify(updatedProfile));

                // Обновляем профиль на сервере (если нужно)
                await axios.put(
                    `${API_BASE_URL}/api/v1/profile/me`,
                    { avatar_url: newAvatarUrl },
                    {
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                        },
                    }
                );
            } catch (err) {
                console.error('Avatar upload error:', err);
            }
        }
    };

    const handleEditSubmit = async (e) => {
        e.preventDefault();
        setEditError('');

        try {
            const response = await axios.put(
                `${API_BASE_URL}/api/v1/profile`,
                editForm,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                }
            );

            console.log('Profile updated:', response.data);
            localStorage.setItem('profile', JSON.stringify(response.data.profile));
            setProfile(response.data.profile);
            setIsEditing(false);
        } catch (err) {
            console.error('Edit profile error:', err);
            setEditError(err.response?.data?.message || 'Failed to update profile');
        }
    };

    if (!profile) {
        return <div>Loading...</div>;
    }

    return (
        <div className="profile-page">
            {/* Боковая панель */}
            <aside className="sidebar">
                <div className="sidebar-header">
                    <h2>Your Playlists</h2>
                </div>
                <ul className="playlist-list">
                    <li>Playlist 1</li>
                    <li>Playlist 2</li>
                    <li>Playlist 3</li>
                    <li className="create-playlist">+ Create Playlist</li>
                </ul>
            </aside>

            {/* Основной контент */}
            <main className="main-content">
                {/* Хедер */}
                <header className="header">
                    <form onSubmit={handleSearch} className="search-bar">
                        <input
                            type="text"
                            placeholder="Search for playlists, profiles, or tracks..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="search-input"
                        />
                        <button type="submit" className="search-button">
                            Search
                        </button>
                    </form>
                    <button onClick={handleLogout} className="logout-button">
                        Log Out
                    </button>
                </header>

                {/* Профиль */}
                <section className="profile-section">
                    <div className="profile-header">
                        <div className="avatar-wrapper">
                            {profile.avatar_url ? (
                                <img
                                    src={profile.avatar_url}
                                    alt="Avatar"
                                    className="profile-avatar"
                                />
                            ) : (
                                <div className="profile-avatar-placeholder">
                                    {profile.display_name?.charAt(0) || 'U'}
                                </div>
                            )}
                            <div className="avatar-overlay">
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleAvatarChange}
                                    className="avatar-input"
                                    id="avatar-upload"
                                />
                                <label htmlFor="avatar-upload" className="avatar-label">
                                    Change Photo
                                </label>
                            </div>
                        </div>
                        <div className="profile-info">
                            <h1>{profile.display_name || 'User'}</h1>
                            <p>
                                <strong>First Name:</strong> {profile.first_name || 'Not set'}
                            </p>
                            <p>
                                <strong>Last Name:</strong> {profile.last_name || 'Not set'}
                            </p>
                            <button onClick={openEditModal} className="edit-button">
                                Edit Profile
                            </button>
                        </div>
                    </div>
                </section>

                {/* Панели */}
                <section className="content-panels">
                    <div className="panel">
                        <h3>Subscribed Playlists</h3>
                        <ul className="panel-list">
                            <li>Subscribed Playlist 1</li>
                            <li>Subscribed Playlist 2</li>
                            <li>Subscribed Playlist 3</li>
                        </ul>
                    </div>
                    <div className="panel">
                        <h3>Followed Profiles</h3>
                        <ul className="panel-list">
                            <li>Profile 1</li>
                            <li>Profile 2</li>
                            <li>Profile 3</li>
                        </ul>
                    </div>
                    <div className="panel">
                        <h3>Top Tracks</h3>
                        <ul className="panel-list">
                            <li>Track 1</li>
                            <li>Track 2</li>
                            <li>Track 3</li>
                        </ul>
                    </div>
                </section>
            </main>

            {/* Модал редактирования */}
            {isEditing && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Edit Profile</h2>
                        <form onSubmit={handleEditSubmit}>
                            <div className="form-group">
                                <label>Display Name</label>
                                <input
                                    type="text"
                                    name="display_name"
                                    value={editForm.display_name}
                                    onChange={handleEditChange}
                                    placeholder="Enter display name"
                                />
                            </div>
                            <div className="form-group">
                                <label>First Name</label>
                                <input
                                    type="text"
                                    name="first_name"
                                    value={editForm.first_name}
                                    onChange={handleEditChange}
                                    placeholder="Enter first name"
                                />
                            </div>
                            <div className="form-group">
                                <label>Last Name</label>
                                <input
                                    type="text"
                                    name="last_name"
                                    value={editForm.last_name}
                                    onChange={handleEditChange}
                                    placeholder="Enter last name"
                                />
                            </div>
                            {editError && <p className="error">{editError}</p>}
                            <div className="modal-buttons">
                                <button type="submit" className="save-button">
                                    Save
                                </button>
                                <button
                                    type="button"
                                    onClick={closeEditModal}
                                    className="cancel-button"
                                >
                                    Cancel
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}

export default Profile;