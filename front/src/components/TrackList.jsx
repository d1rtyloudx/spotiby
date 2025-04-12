import React, { useEffect, useState } from "react";
import axios from "axios";
import { usePlayer } from "../context/PlayerContext";

const TrackList = () => {
    const [tracks, setTracks] = useState([]);
    const { setCurrentTrack } = usePlayer();

    useEffect(() => {
        axios.get("http://localhost:8083/api/v1/track/").then(res => setTracks(res.data.content));
    }, []);

    return (
        <div className="p-4 space-y-4">
            {tracks.map(track => (
                <div
                    key={track.id}
                    className="flex items-center gap-4 p-2 hover:bg-gray-100 rounded cursor-pointer"
                    onClick={() => setCurrentTrack(track)}
                >
                    <img src={track.cover_url} alt={track.title} className="w-16 h-16 rounded" />
                    <div>
                        <div className="font-medium">{track.title}</div>
                        <div className="text-sm text-gray-500">{formatDuration(track.duration_ms)}</div>
                    </div>
                </div>
            ))}
        </div>
    );
};

function formatDuration(ms) {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
}

export default TrackList;
