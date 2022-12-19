package com.learnciper;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.GeneralSecurityException;
import java.util.Arrays;
import java.util.Base64;

public class DecryptDemo {

    public byte[] decrypt(String key, byte[] crypted) throws GeneralSecurityException {
        byte[] keyBytes = getKeyBytes(key);
        byte[] buf = new byte[16];
        System.arraycopy(keyBytes, 0, buf, 0, Math.max(keyBytes.length, buf.length));
        Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
        cipher.init(Cipher.DECRYPT_MODE, new SecretKeySpec(buf, "AES"), new IvParameterSpec(keyBytes));
        return cipher.doFinal(crypted);
    }

    private byte[] getKeyBytes(String key) {
        byte[] bytes = key.getBytes();
        return bytes.length == 16 ? bytes : Arrays.copyOf(bytes, 16);
    }

    //加密
    public String encrypt(String key, String val) {
        try {
            byte[] origData = val.getBytes();
            byte[] crafted = encrypt(key, origData);
            return Base64.getUrlEncoder().withoutPadding().encodeToString(crafted);
        } catch (Exception e) {
            return "";
        }
    }

    public String decrypt(String key, String val) throws GeneralSecurityException {
        byte[] crypted = Base64.getUrlDecoder().decode(val);
        byte[] origData = decrypt(key, crypted);
        return new String(origData);
    }

    public byte[] encrypt(String key, byte[] origData) throws GeneralSecurityException {
        byte[] keyBytes = getKeyBytes(key);
        byte[] buf = new byte[16];
        System.arraycopy(keyBytes, 0, buf, 0, Math.max(keyBytes.length, buf.length));
        Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
        cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(buf, "AES"), new IvParameterSpec(keyBytes));
        return cipher.doFinal(origData);

    }

    public static void main(String[] args) throws GeneralSecurityException {
        String key = "3wjyxqDPNyrd4QrhxTycRMU4dFN2lCm4";
        DecryptDemo decryptClz = new DecryptDemo();
        String encryptedStr = decryptClz.encrypt(key, "37b63ec62ebf8b2e790b8d9829da2ec26f1fad67");
        // 加密前的纯串 37b63ec62ebf8b2e790b8d9829da2ec26f1fad67
        String pureStr = decryptClz.decrypt(key, "2D_24XH_IMxtYLPGdgZ37X_5_a4eFTpqT7v1aNwq-74Ko6XDpE54LRz4e5s4ehMc");

        String es = String.format("加密后的串: %s", encryptedStr);
        String ds = String.format("解密后的串: %s", pureStr);
        System.out.println(es);
        System.out.println(ds);
    }
}
