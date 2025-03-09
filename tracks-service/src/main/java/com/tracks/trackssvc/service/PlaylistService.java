package com.tracks.trackssvc.service;

import com.tracks.trackssvc.model.Playlist;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

public interface PlaylistService {
    PlaylistDto createPlaylist(PlaylistDto playlist);
    Page<PlaylistDto> findAllPlaylists(Pageable pageable);
}
