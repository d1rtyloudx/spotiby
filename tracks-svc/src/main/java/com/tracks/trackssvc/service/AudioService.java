package com.tracks.trackssvc.service;

import org.springframework.web.multipart.MultipartFile;

import java.util.UUID;

public interface AudioService {
    public String upload(MultipartFile file, String trackId);
}
