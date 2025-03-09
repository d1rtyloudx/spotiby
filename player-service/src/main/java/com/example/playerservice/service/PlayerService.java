package com.example.playerservice.service;

import com.example.playerservice.model.Player;

public interface PlayerService {
    Player createPlayer(String userId);
    Player getPlayer(String playerId);
    Player updatePlayer(Player player);
    Player playTrack(String trackId, String playerId);
}
