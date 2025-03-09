package com.tracks.trackssvc.service.impl;

import com.tracks.trackssvc.service.AudioService;
import com.tracks.trackssvc.service.props.MinioProperties;
import io.minio.BucketExistsArgs;
import io.minio.MakeBucketArgs;
import io.minio.MinioClient;
import io.minio.PutObjectArgs;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Path;

@Service
@RequiredArgsConstructor
public class AudioServiceImpl implements AudioService {
    private final MinioProperties minioProperties;
    private final MinioClient minioClient;

    @Override
    public String upload(MultipartFile file, String trackId) {
        try {
            File source = File.createTempFile("source", ".mp3");
            file.transferTo(source);
            File dashDir = Files.createTempDirectory("dash_output").toFile();

            String manifestPath = new File(dashDir, "manifest.mpd").getAbsolutePath();

            ProcessBuilder pb = new ProcessBuilder(
                    "ffmpeg",
                    "-i", source.getAbsolutePath(),
                    "-c:a", "aac",
                    "-b:a", "128k",
                    "-ac", "2",
                    "-ar", "44100",
                    "-seg_duration", "10",
                    "-f", "dash",
                    manifestPath
            );
            pb.redirectErrorStream(true);
            Process process = pb.start();
            int exitCode = process.waitFor();
            if (exitCode != 0) {
                throw new RuntimeException("FFmpeg failed with code " + exitCode);
            }

            uploadDashFiles(dashDir, trackId);

            source.delete();
            Files.walk(dashDir.toPath()).map(Path::toFile).forEach(File::delete);

            return "success";
        } catch (Exception e) {
            throw new RuntimeException("DASH conversion and upload failed: " + e.getMessage());
        }
    }

    public void uploadDashFiles(File dashDir, String trackId) {
        try {
            createBucket();

            File[] files = dashDir.listFiles();
            if (files == null) {
                throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "DASH directory is empty or not found.");
            }

            for (File file : files) {
                String objectName = trackId + "/" + file.getName();  // tracks/12345/manifest.mpd
                try (InputStream inputStream = new FileInputStream(file)) {
                    minioClient.putObject(
                            PutObjectArgs.builder()
                                    .bucket(minioProperties.getBucket())
                                    .object(objectName)
                                    .stream(inputStream, file.length(), -1)
                                    .build()
                    );
                }
            }
        } catch (Exception e) {
            throw new RuntimeException("Error uploading DASH files to MinIO: " + e.getMessage(), e);
        }
    }

    @SneakyThrows
    private void createBucket() {
        boolean found = minioClient.bucketExists(BucketExistsArgs.builder()
                .bucket(minioProperties.getBucket())
                .build());
        if (!found) {
            minioClient.makeBucket(MakeBucketArgs.builder()
                    .bucket(minioProperties.getBucket())
                    .build());
        }
    }
}
