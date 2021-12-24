<?php

/**
 * 注册GateWay
 */

use \Workerman\Worker;
use \GatewayWorker\Gateway;

require_once __DIR__ . '/../../vendor/autoload.php';

$gateway = new Gateway("websocket://127.0.0.1:8282");
$gateway->name = 'Server_GateWay';
// gateway进程数
$gateway->count = 2;
$gateway->lanIp = '127.0.0.1';
$gateway->startPort = 2900;
$gateway->registerAddress = '127.0.0.1:1238';

// 心跳间隔
//$gateway->pingInterval = 10;
// 心跳数据
//$gateway->pingData = '{"type":"ping"}';

// 如果不是在根目录启动，则运行runAll方法
if (!defined('GLOBAL_START')) {
    Worker::runAll();
}
