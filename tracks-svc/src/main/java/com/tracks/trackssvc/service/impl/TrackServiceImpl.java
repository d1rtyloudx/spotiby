package com.tracks.trackssvc.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.repository.TrackRepository;
import com.tracks.trackssvc.service.TrackService;
import com.tracks.trackssvc.web.dto.TrackUploadDto;
import com.tracks.trackssvc.web.dto.UpdateTrackCoverDto;
import lombok.RequiredArgsConstructor;
import org.jetbrains.annotations.NotNull;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.data.domain.*;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;
import ws.schild.jave.EncoderException;
import ws.schild.jave.MultimediaObject;
import ws.schild.jave.info.MultimediaInfo;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.Date;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class TrackServiceImpl implements TrackService {
    private final TrackRepository trackRepository;
    private final AudioServiceImpl audioService;
    private final ObjectMapper objectMapper;

    @Override
    @Transactional
    public Track addTrack(TrackUploadDto trackUploadDto) {
        Track trackToAdd = constructTrack(trackUploadDto);
        Track track = trackRepository.save(trackToAdd);
        audioService.upload(trackUploadDto.getAudioFile(), track.getId());
        return track;
    }

    public Page<Track> getTracks(Pageable pageable) {
        return trackRepository.findAll(pageable);
    }

    @NotNull
    private static Track constructTrack(TrackUploadDto trackUploadDto) {
        Track trackToAdd = new Track();
        if(trackUploadDto.getAudioFile() == null) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "File is required");
        }
        File audioFile;
        try {
            audioFile = File.createTempFile(
                    trackUploadDto.getAudioFile().getOriginalFilename().split("\\.")[0],
                    "." + trackUploadDto.getAudioFile().getOriginalFilename().split("\\.")[1]
            );
            Files.copy(trackUploadDto.getAudioFile().getInputStream(), audioFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Could not create temp file " +  e.getMessage());
        }
        MultimediaObject multimediaObject = new MultimediaObject(audioFile);
        try {
            MultimediaInfo multimediaInfo = multimediaObject.getInfo();
            trackToAdd.setDurationMs(multimediaInfo.getDuration());
        } catch (EncoderException e) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Could not obtain audio info", e);
        }
        audioFile.delete();
        trackToAdd.setAuthorId(trackUploadDto.getAuthorId());
        trackToAdd.setTitle(trackUploadDto.getTitle());
        trackToAdd.setUploadDate(new Date());
        return trackToAdd;
    }

    @Transactional
    @RabbitListener(queues = "track_image")
    public void listenImageQueue(Message message) {
        String body = new String(message.getBody(), StandardCharsets.UTF_8);
        System.out.println(body);
        UpdateTrackCoverDto dto = null;
        try {
            dto = objectMapper.readValue(body, UpdateTrackCoverDto.class);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
        Optional<Track> track = trackRepository.findById(dto.getId());
        if(track.isPresent()) {
            Track foundTrack = track.get();
            foundTrack.setCoverUrl(dto.getAvatarUrl());
            trackRepository.save(foundTrack);
        } else return;
    }
}
