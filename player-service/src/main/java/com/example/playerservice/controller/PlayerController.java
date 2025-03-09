package com.example.playerservice.controller;

import com.example.playerservice.model.Player;
import com.example.playerservice.service.PlayerService;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("api/v1/player")
@RequiredArgsConstructor
public class PlayerController {
    private final PlayerService playerService;

    @PostMapping("/{userId}")
    public Player createPlayer(@PathVariable String userId) {
        return playerService.createPlayer(userId);
    }
    @GetMapping("/{userId}")
    public Player getPlayer(@PathVariable String userId) {
        return playerService.getPlayer(userId);
    }
    @PostMapping("/{userId}/play/{trackId}")
    public Player play(@PathVariable String userId, @PathVariable String trackId) {
        return playerService.playTrack(trackId, userId);
    }
}
