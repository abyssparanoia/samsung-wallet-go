/***********************************************************************************************************************
 * Copyright (c) 2022. Samsung Electronics - All Rights Reserved
 * <p>
 * Unauthorized copying or redistribution of this file in source and binary forms via any medium is strictly prohibited.
 **********************************************************************************************************************/

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.util.Objects;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

/**
 * This 'PayLoadUtil' class is only used to generate sample JSONs.
 * Please generate on your own JSON.
 */
public class PayloadUtil {
    private static final String DEFAULT_LANGUAGE = "en";
    private String json;

    public PayloadUtil(String path) {
        this.json = getStringFromFile(path);
    }

    public static String getStringFromFile(String path) {
        try {
            File file = new File(Objects.requireNonNull(ClassLoader.getSystemClassLoader().getResource(path)).getFile());
            return new String(Files.readAllBytes(file.toPath()));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public String getJson() {
        return json;
    }

    public String getSampleBoardingPass() {
        return this.setReferenceId()
                .setLanguage()
                .setDefaultTimestamp()
                .setTimestampForSampleBoardingPass()
                .getJson();
    }

    public String getSampleCoupon() {
        return this.setReferenceId()
                .setLanguage()
                .setDefaultTimestamp()
                .setTimestampForSampleCoupon()
                .getJson();
    }

    public String getSampleTicket() {
        return this.setReferenceId()
                .setLanguage()
                .setDefaultTimestamp()
                .setTimestampForSampleTicket()
                .getJson();
    }

    private PayloadUtil setReferenceId() {
        this.json = this.json.replace("{refId}", UUID.randomUUID().toString());
        return this;
    }

    private PayloadUtil setLanguage() {
        return setLanguage(DEFAULT_LANGUAGE);
    }

    private PayloadUtil setLanguage(String language) {
        this.json = this.json.replace("{language}", language);
        return this;
    }

    private PayloadUtil setDefaultTimestamp() {
        long currentTimeMillis = System.currentTimeMillis();

        this.json = this.json
                .replace("{createdAt}", String.valueOf(currentTimeMillis))
                .replace("{updatedAt}", String.valueOf(currentTimeMillis));

        return this;
    }

    private PayloadUtil setTimestampForSampleBoardingPass() {
        long currentTimeMillis = System.currentTimeMillis();

        this.json = this.json
                .replace("{boardingTime}", String.valueOf(currentTimeMillis
                        + TimeUnit.DAYS.toMillis(1)))
                .replace("{gateClosingTime}",
                        String.valueOf(currentTimeMillis
                                + TimeUnit.DAYS.toMillis(1)
                                + TimeUnit.MINUTES.toMillis(50)))
                .replace("{estimatedOrActualStartDate}",
                        String.valueOf(currentTimeMillis
                                + TimeUnit.DAYS.toMillis(1)
                                + TimeUnit.HOURS.toMillis(1)))
                .replace("{estimatedOrActualEndDate}",
                        String.valueOf(currentTimeMillis
                                + TimeUnit.DAYS.toMillis(1)
                                + TimeUnit.HOURS.toMillis(12)));

        return this;
    }

    private PayloadUtil setTimestampForSampleCoupon() {
        long currentTimeMillis = System.currentTimeMillis();

        this.json = this.json
                .replace("{issueDate}", String.valueOf(currentTimeMillis))
                .replace("{expiry}", String.valueOf(currentTimeMillis
                        + TimeUnit.DAYS.toMillis(30)));

        return this;
    }

    private PayloadUtil setTimestampForSampleTicket() {
        long currentTimeMillis = System.currentTimeMillis();

        this.json = this.json
                .replace("{issueDate}", String.valueOf(currentTimeMillis))
                .replace("{startDate}", String.valueOf(currentTimeMillis
                        + TimeUnit.DAYS.toMillis(1)))
                .replace("{endDate}", String.valueOf(currentTimeMillis
                        + TimeUnit.DAYS.toMillis(1) + + TimeUnit.HOURS.toMillis(2)));

        return this;
    }
}
