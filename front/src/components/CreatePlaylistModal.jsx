import React, { useState } from 'react';
import axios from 'axios';

const CreatePlaylistModal = ({ isOpen, onClose, onCreated, profile }) => {
    const [playlistData, setPlaylistData] = useState({
        title: '',
        cover: null,
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setPlaylistData((prev) => ({ ...prev, [name]: value }));
    };

    const handleCoverChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setPlaylistData((prev) => ({ ...prev, cover: file }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!playlistData.title || !playlistData.cover) {
            alert('Заполните все поля!');
            return;
        }

        try {
            const formData = new FormData();
            formData.append('title', playlistData.title);
            formData.append('authorId', profile.id);

            // создаём плейлист
            const response = await axios.post('/api/v1/playlist', formData, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            const playlistId = response.data.id;

            // загружаем обложку отдельно
            const coverFormData = new FormData();
            coverFormData.append('file', playlistData.cover);

            await axios.post(`/api/v1/image/playlist/${playlistId}`, coverFormData, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('access_token')}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            onCreated(); // callback чтобы обновить список
            onClose(); // закрыть окно
        } catch (error) {
            console.error('Ошибка при создании плейлиста:', error.response?.data || error.message);
            alert('Не удалось создать плейлист');
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
            <div className="bg-white p-6 rounded-lg w-full max-w-md shadow-lg">
                <h2 className="text-xl font-bold mb-4">Создать плейлист</h2>
                <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                    <input
                        type="text"
                        name="title"
                        placeholder="Название плейлиста"
                        value={playlistData.title}
                        onChange={handleChange}
                        className="border p-2 rounded"
                        required
                    />
                    <input
                        type="file"
                        accept="image/*"
                        onChange={handleCoverChange}
                        className="border p-2 rounded"
                        required
                    />
                    <div className="flex justify-end gap-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
                        >
                            Отмена
                        </button>
                        <button
                            type="submit"
                            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                        >
                            Создать
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default CreatePlaylistModal;