uint8_t lamps[26] = {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};

void myNoteOn(byte channel, byte note, byte velocity) {
  note_to_led(channel, note, true);
}
void myNoteOff(byte channel, byte note, byte velocity) {
  note_to_led(channel, note, false);
}

void note_to_led(uint8_t channel, uint8_t note, bool on) {
  uint8_t lamp_byte = 0xFF;
  uint8_t lamp_bit;
  if (channel == 3 && note == 0x05) {
    lamp_byte = 0;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x0D) {
    lamp_byte = 0;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x15) {
    lamp_byte = 0;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x1D) {
    lamp_byte = 0;
    lamp_bit = 64;
  }
  if (channel == 3 && note == 0x25) {
    lamp_byte = 1;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x2D) {
    lamp_byte = 1;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x35) {
    lamp_byte = 1;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x3D) {
    lamp_byte = 1;
    lamp_bit = 64;
  }
  if (channel == 3 && note == 0x06) {
    lamp_byte = 2;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x0E) {
    lamp_byte = 2;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x1E) {
    lamp_byte = 2;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x26) {
    lamp_byte = 3;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x2E) {
    lamp_byte = 3;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x36) {
    lamp_byte = 3;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x3E) {
    lamp_byte = 3;
    lamp_bit = 64;
  }
  if (channel == 3 && note == 0x07) {
    lamp_byte = 4;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x0F) {
    lamp_byte = 4;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x17) {
    lamp_byte = 4;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x1F) {
    lamp_byte = 4;
    lamp_bit = 64;
  }
  if (channel == 3 && note == 0x27) {
    lamp_byte = 5;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x00) {
    lamp_byte = 6;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x08) {
    lamp_byte = 6;
    lamp_bit = 2;
  }
  if (channel == 1 && note == 0x10) {
    lamp_byte = 6;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x18) {
    lamp_byte = 6;
    lamp_bit = 8;
  }
  if (channel == 1 && note == 0x20) {
    lamp_byte = 6;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x28) {
    lamp_byte = 6;
    lamp_bit = 32;
  }
  if (channel == 1 && note == 0x30) {
    lamp_byte = 6;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x38) {
    lamp_byte = 6;
    lamp_bit = 128;
  }
  if (channel == 1 && note == 0x01) {
    lamp_byte = 7;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x09) {
    lamp_byte = 7;
    lamp_bit = 2;
  }
  if (channel == 1 && note == 0x11) {
    lamp_byte = 7;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x19) {
    lamp_byte = 7;
    lamp_bit = 8;
  }
  if (channel == 1 && note == 0x21) {
    lamp_byte = 7;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x33) {
    lamp_byte = 7;
    lamp_bit = 32;
  }
  if (channel == 1 && note == 0x0C) {
    lamp_byte = 7;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x04) {
    lamp_byte = 7;
    lamp_bit = 128;
  }
  if (channel == 1 && note == 0x3B) {
    lamp_byte = 8;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x31) {
    lamp_byte = 8;
    lamp_bit = 2;
  }
  if (channel == 1 && note == 0x39) {
    lamp_byte = 8;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x02) {
    lamp_byte = 9;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x0A) {
    lamp_byte = 9;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x12) {
    lamp_byte = 9;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x1A) {
    lamp_byte = 9;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x22) {
    lamp_byte = 10;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x2A) {
    lamp_byte = 10;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x32) {
    lamp_byte = 10;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x3A) {
    lamp_byte = 10;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x03) {
    lamp_byte = 11;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x0B) {
    lamp_byte = 11;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x00) {
    lamp_byte = 12;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x08) {
    lamp_byte = 12;
    lamp_bit = 2;
  }
  if (channel == 3 && note == 0x10) {
    lamp_byte = 12;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x18) {
    lamp_byte = 12;
    lamp_bit = 8;
  }
  if (channel == 3 && note == 0x20) {
    lamp_byte = 12;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x28) {
    lamp_byte = 12;
    lamp_bit = 32;
  }
  if (channel == 3 && note == 0x30) {
    lamp_byte = 12;
    lamp_bit = 64;
  }
  if (channel == 3 && note == 0x38) {
    lamp_byte = 12;
    lamp_bit = 128;
  }
  if (channel == 3 && note == 0x01) {
    lamp_byte = 13;
    lamp_bit = 1;
  }
  if (channel == 3 && note == 0x09) {
    lamp_byte = 13;
    lamp_bit = 2;
  }
  if (channel == 3 && note == 0x11) {
    lamp_byte = 13;
    lamp_bit = 4;
  }
  if (channel == 3 && note == 0x19) {
    lamp_byte = 13;
    lamp_bit = 8;
  }
  if (channel == 3 && note == 0x21) {
    lamp_byte = 13;
    lamp_bit = 16;
  }
  if (channel == 3 && note == 0x29) {
    lamp_byte = 13;
    lamp_bit = 32;
  }
  if (channel == 3 && note == 0x31) {
    lamp_byte = 13;
    lamp_bit = 64;
  }
  if (channel == 3 && note == 0x39) {
    lamp_byte = 13;
    lamp_bit = 128;
  }
  if (channel == 1 && note == 0x2D) {
    lamp_byte = 14;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x25) {
    lamp_byte = 14;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x1D) {
    lamp_byte = 14;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x36) {
    lamp_byte = 14;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x05) {
    lamp_byte = 15;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x0D) {
    lamp_byte = 15;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x3C) {
    lamp_byte = 15;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x34) {
    lamp_byte = 16;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x2C) {
    lamp_byte = 16;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x24) {
    lamp_byte = 16;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x1C) {
    lamp_byte = 16;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x14) {
    lamp_byte = 17;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x26) {
    lamp_byte = 17;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x1E) {
    lamp_byte = 17;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x16) {
    lamp_byte = 17;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x0E) {
    lamp_byte = 18;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x3D) {
    lamp_byte = 18;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x35) {
    lamp_byte = 18;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x06) {
    lamp_byte = 18;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x07) {
    lamp_byte = 19;
    lamp_bit = 1;
  }
  if (channel == 1 && note == 0x0F) {
    lamp_byte = 19;
    lamp_bit = 2;
  }
  if (channel == 1 && note == 0x17) {
    lamp_byte = 19;
    lamp_bit = 4;
  }
  if (channel == 1 && note == 0x1F) {
    lamp_byte = 19;
    lamp_bit = 8;
  }
  if (channel == 1 && note == 0x27) {
    lamp_byte = 19;
    lamp_bit = 16;
  }
  if (channel == 1 && note == 0x2F) {
    lamp_byte = 19;
    lamp_bit = 32;
  }
  if (channel == 1 && note == 0x37) {
    lamp_byte = 19;
    lamp_bit = 64;
  }
  if (channel == 1 && note == 0x3F) {
    lamp_byte = 19;
    lamp_bit = 128;
  }
  if (channel == 2 && note == 0x07) {
    lamp_byte = 21;
    lamp_bit = 1;
  }
  if (channel == 2 && note == 0x05) {
    lamp_byte = 21;
    lamp_bit = 2;
  }
  if (channel == 2 && note == 0x0D) {
    lamp_byte = 21;
    lamp_bit = 4;
  }
  if (channel == 2 && note == 0x15) {
    lamp_byte = 21;
    lamp_bit = 8;
  }
  if (channel == 2 && note == 0x1D) {
    lamp_byte = 21;
    lamp_bit = 16;
  }
  if (channel == 2 && note == 0x25) {
    lamp_byte = 21;
    lamp_bit = 32;
  }
  if (channel == 2 && note == 0x2D) {
    lamp_byte = 21;
    lamp_bit = 64;
  }
  if (channel == 2 && note == 0x35) {
    lamp_byte = 21;
    lamp_bit = 128;
  }
  if (channel == 2 && note == 0x3D) {
    lamp_byte = 22;
    lamp_bit = 1;
  }
  if (channel == 2 && note == 0x0F) {
    lamp_byte = 22;
    lamp_bit = 2;
  }
  if (channel == 2 && note == 0x17) {
    lamp_byte = 22;
    lamp_bit = 4;
  }
  if (channel == 2 && note == 0x06) {
    lamp_byte = 22;
    lamp_bit = 8;
  }
  if (channel == 2 && note == 0x0E) {
    lamp_byte = 22;
    lamp_bit = 16;
  }
  if (channel == 2 && note == 0x16) {
    lamp_byte = 22;
    lamp_bit = 32;
  }
  if (channel == 2 && note == 0x1E) {
    lamp_byte = 22;
    lamp_bit = 64;
  }
  if (channel == 2 && note == 0x26) {
    lamp_byte = 22;
    lamp_bit = 128;
  }
  if (channel == 2 && note == 0x2E) {
    lamp_byte = 23;
    lamp_bit = 1;
  }
  if (channel == 2 && note == 0x36) {
    lamp_byte = 23;
    lamp_bit = 2;
  }
  if (channel == 2 && note == 0x3E) {
    lamp_byte = 23;
    lamp_bit = 4;
  }
  if (channel == 2 && note == 0x1F) {
    lamp_byte = 23;
    lamp_bit = 8;
  }
  if (channel == 2 && note == 0x3B) {
    lamp_byte = 24;
    lamp_bit = 1;
  }
  if (channel == 2 && note == 0x04) {
    lamp_byte = 24;
    lamp_bit = 2;
  }
  if (channel == 2 && note == 0x0C) {
    lamp_byte = 24;
    lamp_bit = 4;
  }
  if (channel == 2 && note == 0x14) {
    lamp_byte = 24;
    lamp_bit = 8;
  }
  if (channel == 2 && note == 0x1C) {
    lamp_byte = 24;
    lamp_bit = 16;
  }
  if (channel == 2 && note == 0x24) {
    lamp_byte = 24;
    lamp_bit = 32;
  }
  if (channel == 2 && note == 0x2C) {
    lamp_byte = 24;
    lamp_bit = 64;
  }
  if (channel == 2 && note == 0x34) {
    lamp_byte = 24;
    lamp_bit = 128;
  }
  if (channel == 2 && note == 0x33) {
    lamp_byte = 25;
    lamp_bit = 1;
  }
  if (channel == 2 && note == 0x23) {
    lamp_byte = 25;
    lamp_bit = 2;
  }
  if (channel == 2 && note == 0x2B) {
    lamp_byte = 25;
    lamp_bit = 4;
  }

  if (lamp_byte != 0xFF) {
    if (on) {
      lamps[lamp_byte] |= lamp_bit;
    } else {
      lamps[lamp_byte] &= ~lamp_bit;
    }
    Serial2.write(0x80);
    Serial2.write(lamp_byte);
    Serial2.write(lamps[lamp_byte]);
  }
}

