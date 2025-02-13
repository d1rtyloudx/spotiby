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
import ws.schild.jave.Encoder;
import ws.schild.jave.MultimediaObject;
import ws.schild.jave.encode.AudioAttributes;
import ws.schild.jave.encode.EncodingAttributes;


import java.io.*;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class AudioServiceImpl implements AudioService {
    private final MinioProperties minioProperties;
    private final MinioClient minioClient;


    @Override
    public String upload(MultipartFile file, String trackId) {
        try {
            createBucket();
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Image upload failed " + e.getMessage());
        }
        if (file.isEmpty() || file.getOriginalFilename() == null || (!file.getOriginalFilename().endsWith(".mp3") && !file.getOriginalFilename().endsWith(".mp4"))) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Wrong file name or extension");
        }
        String fileName = trackId + ".mp4";

        if(file.getOriginalFilename().endsWith(".mp3")) {
            File convertedFile = convertMp3ToAac(file);
            try {
                saveAudio(new FileInputStream(convertedFile), fileName);
            } catch (FileNotFoundException e) {
                throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "File invalid" + e.getMessage());
            } finally {
                convertedFile.delete();
            }
            return fileName;
        }

        InputStream inputStream;
        try {
            inputStream = file.getInputStream();
        } catch (IOException e) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "File invalid " + e.getMessage());
        }
        saveAudio(inputStream, fileName);

        return fileName;
    }

    @SneakyThrows
    private File convertMp3ToAac(MultipartFile file) {

        File source = File.createTempFile("source", ".mp3");
        file.transferTo(source);

        File target = File.createTempFile("target", ".aac");
        System.setProperty("jave.ffmpeg.location", "/usr/bin/ffmpeg");
        System.setProperty("jave.disable.ffmpeg.extractor", "true");
        try {
            AudioAttributes audio = new AudioAttributes();
            audio.setCodec("aac");
            audio.setBitRate(128000);
            audio.setChannels(2);
            audio.setSamplingRate(44100);

            EncodingAttributes attrs = new EncodingAttributes();
            attrs.setOutputFormat("mp4");
            attrs.setAudioAttributes(audio);
            Encoder encoder = new Encoder(() -> "/usr/bin/ffmpeg");
            encoder.encode(new MultimediaObject(source), target, attrs);
            return target;
        } catch (Exception e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Error while converting mp3 to AAC " + e.getMessage());
        } finally {
            source.delete();
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

    @SneakyThrows
    private void saveAudio(InputStream inputStream, String fileName) {
        minioClient.putObject(PutObjectArgs.builder()
                .stream(inputStream, inputStream.available(), -1)
                .bucket(minioProperties.getBucket())
                .object(fileName)
                .build());
        inputStream.close();
    }
}
