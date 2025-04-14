package com.tracks.trackssvc.repository;

import com.tracks.trackssvc.model.Playlist;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface PlaylistRepository extends JpaRepository<Playlist, String> {
    public Optional<List<Playlist>> findByAuthorId(String userId);
}
