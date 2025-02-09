package com.tracks.trackssvc.service;

import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.web.dto.TrackUploadDto;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Slice;


public interface TrackService {
    Track addTrack(TrackUploadDto track);
    public Page<Track> getTracks(Pageable pageable);
}
