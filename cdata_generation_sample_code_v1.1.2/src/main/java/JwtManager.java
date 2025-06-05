/***********************************************************************************************************************
 * Copyright (c) 2022. Samsung Electronics - All Rights Reserved
 * <p>
 * Unauthorized copying or redistribution of this file in source and binary forms via any medium is strictly prohibited.
 **********************************************************************************************************************/

import com.nimbusds.jose.EncryptionMethod;
import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWEAlgorithm;
import com.nimbusds.jose.JWEHeader;
import com.nimbusds.jose.JWEObject;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSObject;
import com.nimbusds.jose.JWSSigner;
import com.nimbusds.jose.Payload;
import com.nimbusds.jose.crypto.RSAEncrypter;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.jwk.RSAKey;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.nio.charset.StandardCharsets;
import java.security.KeyFactory;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemReader;

/**
 * This class does not need to be modified.
 */
public class JwtManager {
    private static final String CARD = "CARD";
    private static final String VERSION = "3";
    private static final String PARTNER_ID_KEY = "partnerId";
    private static final String VERSION_KEY = "ver";
    private static final String UTC_KEY = "utc";
    private static final String RSA = "RSA";
    private static final String CERTIFICATE_ID = "YMtt"; 
    private static final String CERTIFICATE_ID_KEY = "certificateId";

    public static String generate(String partnerId, PublicKey samsungPublicKey, PublicKey partnerPublicKey,
                                  PrivateKey partnerPrivateKey, String data) {
        EncryptionMethod jweEnc = EncryptionMethod.A128GCM;
        JWEAlgorithm jweAlg = JWEAlgorithm.RSA1_5;
        JWEHeader jweHeader = new JWEHeader.Builder(jweAlg, jweEnc).build();
        RSAEncrypter encryptor = new RSAEncrypter((RSAPublicKey) samsungPublicKey);
        JWEObject jwe = new JWEObject(jweHeader, new Payload(data));
        try {
            jwe.encrypt(encryptor);
        } catch (JOSEException e) {
            e.printStackTrace();
        }
        String payload = jwe.serialize();

        JWSAlgorithm jwsAlg = JWSAlgorithm.RS256;
        Long utc = System.currentTimeMillis();

        JWSHeader jwsHeader = new JWSHeader.Builder(jwsAlg)
                .contentType(CARD)
                .customParam(PARTNER_ID_KEY, partnerId)
                .customParam(VERSION_KEY, VERSION)
				.customParam(CERTIFICATE_ID_KEY, CERTIFICATE_ID)
                .customParam(UTC_KEY, utc)
                .build();

        JWSObject jwsObj = new JWSObject(jwsHeader, new Payload(payload));

        RSAKey rsaJWK
                = new RSAKey.Builder((RSAPublicKey) partnerPublicKey)
                .privateKey(partnerPrivateKey)
                .build();
        JWSSigner signer;
        try {
            signer = new RSASSASigner(rsaJWK);
            jwsObj.sign(signer);
        } catch (JOSEException e) {
            e.printStackTrace();
        }
        return jwsObj.serialize();
    }

    public static PrivateKey readPrivateKey(String key) {

        final String PKCS_1_PEM_HEADER = "-----BEGIN RSA PRIVATE KEY-----";
        final String PKCS_8_PEM_HEADER = "-----BEGIN PRIVATE KEY-----";

        byte[] keyByte = readKeyByte(key);
        PrivateKey privateKey = null;
        if (key.contains(PKCS_1_PEM_HEADER)) {
            byte[] pkcs1Bytes = keyByte;
            int pkcs1Length = pkcs1Bytes.length;
            int totalLength = pkcs1Length + 22;
            byte[] pkcs8Header = new byte[]{
                    0x30, (byte) 0x82, (byte) ((totalLength >> 8) & 0xff), (byte) (totalLength & 0xff),
                    // Sequence + total length
                    0x2, 0x1, 0x0, // Integer (0)
                    0x30, 0xD, 0x6, 0x9, 0x2A, (byte) 0x86, 0x48, (byte) 0x86, (byte) 0xF7, 0xD, 0x1, 0x1, 0x1, 0x5,
                    0x0, // Sequence: 1.2.840.113549.1.1.1, NULL
                    0x4, (byte) 0x82, (byte) ((pkcs1Length >> 8) & 0xff), (byte) (pkcs1Length & 0xff)
                    // Octet string + length
            };
            keyByte = join(pkcs8Header, pkcs1Bytes);
        }
        PKCS8EncodedKeySpec pkcs8Spec = new PKCS8EncodedKeySpec(keyByte);
        try {
            KeyFactory kf = KeyFactory.getInstance(RSA);
            privateKey = kf.generatePrivate(pkcs8Spec);
        } catch (InvalidKeySpecException | NoSuchAlgorithmException e) {
            e.printStackTrace();
        }

        return privateKey;
    }

    public static PublicKey readPublicKey(String key) {
        PublicKey publicKey = null;
        byte[] keyByte = readKeyByte(key);
        X509EncodedKeySpec x509Spec = new X509EncodedKeySpec(keyByte);
        try {
            KeyFactory kf = KeyFactory.getInstance(RSA);
            publicKey = kf.generatePublic(x509Spec);
        } catch (NoSuchAlgorithmException | InvalidKeySpecException e) {
            e.printStackTrace();
        }
        return publicKey;
    }

    public static PublicKey readCertificate(String cert) {
        Certificate certificate = null;
        byte[] keyByte = readKeyByte(cert);
        InputStream is = new ByteArrayInputStream(keyByte);
        try {
            CertificateFactory cf = CertificateFactory.getInstance("X.509");
            certificate = cf.generateCertificate(is);
        } catch (CertificateException e) {
            e.printStackTrace();
        }
        assert certificate != null;
        return certificate.getPublicKey();
    }

    public static byte[] readKeyByte(String key) {
        byte[] keyByte;
        ByteArrayInputStream bais = new ByteArrayInputStream(key.getBytes(StandardCharsets.UTF_8));
        Reader reader = new InputStreamReader(bais, StandardCharsets.UTF_8);
        PemReader pemReader = new PemReader(reader);
        PemObject pemObject = null;
        try {
            pemObject = pemReader.readPemObject();
        } catch (IOException e) {
            e.printStackTrace();
        }
        if (pemObject == null) {
            keyByte = Base64.getDecoder().decode(key);
        } else {
            keyByte = pemObject.getContent();
        }

        return keyByte;
    }

    private static byte[] join(byte[] byteArray1, byte[] byteArray2) {
        byte[] bytes = new byte[byteArray1.length + byteArray2.length];
        System.arraycopy(byteArray1, 0, bytes, 0, byteArray1.length);
        System.arraycopy(byteArray2, 0, bytes, byteArray1.length, byteArray2.length);
        return bytes;
    }
}
