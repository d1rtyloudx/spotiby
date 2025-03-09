package com.tracks.trackssvc.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.tracks.trackssvc.service.props.MinioProperties;
import io.minio.MinioClient;
import lombok.RequiredArgsConstructor;
import org.apache.kafka.clients.admin.NewTopic;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.serialization.StringSerializer;
import org.springframework.amqp.core.*;
import org.springframework.boot.autoconfigure.kafka.KafkaProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;
import org.springframework.kafka.core.DefaultKafkaProducerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.core.ProducerFactory;

import java.util.HashMap;
import java.util.Map;


@Configuration
@RequiredArgsConstructor
public class ApplicationConfig {
    private final MinioProperties minioProperties;
    private final KafkaProperties kafkaProperties;


    @Bean
    public MinioClient minioClient() {
        return MinioClient.builder()
                .endpoint(minioProperties.getUrl())
                .credentials(minioProperties.getAccessKey(), minioProperties.getSecretKey())
                .build();
    }

    //RabbitMQ beans

    @Bean
    public Queue imageTrackQueue() {
        String queueName = "track_image";
        return new Queue(queueName, true);
    }

    @Bean
    public Queue imagePlaylistQueue() {
        String queueName = "playlist_image";
        return new Queue(queueName, true);
    }

    @Bean
    Exchange exchange() {
        String exchangeName = "image";
        return new DirectExchange(exchangeName, true, false);
    }

    @Bean
    Binding bindingTrack(Queue imageTrackQueue, Exchange exchange) {
        String trackRoutingKey = "track";
        return BindingBuilder.bind(imageTrackQueue).to(exchange).with(trackRoutingKey).noargs();
    }


    @Bean
    Binding bindingPlaylist(Queue imagePlaylistQueue, Exchange exchange) {
        String trackRoutingKey = "playlist";
        return BindingBuilder.bind(imagePlaylistQueue).to(exchange).with(trackRoutingKey).noargs();
    }


    // Json object mapper
    @Bean
    ObjectMapper objectMapper() {
        return new ObjectMapper();
    }



    // Kafka beans and properties
    @Bean
    public ProducerFactory<String, String> producerFactory() {
        Map<String, Object> configProps = new HashMap<>();
        configProps.put(
                ProducerConfig.BOOTSTRAP_SERVERS_CONFIG,
                kafkaProperties.getBootstrapServers());
        configProps.put(
                ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,
                StringSerializer.class);
        configProps.put(
                ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG,
                StringSerializer.class);
        return new DefaultKafkaProducerFactory<>(configProps);
    }

    @Bean
    public KafkaTemplate<String, String> kafkaTemplate() {
        return new KafkaTemplate<>(producerFactory());
    }

    @Bean
    public NewTopic topicTracks() {
        return TopicBuilder.name("track_create").build();
    }

    @Bean
    public NewTopic topicPlaylists() {
        return TopicBuilder.name("playlist_create").build();
    }
}
