package com.tracks.trackssvc.repository;

import com.tracks.trackssvc.model.Track;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface TrackRepository extends JpaRepository<Track, String> {
    public Optional<List<Track>> findByAuthorId(String userId);
    public Optional<Track> findByTitle(String userId);
}
