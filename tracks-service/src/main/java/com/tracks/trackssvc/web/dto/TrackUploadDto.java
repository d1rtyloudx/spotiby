package com.tracks.trackssvc.web.dto;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.web.multipart.MultipartFile;

@Getter
@Setter
@NoArgsConstructor
public class TrackUploadDto {
    private String title;
    private String authorId;
    private MultipartFile audioFile;
}
