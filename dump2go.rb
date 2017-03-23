#!/usr/bin/env ruby

ARGF.each do |line|
  parts = line.split(" ")
  puts "{0x#{parts[2]}, 0x#{parts[3]}, 0x#{parts[4]}, 0x#{parts[5]} }, //#{parts[6..-1].join(" ")}"
end