#!/usr/bin/env ruby

File.open(ARGV.last) do |f|
  f.each_line do |l|
    parts = l.split
    puts "\t#{parts[0]}: \"#{parts[1]}\","
  end
end
