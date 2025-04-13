import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/Profile.css';

const LoadingSpinner = () => (
    <div className="loading-spinner">
        <div className="spinner"></div>
        <p>Updating profile...</p>
    </div>
);

function Profile() {
    const navigate = useNavigate();
    const [profile, setProfile] = useState(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [isEditing, setIsEditing] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [editForm, setEditForm] = useState({
        display_name: '',
        first_name: '',
        last_name: '',
    });
    const [editError, setEditError] = useState('');
    const [avatarSrc, setAvatarSrc] = useState('');
    const [subscriptionAvatarUrls, setSubscriptionAvatarUrls] = useState({});

    const [playlistsPage, setPlaylistsPage] = useState(1);
    const [subscriptions, setSubscriptions] = useState([]);
    const [subscriptionsPage, setSubscriptionsPage] = useState(1);
    const [topTracksPage, setTopTracksPage] = useState(1);
    const itemsPerPage = 4;

    const playlists = ['Плейлист 1', 'Плейлист 2', 'Плейлист 3', 'Плейлист 4', 'Плейлист 5'];
    const topTracks = ['Трек 1', 'Трек 2', 'Трек 3', 'Трек 4', 'Трек 5', 'Трек 6'];

    const fetchAvatar = async (filename) => {
        if (!filename) {
            setAvatarSrc('');
            return;
        }

        try {
            const response = await axios.get(
                `${API_BASE_URL}/api/v1/image/profile/${filename}`,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                    responseType: 'blob',
                }
            );

            const imageUrl = URL.createObjectURL(response.data);
            setAvatarSrc(imageUrl);
        } catch (err) {
            console.error('Failed to fetch avatar:', err);
            setAvatarSrc('');
        }
    };

    const fetchSubscriptionAvatar = async (filename, profileId) => {
        if (!filename) {
            setSubscriptionAvatarUrls((prev) => ({ ...prev, [profileId]: '' }));
            return;
        }

        try {
            const response = await axios.get(
                `${API_BASE_URL}/api/v1/image/profile/${filename}`,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                    responseType: 'blob',
                }
            );

            const imageUrl = URL.createObjectURL(response.data);
            setSubscriptionAvatarUrls((prev) => ({ ...prev, [profileId]: imageUrl }));
        } catch (err) {
            console.error(`Failed to fetch avatar for subscription ${profileId}:`, err);
            setSubscriptionAvatarUrls((prev) => ({ ...prev, [profileId]: '' }));
        }
    };

    const fetchSubscriptions = async () => {
        try {
            let allProfiles = [];
            let currentPage = 1;
            let hasMore = true;

            while (hasMore) {
                const response = await axios.get(`${API_BASE_URL}/api/v1/profile/me/follows/`, {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                    params: {
                        page: currentPage,
                        limit: 1000,
                    },
                });

                const { profiles, pagination } = response.data.items || { profiles: [], pagination: {} };
                allProfiles = [...allProfiles, ...profiles];
                hasMore = pagination.has_more || false;
                currentPage += 1;
            }

            setSubscriptions(allProfiles);

            allProfiles.forEach((profile) => {
                if (profile.avatar_url) {
                    fetchSubscriptionAvatar(profile.avatar_url, profile.id);
                } else {
                    setSubscriptionAvatarUrls((prev) => ({ ...prev, [profile.id]: '' }));
                }
            });
        } catch (err) {
            console.error('Failed to fetch subscriptions:', err);
            setSubscriptions([]);
        }
    };

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

            fetchAvatar(parsedProfile.avatar_url);
            fetchSubscriptions();
        } else {
            navigate('/auth');
        }

        return () => {
            if (avatarSrc) {
                URL.revokeObjectURL(avatarSrc);
            }
            Object.values(subscriptionAvatarUrls).forEach((url) => {
                if (url) URL.revokeObjectURL(url);
            });
        };
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
            setIsLoading(true);
            try {
                const formData = new FormData();
                formData.append('file', file);

                const uploadResponse = await axios.post(
                    `${API_BASE_URL}/api/v1/image/profile/${profile.id}`,
                    formData,
                    {
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                            'Content-Type': 'multipart/form-data',
                        },
                    }
                );

                const newAvatarFilename = uploadResponse.data.image;
                const updatedProfile = { ...profile, avatar_url: newAvatarFilename };
                setProfile(updatedProfile);
                localStorage.setItem('profile', JSON.stringify(updatedProfile));

                await fetchAvatar(newAvatarFilename);
            } catch (err) {
                console.error('Avatar upload error:', err);
                setEditError('Не удалось загрузить аватар');
            } finally {
                setIsLoading(false);
            }
        }
    };

    const handleEditSubmit = async (e) => {
        e.preventDefault();
        setEditError('');
        setIsLoading(true);

        try {
            const response = await axios.put(
                `${API_BASE_URL}/api/v1/profile/me`,
                editForm,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                }
            );

            console.log('Profile updated:', response.data);
            const updatedProfile = response.data.profile;
            localStorage.setItem('profile', JSON.stringify(updatedProfile));
            setProfile(updatedProfile);
            setIsEditing(false);

            fetchAvatar(updatedProfile.avatar_url);
        } catch (err) {
            console.error('Edit profile error:', err);
            setEditError(err.response?.data?.message || 'Не удалось обновить профиль');
        } finally {
            setIsLoading(false);
        }
    };

    const getPaginatedItems = (items, page) => {
        const startIndex = (page - 1) * itemsPerPage;
        return items.slice(startIndex, startIndex + itemsPerPage);
    };

    const handleNextPage = (setPage, currentPage, totalItems) => {
        const maxPage = Math.ceil(totalItems / itemsPerPage);
        if (currentPage < maxPage) {
            setPage(currentPage + 1);
        }
    };

    const handlePrevPage = (setPage, currentPage) => {
        if (currentPage > 1) {
            setPage(currentPage - 1);
        }
    };

    const handleSubscriptionClick = (profileId) => {
        navigate(`/profile/${profileId}`);
    };

    if (!profile) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="profile-page">
            <aside className="sidebar">
                <div className="sidebar-header">
                    <h2>Ваши плейлисты</h2>
                </div>
                <ul className="playlist-list">
                    <li>Плейлист 1</li>
                    <li>Плейлист 2</li>
                    <li>Плейлист 3</li>
                    <li className="create-playlist">+ Создать плейлист</li>
                </ul>
            </aside>

            <main className="main-content">
                <header className="header">
                    <form onSubmit={handleSearch} className="search-bar">
                        <input
                            type="text"
                            placeholder="Поиск плейлистов, профилей или треков..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="search-input"
                        />
                        <button type="submit" className="search-button">
                            Поиск
                        </button>
                    </form>
                    <div className="header-buttons">
                        <button
                            onClick={() => navigate('/profiles')}
                            className="profiles-list-button"
                        >
                            Все профили
                        </button>
                        <button onClick={handleLogout} className="logout-button">
                            Выйти
                        </button>
                    </div>
                </header>

                <section className="profile-section">
                    <div className="profile-header-centered">
                        <div className="avatar-wrapper">
                            {avatarSrc ? (
                                <img
                                    src={avatarSrc}
                                    alt="Аватар"
                                    className="profile-avatar"
                                    onError={() => setAvatarSrc('')}
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
                                    disabled={isLoading}
                                />
                                <label htmlFor="avatar-upload" className="avatar-label">
                                    Изменить фото
                                </label>
                            </div>
                        </div>
                        <div className="profile-info-centered">
                            <h1>{profile.display_name || 'Пользователь'}</h1>
                            <p>
                                <strong>Имя:</strong> {profile.first_name || 'Не указано'}
                            </p>
                            <p>
                                <strong>Фамилия:</strong> {profile.last_name || 'Не указано'}
                            </p>
                            <button onClick={openEditModal} className="edit-button">
                                Редактировать профиль
                            </button>
                        </div>
                    </div>
                </section>

                <section className="content-lists">
                    <div className="list-section">
                        <h3>Подписанные плейлисты</h3>
                        <div className="horizontal-list-wrapper">
                            <ul className="horizontal-list">
                                {getPaginatedItems(playlists, playlistsPage).map((item, index) => (
                                    <li key={index} className="horizontal-item">
                                        {item}
                                    </li>
                                ))}
                            </ul>
                            <div className="pagination-controls">
                                <button
                                    onClick={() => handlePrevPage(setPlaylistsPage, playlistsPage)}
                                    disabled={playlistsPage === 1}
                                    className="pagination-button"
                                >
                                    ←
                                </button>
                                <button
                                    onClick={() =>
                                        handleNextPage(setPlaylistsPage, playlistsPage, playlists.length)
                                    }
                                    disabled={playlistsPage >= Math.ceil(playlists.length / itemsPerPage)}
                                    className="pagination-button"
                                >
                                    →
                                </button>
                            </div>
                        </div>
                    </div>
                    <div className="list-section">
                        <h3>Подписки на профили</h3>
                        <div className="horizontal-list-wrapper">
                            <ul className="horizontal-list">
                                {subscriptions.length > 0 ? (
                                    getPaginatedItems(subscriptions, subscriptionsPage).map((profile, index) => (
                                        <li
                                            key={profile.id || index}
                                            className="horizontal-item subscription-item"
                                            onClick={() => handleSubscriptionClick(profile.id)}
                                        >
                                            <div className="subscription-avatar">
                                                {subscriptionAvatarUrls[profile.id] ? (
                                                    <img
                                                        src={subscriptionAvatarUrls[profile.id]}
                                                        alt="Аватар"
                                                        className="avatar-img"
                                                        onError={() =>
                                                            setSubscriptionAvatarUrls((prev) => ({
                                                                ...prev,
                                                                [profile.id]: '',
                                                            }))
                                                        }
                                                    />
                                                ) : (
                                                    <div className="avatar-placeholder">
                                                        {profile.display_name?.charAt(0) || 'U'}
                                                    </div>
                                                )}
                                            </div>
                                            <p className="subscription-name">
                                                {profile.display_name || 'Пользователь'}
                                            </p>
                                        </li>
                                    ))
                                ) : (
                                    <li className="horizontal-item">Нет подписок</li>
                                )}
                            </ul>
                            <div className="pagination-controls">
                                <button
                                    onClick={() => handlePrevPage(setSubscriptionsPage, subscriptionsPage)}
                                    disabled={subscriptionsPage === 1}
                                    className="pagination-button"
                                >
                                    ←
                                </button>
                                <button
                                    onClick={() =>
                                        handleNextPage(
                                            setSubscriptionsPage,
                                            subscriptionsPage,
                                            subscriptions.length
                                        )
                                    }
                                    disabled={
                                        subscriptionsPage >= Math.ceil(subscriptions.length / itemsPerPage)
                                    }
                                    className="pagination-button"
                                >
                                    →
                                </button>
                            </div>
                        </div>
                    </div>
                    <div className="list-section">
                        <h3>Топ треков</h3>
                        <div className="horizontal-list-wrapper">
                            <ul className="horizontal-list">
                                {getPaginatedItems(topTracks, topTracksPage).map((item, index) => (
                                    <li key={index} className="horizontal-item">
                                        {item}
                                    </li>
                                ))}
                            </ul>
                            <div className="pagination-controls">
                                <button
                                    onClick={() => handlePrevPage(setTopTracksPage, topTracksPage)}
                                    disabled={topTracksPage === 1}
                                    className="pagination-button"
                                >
                                    ←
                                </button>
                                <button
                                    onClick={() =>
                                        handleNextPage(setTopTracksPage, topTracksPage, topTracks.length)
                                    }
                                    disabled={topTracksPage >= Math.ceil(topTracks.length / itemsPerPage)}
                                    className="pagination-button"
                                >
                                    →
                                </button>
                            </div>
                        </div>
                    </div>
                </section>
            </main>

            {isEditing && (
                <div className="modal">
                    <div className="modal-content">
                        {isLoading ? (
                            <LoadingSpinner />
                        ) : (
                            <>
                                <h2>Редактировать профиль</h2>
                                <form onSubmit={handleEditSubmit}>
                                    <div className="form-group">
                                        <label>Отображаемое имя</label>
                                        <input
                                            type="text"
                                            name="display_name"
                                            value={editForm.display_name}
                                            onChange={handleEditChange}
                                            placeholder="Введите отображаемое имя"
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label>Имя</label>
                                        <input
                                            type="text"
                                            name="first_name"
                                            value={editForm.first_name}
                                            onChange={handleEditChange}
                                            placeholder="Введите имя"
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label>Фамилия</label>
                                        <input
                                            type="text"
                                            name="last_name"
                                            value={editForm.last_name}
                                            onChange={handleEditChange}
                                            placeholder="Введите фамилию"
                                        />
                                    </div>
                                    {editError && <p className="error">{editError}</p>}
                                    <div className="modal-buttons">
                                        <button type="submit" className="save-button">
                                            Сохранить
                                        </button>
                                        <button
                                            type="button"
                                            onClick={closeEditModal}
                                            className="cancel-button"
                                        >
                                            Отмена
                                        </button>
                                    </div>
                                </form>
                            </>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}

export default Profile;