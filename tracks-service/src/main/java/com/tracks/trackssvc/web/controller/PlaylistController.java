package com.tracks.trackssvc.web.controller;

import com.tracks.trackssvc.service.PlaylistService;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import com.tracks.trackssvc.web.dto.TrackDto;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("api/v1/playlist")
public class PlaylistController {
    private final PlaylistService playlistService;

    public PlaylistController(PlaylistService playlistService) {
        this.playlistService = playlistService;
    }

    @PostMapping
    public ResponseEntity<PlaylistDto> createPlaylist(@RequestBody PlaylistDto playlistDto) {
        PlaylistDto created = playlistService.createPlaylist(playlistDto);
        return new ResponseEntity<>(created, HttpStatus.CREATED);
    }

    @GetMapping("/{id}")
    public ResponseEntity<PlaylistDto> getPlaylistById(@PathVariable String id) {
        PlaylistDto dto = playlistService.getPlaylistById(id);
        return ResponseEntity.ok(dto);
    }

    @GetMapping
    public ResponseEntity<List<PlaylistDto>> getAllPlaylists() {
        List<PlaylistDto> playlists = playlistService.getAllPlaylists();
        return ResponseEntity.ok(playlists);
    }

    @GetMapping("/user/{id}")
    public ResponseEntity<List<PlaylistDto>> getUserPlaylists(@PathVariable String id) {
        List<PlaylistDto> playlists = playlistService.getPlaylistsByUserId(id);
        return ResponseEntity.ok(playlists);
    }

    @PutMapping("/{id}")
    public ResponseEntity<PlaylistDto> updatePlaylist(@PathVariable String id, @RequestBody PlaylistDto playlistDto) {
        PlaylistDto updated = playlistService.updatePlaylist(id, playlistDto);
        return ResponseEntity.ok(updated);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deletePlaylist(@PathVariable String id) {
        playlistService.deletePlaylist(id);
        return ResponseEntity.noContent().build();
    }

    @PostMapping("/{playlistId}/tracks")
    public ResponseEntity<PlaylistDto> addTrackToPlaylist(@PathVariable String playlistId,
                                                          @RequestBody TrackDto trackDto) {
        PlaylistDto updated = playlistService.addTrackToPlaylist(playlistId, trackDto);
        return ResponseEntity.ok(updated);
    }

    @DeleteMapping("/{playlistId}/tracks/{trackId}")
    public ResponseEntity<PlaylistDto> removeTrackFromPlaylist(@PathVariable String playlistId,
                                                               @PathVariable String trackId) {
        PlaylistDto updated = playlistService.removeTrackFromPlaylist(playlistId, trackId);
        return ResponseEntity.ok(updated);
    }
}
