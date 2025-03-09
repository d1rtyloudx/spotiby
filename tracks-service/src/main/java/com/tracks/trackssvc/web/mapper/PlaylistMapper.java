package com.tracks.trackssvc.web.mapper;

import com.tracks.trackssvc.model.Playlist;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring", uses = {TrackMapper.class})
public interface PlaylistMapper {
    PlaylistDto toDto(Playlist playlist);
    Playlist toEntity(PlaylistDto playlist);
}
