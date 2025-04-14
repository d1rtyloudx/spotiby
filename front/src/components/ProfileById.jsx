import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/Profile.css';
import { usePlayer } from "../context/PlayerContext";

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
    const itemsPerPage = 5;

    const [trackCoverUrls, setTrackCoverUrls] = useState({}); // Новый стейт для обложек треков
    const [playlistCoverUrls, setPlaylistCoverUrls] = useState({});
    const [userPlaylists, setUserPlaylists] = useState([]);
    const [userTopTracks, setUserTopTracks] = useState([]);
    const [selectedPlaylist, setSelectedPlaylist] = useState(null);
    const { setCurrentTrack } = usePlayer();
    const [isPlaylistModalOpen, setIsPlaylistModalOpen] = useState(false);


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

    const fetchTrackCover = async (filename, trackId) => {
        if (!filename) {
            setTrackCoverUrls((prev) => ({ ...prev, [trackId]: '' }));
            return;
        }
        try {
            const response = await axios.get(
                `${API_BASE_URL}/api/v1/image/track/${filename}`,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                    responseType: 'blob',
                }
            );
            const imageUrl = URL.createObjectURL(response.data);
            setTrackCoverUrls((prev) => ({ ...prev, [trackId]: imageUrl }));
        } catch (err) {
            console.error(`Failed to fetch cover for track ${trackId}:`, err);
            setTrackCoverUrls((prev) => ({ ...prev, [trackId]: '' }));
        }
    };

    const fetchTopTracks = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/api/v1/track/user/${profile.id}`, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                },
            });
            setUserTopTracks(response.data);
            response.data.forEach((track) => {
                if (track.cover_url) {
                    fetchTrackCover(track.cover_url, track.id);
                } else {
                    setTrackCoverUrls((prev) => ({ ...prev, [track.id]: '' }));
                }
            });
        } catch (err) {
            console.error('Failed to fetch top tracks:', err);
            setUserTopTracks([]);
        }
    };

    const fetchPlaylistCover = async (filename, playlistId) => {
        if (!filename) {
            setPlaylistCoverUrls((prev) => ({ ...prev, [playlistId]: '' }));
            return;
        }

        try {
            const response = await axios.get(
                `${API_BASE_URL}/api/v1/image/playlist/${filename}`,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    },
                    responseType: 'blob',
                }
            );
            const imageUrl = URL.createObjectURL(response.data);
            setPlaylistCoverUrls((prev) => ({ ...prev, [playlistId]: imageUrl }));
        } catch (err) {
            console.error(`Failed to fetch cover for playlist ${playlistId}:`, err);
            setPlaylistCoverUrls((prev) => ({ ...prev, [playlistId]: '' }));
        }
    };

    const fetchPlaylists = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/api/v1/playlist/user/${profile.id}`, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                },
            });

            const playlists = response.data;
            setUserPlaylists(playlists);

            playlists.forEach((playlist) => {
                // Загружаем обложку плейлиста
                if (playlist.cover_url) {
                    fetchPlaylistCover(playlist.cover_url, playlist.id);
                } else {
                    setPlaylistCoverUrls((prev) => ({ ...prev, [playlist.id]: '' }));
                }

                // Загружаем обложки треков внутри плейлиста
                if (playlist.tracks && playlist.tracks.length > 0) {
                    playlist.tracks.forEach((track) => {
                        if (track.cover_url) {
                            fetchTrackCover(track.cover_url, track.id);
                        } else {
                            setTrackCoverUrls((prev) => ({ ...prev, [track.id]: '' }));
                        }
                    });
                }
            });
        } catch (err) {
            console.error('Failed to fetch playlists:', err);
            setUserPlaylists([]);
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

    useEffect(() => {
        if (profile) {
            fetchPlaylists();
            fetchTopTracks();
        }
    }, [profile]);

    useEffect(() => {
        userPlaylists.forEach((playlist) => {
            if (playlist.cover_url && !playlistCoverUrls[playlist.id]) {
                fetchPlaylistCover(playlist.cover_url, playlist.id);
            }
        });
    }, [userPlaylists]);

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

    const handlePlaylistClick = (playlist) => {
        setSelectedPlaylist(playlist);
        setIsPlaylistModalOpen(true);
    };

    const formatDuration = (ms) => {
        const totalSeconds = Math.floor(ms / 1000);
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
    };

    const handleTrackClick = (track) => {
        console.log("Profile: Клик по треку:", track); // Отладка
        setCurrentTrack({
            id: track.id,
            title: track.title,
            coverUrl: trackCoverUrls[track.id] || "",
        });
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
                        <h3 className="list-name text-xl font-semibold mb-4">Плейлисты</h3>
                        <div className="horizontal-list-wrapper overflow-x-auto">
                            <ul className="horizontal-list flex space-x-6">
                                {userPlaylists.length > 0 ? (
                                    getPaginatedItems(userPlaylists, playlistsPage).map((playlist) => (
                                        <li key={playlist.id}>
                                            <div className="playlist-card" onClick={() => handlePlaylistClick(playlist)}>
                                                <img
                                                    src={playlistCoverUrls[playlist.id] || '/placeholder.png'}
                                                    alt={playlist.title}
                                                    className="playlist-cover"
                                                    onError={() =>
                                                        setPlaylistCoverUrls((prev) => ({ ...prev, [playlist.id]: '' }))
                                                    }
                                                />
                                                <div className="playlist-info">
                                                    <p className="playlist-title">{playlist.title}</p>
                                                </div>
                                            </div>
                                        </li>
                                    ))
                                ) : (
                                    <li className="horizontal-item text-gray-500">Нет плейлистов</li>
                                )}
                            </ul>
                        </div>
                        <div className="pagination-controls flex justify-center space-x-3 mt-4">
                            <button
                                onClick={() => handlePrevPage(setPlaylistsPage, playlistsPage)}
                                disabled={playlistsPage === 1}
                                className="pagination-button px-3 py-1 bg-gray-300 hover:bg-gray-400 rounded disabled:opacity-50"
                            >
                                ←
                            </button>
                            <button
                                onClick={() => handleNextPage(setPlaylistsPage, playlistsPage, userPlaylists.length)}
                                disabled={playlistsPage >= Math.ceil(userPlaylists.length / itemsPerPage)}
                                className="pagination-button px-3 py-1 bg-gray-300 hover:bg-gray-400 rounded disabled:opacity-50"
                            >
                                →
                            </button>
                        </div>
                    </div>
                    <div className="list-section">
                        <h3 className="list-name text-xl font-semibold mb-4">Треки</h3>
                        <div className="horizontal-list-wrapper">
                            <ul className="horizontal-list">
                                {userTopTracks.length > 0 ? (
                                    getPaginatedItems(userTopTracks, topTracksPage).map((item) => (
                                        <li
                                            key={item.id}
                                            className="horizontal-item track-item cursor-pointer hover:bg-neutral-800" // Добавляем стили для интерактивности
                                            onClick={() => handleTrackClick(item)} // Обработчик клика
                                        >
                                            {trackCoverUrls[item.id] ? (
                                                <img
                                                    src={trackCoverUrls[item.id]}
                                                    alt={item.title}
                                                    className="track-cover"
                                                    onError={() =>
                                                        setTrackCoverUrls((prev) => ({
                                                            ...prev,
                                                            [item.id]: "",
                                                        }))
                                                    }
                                                />
                                            ) : null}
                                            <div className="track-info">
                                                <span className="track-title">{item.title}</span>
                                                <span className="track-duration">{formatDuration(item.duration_ms)}</span>
                                            </div>
                                        </li>
                                    ))
                                ) : (
                                    <li className="horizontal-item no-tracks">
                                        <span>Нет треков</span>
                                    </li>
                                )}
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
                                    onClick={() => handleNextPage(setTopTracksPage, topTracksPage, userTopTracks.length)}
                                    disabled={topTracksPage >= Math.ceil(userTopTracks.length / itemsPerPage)}
                                    className="pagination-button"
                                >
                                    →
                                </button>
                            </div>
                        </div>
                    </div>
                </section>
            </main>
            {isPlaylistModalOpen && selectedPlaylist && (
                <div className="modal">
                    <div className="modal-content playlist-modal">
                        <button className="modal-close" onClick={() => setIsPlaylistModalOpen(false)}>×</button>

                        <div className="playlist-header">
                            {playlistCoverUrls[selectedPlaylist.id] && (
                                <img
                                    src={playlistCoverUrls[selectedPlaylist.id]}
                                    alt="cover"
                                    className="playlist-modal-cover"
                                />
                            )}
                            <h2>{selectedPlaylist.title}</h2>
                        </div>

                        <ul className="playlist-track-list">
                            {selectedPlaylist.tracks.map((track) => (
                                <li key={track.id} className="playlist-track-item" onClick={() => handleTrackClick(track)}>
                                    {trackCoverUrls[track.id] && (
                                        <img src={trackCoverUrls[track.id]} alt="track" className="track-cover-small" />
                                    )}
                                    <span className="track-title">{track.title}</span>
                                    <span className="track-duration">{formatDuration(track.duration_ms)}</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            )}
        </div>
    );
}

export default ProfileById;