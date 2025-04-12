import React, { useRef, useEffect, useState } from "react";
import * as dashjs from "dashjs";
import { usePlayer } from "../context/PlayerContext";

const Player = () => {
    const { currentTrack } = usePlayer();
    const videoRef = useRef(null);
    const dashPlayerRef = useRef(null);
    const [volume, setVolume] = useState(0.5);
    const [isPlaying, setIsPlaying] = useState(false);
    const [progress, setProgress] = useState(0); // текущее время воспроизведения
    const [duration, setDuration] = useState(0); // длительность трека

    // Инициализация плеера только при изменении currentTrack
    useEffect(() => {
        if (!currentTrack || !currentTrack.id) {
            return;
        }

        const track = currentTrack;

        const initPlayer = (trackToPlay) => {
            try {
                if (dashPlayerRef.current) {
                    dashPlayerRef.current.reset();
                    dashPlayerRef.current = null;
                }
                const player = dashjs.MediaPlayer().create();

                if (!videoRef.current) {
                    console.error("videoRef не найден перед инициализацией плеера.");
                    return;
                }

                const streamUrl = `http://localhost:8085/api/v1/stream/${trackToPlay.id}/manifest.mpd`;
                player.initialize(videoRef.current, streamUrl, true);
                player.setVolume(volume);
                dashPlayerRef.current = player;
                setIsPlaying(true);

                // При запуске, если возможно, установим длительность
                const updateDuration = () => {
                    if (videoRef.current && videoRef.current.duration) {
                        setDuration(videoRef.current.duration);
                    }
                };
                // Попробуем обновить длительность через небольшой таймаут,
                // чтобы videoRef успел загрузить метаданные.
                setTimeout(updateDuration, 500);
            } catch (err) {
                console.error("Ошибка инициализации dash-плеера:", err);
            }
        };

        if (!videoRef.current) {
            requestAnimationFrame(() => {
                if (videoRef.current) {
                    initPlayer(track);
                } else {
                    console.error("videoRef не доступен даже после ожидания.");
                }
            });
        } else {
            initPlayer(track);
        }

        return () => {
            if (dashPlayerRef.current) {
                dashPlayerRef.current.reset();
                dashPlayerRef.current = null;
            }
        };
    }, [currentTrack]); // зависим только от currentTrack

    // Отдельный эффект для обновления громкости (без перезапуска плеера)
    useEffect(() => {
        if (dashPlayerRef.current) {
            dashPlayerRef.current.setVolume(volume);
        }
    }, [volume]);

    // Обновляем прогресс воспроизведения каждую секунду
    useEffect(() => {
        const interval = setInterval(() => {
            if (videoRef.current && !videoRef.current.paused) {
                setProgress(videoRef.current.currentTime);
                if (videoRef.current.duration) {
                    setDuration(videoRef.current.duration);
                }
            }
        }, 1000);

        return () => clearInterval(interval);
    }, []);

    // Функция для форматирования времени (мм:сс)
    const formatTime = (time) => {
        if (isNaN(time)) return "00:00";
        const minutes = Math.floor(time / 60);
        const seconds = Math.floor(time % 60);
        return `${minutes.toString().padStart(2, "0")}:${seconds
            .toString()
            .padStart(2, "0")}`;
    };

    const togglePlay = () => {
        const video = videoRef.current;
        if (!video) return;
        if (video.paused) {
            video.play();
            setIsPlaying(true);
        } else {
            video.pause();
            setIsPlaying(false);
        }
    };

    // Обработчик перемотки: позволяет пользователю перетаскивать прогресс-бар
    const handleProgressChange = (e) => {
        const newTime = parseFloat(e.target.value);
        if (videoRef.current) {
            videoRef.current.currentTime = newTime;
            setProgress(newTime);
        }
    };

    return (
        <div className="fixed bottom-0 left-0 right-0 bg-neutral-900 text-white p-4 flex flex-col gap-2 z-50">
            {/* Скрытый video элемент для dash.js */}
            <video ref={videoRef} className="hidden" preload="auto" />

            {currentTrack ? (
                <>
                    <div className="flex items-center gap-4 w-full">
                        <img
                            src={currentTrack.coverUrl}
                            alt={currentTrack.title}
                            className="w-12 h-12 rounded"
                        />
                        <div className="flex-1 overflow-hidden">
                            <div className="text-sm truncate">{currentTrack.title}</div>
                        </div>
                        <button
                            onClick={togglePlay}
                            className="text-xl hover:scale-110 transition-transform"
                        >
                            {isPlaying ? "⏸️" : "▶️"}
                        </button>
                        <input
                            type="range"
                            min={0}
                            max={1}
                            step={0.01}
                            value={volume}
                            onChange={(e) => setVolume(parseFloat(e.target.value))}
                            className="w-32"
                        />
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="text-xs">{formatTime(progress)}</span>
                        <input
                            type="range"
                            min={0}
                            max={duration || 0}
                            step={0.1}
                            value={progress}
                            onChange={handleProgressChange}
                            className="w-full"
                        />
                        <span className="text-xs">{formatTime(duration)}</span>
                    </div>
                </>
            ) : (
                <div className="text-sm text-gray-400">Трек не выбран</div>
            )}
        </div>
    );
};

export default Player;
