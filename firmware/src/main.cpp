#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include "config.h"

String tripId = "";
bool tripActive = false;

void connectWiFi() {
    Serial.printf("Connecting to %s", WIFI_SSID);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.printf("Connected! IP: %s\n", WiFi.localIP().toString().c_str());
}

String apiUrl(const char* path) {
    return String("http://") + API_HOST + ":" + API_PORT + path;
}

int apiPost(const char* path, const String& body) {
    HTTPClient http;
    http.begin(apiUrl(path));
    http.addHeader("X-API-Key", API_KEY);
    http.addHeader("Content-Type", "application/json");

    int code = http.POST(body);
    if (code > 0) {
        String response = http.getString();
        Serial.printf("[%d] %s\n", code, response.c_str());

        // Парсим trip_id из ответа на trip/start
        if (String(path).endsWith("trip/start") && code == 200) {
            JsonDocument doc;
            if (deserializeJson(doc, response) == DeserializationError::Ok) {
                tripId = doc["trip_id"].as<String>();
                Serial.printf("Trip started: %s\n", tripId.c_str());
            }
        }
    } else {
        Serial.printf("HTTP error: %s\n", http.errorToString(code).c_str());
    }

    http.end();
    return code;
}

void startTrip() {
    Serial.println("Starting trip...");
    int code = apiPost("/api/device/trip/start", "{}");
    if (code == 200) {
        tripActive = true;
    }
}

void endTrip() {
    Serial.println("Ending trip...");
    apiPost("/api/device/trip/end", "{}");
    tripActive = false;
    tripId = "";
}

void sendFakeLocation() {
    if (!tripActive || tripId.isEmpty()) {
        return;
    }

    // Фейковые координаты: Москва + небольшой случайный сдвиг
    float lat = 55.7558 + random(-100, 100) / 100000.0;
    float lon = 37.6173 + random(-100, 100) / 100000.0;
    float speed = random(0, 80);
    float heading = random(0, 360);
    int satellites = random(6, 12);

    JsonDocument doc;
    doc["trip_id"] = tripId;
    doc["lat"] = lat;
    doc["lon"] = lon;
    doc["speed"] = speed;
    doc["heading"] = heading;
    doc["satellites"] = satellites;

    String body;
    serializeJson(doc, body);

    Serial.printf("Sending: lat=%.5f lon=%.5f speed=%.0f\n", lat, lon, speed);
    apiPost("/api/device/location", body);
}

void setup() {
    Serial.begin(115200);
    delay(1000);
    Serial.println("\n=== Car Tracking ESP32 ===");

    connectWiFi();

    // Автоматически стартуем тестовую поездку
    startTrip();
}

void loop() {
    // Переподключение WiFi если потеряно
    if (WiFi.status() != WL_CONNECTED) {
        Serial.println("WiFi lost, reconnecting...");
        connectWiFi();
    }

    sendFakeLocation();
    delay(SEND_INTERVAL_MS);
}
