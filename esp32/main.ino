//Koding RTC -------------------------------------->
#include <Wire.h>
#include <RTClib.h>
RTC_DS1307 rtc;

//Koding LCD --------------------------------------->
#include <LiquidCrystal_I2C.h>
LiquidCrystal_I2C lcd(0x27, 20, 4);

//Koding OLED --------------------------------------->
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

#define SCREEN_WIDTH 128 // OLED width,  in pixels
#define SCREEN_HEIGHT 64 // OLED height, in pixels

// create an OLED display object connected to I2C
Adafruit_SSD1306 oled(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, -1);

//Koding Servo --------------------------------------->
#include <ESP32Servo.h>
const int servoPin = 5;
Servo mekanik;

//Koding Sensor dan LED --------------------------------------->
// proses include library
#include "DHT.h"
// deklarasi variable
// set pin yang digunakan
#define DHTPIN 15
#define DHTTYPE DHT22 // DHT 22 (AM2302), AM2321
DHT dht(DHTPIN, DHTTYPE);

//pin analog
const int pinSensor = A0;
int adcValue = 0;

#include "HX711.h"
const int pinDOUT = 2;
const int pinSCK = 4;
HX711 scale;

//pin LED RGB 
#define LEDPIN1 23
#define LEDPIN2 18
#define LEDPIN3 19

//inisialisasi milis untuk real time
unsigned long sebelum = 0;

// Publisher ke web
// wifi
#include <WiFi.h>
#include <PubSubClient.h>     // Library untuk mengkoneksikan ESP32 ke MQTT broker

const char* ssid = "Wokwi-GUEST";
const char* password = "";
char *mqttServer = "broker.hivemq.com";
int mqttPort = 1883;

char clientId[50];

WiFiClient wifiClient;                // Membuat objek wifiClient
PubSubClient mqttClient(wifiClient);  // Membuat  objek mqttClient dengan konstruktor objek WiFiClient (Permintaan dari Lib)


// mengatur dan menginisialisasi koneksi ke broker MQTT 
// serta menetapkan fungsi callback yang akan dipanggil 
// ketika pesan diterima oleh klien dari langganan yang dibuat.
void setupMQTT() {
  mqttClient.setServer(mqttServer, mqttPort); // Mengatur detail broker target  yang digunakan
  mqttClient.setCallback(callback);           // jika kita ingin menerima pesan untuk langganan yang dibuat oleh klien
}

// untuk melakukan koneksi ulang (reconnect) ke broker MQTT 
// jika klien kehilangan koneksi dengan broker 
// atau gagal melakukan koneksi saat pertama kali menjalankan program.
// ESP32 Reconnect to broker
void reconnect() {
  Serial.println("Connecting to MQTT Broker...");
  while (!mqttClient.connected()) {
      Serial.println("Reconnecting to MQTT Broker..");
      String clientId = "ESP32Client-";
      clientId += String(random(0xffff), HEX);
      
      if (mqttClient.connect(clientId.c_str())) {
        Serial.println("Connected.");
      }      
  }
}

//  fungsi yang akan dipanggil ketika klien menerima pesan dari broker 
//  MQTT yang sesuai dengan langganan yang telah dibuat sebelumnya. 
void callback(char* topic, byte* message, unsigned int length) {
  Serial.print("Callback - ");
  Serial.print("Message:");
  for (int i = 0; i < length; i++) {
    Serial.print((char)message[i]);
  }
}

void connectToInternet(){
  WiFi.begin(ssid, password);   // Mencoba connect ke Wifi

  // Melakukan pengecekan terhadap status koneksi ke WI-Fi
  while (WiFi.status() != WL_CONNECTED) {
      delay(500);
      Serial.print(".");
    } 
    Serial.println("");
    Serial.println("Connected to Wi-Fi");
}


void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  Serial.println("EDSPERT IoT Bootcamp kelompok 3");
  // inisiasi sensor DHT
  dht.begin();

  connectToInternet();
  setupMQTT(); // setup koneksi ke broker

  // inisiasi pin led rgb
  pinMode(LEDPIN1, OUTPUT);
  pinMode(LEDPIN2, OUTPUT);
  pinMode(LEDPIN3, OUTPUT);

  // inisiasi pin sensor massa
  scale.begin(pinDOUT, pinSCK);
  scale.set_scale(0.42);
  scale.tare();

  // inisiasi lcd20x4
  lcd.init();
  lcd.backlight();

  //Koding Servo --------------------------------------->
  mekanik.attach(servoPin);
  mekanik.write(0);

  //Koding RTC -------------------------------------->
  Wire.begin();
  if(!rtc.begin()){
    Serial.println("RTC Tidak Terhubung");
    lcd.setCursor(0,0);
    lcd.print("RTC Tidak Konek");
    while (1);
  }

  lcd.clear();
 
  // initialize OLED display with I2C address 0x3C
  if (!oled.begin(SSD1306_SWITCHCAPVCC, 0x3C)) {
    Serial.println(F("failed to start SSD1306 OLED"));
    while (1);
  }
  lcd.print("Pakan Ayam IoT");  
  delay(2000);         // wait two seconds for initializing
  lcd.clear();
  oled.clearDisplay(); // clear display  
  oled.setTextSize(1);         // set text size
  oled.setTextColor(WHITE);    // set text color
  oled.setCursor(0, 2);       // set position to display (x,y)
  oled.println("Monitor Kandang Ayam"); // set text
  oled.display();
  oled.clearDisplay(); // clear display  
}

void loop() {
  
  //mendefinisikan waktu
  DateTime now = rtc.now();
  int tahun = now.year() % 100;
  int jam = now.hour();
  int menit = now.minute();
  int detik = now.second();

  //mendefinisikan program sensor
  float h = dht.readHumidity();
  float t = dht.readTemperature();

  adcValue = analogRead(pinSensor);
  // put your main code here, to run repeatedly:
  
  long reading = scale.get_units(3);
  float kg = float(reading) / 1000;
  unsigned long sekarang = millis();
  if(sekarang - sebelum >= 1000)
  {
    sebelum = sekarang;  
    sensortampil();

  }

  
/*jadwal makan pagi*/
  if( (jam == 6) && (menit == 15) && (detik == 1) ){    
    makan();
  }

/*jadwal makan sore */
  if( (jam == 15) && (menit == 15) && (detik == 1) ){    
    makan();
  }

/*jadwal makan malam*/
  if( (jam == 19) && (menit == 45) && (detik == 1) ){    
    makan();
  }


  if(t>=40.00){
Serial.println("WARNING");
digitalWrite(LEDPIN1, HIGH);
} else {
Serial.println("NORMAL");
digitalWrite(LEDPIN1, LOW);
digitalWrite(LEDPIN3, HIGH);
}

if(adcValue>=3413){
Serial.println("WARNING");
digitalWrite(LEDPIN2, HIGH);
} else {
Serial.println("NORMAL");
digitalWrite(LEDPIN2, LOW);
digitalWrite(LEDPIN3, HIGH);
}

if (!mqttClient.connected()){
    reconnect();   // Try to connect with broker
  }
  else{
    // Send the Data
    adcValue = analogRead(pinSensor);
    char amoString[8];
    dtostrf(adcValue, 1, 2, amoString);  //Convert float to String
    Serial.print("NH3: ");
    Serial.println(amoString);
    mqttClient.publish("edspertkel3/amo", amoString);


    float temperature = dht.readTemperature();
    char tempString[8];
    dtostrf(temperature, 1, 2, tempString);  //Convert float to String
    Serial.print("Temperature: ");
    Serial.println(tempString);
    mqttClient.publish("edspertkel3/temp", tempString);

  
    float humidity = dht.readHumidity();
    char humString[8];
    dtostrf(humidity, 1, 2, humString);
    Serial.print("Humidity: ");
    Serial.println(humString);
    mqttClient.publish("edspertkel3/hum", humString);

    long reading = scale.get_units(3);
    float kg = float(reading) / 1000;
    char massString[8];
    dtostrf(kg, 1, 2, massString);
    Serial.print("Massa: ");
    Serial.println(massString);
    mqttClient.publish("edspertkel3/mass", massString);

    delay(2000);
  }

}


void sensortampil(){
    //mendefinisikan program sensor
  float h = dht.readHumidity();
  float t = dht.readTemperature();

  adcValue = analogRead(pinSensor);
  // put your main code here, to run repeatedly:
  
  long reading = scale.get_units(3);
  float kg = float(reading) / 1000;

  DateTime now = rtc.now();
  int tahun = now.year() % 100;

  lcd.setCursor(0,0);
  lcd.print("Pakan Ayam IoT");
  lcd.setCursor(0,1);
  lcd.print(String() + now.day() + "/" + now.month() + "/" + tahun);
  lcd.print(" ");
  lcd.print(String() + now.hour() + ":" + now.minute() + ":" + now.second());
  lcd.print(" ");
  lcd.setCursor(0,2);
  lcd.print("Massa: " + String(kg, 2) + "kg");

  oled.setCursor(0, 2);       // set position to display (x,y)
  oled.println("Monitoring Sensor"); // set text
  oled.setCursor(0, 11);       // set position to display (x,y)
  oled.println("T: " + String(t, 2) + char(247) +"C"); // set text
  oled.setCursor(0, 20);       // set position to display (x,y)
  oled.println("%RH : " + String(h, 2) + "%"); // set text
  oled.setCursor(0, 29);       // set position to display (x,y)
  oled.println("NH3: " + String(adcValue, 2) + "PPM"); // set text
  oled.display();              // display on OLED
}

void makan(){
  lcd.setCursor(0,3);
    lcd.print("Feeding Time...");
    mekanik.write(180);
    delay(3000);
    mekanik.write(0);
    lcd.clear();
}