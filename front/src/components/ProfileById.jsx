import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/Profile.css';

const LoadingSpinner = () => (
    <div className="loading-spinner">
        <div className="spinner"></div>
        <p>Loading profile...</p>
    </div>
);

function ProfileById() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [profile, setProfile] = useState(null);
    const [avatarSrc, setAvatarSrc] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');

    const [playlistsPage, setPlaylistsPage] = useState(1);
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

    const fetchProfile = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/api/v1/profile/${id}`, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                },
            });

            const profileData = response.data.profile;
            setProfile(profileData);
            fetchAvatar(profileData.avatar_url);
        } catch (err) {
            console.error('Failed to fetch profile:', err);
            setError('Не удалось загрузить профиль');
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchProfile();

        return () => {
            if (avatarSrc) {
                URL.revokeObjectURL(avatarSrc);
            }
        };
    }, [id]);

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

    if (isLoading) {
        return <LoadingSpinner />;
    }

    if (error) {
        return <div className="error">{error}</div>;
    }

    if (!profile) {
        return <div>Профиль не найден</div>;
    }

    return (
        <div className="profile-page">
            <aside className="sidebar">
                <div className="sidebar-header">
                    <h2>Плейлисты пользователя</h2>
                </div>
                <ul className="playlist-list">
                    <li>Плейлист 1</li>
                    <li>Плейлист 2</li>
                    <li>Плейлист 3</li>
                </ul>
            </aside>

            <main className="main-content">
                <header className="header">
                    <div className="header-buttons">
                        <button
                            onClick={() => navigate('/profiles')}
                            className="profiles-list-button"
                        >
                            Все профили
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
                        </div>
                        <div className="profile-info-centered">
                            <h1>{profile.display_name || 'Пользователь'}</h1>
                            <p>
                                <strong>Имя:</strong> {profile.first_name || 'Не указано'}
                            </p>
                            <p>
                                <strong>Фамилия:</strong> {profile.last_name || 'Не указано'}
                            </p>
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
        </div>
    );
}

export default ProfileById;