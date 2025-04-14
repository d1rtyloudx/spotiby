import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/Profile.css';
import { usePlayer } from "../context/PlayerContext";
import CreatePlaylistModal from './CreatePlaylistModal';

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
    const [isAddingTrack, setIsAddingTrack] = useState(false);
    const [trackData, setTrackData] = useState({
        track_name: '',
        cover: null,
        audio: null,
    });
    const [isCreatingPlaylist, setIsCreatingPlaylist] = useState(false);
    const [playlistData, setPlaylistData] = useState({
        title: '',
        cover: null,
    });
    const [editError, setEditError] = useState('');
    const [avatarSrc, setAvatarSrc] = useState('');
    const [subscriptionAvatarUrls, setSubscriptionAvatarUrls] = useState({});
    const [trackCoverUrls, setTrackCoverUrls] = useState({}); // Новый стейт для обложек треков
    const [playlistCoverUrls, setPlaylistCoverUrls] = useState({});
    const [selectedPlaylist, setSelectedPlaylist] = useState(null);
    const [isPlaylistModalOpen, setIsPlaylistModalOpen] = useState(false);
    const [searchTrack, setSearchTrack] = useState('');
    const [foundTrack, setFoundTrack] = useState(null);
    const [userPlaylists, setUserPlaylists] = useState([]);
    const [userTopTracks, setUserTopTracks] = useState([]);
    const [playlistsPage, setPlaylistsPage] = useState(1);
    const [topTracksPage, setTopTracksPage] = useState(1);
    const itemsPerPage = 5;

    // Существующие подписки
    const [subscriptions, setSubscriptions] = useState([]);
    const [subscriptionsPage, setSubscriptionsPage] = useState(1);

    const { setCurrentTrack } = usePlayer();
    // Функция получения аватара
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

    // Функция получения аватара для подписок
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


    // Получение подписок на профили
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

    // Получение плейлистов пользователя
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

    // Получение топ треков пользователя
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
            Object.values(trackCoverUrls).forEach((url) => {
                if (url) URL.revokeObjectURL(url);
            });
        };
    }, [navigate]);

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

    // Остальные функции (handleLogout, handleSearch, и т.д.) остаются без изменений
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

    const handlePlaylistClick = (playlist) => {
        setSelectedPlaylist(playlist);
        setIsPlaylistModalOpen(true);
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

    const formatDuration = (ms) => {
        const totalSeconds = Math.floor(ms / 1000);
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
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

    const handleAddTrack = () => {
        setIsAddingTrack(true);
    };

    const closeAddTrackModal = () => {
        setIsAddingTrack(false);
        setTrackData({ track_name: '', cover: null, audio: null });
    };

    const handleTrackChange = (e) => {
        const { name, value } = e.target;
        setTrackData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    };

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        console.log('Выбранный файл обложки:', file); // Отладка
        if (file) {
            setTrackData((prevData) => ({
                ...prevData,
                cover: file,
            }));
        } else {
            console.warn('Файл обложки не выбран');
        }
    };

    const handleAudioChange = (e) => {
        const file = e.target.files[0];
        setTrackData((prevData) => ({
            ...prevData,
            audio: file,
        }));
    };

    const handleTrackSubmit = async (e) => {
        e.preventDefault();

        if (!trackData.track_name || !trackData.audio || !trackData.cover) {
            alert('Все поля обязательны для заполнения!');
            return;
        }

        try {
            const trackFormData = new FormData();
            trackFormData.append('title', trackData.track_name);
            trackFormData.append('authorId', profile.id);
            trackFormData.append('audioFile', trackData.audio);

            const trackResponse = await axios.post(`${API_BASE_URL}/api/v1/track/add`, trackFormData, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            console.log('Трек успешно загружен:', trackResponse.data);

            const coverFormData = new FormData();
            coverFormData.append('file', trackData.cover);

            const coverResponse = await axios.post(`${API_BASE_URL}/api/v1/image/track/${trackResponse.data.id}`, coverFormData, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            console.log('Обложка успешно загружена:', coverResponse.data);
            closeAddTrackModal();
            // Обновляем список треков после добавления
            fetchTopTracks();
        } catch (error) {
            console.error('Ошибка:', error.response?.data || error.message);
            alert('Произошла ошибка при загрузке данных');
        }
    };

    const handleTrackClick = (track) => {
        console.log("Profile: Клик по треку:", track); // Отладка
        setCurrentTrack({
            id: track.id,
            title: track.title,
            coverUrl: trackCoverUrls[track.id] || "",
        });
    };

    const handleCreatePlaylist = () => {
        setIsCreatingPlaylist(true);
    };

    const closeCreatePlaylistModal = () => {
        setIsCreatingPlaylist(false);
        setPlaylistData({ title: '', cover: null });
    };

    const handlePlaylistChange = (e) => {
        const { name, value } = e.target;
        setPlaylistData((prev) => ({ ...prev, [name]: value }));
    };

    const handlePlaylistCoverChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setPlaylistData((prev) => ({ ...prev, cover: file }));
        }
    };

    const handlePlaylistSubmit = async (e) => {
        e.preventDefault();

        if (!playlistData.title || !playlistData.cover) {
            alert('Все поля обязательны!');
            return;
        }

        try {
            const payload = {
                title: playlistData.title,
                author_id: profile.id,
            };

            const playlistResponse = await axios.post(`${API_BASE_URL}/api/v1/playlist`, payload, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                },
            });

            const playlistId = playlistResponse.data.id;

            const coverForm = new FormData();
            coverForm.append('file', playlistData.cover);

            await axios.post(`${API_BASE_URL}/api/v1/image/playlist/${playlistId}`, coverForm, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            closeCreatePlaylistModal();
            fetchPlaylists()
        } catch (error) {
            console.error('Ошибка при создании плейлиста:', error.response?.data || error.message);
            alert('Не удалось создать плейлист');
        }
    };

    const handleFindTrack = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/api/v1/track/search/`, {
                params: { title: searchTrack },
                headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
            });
            setFoundTrack(response.data);
        } catch (err) {
            console.error("Ошибка при поиске трека:", err);
            setFoundTrack(null);
        }
    };

    const handleAddTrackToPlaylist = async (trackDto, playlistId) => {
        try {
            await axios.post(
                `${API_BASE_URL}/api/v1/playlist/${playlistId}/tracks`, trackDto,
                { headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` } }
            );
            fetchPlaylists()
            setIsPlaylistModalOpen(false);
        } catch (err) {
            console.error('Ошибка при добавлении трека в плейлист:', err);
        }
    };

    if (!profile) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="profile-page">
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
                        <h3 className="list-name text-xl font-semibold mb-4">Ваши плейлисты</h3>
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
                            <button onClick={handleCreatePlaylist} className="add-button small px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700">
                                + Создать
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
                        <h3 className="list-name text-xl font-semibold mb-4">Подписки</h3>
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
                                <button onClick={handleAddTrack} className="add-button small">
                                    + Добавить
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
            {isAddingTrack && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Добавить трек</h2>
                        <form onSubmit={handleTrackSubmit}>
                            <div className="form-group">
                                <label>Имя трека</label>
                                <input
                                    type="text"
                                    name="track_name"
                                    value={trackData.track_name}
                                    onChange={handleTrackChange}
                                    placeholder="Введите имя трека"
                                />
                            </div>
                            <div className="form-group">
                                <label>Обложка</label>
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleImageChange}
                                />
                            </div>
                            <div className="form-group">
                                <label>Аудиофайл</label>
                                <input
                                    type="file"
                                    accept="audio/*"
                                    onChange={handleAudioChange}
                                />
                            </div>
                            <div className="modal-buttons">
                                <button type="submit" className="save-button">
                                    Добавить
                                </button>
                                <button
                                    type="button"
                                    onClick={closeAddTrackModal}
                                    className="cancel-button"
                                >
                                    Отмена
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
            {isCreatingPlaylist && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Создать плейлист</h2>
                        <form onSubmit={handlePlaylistSubmit}>
                            <div className="form-group">
                                <label>Название</label>
                                <input
                                    type="text"
                                    name="title"
                                    value={playlistData.title}
                                    onChange={handlePlaylistChange}
                                    placeholder="Введите название плейлиста"
                                />
                            </div>
                            <div className="form-group">
                                <label>Обложка</label>
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handlePlaylistCoverChange}
                                />
                            </div>
                            <div className="modal-buttons">
                                <button type="submit" className="save-button">
                                    Создать
                                </button>
                                <button
                                    type="button"
                                    onClick={closeCreatePlaylistModal}
                                    className="cancel-button"
                                >
                                    Отмена
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
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

                        <div className="track-search-section">
                            <input
                                type="text"
                                placeholder="Введите название трека"
                                value={searchTrack}
                                onChange={(e) => setSearchTrack(e.target.value)}
                                className="track-search-input"
                            />
                            <button onClick={handleFindTrack} className="track-search-button">Найти</button>
                            {foundTrack && (
                                <div className="found-track-info">
                                    <span>{foundTrack.title}</span>
                                    <button className="add-track-button"  onClick={() => handleAddTrackToPlaylist(foundTrack, selectedPlaylist.id)}>
                                        Добавить
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

export default Profile;