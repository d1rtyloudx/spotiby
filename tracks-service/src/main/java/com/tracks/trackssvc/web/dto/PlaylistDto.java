package com.tracks.trackssvc.web.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.*;

import java.util.Date;
import java.util.Set;

@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
@Data
public class PlaylistDto {
    private String id;
    private String title;
    @JsonProperty("author_id")
    private String authorId;
    @JsonProperty("cover_url")
    private String coverUrl;
    @JsonProperty("created_at")
    private Date createdAt;
    private Set<TrackDto> tracks;
}
