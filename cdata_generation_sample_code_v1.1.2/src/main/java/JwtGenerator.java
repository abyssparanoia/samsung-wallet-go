/***********************************************************************************************************************
 * Copyright (c) 2022. Samsung Electronics - All Rights Reserved
 * <p>
 * Unauthorized copying or redistribution of this file in source and binary forms via any medium is strictly prohibited.
 **********************************************************************************************************************/

import java.security.PrivateKey;
import java.security.PublicKey;

public class JwtGenerator {
    private static final String PARTNER_ID = "4059557693262156416"; // partnerCode
    private static final String BOARDING_PASS_CARD_ID = "3glji6d406000";
    private static final String TICKET_CARD_ID = "3gljicv6ssr00";
    private static final String COUPON_CARD_ID = "3gljl0iagsd00";

    public static void main(String[] args) {

        PublicKey samsungPublicKey =
                JwtManager.readCertificate(PayloadUtil.getStringFromFile("sample/securities/Samsung.crt"));
        PublicKey partnerPublicKey =
                JwtManager.readCertificate(PayloadUtil.getStringFromFile("sample/securities/Partner.crt"));
        PrivateKey partnerPrivateKey =
                JwtManager.readPrivateKey(PayloadUtil.getStringFromFile("sample/securities/Partner.key"));

        // Boarding pass
        String plainData = new PayloadUtil("sample/payload/BoardingPass.json").getSampleBoardingPass();

        // Coupon
        // String plainData = new PayloadUtil("sample/payload/Coupon.json").getSampleCoupon();

        // Ticket
        // String plainData = new PayloadUtil("sample/payload/Ticket.json").getSampleTicket();

        String cdata = JwtManager.generate(PARTNER_ID, samsungPublicKey, partnerPublicKey, partnerPrivateKey, plainData);

        System.out.println("Generated CDATA: " + cdata);
        System.out.println("The sample CDATA is never being expired and can be used for Code Lab testing only. " +
                "In real-life application, CDATA must be expired in 30 seconds after creation due to security purposes.");
    }
}
