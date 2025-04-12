package com.tracks.trackssvc.web.controller;

import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.service.TrackService;
import com.tracks.trackssvc.web.dto.TrackDto;
import com.tracks.trackssvc.web.dto.TrackUploadDto;
import com.tracks.trackssvc.web.mapper.TrackMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Slice;
import org.springframework.data.web.PagedModel;
import org.springframework.data.web.PagedResourcesAssembler;
import org.springframework.web.bind.annotation.*;

@RequiredArgsConstructor
@RestController
@RequestMapping("api/v1/track")
@CrossOrigin(origins = "http://localhost:5173")
public class TrackController {
    private final TrackService trackService;
    private final TrackMapper trackMapper;

    @PostMapping(value = "/add", consumes = { "multipart/form-data" })
    public TrackDto addTrack(@ModelAttribute TrackUploadDto trackUploadDto) {
        return trackService.addTrack(trackUploadDto);
    }
    @GetMapping("/")
    public Page<TrackDto> getTracks(Pageable pageable) {
        return trackService.getTracks(pageable);
    }

}
