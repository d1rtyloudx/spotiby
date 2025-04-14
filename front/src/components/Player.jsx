import React, { useRef, useEffect, useState } from "react";
import * as dashjs from "dashjs";
import { usePlayer } from "../context/PlayerContext";
import "../styles/Player.css";

const Player = () => {
    const { currentTrack, setCurrentTrack } = usePlayer();
    const videoRef = useRef(null);
    const dashPlayerRef = useRef(null);
    const [volume, setVolume] = useState(0.5);
    const [isPlaying, setIsPlaying] = useState(false);
    const [progress, setProgress] = useState(0);
    const [duration, setDuration] = useState(0);

    // Функция форматирования времени в формате mm:ss
    const formatTime = (time) => {
        if (isNaN(time)) return "00:00";
        const minutes = Math.floor(time / 60);
        const seconds = Math.floor(time % 60);
        return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
    };

    // Инициализация плеера при изменении currentTrack
    useEffect(() => {
        if (!currentTrack) return;

        const initPlayer = (trackToPlay) => {
            try {
                // Если плеер уже существует — сбросим его
                if (dashPlayerRef.current) {
                    dashPlayerRef.current.reset();
                    dashPlayerRef.current = null;
                }
                const player = dashjs.MediaPlayer().create();

                if (!videoRef.current) {
                    console.error("Player: videoRef не найден.");
                    return;
                }

                const token = localStorage.getItem("access_token");
                if (!token) {
                    console.error("Player: Токен авторизации отсутствует");
                    alert("Вы не авторизованы. Пожалуйста, войдите снова.");
                    return;
                }

                player.addRequestInterceptor((request) => {
                    request.headers = request.headers || {};
                    request.headers["Authorization"] = `Bearer ${token}`;
                    return request;
                });

                const streamUrl = `http://localhost:8080/api/v1/stream/${trackToPlay.id}/manifest.mpd`;
                player.initialize(videoRef.current, streamUrl, false);
                player.setVolume(volume);
                dashPlayerRef.current = player;

                videoRef.current.play()
                    .then(() => setIsPlaying(true))
                    .catch((err) => {
                        console.error("Player: Ошибка воспроизведения:", err);
                        setIsPlaying(false);
                    });

                // Обновляем длительность трека с небольшой задержкой
                setTimeout(() => {
                    if (videoRef.current && videoRef.current.duration) {
                        setDuration(videoRef.current.duration);
                    }
                }, 500);
            } catch (err) {
                console.error("Player: Ошибка инициализации dash-плеера:", err);
            }
        };

        if (!videoRef.current) {
            // Если videoRef ещё не установлен – ждём следующего рендера
            requestAnimationFrame(() => {
                if (videoRef.current) {
                    initPlayer(currentTrack);
                } else {
                    console.error("Player: videoRef по-прежнему не доступен.");
                }
            });
        } else {
            initPlayer(currentTrack);
        }

        // Очистка плеера при смене трека или размонтировании компонента
        return () => {
            if (dashPlayerRef.current) {
                dashPlayerRef.current.reset();
                dashPlayerRef.current = null;
            }
        };
    }, [currentTrack]);

    // Обновление прогресса воспроизведения каждую секунду
    useEffect(() => {
        const interval = setInterval(() => {
            if (videoRef.current && !videoRef.current.paused) {
                setProgress(videoRef.current.currentTime);
                setDuration(videoRef.current.duration || 0);
            }
        }, 1000);
        return () => clearInterval(interval);
    }, []);

    // Обновление громкости без перезапуска трека
    useEffect(() => {
        if (dashPlayerRef.current) {
            dashPlayerRef.current.setVolume(volume);
        }
    }, [volume]);

    // Переключение воспроизведения
    const togglePlay = () => {
        if (!videoRef.current) return;
        if (videoRef.current.paused) {
            videoRef.current.play()
                .then(() => setIsPlaying(true))
                .catch((err) => console.error("Player: Ошибка воспроизведения:", err));
        } else {
            videoRef.current.pause();
            setIsPlaying(false);
        }
    };

    // Изменение позиции воспроизведения
    const handleProgressChange = (e) => {
        const newTime = parseFloat(e.target.value);
        if (videoRef.current) {
            videoRef.current.currentTime = newTime;
            setProgress(newTime);
        }
    };

    // Закрытие плеера
    const closePlayer = () => {
        setCurrentTrack(null);
        if (dashPlayerRef.current) {
            dashPlayerRef.current.reset();
            dashPlayerRef.current = null;
        }
    };

    if (!currentTrack) return null;

    return (
        <div className="player-container">
            <video ref={videoRef} style={{ display: "none" }} />
            <div className="player-center">
                <div className="player-info">
                    <img
                        src={currentTrack.coverUrl}
                        alt={currentTrack.title}
                        className="player-cover"
                    />
                    <div className="player-track-info">
                        <div className="player-track-title">{currentTrack.title}</div>
                    </div>
                    <button onClick={togglePlay} className="player-play-pause">
                        {isPlaying ? " ⏸ " : " ▶ "}
                    </button>
                    <input
                        type="range"
                        min={0}
                        max={1}
                        step={0.01}
                        value={volume}
                        onChange={(e) => setVolume(parseFloat(e.target.value))}
                        className="player-volume-slider"
                    />
                    <button onClick={closePlayer} className="player-close-button">
                        X
                    </button>
                </div>
                <div className="player-progress">
                    <span className="player-current-time">{formatTime(progress)}</span>
                    <input
                        type="range"
                        min={0}
                        max={duration || 0}
                        step={0.1}
                        value={progress}
                        onChange={handleProgressChange}
                        className="player-progress-slider"
                    />
                    <span className="player-duration">{formatTime(duration)}</span>
                </div>
            </div>
        </div>
    );
};

export default Player;
