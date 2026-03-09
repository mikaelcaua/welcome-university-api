package com.welcomeuniversity.provas.service;

import java.awt.Graphics2D;
import java.awt.RenderingHints;
import java.awt.image.BufferedImage;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.Iterator;
import java.util.Locale;

import javax.imageio.IIOImage;
import javax.imageio.ImageIO;
import javax.imageio.ImageWriteParam;
import javax.imageio.ImageWriter;
import javax.imageio.stream.ImageOutputStream;

import org.apache.pdfbox.Loader;
import org.apache.pdfbox.pdfwriter.compress.CompressParameters;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;

@Service
public class UploadCompressionService {

    private static final Logger LOG = LoggerFactory.getLogger(UploadCompressionService.class);
    private static final int MAX_IMAGE_DIMENSION = 2200;

    public S3Service.UploadPayload compressForUpload(MultipartFile file) {
        String originalFilename = file.getOriginalFilename() == null ? "exam" : file.getOriginalFilename();
        String contentType = file.getContentType() == null ? "" : file.getContentType().toLowerCase(Locale.ROOT);

        byte[] originalBytes;
        try {
            originalBytes = file.getBytes();
        } catch (IOException ex) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Falha ao ler arquivo enviado.");
        }

        String lowerFilename = originalFilename.toLowerCase(Locale.ROOT);
        if ("application/pdf".equals(contentType) || lowerFilename.endsWith(".pdf")) {
            byte[] compressedPdf = tryCompressPdf(originalBytes);
            return new S3Service.UploadPayload(originalFilename, "application/pdf", compressedPdf);
        }

        if (contentType.startsWith("image/") || hasImageExtension(lowerFilename)) {
            String format = resolveImageFormat(contentType, lowerFilename);
            if (format.isBlank()) {
                return new S3Service.UploadPayload(
                    originalFilename,
                    contentType.isBlank() ? "application/octet-stream" : contentType,
                    originalBytes
                );
            }
            byte[] compressedImage = tryCompressImage(originalBytes, format);
            String resolvedContentType = switch (format) {
                case "png" -> "image/png";
                case "jpg", "jpeg" -> "image/jpeg";
                default -> contentType.isBlank() ? "application/octet-stream" : contentType;
            };
            return new S3Service.UploadPayload(originalFilename, resolvedContentType, compressedImage);
        }

        return new S3Service.UploadPayload(
            originalFilename,
            contentType.isBlank() ? "application/octet-stream" : contentType,
            originalBytes
        );
    }

    private byte[] tryCompressPdf(byte[] originalBytes) {
        try (PDDocument document = Loader.loadPDF(originalBytes);
             ByteArrayOutputStream output = new ByteArrayOutputStream()) {
            document.save(output, CompressParameters.DEFAULT_COMPRESSION);
            byte[] compressed = output.toByteArray();
            return compressed.length < originalBytes.length ? compressed : originalBytes;
        } catch (IOException ex) {
            LOG.warn("Falha ao comprimir PDF, enviando original. motivo={}", ex.getMessage());
            return originalBytes;
        }
    }

    private byte[] tryCompressImage(byte[] originalBytes, String format) {
        try {
            BufferedImage source = ImageIO.read(new ByteArrayInputStream(originalBytes));
            if (source == null) {
                return originalBytes;
            }

            BufferedImage resized = resizeIfNeeded(source);
            byte[] compressed = writeImage(resized, format);
            return compressed.length < originalBytes.length ? compressed : originalBytes;
        } catch (IOException ex) {
            LOG.warn("Falha ao comprimir imagem, enviando original. motivo={}", ex.getMessage());
            return originalBytes;
        }
    }

    private BufferedImage resizeIfNeeded(BufferedImage source) {
        int width = source.getWidth();
        int height = source.getHeight();
        int largestDimension = Math.max(width, height);
        if (largestDimension <= MAX_IMAGE_DIMENSION) {
            return source;
        }

        double ratio = (double) MAX_IMAGE_DIMENSION / largestDimension;
        int targetWidth = Math.max(1, (int) Math.round(width * ratio));
        int targetHeight = Math.max(1, (int) Math.round(height * ratio));
        int imageType = source.getColorModel().hasAlpha() ? BufferedImage.TYPE_INT_ARGB : BufferedImage.TYPE_INT_RGB;
        BufferedImage resized = new BufferedImage(targetWidth, targetHeight, imageType);
        Graphics2D graphics = resized.createGraphics();
        graphics.setRenderingHint(RenderingHints.KEY_INTERPOLATION, RenderingHints.VALUE_INTERPOLATION_BILINEAR);
        graphics.setRenderingHint(RenderingHints.KEY_RENDERING, RenderingHints.VALUE_RENDER_QUALITY);
        graphics.drawImage(source, 0, 0, targetWidth, targetHeight, null);
        graphics.dispose();
        return resized;
    }

    private byte[] writeImage(BufferedImage image, String format) throws IOException {
        String normalizedFormat = "jpg".equals(format) ? "jpeg" : format;
        BufferedImage outputImage = image;
        if ("jpeg".equals(normalizedFormat) && image.getColorModel().hasAlpha()) {
            outputImage = new BufferedImage(image.getWidth(), image.getHeight(), BufferedImage.TYPE_INT_RGB);
            Graphics2D graphics = outputImage.createGraphics();
            graphics.drawImage(image, 0, 0, null);
            graphics.dispose();
        }

        Iterator<ImageWriter> writers = ImageIO.getImageWritersByFormatName(normalizedFormat);
        if (!writers.hasNext()) {
            return toFormatWithDefaultWriter(outputImage, normalizedFormat);
        }

        ImageWriter writer = writers.next();
        try (ByteArrayOutputStream output = new ByteArrayOutputStream();
             ImageOutputStream imageOutput = ImageIO.createImageOutputStream(output)) {
            writer.setOutput(imageOutput);
            ImageWriteParam param = writer.getDefaultWriteParam();
            if (param.canWriteCompressed()) {
                param.setCompressionMode(ImageWriteParam.MODE_EXPLICIT);
                if ("jpeg".equals(normalizedFormat)) {
                    param.setCompressionQuality(0.78f);
                } else if ("png".equals(normalizedFormat)) {
                    param.setCompressionQuality(0.35f);
                }
            }
            writer.write(null, new IIOImage(outputImage, null, null), param);
            writer.dispose();
            return output.toByteArray();
        }
    }

    private byte[] toFormatWithDefaultWriter(BufferedImage image, String format) throws IOException {
        try (ByteArrayOutputStream output = new ByteArrayOutputStream()) {
            if (!ImageIO.write(image, format, output)) {
                return new byte[0];
            }
            return output.toByteArray();
        }
    }

    private String resolveImageFormat(String contentType, String lowerFilename) {
        if (contentType.equals("image/png") || lowerFilename.endsWith(".png")) {
            return "png";
        }
        if (contentType.equals("image/jpg")
            || contentType.equals("image/jpeg")
            || lowerFilename.endsWith(".jpg")
            || lowerFilename.endsWith(".jpeg")) {
            return "jpeg";
        }
        return "";
    }

    private boolean hasImageExtension(String lowerFilename) {
        return lowerFilename.endsWith(".png")
            || lowerFilename.endsWith(".jpg")
            || lowerFilename.endsWith(".jpeg")
            || lowerFilename.endsWith(".gif")
            || lowerFilename.endsWith(".webp");
    }
}
