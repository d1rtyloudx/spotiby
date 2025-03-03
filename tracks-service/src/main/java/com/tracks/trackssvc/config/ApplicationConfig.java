package com.tracks.trackssvc.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracks.trackssvc.service.props.MinioProperties;
import com.tracks.trackssvc.web.dto.UpdateTrackCoverDto;
import io.minio.MinioClient;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.core.*;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;


@Configuration
@RequiredArgsConstructor
public class ApplicationConfig {
    private final MinioProperties minioProperties;
    private final String queueName = "track_image";
    private final String exchangeName = "image";
    private final String trackRoutingKey = "track";


    @Bean
    public MinioClient minioClient() {
        return MinioClient.builder()
                .endpoint(minioProperties.getUrl())
                .credentials(minioProperties.getAccessKey(), minioProperties.getSecretKey())
                .build();
    }

    @Bean
    public Queue imageQueue() {
        return new Queue(queueName, true);
    }

    @Bean
    Exchange exchange() {
        return new DirectExchange(exchangeName, true, false);
    }

    @Bean
    Binding binding(Queue imageQueue, Exchange exchange) {
        return BindingBuilder.bind(imageQueue).to(exchange).with(trackRoutingKey).noargs();
    }

    @Bean
    ObjectMapper objectMapper() {
        return new ObjectMapper();
    }
}
