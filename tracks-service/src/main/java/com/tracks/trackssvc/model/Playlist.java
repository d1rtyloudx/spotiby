package com.tracks.trackssvc.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;
import java.util.HashSet;
import java.util.Set;

@Entity
@Table(name = "playlist", schema = "public")
@NoArgsConstructor
@Getter
@Setter
public class Playlist {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;
    private String title;
    private String authorId;
    private String coverUrl;
    @Temporal(TemporalType.TIMESTAMP)
    private Date createdAt;
    @ManyToMany(cascade = {CascadeType.PERSIST, CascadeType.MERGE})
    @JoinTable(
            name = "playlist_track",
            joinColumns = @JoinColumn(name = "playlist_id"),
            inverseJoinColumns = @JoinColumn(name = "track_id")
    )
    private Set<Track> tracks = new HashSet<>();
}