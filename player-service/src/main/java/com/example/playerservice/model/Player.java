package com.example.playerservice.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.Getter;
import lombok.Setter;

@Entity
@Getter
@Setter
public class Player {
    @Id
    String id;
    @Column(name = "is_playing")
    Boolean isPlaying;
    @Column(name = "playback_time")
    Long playbackTime;
    @Column(name = "current_track_id")
    String currentTrackId;
    @Column(name = "volume")
    Integer volume;
    @Column(name = "manifest_url")
    String manifestUrl;
}
