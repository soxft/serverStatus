<?php

namespace Class;

class Tool
{
    public static function outPut($msg): void
    {
        print_r(date("Y-m-d H:i:s") . " " . $msg);
        echo "\r\n";
    }
}
