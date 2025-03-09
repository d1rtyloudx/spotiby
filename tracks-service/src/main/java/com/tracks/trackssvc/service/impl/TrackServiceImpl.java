package com.tracks.trackssvc.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.repository.TrackRepository;
import com.tracks.trackssvc.service.TrackService;
import com.tracks.trackssvc.web.dto.TrackDto;
import com.tracks.trackssvc.web.dto.TrackUploadDto;
import com.tracks.trackssvc.web.dto.UpdateCoverDto;
import com.tracks.trackssvc.web.mapper.TrackMapper;
import lombok.RequiredArgsConstructor;
import org.jetbrains.annotations.NotNull;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.data.domain.*;
import org.springframework.http.HttpStatus;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.io.BufferedReader;
import java.io.File;
import java.io.IOException;
import java.io.InputStreamReader;
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
    private final TrackMapper trackMapper;
    private final KafkaTemplate<String, String> kafkaTemplate;

    @Override
    @Transactional
    public TrackDto addTrack(TrackUploadDto trackUploadDto) {
        Track trackToAdd = constructTrack(trackUploadDto);
        Track track = trackRepository.save(trackToAdd);
        audioService.upload(trackUploadDto.getAudioFile(), track.getId());
        TrackDto trackDto = trackMapper.toDto(track);
        try {
            kafkaTemplate.send("track_create", (objectMapper.writeValueAsString(trackDto)));
        } catch (JsonProcessingException e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Failed to serialize track", e);
        }
        return trackDto;
    }

    public Page<TrackDto> getTracks(Pageable pageable) {
        return trackRepository.findAll(pageable).map(trackMapper::toDto);
    }

    @NotNull
    private static Track constructTrack(TrackUploadDto trackUploadDto) {
        Track trackToAdd = new Track();

        if (trackUploadDto.getAudioFile() == null) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "File is required");
        }

        File audioFile;
        try {
            String originalFilename = trackUploadDto.getAudioFile().getOriginalFilename();
            String baseName = originalFilename.contains(".")
                    ? originalFilename.substring(0, originalFilename.lastIndexOf('.'))
                    : originalFilename;
            String extension = originalFilename.contains(".")
                    ? originalFilename.substring(originalFilename.lastIndexOf('.') + 1)
                    : "";
            audioFile = File.createTempFile(baseName, "." + extension);
            Files.copy(trackUploadDto.getAudioFile().getInputStream(), audioFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Could not create temp file " + e.getMessage());
        }

        try {
            long durationMs = getDurationMs(audioFile);
            trackToAdd.setDurationMs(durationMs);
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Could not obtain audio info: " + e.getMessage(), e);
        } finally {
            if (audioFile.exists()) {
                audioFile.delete();
            }
        }

        trackToAdd.setAuthorId(trackUploadDto.getAuthorId());
        trackToAdd.setTitle(trackUploadDto.getTitle());
        trackToAdd.setUploadDate(new Date());
        return trackToAdd;
    }

    private static long getDurationMs(File audioFile) throws IOException, InterruptedException {
        ProcessBuilder pb = new ProcessBuilder(
                "ffprobe",
                "-v", "error",
                "-show_entries", "format=duration",
                "-of", "default=noprint_wrappers=1:nokey=1",
                audioFile.getAbsolutePath()
        );
        pb.redirectErrorStream(true);
        Process process = pb.start();

        BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
        String durationStr = reader.readLine();

        int exitCode = process.waitFor();
        if (exitCode != 0 || durationStr == null) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "ffprobe returned non-zero exit code or empty duration");
        }

        double durationSec = Double.parseDouble(durationStr);
        return (long) (durationSec * 1000);
    }

    @Transactional
    @RabbitListener(queues = "track_image")
    public void listenImageQueue(Message message) {
        String body = new String(message.getBody(), StandardCharsets.UTF_8);
        System.out.println(body);
        UpdateCoverDto dto;
        try {
            dto = objectMapper.readValue(body, UpdateCoverDto.class);
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
