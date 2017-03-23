#!/usr/bin/env ruby

ARGF.each do |line|
  parts = line.split(" ")
  puts "{0x#{parts[2..5].join(", 0x")} }, //#{parts[6..-1].join(" ")}"
end