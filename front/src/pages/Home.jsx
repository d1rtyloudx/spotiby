import React, { useState } from "react";
import DashPlayer from "../components/DashPlayer";

const Home = () => {
    const [trackId, setTrackId] = useState("41005a3f-99f1-4700-866a-ebda82a8d94d");

    return (
        <div>
            <h1>Spotiby player test</h1>
            <DashPlayer trackId={trackId} />

            <button onClick={() => setTrackId()}>Play Next Track</button>
        </div>
    );
};

export default Home;