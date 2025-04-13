import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/Profile.css';

// Компонент спиннера для загрузки
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
    const [isLoading, setIsLoading] = useState(false); // Новое состояние для загрузки
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
            } catch (err) {
                console.error('Avatar upload error:', err);
            }
        }
    };

    const handleEditSubmit = async (e) => {
        e.preventDefault();
        setEditError('');
        setIsLoading(true); // Включаем индикатор загрузки

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
            localStorage.setItem('profile', JSON.stringify(response.data.profile));
            setProfile(response.data.profile);
            setIsEditing(false);
        } catch (err) {
            console.error('Edit profile error:', err);
            setEditError(err.response?.data?.message || 'Не удалось обновить профиль');
        } finally {
            setIsLoading(false); // Выключаем индикатор загрузки
        }
    };

    if (!profile) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="profile-page">
            {/* Боковая панель */}
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

            {/* Основной контент */}
            <main className="main-content">
                {/* Хедер */}
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
                    <button onClick={handleLogout} className="logout-button">
                        Выйти
                    </button>
                </header>

                {/* Профиль */}
                <section className="profile-section">
                    <div className="profile-header">
                        <div className="avatar-wrapper">
                            {profile.avatar_url ? (
                                <img
                                    src={profile.avatar_url}
                                    alt="Аватар"
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
                                    Изменить фото
                                </label>
                            </div>
                        </div>
                        <div className="profile-info">
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

                {/* Панели */}
                <section className="content-panels">
                    <div className="panel">
                        <h3>Подписанные плейлисты</h3>
                        <ul className="panel-list">
                            <li>Подписанный плейлист 1</li>
                            <li>Подписанный плейлист 2</li>
                            <li>Подписанный плейлист 3</li>
                        </ul>
                    </div>
                    <div className="panel">
                        <h3>Подписки на профили</h3>
                        <ul className="panel-list">
                            <li>Профиль 1</li>
                            <li>Профиль 2</li>
                            <li>Профиль 3</li>
                        </ul>
                    </div>
                    <div className="panel">
                        <h3>Топ треков</h3>
                        <ul className="panel-list">
                            <li>Трек 1</li>
                            <li>Трек 2</li>
                            <li>Трек 3</li>
                        </ul>
                    </div>
                </section>
            </main>

            {/* Модал редактирования */}
            {isEditing && (
                <div className="modal">
                    <div className="modal-content">
                        {isLoading ? (
                            <LoadingSpinner /> // Показываем спиннер во время загрузки
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