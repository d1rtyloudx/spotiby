package com.tracks.trackssvc.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracks.trackssvc.model.Playlist;
import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.repository.PlaylistRepository;
import com.tracks.trackssvc.repository.TrackRepository;
import com.tracks.trackssvc.service.PlaylistService;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import com.tracks.trackssvc.web.dto.TrackDto;
import com.tracks.trackssvc.web.dto.UpdateCoverDto;
import com.tracks.trackssvc.web.mapper.PlaylistMapper;
import com.tracks.trackssvc.web.mapper.TrackMapper;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Service
public class PlaylistServiceImpl implements PlaylistService {
    private final ObjectMapper objectMapper;
    private final PlaylistRepository playlistRepository;
    private final TrackRepository trackRepository;
    private final PlaylistMapper playlistMapper;
    private final TrackMapper trackMapper;

    public PlaylistServiceImpl(PlaylistRepository playlistRepository,
                               TrackRepository trackRepository,
                               PlaylistMapper playlistMapper,
                               TrackMapper trackMapper,
                               ObjectMapper objectMapper) {
        this.playlistRepository = playlistRepository;
        this.trackRepository = trackRepository;
        this.playlistMapper = playlistMapper;
        this.trackMapper = trackMapper;
        this.objectMapper = objectMapper;
    }

    @Override
    public PlaylistDto createPlaylist(PlaylistDto playlistDto) {
        Playlist playlist = playlistMapper.toEntity(playlistDto);
        playlist.setCreatedAt(new Date());
        Playlist saved = playlistRepository.save(playlist);
        return playlistMapper.toDto(saved);
    }

    @Override
    public PlaylistDto getPlaylistById(String id) {
        Playlist playlist = playlistRepository.findById(id)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Playlist not found"));
        return playlistMapper.toDto(playlist);
    }

    @Override
    public List<PlaylistDto> getAllPlaylists() {
        return playlistRepository.findAll()
                .stream()
                .map(playlistMapper::toDto)
                .collect(Collectors.toList());
    }

    @Override
    public PlaylistDto updatePlaylist(String id, PlaylistDto playlistDto) {
        Playlist playlist = playlistRepository.findById(id)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Playlist not found"));
        playlist.setTitle(playlistDto.getTitle());
        playlist.setCoverUrl(playlistDto.getCoverUrl());
        Playlist updated = playlistRepository.save(playlist);
        return playlistMapper.toDto(updated);
    }

    @Override
    public void deletePlaylist(String id) {
        Playlist playlist = playlistRepository.findById(id)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Playlist not found"));
        playlistRepository.delete(playlist);
    }

    @Override
    public PlaylistDto addTrackToPlaylist(String playlistId, TrackDto trackDto) {
        Playlist playlist = playlistRepository.findById(playlistId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Playlist not found"));
        Track track;
        if (trackDto.getId() != null) {
            track = trackRepository.findById(trackDto.getId()).orElse(null);
            if (track == null) {
                track = trackMapper.toEntity(trackDto);
                track = trackRepository.save(track);
            }
        } else {
            track = trackMapper.toEntity(trackDto);
            track = trackRepository.save(track);
        }
        playlist.getTracks().add(track);
        Playlist updated = playlistRepository.save(playlist);
        return playlistMapper.toDto(updated);
    }

    @Override
    public PlaylistDto removeTrackFromPlaylist(String playlistId, String trackId) {
        Playlist playlist = playlistRepository.findById(playlistId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Playlist not found"));
        Track track = trackRepository.findById(trackId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Playlist not found"));
        playlist.getTracks().remove(track);
        Playlist updated = playlistRepository.save(playlist);
        return playlistMapper.toDto(updated);
    }

    @Override
    public List<PlaylistDto> getPlaylistsByUserId(String authorId) {
        return playlistRepository.findByAuthorId(authorId).orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "No playlists"))
                .stream().map(playlistMapper::toDto).toList();
    }

    @Transactional
    @RabbitListener(queues = "playlist_image")
    public void listenImageQueue(Message message) {
        String body = new String(message.getBody(), StandardCharsets.UTF_8);
        System.out.println(body);
        UpdateCoverDto dto;
        try {
            dto = objectMapper.readValue(body, UpdateCoverDto.class);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
        Optional<Playlist> playlist = playlistRepository.findById(dto.getId());
        if(playlist.isPresent()) {
            Playlist foundPlaylist = playlist.get();
            foundPlaylist.setCoverUrl(dto.getAvatarUrl());
            playlistRepository.save(foundPlaylist);
        } else return;
    }
}