package com.tracks.trackssvc.repository;

import com.tracks.trackssvc.model.Track;
import org.springframework.data.jpa.repository.JpaRepository;

public interface TrackRepository extends JpaRepository<Track, Long> {

}
