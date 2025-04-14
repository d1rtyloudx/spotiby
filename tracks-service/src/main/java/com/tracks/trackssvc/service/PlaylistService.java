package com.tracks.trackssvc.service;

import com.tracks.trackssvc.model.Playlist;
import com.tracks.trackssvc.web.dto.PlaylistDto;
import com.tracks.trackssvc.web.dto.TrackDto;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface PlaylistService {
    PlaylistDto createPlaylist(PlaylistDto playlistDto);
    PlaylistDto getPlaylistById(String id);
    List<PlaylistDto> getAllPlaylists();
    PlaylistDto updatePlaylist(String id, PlaylistDto playlistDto);
    void deletePlaylist(String id);
    PlaylistDto addTrackToPlaylist(String playlistId, TrackDto trackDto);
    PlaylistDto removeTrackFromPlaylist(String playlistId, String trackId);

    List<PlaylistDto> getPlaylistsByUserId(String authorId);
}
