package com.tracks.trackssvc.web.controller;

import com.tracks.trackssvc.service.PlaylistService;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RequiredArgsConstructor
@RestController
@RequestMapping("api/v1/playlist")
public class PlaylistController {
    private final PlaylistService playlistService;


    @PostMapping("/create")
    public PlaylistDto create(@RequestBody PlaylistDto playlistDto) {
        return playlistService.createPlaylist(playlistDto);
    }

    @GetMapping("/")
    public Page<PlaylistDto> getAll(Pageable pageable) {
        return playlistService.findAllPlaylists(pageable);
    }
}
