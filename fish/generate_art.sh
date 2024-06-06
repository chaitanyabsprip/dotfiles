#! /usr/bin/fish

set random_index (random 1 2)
set cond 1

if test $cond -eq $random_index
  colorscript -r
else
  fortune | cowsay
end
