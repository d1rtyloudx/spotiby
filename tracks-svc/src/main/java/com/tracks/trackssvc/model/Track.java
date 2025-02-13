package com.tracks.trackssvc.model;


import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;

@Entity
@Table(name = "track", schema = "public")
@NoArgsConstructor
@Getter
@Setter
public class Track {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;
    private String title;
    private Long durationMs;
    private String authorId;
    private String coverUrl; //todo broker: bind queue with exchange
    @Temporal(TemporalType.TIMESTAMP)
    private Date uploadDate;
}
