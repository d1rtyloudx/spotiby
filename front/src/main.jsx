import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.jsx';
import './index.css';
import { PlayerProvider } from './context/PlayerContext.jsx';
import Player from "./components/Player.jsx";

ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        <PlayerProvider>
            <App />
            <Player />
        </PlayerProvider>
    </React.StrictMode>,
);