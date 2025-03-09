import React, { useState, useEffect, useRef } from "react";
import * as dashjs from "dashjs";

const DashAudioPlayer = ({ trackId }) => {
    const [isPlaying, setIsPlaying] = useState(false);
    const videoRef = useRef(null);

    const handlePlayClick = () => {
        if (videoRef.current) {
            videoRef.current.play().then(() => {
                setIsPlaying(true);
            }).catch((error) => {
                console.error("Audio failed to play:", error);
            });
        }
    };

    useEffect(() => {
        if (!trackId) return;

        console.log("useEffect triggered with trackId:", trackId);

        if (!videoRef.current) {
            console.error("Video element is not ready!");
            return;
        }

        const player = dashjs.MediaPlayer().create();
        const manifestUrl = `http://localhost:4444/api/v1/stream/${trackId}/manifest.mpd`;

        player.initialize(videoRef.current, manifestUrl, true);
        player.
        videoRef.current.muted = false;

        return () => {
            player.reset();
        };
    }, [trackId]);

    return (
        <div>
            <h2>Now Playing: {trackId}</h2>
            <button onClick={handlePlayClick}>Play</button>
            <video ref={videoRef} controls hidden />
        </div>
    );
};

export default DashAudioPlayer;