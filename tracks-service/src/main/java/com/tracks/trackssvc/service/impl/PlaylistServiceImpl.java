package com.tracks.trackssvc.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracks.trackssvc.model.Playlist;
import com.tracks.trackssvc.repository.PlaylistRepository;
import com.tracks.trackssvc.service.PlaylistService;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import com.tracks.trackssvc.web.dto.UpdateCoverDto;
import com.tracks.trackssvc.web.mapper.PlaylistMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class PlaylistServiceImpl implements PlaylistService {
    private final PlaylistRepository playlistRepository;
    private final PlaylistMapper playlistMapper;
    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    @Override
    @Transactional
    public PlaylistDto createPlaylist(PlaylistDto playlist) {
        Playlist playlistEntity = playlistMapper.toEntity(playlist);
        playlistEntity.setCreatedAt(new Date());
        PlaylistDto created_playlist = playlistMapper.toDto(playlistRepository.save(playlistEntity));
        try {
            kafkaTemplate.send("playlist_create", objectMapper.writeValueAsString(created_playlist));
        } catch (JsonProcessingException e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Error while serialising playlist data", e);
        }
        return created_playlist;
    }

    @Override
    public Page<PlaylistDto> findAllPlaylists(Pageable pageable) {
        return playlistRepository.findAll(pageable).map(playlistMapper::toDto);
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
