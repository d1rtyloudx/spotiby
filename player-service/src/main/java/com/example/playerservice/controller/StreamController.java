package com.example.playerservice.controller;

import com.example.playerservice.props.MinioProperties;
import io.minio.GetObjectArgs;
import io.minio.MinioClient;
import io.minio.errors.*;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.InputStreamResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.io.IOException;
import java.io.InputStream;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;

@RestController
@CrossOrigin(origins = "http://localhost:5173")
@RequestMapping("api/v1/stream")
@RequiredArgsConstructor
public class StreamController {

    private final MinioClient minioClient;
    private final MinioProperties minioProperties;
    @GetMapping("/{trackId}/{fileName}")
    public Resource getResource(@PathVariable String fileName, @PathVariable String trackId) {
        try {
            InputStream fileStream = minioClient.getObject(GetObjectArgs.builder()
                    .bucket(minioProperties.getBucket())
                    .object(trackId+"/"+fileName).build());
            return new InputStreamResource(fileStream);
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Unnable to extract object from minio");
        }
    }
}
