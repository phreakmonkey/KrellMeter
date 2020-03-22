#include <Arduino.h>

#define pinA 5
#define pinB 6

const uint8_t maxChars = 20;
char receivedChars[maxChars];   // an array to store the received data
boolean newData = false;
unsigned long lastwrite = millis();
bool idle = false;

void recvData() {
    static byte ndx = 0;
    char rc;
    
    if (Serial.available() > 0) {
        rc = Serial.read();

        if (rc == 'A' || rc == 'B' || (rc > 47 && rc < 58 )) {
              receivedChars[ndx] = rc;
              Serial.print(rc);
              ndx++;
              if (ndx >= maxChars) {
                  ndx = maxChars - 1;
              }
        }
        else if (rc == '\r' || rc == '\n' || rc == ':')
        {
            receivedChars[ndx] = '\0'; // terminate the string
            ndx = 0;
            newData = true;
            Serial.print('\n');
        }
    }
}

void updateMeters() {
    int newVal = 0;
    if (newData == true) {
      lastwrite = millis();
      idle = false;
      Serial.println("Got data");
      newVal = atoi(receivedChars+1);
      if (receivedChars[0] == 'A') {
        analogWrite(pinA, newVal);
        Serial.print("A");
        Serial.println(newVal);
      } else if (receivedChars[0] == 'B') {
        analogWrite(pinB, newVal);
        Serial.print("B");
        Serial.println(newVal);
      }
      newData = false;
    }
}

void setup() {
  delay(5000);
  Serial.begin(115200);
  Serial.println("SerialPWM 1.0");
  pinMode(pinA, OUTPUT);
  pinMode(pinB, OUTPUT);
  digitalWrite(pinA, LOW);
  digitalWrite(pinB, LOW);
  analogWrite(pinA, 512);
  delay(1000);
  analogWrite(pinA, 0);
  analogWrite(pinB, 512);
  delay(1000);
  analogWrite(pinB, 0);
}

void loop() {
  recvData();
  updateMeters();
  if ((!idle) && millis() - lastwrite > 5000) {
    idle = true;
    analogWrite(pinA, 0);
    analogWrite(pinB, 0);
  }
}