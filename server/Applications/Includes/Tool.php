<?php

namespace Includes;

class Tool
{
    public static function out(mixed $msg): void
    {
        print_r(date("Y-m-d H:i:s") . " ");
        print_r($msg);
        echo "\r\n";
    }

    public static function exec_php_file($file)
    {
        \ob_start();
        try {
            require $file;
        } catch (\Exception $e) {
            echo $e;
        }
        return \ob_get_clean();
    }
}
