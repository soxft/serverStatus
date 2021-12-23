<?php

/**
 * 注册服务
 */

use \Workerman\Worker;
use \GatewayWorker\Register;

require_once __DIR__ . '/../../vendor/autoload.php';

$register = new Register('text://0.0.0.0:1238');

if (!defined('GLOBAL_START')) {
    Worker::runAll();
}
