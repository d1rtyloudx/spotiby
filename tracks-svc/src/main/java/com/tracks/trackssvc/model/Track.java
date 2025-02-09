package com.tracks.trackssvc.model;


import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;
import java.util.UUID;

@Entity
@Table(name = "track", schema = "public")
@NoArgsConstructor
@Getter
@Setter
public class Track {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private UUID id;
    private String title;
    private Long durationMs;
    private Long authorId;
    private String coverUri; //todo broker: bind queue with exchange
    @Temporal(TemporalType.TIMESTAMP)
    private Date uploadDate;
}
