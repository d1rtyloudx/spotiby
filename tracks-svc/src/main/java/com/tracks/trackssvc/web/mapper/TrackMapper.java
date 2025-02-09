package com.tracks.trackssvc.web.mapper;

import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.web.dto.TrackDto;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface TrackMapper extends Mappable<Track, TrackDto> {

}
