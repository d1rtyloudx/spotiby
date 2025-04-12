package com.example.playerservice.service.impl;

import com.example.playerservice.model.Player;
import com.example.playerservice.repository.PlayerRepository;
import com.example.playerservice.service.PlayerService;
import io.minio.MinioClient;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

@Service
@RequiredArgsConstructor(onConstructor_ = {@Autowired})
public class PlayerServiceImpl implements PlayerService {
    private final PlayerRepository playerRepository;
    private final MinioClient minioClient;

    @Override
    public Player createPlayer(String userId) {
        if(playerRepository.existsById(userId)) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Player already exists");
        }
        Player player = new Player();
        player.setId(userId);
        player.setVolume(0);
        player.setCurrentTrackId("none");
        player.setIsPlaying(false);
        player.setPlaybackTime(0L);
        player.setManifestUrl("");
        playerRepository.save(player);
        return player;
    }

    @Override
    public Player getPlayer(String playerId) {
        try {
            return playerRepository.getReferenceById(playerId);
        } catch(NullPointerException e) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND);
        }
    }

    @Override
    public Player updatePlayer(Player player) {
        return playerRepository.save(player);
    }

    @Override
    public Player playTrack(String trackId, String playerId) {
        Player player = getPlayer(playerId);
        player.setCurrentTrackId(trackId);
        player.setManifestUrl("stream/" + trackId + "/manifest.mpd");
        return playerRepository.save(player);
    }
}
