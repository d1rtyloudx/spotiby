import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';
import '../styles/ProfilesList.css';

const ProfilesList = () => {
    const navigate = useNavigate();
    const [profiles, setProfiles] = useState([]);
    const [avatarUrls, setAvatarUrls] = useState({});
    const [pagination, setPagination] = useState({
        limit: 1000,
        current_page: 1,
        total_pages: 1,
        has_more: false,
    });
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const [ownProfileId, setOwnProfileId] = useState('');
    const [followingStatus, setFollowingStatus] = useState({});
    const [followedProfileIds, setFollowedProfileIds] = useState([]);

    useEffect(() => {
        const storedProfile = localStorage.getItem('profile');
        if (storedProfile) {
            const parsedProfile = JSON.parse(storedProfile);
            setOwnProfileId(parsedProfile.id);
        } else {
            navigate('/auth');
        }
    }, [navigate]);

    const fetchAvatar = async (filename, profileId) => {
        if (!filename) {
            setAvatarUrls((prev) => ({ ...prev, [profileId]: '' }));
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
            setAvatarUrls((prev) => ({ ...prev, [profileId]: imageUrl }));
        } catch (err) {
            console.error(`Failed to fetch avatar for profile ${profileId}:`, err);
            setAvatarUrls((prev) => ({ ...prev, [profileId]: '' }));
        }
    };

    const fetchFollowedProfiles = async () => {
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

            return allProfiles.map((profile) => profile.id);
        } catch (err) {
            console.error('Failed to fetch followed profiles:', err);
            return [];
        }
    };

    const fetchProfiles = async (page = 1) => {
        setIsLoading(true);
        setError('');
        try {
            const profilesResponse = await axios.get(`${API_BASE_URL}/api/v1/profile/`, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                },
                params: {
                    page,
                    limit: pagination.limit,
                },
            });

            const { profiles, pagination: newPagination } = profilesResponse.data.items;
            const filteredProfiles = profiles.filter((profile) => profile.id !== ownProfileId);
            setProfiles(filteredProfiles);
            setPagination(newPagination);

            filteredProfiles.forEach((profile) => {
                if (profile.avatar_url) {
                    fetchAvatar(profile.avatar_url, profile.id);
                } else {
                    setAvatarUrls((prev) => ({ ...prev, [profile.id]: '' }));
                }
            });

            const followedIds = await fetchFollowedProfiles();
            setFollowedProfileIds(followedIds);

            const statusMap = filteredProfiles.reduce((acc, profile) => {
                acc[profile.id] = followedIds.includes(profile.id);
                return acc;
            }, {});
            setFollowingStatus(statusMap);
        } catch (err) {
            console.error('Failed to fetch profiles:', err);
            setError('Не удалось загрузить профили');
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        if (ownProfileId) {
            fetchProfiles(pagination.current_page);
        }
    }, [pagination.current_page, ownProfileId]);

    useEffect(() => {
        return () => {
            Object.values(avatarUrls).forEach((url) => {
                if (url) URL.revokeObjectURL(url);
            });
        };
    }, [avatarUrls]);

    const handleNextPage = () => {
        if (pagination.has_more && pagination.current_page < pagination.total_pages) {
            setPagination((prev) => ({
                ...prev,
                current_page: prev.current_page + 1,
            }));
        }
    };

    const handlePrevPage = () => {
        if (pagination.current_page > 1) {
            setPagination((prev) => ({
                ...prev,
                current_page: prev.current_page - 1,
            }));
        }
    };

    const handleFollowToggle = async (profileId) => {
        const isFollowing = followingStatus[profileId];
        try {
            if (isFollowing) {
                await axios.put(
                    `${API_BASE_URL}/api/v1/profile/me/unfollow/${profileId}`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                        },
                    }
                );
                setFollowingStatus((prev) => ({ ...prev, [profileId]: false }));
                setFollowedProfileIds((prev) => prev.filter((id) => id !== profileId));
            } else {
                await axios.put(
                    `${API_BASE_URL}/api/v1/profile/me/follow/${profileId}`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                        },
                    }
                );
                setFollowingStatus((prev) => ({ ...prev, [profileId]: true }));
                setFollowedProfileIds((prev) => [...prev, profileId]);
            }
        } catch (err) {
            console.error(`Failed to ${isFollowing ? 'unfollow' : 'follow'} profile:`, err);
            alert(`Не удалось ${isFollowing ? 'отписаться' : 'подписаться'}`);
        }
    };

    const handleProfileClick = (profileId) => {
        navigate(`/profile/${profileId}`);
    };

    return (
        <div className="profiles-list-page">
            <header className="header">
                <h1>Все профили</h1>
            </header>
            <main className="main-content scrollable">
                {isLoading ? (
                    <div className="loading-spinner">
                        <div className="spinner"></div>
                        <p>Загрузка профилей...</p>
                    </div>
                ) : error ? (
                    <p className="error">{error}</p>
                ) : profiles.length === 0 ? (
                    <p className="no-profiles">Нет профилей для отображения</p>
                ) : (
                    <ul className="profiles-list">
                        {profiles.map((profile) => (
                            <li
                                key={profile.id}
                                className="profile-item"
                                onClick={() => handleProfileClick(profile.id)}
                            >
                                <div className="profile-avatar">
                                    {avatarUrls[profile.id] ? (
                                        <img
                                            src={avatarUrls[profile.id]}
                                            alt="Аватар"
                                            className="avatar-img"
                                            onError={() =>
                                                setAvatarUrls((prev) => ({ ...prev, [profile.id]: '' }))
                                            }
                                        />
                                    ) : (
                                        <div className="avatar-placeholder">
                                            {profile.display_name?.charAt(0) || 'U'}
                                        </div>
                                    )}
                                </div>
                                <div className="profile-info">
                                    <h2>{profile.display_name || 'Пользователь'}</h2>
                                    <p>
                                        {profile.first_name || ''} {profile.last_name || ''}
                                    </p>
                                    <p className="description">
                                        {profile.description || 'Нет описания'}
                                    </p>
                                </div>
                                <button
                                    className={`follow-button ${followingStatus[profile.id] ? 'unfollow' : ''}`}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        handleFollowToggle(profile.id);
                                    }}
                                >
                                    {followingStatus[profile.id] ? 'Отписаться' : 'Подписаться'}
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </main>
            <footer className="pagination-footer">
                <div className="pagination">
                    <button
                        onClick={handlePrevPage}
                        disabled={pagination.current_page === 1}
                        className="pagination-button"
                    >
                        Назад
                    </button>
                    <span className="pagination-info">
                        Страница {pagination.current_page} из {pagination.total_pages}
                    </span>
                    <button
                        onClick={handleNextPage}
                        disabled={
                            !pagination.has_more ||
                            pagination.current_page === pagination.total_pages
                        }
                        className="pagination-button"
                    >
                        Вперед
                    </button>
                </div>
            </footer>
        </div>
    );
};

export default ProfilesList;