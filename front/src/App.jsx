import React from "react";
import TrackList from "./components/TrackList";
import Player from "./components/Player";
import { PlayerProvider } from "./context/PlayerContext";

const App = () => {
    return (
        <PlayerProvider>
            <div className="pb-28">
                <TrackList />
            </div>
            <Player />
        </PlayerProvider>
    );
};

export default App;
