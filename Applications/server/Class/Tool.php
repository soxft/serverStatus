<?php

namespace Class;

class Tool
{
    public static function outPut($msg): void
    {
        echo date("Y-m-d H:i:s");
        print_r($msg) . PHP_EOL;
    }
}
