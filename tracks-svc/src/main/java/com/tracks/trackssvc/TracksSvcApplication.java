package com.tracks.trackssvc;

import com.tracks.trackssvc.repository.TrackRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class TracksSvcApplication {

    public static void main(String[] args) {
        SpringApplication.run(TracksSvcApplication.class, args);
        
    }

}
