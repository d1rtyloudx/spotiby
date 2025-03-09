package com.tracks.trackssvc.web.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class UpdateCoverDto {
    @JsonProperty("id")
    String id;
    @JsonProperty("avatar_url")
    String avatarUrl;
}
