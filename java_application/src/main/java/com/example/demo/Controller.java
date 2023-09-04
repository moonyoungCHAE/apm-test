package com.example.demo;

import com.amazonaws.auth.InstanceProfileCredentialsProvider;
import com.amazonaws.client.builder.AwsClientBuilder;
import com.amazonaws.services.dynamodbv2.AmazonDynamoDB;
import com.amazonaws.services.dynamodbv2.model.ScanRequest;
import com.amazonaws.services.dynamodbv2.model.ScanResult;
import org.apache.catalina.authenticator.BasicAuthenticator;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import com.amazonaws.services.dynamodbv2.AmazonDynamoDBClientBuilder;

import java.util.logging.Logger;

@RestController
public class Controller {
    private static Logger logger = Logger.getLogger("mylogger");

    private RestTemplate restTemplate;

    public Controller() {
        this.restTemplate = new RestTemplate();
    }

    @GetMapping("/ping")
    public String Ping() {
        MyFunction1("hello");
        return "pong";
    }

    @GetMapping("/golang")
    public ResponseEntity<String> SendGolangAPI() {
        ResponseEntity<String> response = restTemplate.getForEntity("http://localhost:8081/golang", String.class);
        return response;
    }

    @GetMapping("/golang/db")
    public ResponseEntity<String> SendGolangDBAPI() {
        ResponseEntity<String> response = restTemplate.getForEntity("http://localhost:8081/golang/db", String.class);
        return response;
    }

    @GetMapping("/dynamodb")
    public Integer DynamoDB() {
        AwsClientBuilder.EndpointConfiguration endpointConfiguration = new AwsClientBuilder.EndpointConfiguration("http://localhost:8000", "ap-northeast-2");

        AmazonDynamoDB amazonDynamoDb = AmazonDynamoDBClientBuilder.standard()
                .withCredentials(InstanceProfileCredentialsProvider.getInstance())
                .withEndpointConfiguration(endpointConfiguration).build();;
        ScanResult result = amazonDynamoDb.scan(new ScanRequest("tommoy_test"));
        return result.getCount();
    }

    private void MyFunction1(String arg1) {
        logger.info("my function 1 arg1 : " + arg1);
        MyFunction2(arg1, Math.random());
    }

    private void MyFunction2(String arg1, double arg2) {
        logger.info("my function 2 arg1 : " + arg1 + ", arg2 : " + arg2);
    }
}

