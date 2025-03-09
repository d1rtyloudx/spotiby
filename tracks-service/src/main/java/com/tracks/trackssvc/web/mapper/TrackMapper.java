package com.tracks.trackssvc.web.mapper;

import com.tracks.trackssvc.model.Track;
import com.tracks.trackssvc.web.dto.TrackDto;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface TrackMapper {

    TrackDto toDto(Track track);
    Track toEntity(TrackDto trackDto);
}
