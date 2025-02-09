package com.tracks.trackssvc.web.mapper;

import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.service.impl.AudioServiceImpl;
import com.tracks.trackssvc.web.dto.TrackDto;
import com.tracks.trackssvc.web.dto.TrackUploadDto;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.Date;

@Mapper(componentModel = "spring")
public interface TrackUploadMapper {
}
