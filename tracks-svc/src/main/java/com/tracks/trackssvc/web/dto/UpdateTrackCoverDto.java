package com.tracks.trackssvc.web.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class UpdateTrackCoverDto {
    @JsonProperty("id")
    String id;
    @JsonProperty("avatar_url")
    String avatarUrl;
}
