// The controller wants 38400bbs 8bits, odd parity, 1 stop bit
#define SERIAL_CONFIG 38400, SERIAL_8O1

// OE pins for the rs422 drivers/receivers
static int oe_pins[6] = {A0, A1, A2, A3, A4, A5};

// Ping packets
static uint8_t ping[4] = {0xE0, 0x00, 0x00, 0x00};
static uint8_t pong[4] = {0xE0, 0xFF, 0xFF, 0xF0};

// Valid message start bytes
static uint8_t valid_start[] = {0x80, 0x84,0x86,0xE0};

volatile uint8_t msg[3] = {0,0,0};
volatile uint8_t msg_field = 0;
volatile uint8_t control_cache[7] = {255,255,255,255,255,255,255};


static uint8_t segment_sequence[] = {0xFF - 0x01, 0xFF - 0x02, 0xFF - 0x04, 0xFF - 0x08, 0xFF - 0x10, 0xFF - 0x20};

void setup() {
  
  Serial.begin(57600);
  Serial1.begin(38400);
  Serial2.begin(SERIAL_CONFIG);
  Serial3.begin(SERIAL_CONFIG);

  // Enable the receivers and drivers
  for (uint8_t i = 0; i < 6; i++){
    int pin = oe_pins[i];
    pinMode(pin, OUTPUT);
    digitalWrite(pin, HIGH);
  }
  usbMIDI.setHandleNoteOff(myNoteOff);
  usbMIDI.setHandleNoteOn(myNoteOn);
}

static uint8_t lamp_offset = 0x00;
static uint8_t lamp_bit = 0x00;
static unsigned long last_run = 0;
#define INTERVAL 300

void loop() {

  // Link to the computer
  while (Serial.available()) {
    uint8_t in = Serial.read();
    Serial1.write(in); // Rom2
    // Serial2.write(in); // Rom1
    // Serial3.write(in);
    // lamp_tester();
    // lcd_tester();
  }
  
  while (Serial1.available()) {
    uint8_t in = Serial1.read();
    Serial.write(in);
    parse_midi(in);
  }

  while (Serial2.available()) {
    uint8_t in = Serial2.read();
    // Serial.write(in);
  }

  while (Serial3.available()) {
    uint8_t in = Serial3.read();
    Serial.write(in);
  }
  usbMIDI.read();
  
  unsigned long t;
  t = millis();
  if (last_run + INTERVAL < t) {
    last_run = t;
    // lcd_tester();
    // lamp_tester();
    segment_tester();
  }
}


void parse_midi(uint8_t in) {
  if (msg_field == 0 && (in & 0x80) != 0) {
    msg[0] = in;
    msg_field++;
  } else if (msg_field != 0) {
    msg[msg_field] = in;
    msg_field++;
  }
  if ((msg[0] & 0x90) == 0x90 && msg_field == 2) {
    // Keyboard message
    uint8_t channel = msg[0] - 0x90;
    if ((msg[1] & 0x40) == 0x40) {
      usbMIDI.sendNoteOn(msg[1] - 0x40, 64, channel);
      // note_to_led(channel, msg[1] - 0x40, true);
    } else {
      usbMIDI.sendNoteOff(msg[1], 64, channel);
      // note_to_led(channel, msg[1], false);
    }
    msg_field = 0;
  } else if ((msg[0] & 0xA0) == 0xA0 && msg_field==3) {
    // Analog message
    uint8_t control = msg[0] - 0xA0;
    if (control > 6) {
      // invalid control number
      
    } else if (control_cache[control] != msg[1]) {
      usbMIDI.sendControlChange(control, msg[1], 1);
      control_cache[control] = msg[1];
    }
    msg_field == 0;
  }

  if (msg_field > 2) {
    msg_field = 0;
  }
}

uint8_t segment_number;
void segment_tester() {
  Serial1.write(0x83);
  Serial1.write(segment_sequence[(segment_number+2) % 6]);
  Serial1.write(segment_sequence[(segment_number+1) % 6]);  
  Serial1.write(segment_sequence[segment_number]);
  Serial1.write(0x84);
  Serial1.write(segment_sequence[(segment_number+2) % 6]);
  Serial1.write(segment_sequence[(segment_number+1) % 6]);  
  Serial1.write(segment_sequence[segment_number]);
  segment_number = (segment_number + 1) % 6;
}

uint8_t lcd_char = 'A';
void lcd_tester() {
  Serial2.write(0x81);
  Serial2.write(lcd_char);
  /*
  Serial.print("\n");
  Serial.print("Wrote lcd char: ");
  Serial.write(lcd_char);
  */
  lcd_char++;
  if (lcd_char > 'Z') {
    lcd_char = 'A';
  }
}

void lamp_tester() {
    uint8_t lamp_byte = 0x01 << lamp_bit;
    Serial2.write(0x80);
    Serial2.write(lamp_offset);
    Serial2.write(lamp_byte);
    /*
    Serial.print("\n");
    Serial.print("Sent: ");
    Serial.print(lamp_offset);
    Serial.print(" - ");
    Serial.print(lamp_byte, HEX);
    Serial.print("\n");
    */
    lamp_bit++;
    if (lamp_bit == 0x09) {
      lamp_bit = 0x00;
      lamp_offset++;
      if (lamp_offset == 0x1B) {
        lamp_offset = 0x00;
      }
    }
}

