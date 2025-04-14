package com.tracks.trackssvc.service;

import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.web.dto.TrackDto;
import com.tracks.trackssvc.web.dto.TrackUploadDto;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.util.List;


public interface TrackService {
    TrackDto addTrack(TrackUploadDto track);
    public Page<TrackDto> getTracks(Pageable pageable);
    public TrackDto getTrackByTitle(String filter);
    TrackDto getTrackById(String id);
    public List<TrackDto> getAllTracksByAuthorId(String authorId);
}
