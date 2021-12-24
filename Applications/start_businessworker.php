<?php

/**
 * 注册业务
 */

use \Workerman\Worker;
use \GatewayWorker\BusinessWorker;

// 自动加载类
require_once __DIR__ . '/../vendor/autoload.php';

// bussinessWorker 进程
$worker = new BusinessWorker();
// worker名称
$worker->name = 'Server_Bussiness_Worker';
// bussinessWorker进程数量
$worker->count = 4;
// 服务注册地址
$worker->registerAddress = '127.0.0.1:1238';

// 如果不是在根目录启动，则运行runAll方法
if (!defined('GLOBAL_START')) {
    Worker::runAll();
}
