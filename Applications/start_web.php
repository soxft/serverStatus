<?php

/**
 * web界面
 */

use \Workerman\Worker;
use \Workerman\Connection\TcpConnection;
use \Workerman\Protocols\Http\Request;
use \Workerman\Protocols\Http\Response;
use \Class\Tool;

require_once __DIR__ . '/../vendor/autoload.php';
$http_worker = new Worker("http://0.0.0.0:2345");

$http_worker->count = 1;

define("WEBROOT", __DIR__ . "/web/");

$http_worker->onMessage = function (TcpConnection $connection, Request $request) {
    $path = $request->path();
    $request_method = $request->method(); // 请求方式 post,get etc
    $userIp = $request->header('x-real-ip'); //获取真实ip

    Tool::out("{$userIp} > {$request_method} > {$path}", 'WEB');

    $response = new Response(200);
    $response->withHeaders([
        'Content-Type' => 'text/html; charset=utf-8',
        'Access-Control-Allow-Methods' => 'GET'
    ]);

    if ($path == "/" || $path == "") $path = "index.php"; //默认文档
    if (file_exists(WEBROOT . $path)) {
        $response->withBody(Tool::exec_php_file(WEBROOT . $path));
    } else {
        $response->withStatus(404);
        $response->withBody("<h1>404 Not Found</h1>");
    }
    return $connection->send($response);
};

if (!defined('GLOBAL_START')) {
    Worker::runAll();
}
